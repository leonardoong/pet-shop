# Phase 3: Admin Dashboard — Implementation Plan

## Overview

Phase 3 makes the admin dashboard fully functional. Currently the admin has login + a dashboard
with hardcoded "---" placeholders, and every management page is `<div>coming in Phase 3</div>`.
There are **zero admin CRUD APIs** on the backend — no way to manage products, categories,
orders, or customers without connecting directly to the DB.

---

## Part A: Backend Admin APIs

### A1 — Product Management (products:*)

**Routes** (registered under `/api/v1/admin` in server.go):

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET`    | `/admin/products`          | `products:read`   | List with pagination, search, category filter, active/inactive filter |
| `POST`   | `/admin/products`          | `products:create` | Create new product |
| `GET`    | `/admin/products/:id`      | `products:read`   | Product detail |
| `PUT`    | `/admin/products/:id`      | `products:update` | Update product (fields partial) |
| `DELETE` | `/admin/products/:id`      | `products:delete` | Soft-delete or deactivate |

**New files:**
- `backend/internal/handler/product/admin_handler.go` — Swagger-annotated handlers
- `backend/internal/service/product/admin_service.go` — Business logic (or extend existing)
- `backend/internal/dto/product/admin_dto.go` — Create/Update request + admin response types

**DTOs:**

```go
// CreateProductRequest
type CreateProductRequest struct {
    Name        string `json:"name" validate:"required,min=3,max=200"`
    CategoryID  string `json:"category_id" validate:"required,uuid"`
    Description string `json:"description" validate:"max=2000"`
    Price       float64 `json:"price" validate:"required,min=0"`
    Stock       int     `json:"stock" validate:"required,min=0"`
    SKU         string `json:"sku" validate:"required,min=3,max=50"`
    ImageURL    string `json:"image_url" validate:"omitempty,url"`
    IsActive    bool   `json:"is_active"`
}

// UpdateProductRequest — same fields but all optional, omitempty
// AdminProductResponse — same as public ProductResponse + created_at, updated_at, stock, is_active
```

**Logic details:**
- Create: validate category exists, auto-generate slug from name (append suffix if conflict), set default is_active=true
- Update: partial update (only fields present in request), re-slug if name changes
- Delete: set `is_active = false` (soft-delete, products stay in DB for order history)
- List: paginated (page, limit), sort (name, price, stock, created_at), filter (category_id, is_active, search by name/sku)

---

### A2 — Category Management (categories:*)

**Routes:**

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET`    | `/admin/categories`        | `categories:read`   | List all (optionally with product count) |
| `POST`   | `/admin/categories`        | `categories:create` | Create category |
| `PUT`    | `/admin/categories/:id`    | `categories:update` | Update category |
| `DELETE` | `/admin/categories/:id`    | `categories:delete` | Delete (only if 0 products) |

**New files:**
- `backend/internal/handler/product/admin_handler.go` — (reuse same handler file for categories)
- `backend/internal/service/product/admin_category_service.go`
- `backend/internal/dto/product/admin_dto.go` — Category request/response types

**Logic details:**
- Delete guard: prevent deleting categories that have products
- Auto-generate slug from name on create/update

---

### A3 — Order Management (orders:*)

**Routes:**

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET`    | `/admin/orders`            | `orders:read`   | List all orders (pagination, filters) |
| `GET`    | `/admin/orders/:id`        | `orders:read`   | Order detail |
| `PATCH`  | `/admin/orders/:id/status` | `orders:update` | Update order status |

**New files:**
- `backend/internal/handler/order/admin_handler.go`
- `backend/internal/service/order/admin_service.go`
- `backend/internal/dto/order/admin_dto.go`

**Logic details:**
- List: paginated, filter by status, customer email, date range (from/to)
- Status update: validate valid transitions (`pending → confirmed → processing → shipped → delivered`; any status → `cancelled`)
- Add `status_note` field on status update for audit trail

---

### A4 — Inventory Management (inventory:*)

**Routes:**

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET`    | `/admin/inventory`              | `inventory:read`   | Product stock list |
| `PATCH`  | `/admin/inventory/:productId`   | `inventory:update` | Adjust stock |

**New files:**
- `backend/internal/handler/product/admin_inventory_handler.go`
- `backend/internal/service/product/inventory_service.go`
- `backend/internal/dto/product/admin_dto.go` — Inventory response types

**Logic details:**
- List: paginated, filter low-stock (stock <= threshold), search by name/sku
- Adjust stock: `operation = "add" | "subtract" | "set"`, `quantity`, `note` (audit)

