## Project Overview

Build a production-ready full-stack Pet Shop platform consisting of:

1. Customer-facing E-Commerce Website
2. Admin Dashboard
3. REST API Backend
4. Authentication & Authorization
5. Inventory & Order Management
6. Responsive UI for desktop & mobile
7. Payment Integration

---

# Tech Stack

## Frontend

### Customer Website
- React
- Vite
- TypeScript
- React Router
- TanStack Query (React Query)
- Zustand (state management)
- TailwindCSS
- Shadcn/UI
- Axios
- React Hook Form + Zod
- Framer Motion

### Admin Dashboard
- React
- TypeScript
- TailwindCSS
- Shadcn/UI
- TanStack Table
- Recharts
- React Hook Form
- Zod

---

## Backend

### Core
- Go (latest stable version)
- Gin Framework
- PostgreSQL
- Redis
- JWT Authentication
- GORM

### Infrastructure
- Docker
- Docker Compose
- Nginx
- GitHub Actions CI/CD

---

# Architecture

Clean architecture with layer-first, domain-second organization.

## Backend Structure

```
backend/
├── cmd/api/            # Entry point (main.go)
├── internal/
│   ├── model/          # GORM structs — no business logic
│   │   ├── user/       # Customer, Admin, Role, Permission, Address, RefreshToken
│   │   ├── product/    # Category, Product
│   │   ├── cart/       # Cart, CartItem
│   │   └── order/      # Order, OrderItem
│   ├── dto/            # Request/response types — no DB tags
│   │   ├── auth/
│   │   ├── user/       # Address request types
│   │   ├── product/    # Filter
│   │   ├── cart/
│   │   └── order/
│   ├── repository/     # DB operations only — no business rules
│   │   ├── auth/
│   │   ├── user/       # Address persistence
│   │   ├── product/    # Category + Product (shared domain)
│   │   ├── cart/
│   │   └── order/
│   ├── service/        # Domain logic — testable, no infra coupling
│   │   ├── auth/
│   │   ├── user/       # Address service
│   │   ├── product/    # CategoryService + ProductService
│   │   ├── cart/
│   │   └── order/
│   ├── handler/        # HTTP layer — parse, call service, respond
│   │   ├── auth/
│   │   ├── user/       # Address handler
│   │   ├── product/    # CategoryHandler + ProductHandler
│   │   ├── cart/
│   │   └── order/
│   ├── middleware/     # Auth (RequireCustomer, RequireAdmin), RBAC
│   └── server/         # Route registration, CORS, Swagger
├── pkg/
│   ├── config/         # Env loading
│   ├── jwt/            # JWT Manager (customer + admin issuers)
│   ├── response/       # Standard JSON envelope helpers
│   └── validator/      # Custom Gin validators (indonesian_phone)
├── migrations/         # golang-migrate SQL files (up + down)
└── docs/               # Auto-generated Swagger (swag init)
```

## Import Alias Convention

When the same domain name appears across multiple layers, use short aliases:

```go
// Within a service file, importing own-domain layers:
model  "petshop/internal/model/product"
repo   "petshop/internal/repository/product"
dto    "petshop/internal/dto/product"

// Cross-domain imports use the domain as prefix:
cartrepo    "petshop/internal/repository/cart"
userrepo    "petshop/internal/repository/user"
productrepo "petshop/internal/repository/product"
```

## Domain Grouping

| Domain folder | Covers |
|---|---|
| `auth` | Customer + admin authentication, token lifecycle |
| `user` | Customer profile, admin accounts, shipping addresses |
| `product` | Categories and products (shared layer folder) |
| `cart` | Cart and cart items |
| `order` | Orders and order items, checkout flow |
| `payment` | *(future)* |
| `inventory` | *(future)* |
| `dashboard` | *(future)* |

## Frontend Structure

