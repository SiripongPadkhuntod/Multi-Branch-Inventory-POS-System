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
- Nginx: http://localhost
- PostgreSQL: localhost:5432

Seed users:

- Owner: `owner@example.com` / `password123`
- Employee: `employee@example.com` / `password123`

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

- Add a migration runner such as Goose or Atlas instead of init-only Docker SQL.
- Implement refresh-token persistence/revocation.
- Add rate limiting middleware backed by Redis.
- Complete transfer approval/completion usecases.
- Add receipt PDF/print rendering.
- Add repository and integration tests using a test PostgreSQL container.
- Add CI jobs for Go tests, frontend typecheck, lint, and Docker image build.
