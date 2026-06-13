# Multi-Branch Inventory & POS System

Production-oriented starter for a multi-branch inventory and POS web application.

## Stack

- Frontend: Next.js App Router, TypeScript, TailwindCSS, Zustand, TanStack-ready API layer, React Hook Form, Zod
- Backend: Go, Gin, Clean Architecture, JWT auth, PostgreSQL, bcrypt
- Infrastructure: Docker Compose, Nginx reverse proxy, PostgreSQL migrations and seed data

## Structure

```text
backend/
  cmd/api
  internal/domain
  internal/usecase
  internal/repository
  internal/delivery
  internal/infrastructure
  migrations
  seed
frontend/
  src/app
  src/features
  src/components
  src/services
  src/hooks
  src/stores
  src/types
  src/utils
nginx/
docker-compose.yml
```

## Run With Docker

```bash
docker compose up --build
```

Services:

- Frontend: http://localhost:3000
- Backend: http://localhost:8080/api/v1
- Backend health check: http://localhost:8080/healthz
- Swagger: http://localhost:8080/swagger
- Nginx: http://localhost
- PostgreSQL: localhost:5433

Seed users and mock data:

- Owner: `owner@example.com` / `password123`
- Manager Siam: `manager@example.com` / `password123`
- Manager North: `manager-north@example.com` / `password123`
- Employee Siam: `employee@example.com` / `password123`
- Employee North: `employee-north@example.com` / `password123`
- Employee Phuket: `employee-phuket@example.com` / `password123`
- No-sales employee: `nosales@example.com` / `password123`

The seed includes demo branches, categories, product images, branch stock, low-stock items, sales, refunds, transfer statuses, movement logs, and audit logs.

To refresh mock data in a running local Docker database:

```bash
psql postgresql://pos:pos_password@localhost:5433/pos -f backend/seed/001_seed.sql
```

## Deploy Notes

The backend runs SQL migrations automatically on startup. Existing databases that already have the initial tables are baselined, then newer migrations are applied in order.

Backend architecture notes live in `backend/ARCHITECTURE.md`. The runtime entrypoint is `backend/cmd/server`, with DI in `backend/internal/server`.

```bash
cd backend
go run ./cmd/server
```

Recommended Render backend environment:

```env
APP_ENV=production
RUN_MIGRATIONS=true
MIGRATIONS_PATH=migrations
DATABASE_URL=<Render PostgreSQL internal database URL>
JWT_SECRET=<strong secret>
```

Use `/healthz` as the Render health check path.

## Implemented Core

- JWT login/logout/me endpoints
- Role middleware for owner-only areas
- Branch-aware employee permissions
- Product CRUD with barcode lookup and soft delete
- Inventory list and transaction-safe stock adjustment
- POS sale creation with row-level inventory locking
- Movement logs for stock changes and sales
- Refund inventory restoration endpoint
- PostgreSQL schema with UUIDs, foreign keys, constraints, indexes, and soft-delete columns
- Automatic migration runner for Docker and Render deployments
- Next.js login, owner dashboard, POS screen, employee pages, and owner module pages

## API Format

Success:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Error message"
}
```

## Next Production Steps

- Implement refresh-token persistence/revocation.
- Add rate limiting middleware backed by Redis.
- Complete transfer approval/completion usecases.
- Add receipt PDF/print rendering.
- Add repository and integration tests using a test PostgreSQL container.
- Add CI jobs for Go tests, frontend typecheck, lint, and Docker image build.
