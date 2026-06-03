# Phase 4: Payment, Profiles & Images — Implementation Plan

## Overview

Phase 4 adds the missing pieces that turn the pet shop from a catalog+admin system into a
functional e-commerce platform that can accept real money. It also fills UX gaps (profile editing,
password reset, image uploads) that make the platform feel incomplete.

---

## Part A: Payment Integration (Midtrans)

### Why Midtrans
- Dominant Indonesian payment gateway (GoPay, OVO, bank transfer, credit card)
- Well-documented API, official Go SDK (`github.com/midtrans/midtrans-go`)
- Sandbox environment for testing
- Webhook callbacks for payment status updates

### A1 — Data Model Changes

```
New table: payments
├── id            UUID PK
├── order_id      UUID FK → orders (unique, 1:1)
├── external_id   VARCHAR (Midtrans order_id)
├── amount        NUMERIC(12,2)
├── method        VARCHAR (payment method chosen)
├── status        ENUM: pending, settlement, expired, cancel, deny, refund
├── paid_at       TIMESTAMPTZ
├── raw_callback  JSONB (store full callback for audit)
└── timestamps
```

Modify `orders`:
- Add optional fields: `payment_status`, `payment_url`

### A2 — Midtrans Integration Service

**Package:** `pkg/midtrans/`

```go
type Client interface {
    CreateTransaction(order *OrderPayload) (*TransactionResponse, error)
    VerifySignature(orderID, statusCode, grossAmount, serverKey string, rawBody []byte) bool
}

type OrderPayload struct {
    OrderID     string
    GrossAmount int64
    Customer    CustomerInfo
    Items       []ItemDetail
}
```

**Flow:**
1. Customer submits checkout → order created with status `pending`
2. Backend calls Midtrans Snap API → gets `redirect_url` + `token`
3. Frontend opens Midtrans Snap popup → customer pays
4. Midtrans sends webhook to `/api/v1/payment/callback`
5. Backend validates signature → updates `payments` + `orders` status to `confirmed`

### A3 — New Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/customer/checkout` | Replaces current checkout: creates order + Midtrans transaction |
| `GET`  | `/customer/orders/:id/pay` | Returns payment URL/token for unpaid order |
| `POST` | `/payment/callback` | Midtrans webhook (public, signature-verified) |
| `GET`  | `/admin/payments` | Admin view of all payments |

**Modified `POST /customer/orders` (checkout):**
- Keep existing validation (address, cart, stock)
- After order creation → call Midtrans → return `{order, payment_url}`

### A4 — Frontend Changes

**Customer side:**
- Checkout page: after submit, redirect to Midtrans Snap or show payment button
- Order detail page: show payment status, "Bayar Sekarang" button for pending orders
- Order history: payment status badges

**Admin side:**
- Orders page: payment status column
- New Payments tab or filter showing payment statuses

---

## Part B: Customer Profile Management

### B1 — New Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/customer/me` | Get profile (currently missing from routes) |
| `PUT`  | `/customer/me` | Update name, phone |
| `PUT`  | `/customer/me/email` | Change email (requires password verification) |
| `PUT`  | `/customer/me/password` | Change password (requires old password) |

### B2 — Frontend: Profile Page

**Route:** `/akun/profil`

- Editable fields: full name, phone number
- Email change section (with current password)
- Password change section (old + new + confirm)
- Avatar placeholder (future Phase 5)

**Files:**
- `frontend/src/pages/Profile.tsx`
- `frontend/src/api/profile.ts`

---

## Part C: Password Reset Flow

### C1 — Token-based Reset

```
New table: password_resets
├── id          UUID PK
├── email       VARCHAR
├── token_hash  VARCHAR UNIQUE
├── user_type   ENUM: customer, admin
├── expires_at  TIMESTAMPTZ
├── used        BOOLEAN
└── timestamps
```

### C2 — Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/auth/forgot-password` | Send reset link (email) |
| `POST` | `/auth/reset-password` | Reset password with token |

### C3 — Frontend

- Login page: functional "Lupa password?" link
- `/lupa-password` page: email input form
- `/reset-password?token=xxx` page: new password form

---

## Part D: Product Image Upload