---

### A5 — Customer Management (customers:*)

**Routes:**

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET`    | `/admin/customers`         | `customers:read`   | List customers |
| `GET`    | `/admin/customers/:id`     | `customers:read`   | Customer detail |
| `PATCH`  | `/admin/customers/:id`     | `customers:update` | Toggle active/inactive |

**New files:**
- `backend/internal/handler/user/admin_handler.go`
- `backend/internal/service/user/admin_service.go`
- `backend/internal/dto/user/admin_dto.go`

**Logic details:**
- List: paginated, search by name/email/phone
- Detail: includes order count, total spent
- Toggle: only active ↔ inactive (never delete)

---

### A6 — Dashboard Analytics (dashboard:read)

**Route:**

| Method | Path | Permission | Description |
|--------|------|------------|-------------|
| `GET` | `/admin/dashboard` | `dashboard:read` | Aggregated stats |

**Response shape:**

```json
{
  "total_revenue": 52500000,
  "total_orders": 142,
  "total_products": 87,
  "total_customers": 356,
  "orders_by_status": {
    "pending": 5, "confirmed": 12, "processing": 8,
    "shipped": 23, "delivered": 89, "cancelled": 5
  },
  "recent_orders": [ /* last 5 */ ],
  "low_stock_products": [ /* stock <= 10, top 5 */ ],
  "revenue_by_month": [ /* last 6 months, for chart */ ]
}
```

**New files:**
- `backend/internal/handler/dashboard/handler.go` ** ← NEW DOMAIN `dashboard`
- `backend/internal/service/dashboard/service.go`
- `backend/internal/repository/dashboard/repository.go`
- `backend/internal/dto/dashboard/dto.go`

**Logic details:**
- Revenue: `SUM(order_items.subtotal)` from `orders` with status not `cancelled`
- Revenue by month: group by month for last 6 months
- Low stock: products WHERE stock <= 10, ordered by stock ASC

---

## Part B: Admin Frontend Pages

### Shared Components (new)
- `admin/src/components/DataTable.tsx` — Reusable TanStack Table with search, pagination, column toggle
- `admin/src/components/ConfirmDialog.tsx` — Reusable delete confirmation modal
- `admin/src/components/StatusBadge.tsx` — Colored order status badges
- `admin/src/components/PageHeader.tsx` — Breadcrumb + title + action button pattern

### API Modules (new)
- `admin/src/api/products.ts` — CRUD + list with filters
- `admin/src/api/categories.ts` — CRUD + list
- `admin/src/api/orders.ts` — List, detail, updateStatus
- `admin/src/api/inventory.ts` — List, adjustStock
- `admin/src/api/customers.ts` — List, detail, toggleActive
- `admin/src/api/dashboard.ts` — Get stats

### B1 — Dashboard Page (real data)

Replace hardcoded "---" with live API data:
- 4 KPI cards: total revenue (formatted IDR), total orders, total products, total customers
- Revenue chart (Recharts AreaChart) — last 6 months
- Order status bar chart — count by status
- Recent orders table (last 5, with status badges, clickable)
- Low stock alerts section (products with stock ≤ 10)

**Files to modify:**
- `admin/src/pages/Dashboard.tsx`

---

### B2 — Products Page

Full product management:
- DataTable: name, SKU, category, price, stock, status (active/inactive), actions
- Search bar (name/SKU), filter by category, filter by status (active/inactive)
- "Tambah Produk" button → slide-out drawer with CreateProduct form (React Hook Form + Zod)
- Row actions: Edit (opens drawer pre-filled), Toggle active, Delete (confirm dialog)
- Pagination

**New files:**
- `admin/src/pages/Products.tsx`
- `admin/src/components/ProductForm.tsx` (reused for create + edit)

---

### B3 — Categories Page

Simple CRUD:
- Table: icon, name, slug, product count, actions
- "Tambah Kategori" button → modal with name + image URL fields
- Inline edit (double-click name to edit)
- Delete with confirm dialog (show error if category has products)
- No pagination needed (categories are few)

**New files:**
- `admin/src/pages/Categories.tsx`

---

### B4 — Orders Page

Order management:
- DataTable: order ID, customer name, date, total, status badge, actions
- Filters: status dropdown, date range, search by customer/order ID
- Row actions: View detail (expandable row or drawer), Update status (dropdown)
- Status update with confirmation (e.g., "Konfirmasi Pesanan?")
- Order detail drawer: customer info, shipping address, items table, totals

**New files:**
- `admin/src/pages/Orders.tsx`

---

### B5 — Inventory Page

Stock management:
- Table: product name, SKU, category, current stock, status (in-stock / low / out)
- Low-stock rows highlighted (amber) / out-of-stock (red)
- Search by name/SKU, filter by category
- "Sesuaikan Stok" button → modal with operation (add/subtract/set), quantity, reason
- Stock history log (stretch goal — probably just the adjustment for now)

**New files:**
- `admin/src/pages/Inventory.tsx`

---

### B6 — Customers Page

Customer list:
- Table: name, email, phone, orders count, total spent, status (active/inactive), joined date
- Search by name/email
- Click row → detail drawer with: contact info, addresses (from addresses table), recent orders, totals
- Toggle active/inactive (confirm dialog)

**New files:**
- `admin/src/pages/Customers.tsx`

---

### B7 — Roles & Permissions (stretch goal within Phase 3)

View-only page showing existing roles and their permissions:
- Role cards with permission chips
- Read-only (no CRUD for now — use DB seed + future create UI)
- Could be part of Settings page

---

## Part C: Customer Frontend Bug Fixes

Minor fixes that improve UX:

1. **Fix refresh token**: In `Login.tsx` and `Register.tsx`, the `setAuth` call passes `access_token` as the refresh token. Backend returns a refresh token — capture and store it properly.
2. **Dynamic nav categories**: Replace hardcoded "Anjing"/"Kucing" links in `CustomerLayout.tsx` with fetched categories.
3. **Handle dead links**: Either create `/lupa-password`, `/tentang`, `/kontak` placeholder pages or remove the links temporarily.
4. **Cart stock feedback**: Show error toast when server rejects add-to-cart due to insufficient stock.

---

## Implementation Order

```
Step 1 (Backend): A1 Product Management
Step 2 (Backend): A2 Category Management
Step 3 (Backend): A3 Order Management
Step 4 (Backend): A4 Inventory Management
Step 5 (Backend): A5 Customer Management
Step 6 (Backend): A6 Dashboard Analytics
Step 7 (Frontend): Shared components + API modules
Step 8 (Frontend): B2 Products page
Step 9 (Frontend): B3 Categories page
Step 10 (Frontend): B4 Orders page
Step 11 (Frontend): B5 Inventory page
Step 12 (Frontend): B6 Customers page
Step 13 (Frontend): B1 Dashboard (real data)
Step 14 (Bug Fixes): Part C
Step 15 (Polish): Route registration, Swagger, testing
```

---

## Files To Create/Modify Summary

### New Backend Files (~18 files)
```
backend/internal/handler/product/admin_handler.go
backend/internal/handler/product/admin_inventory_handler.go
backend/internal/handler/order/admin_handler.go
backend/internal/handler/user/admin_handler.go
backend/internal/handler/dashboard/handler.go
backend/internal/service/product/admin_product_service.go
backend/internal/service/product/admin_category_service.go
backend/internal/service/product/inventory_service.go
backend/internal/service/order/admin_service.go
backend/internal/service/user/admin_service.go
backend/internal/service/dashboard/service.go
backend/internal/repository/dashboard/repository.go
backend/internal/dto/product/admin_dto.go
backend/internal/dto/order/admin_dto.go
backend/internal/dto/user/admin_dto.go
backend/internal/dto/dashboard/dto.go
```

### New Admin Frontend Files (~12 files)
```
admin/src/api/products.ts
admin/src/api/categories.ts
admin/src/api/orders.ts
admin/src/api/inventory.ts
admin/src/api/customers.ts
admin/src/api/dashboard.ts
admin/src/pages/Products.tsx
admin/src/pages/Categories.tsx
admin/src/pages/Orders.tsx
admin/src/pages/Inventory.tsx
admin/src/pages/Customers.tsx
admin/src/components/DataTable.tsx
admin/src/components/ConfirmDialog.tsx
admin/src/components/StatusBadge.tsx
admin/src/components/ProductForm.tsx
```

### Modified Files (~8 files)
```
backend/internal/server/server.go          — Register admin routes + RBAC
admin/src/routes/index.tsx                 — Replace placeholders with real pages
admin/src/pages/Dashboard.tsx              — Real data
frontend/src/pages/Login.tsx               — Fix refresh token
frontend/src/pages/Register.tsx            — Fix refresh token
frontend/src/layouts/CustomerLayout.tsx    — Dynamic categories
frontend/src/store/authStore.ts            — Fix token storage
```
