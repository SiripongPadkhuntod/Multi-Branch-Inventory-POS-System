# Architecture

This backend follows the same high-level layout

## Layout

```text
cmd/
  docs/
  server/
configs/
docs/
internal/
  app/
    domain/
      audit.go
      branch.go
      category.go
      dashboard.go
      inventory.go
      product.go
      role.go
      sale.go
      user.go
    handler/
    port/
    repository/
    service/
  constant/
  infrastructure/
    db-client/
    gin-client/
    logger-client/
    middleware-client/
  property/
  router/
  server/
pkg/
  v1/
    dto/
    model/
```

`internal/app/domain` and `internal/app/port` are the source of truth for core models and contracts. PostgreSQL code is an outbound adapter that satisfies `internal/app/port` repository interfaces. Template-style routes under `/pos-service/v1` call `internal/app/handler -> internal/app/service -> internal/app/port` and do not depend on legacy usecase services.

Domain files are intentionally split by business area so inventory, sales, product, user, dashboard, and audit models can be found without scanning one large model file.

The POS implementation still keeps packages under `internal/delivery`, `internal/domain`, `internal/repository`, and `internal/usecase` as compatibility wrappers for the existing `/api/v1` frontend API while routes are migrated endpoint by endpoint. Those packages must not be imported by new app-layer code.

Runtime bootstrapping starts in `cmd/server`, dependency wiring lives in `internal/server`, and contracts live in `internal/app/port`.