```
frontend/
├── src/
│   ├── api/            # Axios instances + per-domain API functions
│   ├── components/     # Shared UI components
│   ├── pages/          # Route-level page components
│   ├── hooks/          # Custom React hooks
│   ├── store/          # Zustand stores
│   ├── layouts/        # CustomerLayout, AdminLayout
│   ├── types/          # TypeScript interfaces matching backend models
│   ├── lib/            # Utilities, helpers
│   └── routes/         # React Router configuration
```

---

## Backend Layer Rules

### Handlers
- Only parse requests, call services, return responses
- No business logic — ever

### Services
- Own all domain logic
- Testable with no infrastructure coupling
- Define their own error sentinel vars (e.g. `ErrProductNotFound`)

### Repositories
- Only handle DB operations
- No business rules
- Interfaces defined in the repository package; services import them directly
- Sentinel errors that arise from DB operations live here (e.g. `ErrInsufficientStock`)

---

## Database Restrictions

### NEVER
- use SELECT *
- skip indexes on searchable fields
- fetch unnecessary columns
- create N+1 query problems
- run expensive queries in loops

### ALWAYS
- paginate list endpoints
- use transactions where needed
- add created_at and updated_at
- validate migrations

---

# API Standards

## Response Format

Always use consistent response structure:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Error format:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": []
}
```

---

# Authentication Rules

Claude must:
- use JWT expiration
- implement refresh token rotation
- hash passwords with bcrypt/argon2
- secure admin routes
- validate authorization on every protected endpoint

Never trust frontend roles.

---

# User & Admin Separation Rules

## Separate Tables

- `customers` table: stores customer accounts only
- `admins` table: stores admin/staff accounts only
- These tables are NEVER merged — do not use a single `users` table with a role column
- Each has its own JWT issuer claim (`iss: customer` vs `iss: admin`) to prevent token cross-use
- Each has its own refresh token scope — a customer token cannot authenticate as admin

## RBAC for Admin Dashboard

- Implement Role-Based Access Control (RBAC) for all admin routes
- `roles` table: named roles (e.g. `super_admin`, `manager`, `inventory_staff`, `support`)
- `permissions` table: granular actions (e.g. `products:create`, `orders:update`, `inventory:read`)
- `role_permissions` table: many-to-many join between roles and permissions
- `admin_roles` table: many-to-many join between admins and roles
- An admin can have multiple roles; a role can have multiple permissions
- Authorization middleware must check the admin's effective permissions on every protected endpoint
- Never derive permissions solely from the frontend — always verify server-side
- New admin roles can be added via the admin dashboard without code changes

---

# Validation Rules

Validate:
- request body
- query params
- path params
- uploaded files

Validation must happen before business logic execution.

---

# Documentation Rules

Claude should generate:
- README
- API documentation
- environment setup guide
- architecture notes

Important business logic should be documented.

---

# UI/UX Rules

Customer-facing pages must:
- load fast
- be mobile responsive
- provide clear feedback states
- have accessible forms

Admin dashboard must:
- prioritize data readability
- support filtering/search
- support pagination

---

# Customer Frontend Language & UI Rules

## Language
- All customer-facing UI text must be in **Bahasa Indonesia**
- This includes: navigation, labels, buttons, placeholders, error messages, empty states, and any copy
- Admin dashboard remains in English

## Color Theme
- Primary color: **warm light green** (hijau muda warm)
- Use a green palette with warm/yellow undertones — not cold/blue greens
- Tailwind custom color name: `primary` (mapped to the warm green scale)
- Accent warmth via cream, soft amber, or warm white backgrounds
- Avoid harsh contrasts; prefer soft, friendly tones that feel welcoming for a pet shop

## Home Page Layout
- Home page is a **catalog-first layout**, NOT a marketing hero-only page
- Structure (top to bottom):
  1. Promotional banner / hero (slim, not full-screen)
  2. Category quick-links (icon + label grid)
  3. Featured / all products grid (main content)
  4. Secondary banner (optional promo)
- Goal: customer can start browsing products immediately without scrolling past heavy hero sections

---

# Final Rule

Claude should always:
- think before generating code
- prefer simple solutions
- explain tradeoffs when relevant
- avoid overengineering
- generate production-grade code by default
- optimize for long-term maintainability