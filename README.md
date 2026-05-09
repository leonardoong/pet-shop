# 🐾 PetShop

Full-stack e-commerce platform for a pet shop — built with Go, React, PostgreSQL, and Redis.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.23, Gin, GORM, PostgreSQL 16, Redis 7 |
| **Auth** | JWT (access + refresh token rotation), bcrypt, RBAC |
| **Customer Frontend** | React 18, Vite, TypeScript, TailwindCSS, TanStack Query, Zustand |
| **Admin Dashboard** | React 18, Vite, TypeScript, TailwindCSS, TanStack Table, Recharts |
| **Infrastructure** | Docker, Docker Compose, Nginx |

---

## Project Structure

```
pet-shop/
├── backend/                  # Go REST API
│   ├── cmd/api/              # Entry point
│   ├── internal/
│   │   ├── auth/             # Auth handler, service, repository, DTOs
│   │   ├── admin/            # Admin model + RBAC
│   │   ├── customer/         # Customer model
│   │   ├── token/            # Refresh token model
│   │   ├── middleware/       # JWT auth + RBAC permission middleware
│   │   └── server/           # Gin server, CORS, route registration
│   ├── pkg/
│   │   ├── config/           # Environment config loader
│   │   ├── database/         # Postgres + Redis connection
│   │   ├── jwt/              # JWT generation & validation
│   │   ├── response/         # Consistent API response helpers
│   │   └── validator/        # Custom validators (Indonesian phone, etc.)
│   ├── migrations/           # SQL migration files (up + down)
│   └── docs/                 # Auto-generated Swagger docs
├── frontend/                 # Customer-facing storefront (Bahasa Indonesia)
├── admin/                    # Admin dashboard
├── nginx/                    # Reverse proxy config
├── docker-compose.yml        # Postgres + Redis for local dev
└── docker-compose.full.yml   # Full stack (all services)
```

---

## Getting Started

### Prerequisites

- [Go 1.23+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker & Docker Compose](https://www.docker.com/)

### 1. Clone and configure environment

```bash
git clone <repo-url>
cd pet-shop
cp .env.example .env
```

Edit `.env` with your secrets (JWT secrets, DB password, etc.).

### 2. Start infrastructure (Postgres + Redis)

```bash
make up
```

### 3. Run database migrations

```bash
make migrate
# or manually:
for f in backend/migrations/*.up.sql; do
  docker exec -i petshop_postgres psql -U petshop -d petshop_db < "$f"
done
```

### 4. Start backend

```bash
make dev-backend
# API available at http://localhost:8080
# Swagger UI at  http://localhost:8080/swagger/index.html
```

### 5. Start frontend apps

```bash
# Customer storefront — http://localhost:3000
make dev-frontend

# Admin dashboard — http://localhost:3001
make dev-admin
```

---

## API Documentation

Interactive Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All protected endpoints require a `Bearer` token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

Access tokens expire in **15 minutes**. Use the refresh endpoint to obtain a new one.

---

### Endpoints

#### Customer Auth

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/customer/auth/register` | Register new customer |
| `POST` | `/customer/auth/login` | Customer login |
| `POST` | `/customer/auth/refresh` | Refresh access token |
| `POST` | `/customer/auth/logout` | Revoke refresh token |

#### Admin Auth

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/admin/auth/login` | Admin login (returns permissions) |
| `POST` | `/admin/auth/refresh` | Refresh admin access token |
| `POST` | `/admin/auth/logout` | Revoke admin refresh token |

#### Protected (examples)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/customer/me` | Customer JWT | Get current customer |
| `GET` | `/admin/me` | Admin JWT | Get current admin + permissions |

---

### Request / Response Format

**Success**
```json
{
  "success": true,
  "message": "Login successful",
  "data": { ... }
}
```

**Error**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": ["field is required"]
}
```

---

## Database Schema

| Table | Description |
|-------|-------------|
| `customers` | Customer accounts |
| `admins` | Admin accounts |
| `roles` | RBAC roles (super_admin, manager, inventory_staff, support) |
| `permissions` | Granular permissions (e.g. `products:create`, `orders:update`) |
| `role_permissions` | Roles ↔ permissions mapping |
| `admin_roles` | Admins ↔ roles mapping |
| `refresh_tokens` | Hashed refresh tokens (customer + admin, with rotation) |
| `categories` | Product categories |
| `products` | Products with stock and SKU |
| `addresses` | Customer shipping addresses |
| `carts` | One cart per customer |
| `cart_items` | Items in cart |
| `orders` | Orders with snapshotted shipping address |
| `order_items` | Line items with snapshotted unit price |

---

## RBAC Roles (seeded)

| Role | Permissions |
|------|-------------|
| `super_admin` | All permissions |
| `manager` | Products, categories, orders, inventory, customers, dashboard |
| `inventory_staff` | Inventory read/update, products read |
| `support` | Orders read, customers read, products read |

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | `development` or `production` | `development` |
| `APP_PORT` | API server port | `8080` |
| `DB_HOST` | Postgres host | `localhost` |
| `DB_PORT` | Postgres port | `5432` |
| `DB_USER` | Postgres user | `petshop` |
| `DB_PASSWORD` | Postgres password | — |
| `DB_NAME` | Database name | `petshop_db` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | — |
| `JWT_ACCESS_SECRET` | Secret for signing access tokens | — |
| `JWT_REFRESH_SECRET` | Secret for signing refresh tokens | — |
| `JWT_ACCESS_EXPIRY_MINUTES` | Access token TTL in minutes | `15` |
| `JWT_REFRESH_EXPIRY_DAYS` | Refresh token TTL in days | `7` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins | `http://localhost:3000,...` |

---

## Make Commands

```bash
make up            # Start Postgres + Redis
make down          # Stop all containers
make build         # Build all Docker images
make dev-backend   # Run Go API locally
make dev-frontend  # Run customer frontend (port 3000)
make dev-admin     # Run admin dashboard (port 3001)
make install       # npm install for frontend + admin
make tidy          # go mod tidy
make env           # Copy .env.example → .env
```

---

## Roadmap

- [x] Phase 1 — Foundation (auth, DB schema, scaffolding)
- [ ] Phase 2 — Product catalog, cart, customer shop pages
- [ ] Phase 3 — Admin dashboard (products, orders, inventory)
- [ ] Phase 4 — Payment integration, user profiles
- [ ] Phase 5 — Reviews, email notifications, CI/CD
