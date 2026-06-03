# Phase 5: Reviews, Notifications, CI/CD — Implementation Plan

## Overview

Phase 5 adalah tahap polish dan production-readiness. Fokus pada fitur yang membuat platform terasa
lengkap sebagai e-commerce: reviews, notifikasi email, Midtrans nyata, dan CI/CD pipeline.

---

## Part A: Product Reviews & Ratings

### A1 — Data Model

```
Tabel reviews:
├── id            UUID PK
├── product_id    UUID FK → products (index)
├── customer_id   UUID FK → customers (index)
├── rating        INT (1-5, CHECK constraint)
├── comment       TEXT
├── is_approved   BOOLEAN (admin moderation)
├── unique(product_id, customer_id) — satu review per customer per produk
└── timestamps
```

### A2 — Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/products/:slug/reviews` | Public | List approved reviews (pagination) |
| `POST` | `/customer/reviews` | Customer JWT | Create review (only if ordered) |
| `GET` | `/admin/reviews` | Admin | List all reviews (filter pending) |
| `PATCH` | `/admin/reviews/:id/approve` | Admin | Toggle approval |

### A3 — Frontend

- Product detail page: star rating display + review list
- "Tulis Review" form di order detail page (setelah delivered)
- Admin: moderation queue di sidebar

---

## Part B: Midtrans Real Driver

### B1 — Midtrans Client

**Package:** `pkg/payment/midtrans.go` — implementasi interface `payment.Driver`

```go
type MidtransDriver struct {
    client snap.Client
}

func (d *MidtransDriver) CreateTransaction(req CreateRequest) (*Response, error) {
    snapReq := &snap.Request{
        TransactionDetails: midtrans.TransactionDetails{
            OrderID:  req.OrderID,
            GrossAmt: int64(req.GrossAmount),
        },
        CustomerDetail: &midtrans.CustomerDetail{
            FName: req.CustomerName,
            Email: req.CustomerEmail,
            Phone: req.CustomerPhone,
        },
        Items: convertItems(req.Items),
    }
    resp, err := d.client.CreateTransaction(snapReq)
    return &Response{
        TransactionID: resp.OrderID,
        PaymentToken:  resp.Token,
        RedirectURL:   resp.RedirectURL,
        PaymentURL:    resp.RedirectURL,
    }, nil
}
```

### B2 — Webhook Handler

- `POST /api/v1/payment/callback` — Midtrans POST notification
- Validasi signature dengan server key
- Update `orders.payment_status` → settlement/cancel/expire
- Update `orders.status` → confirmed (jika settlement)
- Log seluruh callback body

### B3 — Config

```env
PAYMENT_DRIVER=midtrans
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxx
MIDTRANS_MERCHANT_ID=G123456
MIDTRANS_IS_PRODUCTION=false
```

---

## Part C: Email Notifications (SMTP)

### C1 — Email Service

**Package:** `pkg/email/email.go`

```go
type Client interface {
    Send(to, subject, htmlBody string) error
}
```

- SMTP driver via `net/smtp` (tanpa dependency eksternal)
- HTML templates via `html/template`
- Debug driver (log to console, like current reset token)

### C2 — Trigger Points

| Event | Email |
|-------|-------|
| Order status change | "Pesanan Anda [Confirmed/Shipped/Delivered]" |
| Password reset | Link reset (currently console log) |
| Payment confirmation | "Pembayaran Diterima" |

### C3 — Templates

```
backend/templates/
├── order-confirmed.html
├── order-shipped.html
├── order-delivered.html
├── password-reset.html
├── payment-received.html
```

---

## Part D: Customer Avatar

### D1 — Backend

- Field `avatar_url` di `customers` table
- Endpoint `POST /customer/me/avatar` (multipart, crop to 200x200)
- Serve via static `/api/v1/uploads/avatars/`

### D2 — Frontend

- Avatar display di nav header (ganti inisial)
- Upload di halaman profil

---

## Part E: CI/CD Pipeline

### E1 — GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  backend:
    - go test ./...
    - go build ./cmd/api
  
  frontend:
    - npm run build (frontend & admin)
  
  deploy:
    - docker build & push to registry
    - SSH deploy to VPS (optional)
```

### E2 — Docker Production Setup

- Multi-stage Dockerfile untuk Go (size ~20MB)
- docker-compose.prod.yml dengan Nginx reverse proxy
- Let's Encrypt SSL via certbot

---

## Part F: Essential Polish

### F1 — Rate Limiting (Redis)

```go
// middleware/rate_limiter.go
func RateLimit(redis *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
    // Sliding window rate limiter per IP
    // Apply to: login, register, forgot-password
}
```

### F2 — Product Variants (stretch)

- `product_variants` table: size, color, price_adjustment, stock
- Variant selector di product detail
- Update cart/checkout untuk variant-based stock

### F3 — Coupon/Discount System (stretch)

- `coupons` table: code, type (fixed/percentage), value, min_order, expiry
- Apply coupon di checkout
- Track usage per customer

---

## Implementation Order

```
Step 1:  A — Product Reviews (model, endpoints, frontend)
Step 2:  B — Midtrans Real Driver
Step 3:  C — Email SMTP + Templates
Step 4:  D — Customer Avatar
Step 5:  E — CI/CD Pipeline
Step 6:  F1 — Rate Limiting
Step 7:  F2 — Product Variants (optional)
Step 8:  F3 — Coupon System (optional)
```

---

## Files Summary

### New Backend (~30 files)
```
backend/internal/model/product/review.go
backend/internal/dto/product/review_dto.go
backend/internal/repository/product/review_repository.go
backend/internal/service/product/review_service.go
backend/internal/handler/product/review_handler.go

backend/pkg/payment/midtrans.go              — Midtrans Snap driver
backend/internal/handler/payment/webhook.go  — Midtrans callback handler

backend/pkg/email/email.go                   — Email client interface
backend/pkg/email/smtp.go                    — SMTP driver
backend/pkg/email/debug.go                   — Console logger driver
backend/templates/*.html                     — Email templates
backend/internal/service/notification/service.go

backend/internal/handler/user/avatar_handler.go
backend/internal/middleware/rate_limiter.go

backend/internal/model/product/variant.go    — (stretch)
backend/internal/model/order/coupon.go       — (stretch)

backend/migrations/000013_create_reviews.up.sql
backend/migrations/000014_add_avatar_url.up.sql
```

### New Frontend (~10 files)
```
frontend/src/components/ReviewStars.tsx
frontend/src/components/ReviewForm.tsx
frontend/src/components/ReviewList.tsx
admin/src/pages/ReviewModeration.tsx

frontend/src/components/PaymentButton.tsx    — Midtrans Snap popup
frontend/src/pages/Profile.tsx               — Update (avatar upload)
```

---

## Dependencies to Add

```
# Go
github.com/midtrans/midtrans-go    — Midtrans SDK (real)
golang.org/x/time/rate             — Rate limiting

# Env additions
PAYMENT_DRIVER=midtrans
MIDTRANS_SERVER_KEY=SB-Mid-server-xxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxx
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@petshop.com
SMTP_PASSWORD=app-password
SMTP_FROM=noreply@petshop.com
APP_URL=http://localhost:3000
```

---

## Estimasi Total: ~40 files, ~2-3x effort dari Phase 4