### D1 — Backend

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/admin/upload` | Upload image (multipart) → return URL |
| `POST` | `/admin/upload/bulk` | Multiple images |

- Store files in `uploads/` directory (dev) or S3/MinIO (prod)
- Configurable storage backend via `UPLOAD_DRIVER` env var
- File validation: max 5MB, JPEG/PNG/WebP only
- Return publicly accessible URL

### D2 — Frontend

- ProductForm: replace image_url text input with file picker + preview
- Upload on file select, then fill image_url automatically

---

## Part E: Email Notifications (lite)

Basic transactional emails without heavy infrastructure:

- Use Go's `net/smtp` or a small lib like `gomail`
- Templates: Go `html/template`
- Events:
  - Order status changes → email customer
  - Password reset → email reset link

**Config:**
```env
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=587
SMTP_USER=xxx
SMTP_PASSWORD=xxx
SMTP_FROM=noreply@petshop.com
```

---

## Implementation Order

```
Step 1:  A2 — Midtrans Go client (pkg/midtrans/)
Step 2:  A1 — Payment migration + model
Step 3:  A3 — Payment endpoints + webhook handler
Step 4:  A4 — Frontend: Checkout flow with payment
Step 5:  B  — Customer profile endpoints + page
Step 6:  C  — Password reset flow
Step 7:  D  — File upload for product images
Step 8:  E  — Email notifications (order status, password reset)
```

---

## Files Summary

### New Backend (~25 files)
```
backend/pkg/midtrans/client.go          — Midtrans API wrapper
backend/pkg/midtrans/config.go          — Midtrans env config
backend/internal/model/payment/model.go — Payment GORM model
backend/internal/dto/payment/dto.go     — Payment request/response types
backend/internal/repository/payment/repository.go
backend/internal/service/payment/service.go
backend/internal/handler/payment/handler.go
backend/internal/handler/payment/webhook.go
backend/internal/dto/user/profile_dto.go   — Profile update DTOs
backend/internal/service/user/profile_service.go
backend/internal/handler/user/profile_handler.go
backend/internal/model/user/password_reset.go
backend/internal/dto/auth/reset_dto.go
backend/internal/service/auth/reset_service.go
backend/internal/handler/auth/reset_handler.go
backend/internal/handler/upload/handler.go — File upload handler
backend/pkg/upload/upload.go               — File storage abstraction
backend/internal/service/notification/email.go
backend/migrations/000010_create_payments.up.sql
backend/migrations/000011_create_password_resets.up.sql
backend/migrations/xxxx_down.sql counterparts
```

### New Frontend (~8 files)
```
frontend/src/pages/Profile.tsx
frontend/src/pages/ForgotPassword.tsx
frontend/src/pages/ResetPassword.tsx
frontend/src/api/profile.ts
frontend/src/api/payment.ts
admin/src/components/FileUpload.tsx       — Drag-drop image upload
frontend/src/components/PaymentButton.tsx — Midtrans Snap wrapper
```

### Modified Files (~12 files)
```
backend/internal/server/server.go         — Register new routes
backend/internal/model/order/model.go     — Add payment_status, payment_url
backend/internal/service/order/service.go — Integrate payment creation
backend/migrations/000009_seed_rbac.up.sql — Add customers:update (already done)
frontend/src/pages/Login.tsx              — Enable lupa password link
frontend/src/pages/Checkout.tsx           — Payment flow integration
frontend/src/pages/OrderDetail.tsx        — Payment status + pay button
frontend/src/routes/index.tsx             — New routes
admin/src/components/ProductForm.tsx      — File upload instead of URL input
admin/src/pages/Orders.tsx                — Payment status column
admin/src/routes/index.tsx                — Payments route (optional)
backend/go.mod                             — midtrans-go dependency
```

---

## Dependencies to Add

```
# Go
github.com/midtrans/midtrans-go    — Midtrans SDK
github.com/go-gomail/gomail         — Email sending (or net/smtp)

# Env
MIDTRANS_SERVER_KEY=SB-Mid-server-xxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxx
MIDTRANS_MERCHANT_ID=G123456
MIDTRANS_IS_PRODUCTION=false
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=587
SMTP_USER=xxx
SMTP_PASSWORD=xxx
SMTP_FROM=noreply@petshop.com
UPLOAD_DIR=./uploads
UPLOAD_MAX_SIZE=5242880
```

---

## What Phase 4 Does NOT Include (deferred to Phase 5)

- Product reviews & ratings
- Advanced email templates (welcome, promotion, abandoned cart)
- CI/CD pipeline (GitHub Actions)
- Rate limiting (Redis)
- Admin role/permission CRUD UI
- Stock history audit log
- Analytics export (CSV/Excel)
- Product variants (size, color)
- Wishlist
- Coupon/discount system
- SEO meta tags
- Sitemap
