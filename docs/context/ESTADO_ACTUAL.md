# Estado Actual del Proyecto — Sales Insight Backend

**Última actualización:** 2026-05-22  
**Sesión:** Sincronización de documentación con estado real del código (Fases 1-5 implementadas)  
**Branch:** `development`  

---

## 1. Resumen Ejecutivo

El backend ha sido completamente reconstruido con arquitectura **feature-based**, **DI manual** y **repository pattern sobre sqlc**. Las Fases 1 a 5 del MVP están implementadas, compilando y con tests pasando. Las Fases 6 (Frontend React) y 7 (Integración / Docker Compose) están pendientes.

**Estado de compilación:** ✅ `go build ./cmd/api/` funciona  
**Estado de tests:** ✅ `go test ./...` pasa (tests unitarios con mocks en: `merchant/`, `analytics/`, `dashboard/`, `employees/`, `orders/`, `payments/`, `products/`; packages `platform/*`, `auth/`, `clover/`, `sync/`, `token_cache/` pendientes de tests)  
**Servidor arranca:** ✅ Lee `.env` desde cualquier directorio  

---

## 2. Arquitectura Implementada

### Layout Feature-Based (DECISIÓN PERMANENTE)
```
backend/
├── cmd/
│   ├── api/           # Entrypoint unificado (API + Scheduler en MVP)
│   └── worker/        # Placeholder para futura separación
├── internal/
│   ├── platform/      # Infra transversal (sin lógica de negocio)
│   │   ├── config/    # Viper + godotenv (lectura manual, no Unmarshal)
│   │   ├── db/        # PostgreSQL pool (postgres.go)
│   │   ├── logger/    # slog estructurado
│   │   ├── scheduler/ # Wrapper robfig/cron/v3
│   │   └── security/  # AES-256-GCM
│   ├── clover/        # Integración externa
│   │   ├── client.go       # REST API client (con rate limiter)
│   │   ├── oauth_client.go # OAuth endpoints separados
│   │   ├── rate_limiter.go # Token bucket + semaphore + retry 429
│   │   └── dto.go          # Structs de respuesta Clover
│   ├── merchant/      # Feature completo con sqlc
│   │   ├── model.go
│   │   ├── repository.go      # INTERFAZ
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── sqlc/
│   │       ├── queries.sql      # SQL para sqlc generate
│   │       ├── db.go            # GENERADO (no tocar)
│   │       ├── models.go        # GENERADO (no tocar)
│   │       ├── queries.sql.go   # GENERADO (no tocar)
│   │       └── repository.go    # Wrapper manual → implementa interfaz
│   ├── orders/        # Feature completo (modelo + interfaz + sqlc + service + handler + tests)
│   ├── products/      # Feature completo (modelo + interfaz + sqlc + service + handler + tests)
│   ├── employees/     # Feature completo (modelo + interfaz + sqlc + service + handler + tests)
│   ├── payments/      # Feature completo (modelo + interfaz + sqlc + service + tests)
│   ├── auth/          # Feature completo (OAuth + JWT middleware)
│   ├── sync/          # Engine de sincronización + API de status
│   │   ├── engine.go    # Orquestador con DI
│   │   ├── orders.go    # Extracción incremental
│   │   ├── items.go     # Extracción productos/categorías
│   │   ├── employees.go # Extracción empleados
│   │   ├── payments.go  # Extracción pagos
│   │   ├── transform.go # Utilidades de normalización
│   │   ├── model.go     # Structs Log, Error, SyncEntity
│   │   ├── repository.go # Interfaces LogRepository, ErrorRepository
│   │   ├── handler.go   # Endpoints /sync/status, /sync/logs, /sync/trigger, /sync/errors
│   │   └── sqlc/        # Implementación sqlc de sync logs y errors
│   ├── analytics/     # Feature completo (modelo + interfaz + sqlc + service + tests; sin handler propio)
│   ├── dashboard/     # Feature completo (modelo + service + handler + tests; orquesta analytics)
│   └── token_cache/   # Cache en memoria de tokens
├── migrations/        # golang-migrate
├── sqlc.yaml          # Configuración por feature
└── .env               # Desarrollo local (no commitear)
```

### Patrón por Feature (DECISIÓN PERMANENTE)
Cada feature sigue:
```
feature/
├── model.go              # Structs de dominio
├── repository.go         # INTERFAZ (contrato)
├── service.go            # Lógica de negocio (depende de interfaz)
├── handler.go            # HTTP handlers
└── sqlc/
    ├── queries.sql       # SQL para sqlc generate
    └── repository.go     # Implementación concreta (wrapper sqlc)
```

### Dependency Injection (DECISIÓN PERMANENTE)
- **DI manual puro** en `cmd/api/main.go`
- Cada struct recibe dependencias por constructor
- **NO usar** samber/do, uber/fx, ni google/wire (a menos que el grafo crezca >20 servicios)
- `main.go` es el único lugar con wiring

### Database Access (DECISIÓN PERMANENTE)
- **sqlc** para generación de código tipado desde SQL
- **Repository pattern**: interfaz en `feature/repository.go`, implementación en `feature/sqlc/repository.go`
- **NO usar ORM** (GORM/ent)
- Transacciones manuales con `*sql.Tx`
- Batch size default: 100 registros

---

## 3. Decisiones Permanentes Tomadas

| # | Decisión | Justificación | Archivos afectados |
|---|----------|---------------|-------------------|
| D1 | Feature-based layout | Escalabilidad cognitiva, menos conflictos de merge | Toda la estructura `internal/` |
| D2 | DI manual | <15 servicios, explícito, fácil de testear | `cmd/api/main.go` |
| D3 | Repository pattern sobre sqlc | Desacopla negocio de infraestructura, testeable | `*/repository.go` + `*/sqlc/repository.go` |
| D4 | API + Worker unificados en MVP | Reduce complejidad operativa en local | `cmd/api/main.go` (scheduler dentro) |
| D5 | OAuth client separado de REST client | URLs diferentes, headers diferentes | `clver/oauth_client.go` vs `clover/client.go` |
| D6 | Token cache en memoria (1 réplica) | MVP: 1 réplica. Futuro: Redis/external | `token_cache/cache.go` |
| D7 | JWT para sesiones de dashboard | Clover es el IdP, nosotros generamos JWT interno | `auth/middleware.go` |
| D8 | Config loader con godotenv + búsqueda recursiva | Viper Unmarshal con squash no funciona con env vars | `platform/config/config.go` |
| D9 | Graceful shutdown: `e.Close()` en dev, `e.Shutdown()` en prod | Evita TIME_WAIT problemático en desarrollo | `cmd/api/main.go` |
| D10 | ~~Stub repositories para features no implementadas~~ *(resuelto en Fase 4)* | ~~Permitía compilar el sync engine sin sqlc completo~~ | ~~`orders/stub.go`, etc.~~ — **Eliminados** |

---

## 4. Correcciones Temporales (Deuda Técnica)

### ~~CT1: Stub repositories en dominio~~ ✅ RESUELTO
**Ubicación:** ~~`orders/stub.go`, `products/stub.go`, `employees/stub.go`, `payments/stub.go`~~  
**Estado:** Resuelto en Fase 4. Todos los features de dominio (`orders`, `products`, `employees`, `payments`, `analytics`, `sync`) ahora tienen implementaciones sqlc reales con modelos, queries, repositorios, servicios, handlers y tests. Los stubs fueron eliminados y `cmd/api/main.go` realiza el wiring de repos reales. El sync engine opera con implementaciones completas.

### CT2: `models.go` generado por sqlc tiene todas las tablas
**Ubicación:** `merchant/sqlc/models.go`  
**Problema:** sqlc lee todo el schema desde `migrations/` y genera modelos para todas las tablas (incluyendo `employees`, `orders`, etc.) en cada package.  
**Solución temporal:** Aceptar que `merchant/sqlc/models.go` tiene ~15 structs, solo usamos `Restaurant`.  
**Cuándo corregir:** No es crítico. Cuando haya muchos features, podríamos usar un package `internal/models` central o scripts de post-procesado.  
**Impacto:** Ruido en código generado, pero no afecta funcionalidad.

### CT3: `db.go` duplicado por feature
**Ubicación:** `*/sqlc/db.go`  
**Problema:** sqlc genera `DBTX` interface en cada package.  
**Solución temporal:** Aceptar la duplicación.  
**Cuándo corregir:** Si crece mucho, definir `platform/db.DBTX` y configurar sqlc para no generar `db.go`.  
**Impacto:** Mínimo, es código generado.

### ~~CT4: Cursor inicial en scheduler usa `time.Time{}`~~ ✅ RESUELTO
**Ubicación:** `platform/scheduler/cron.go`  
**Estado:** Resuelto. El scheduler ahora:
1. Lee el último cursor desde `sync_logs` vía `syncLogRepo.GetLatestByEntity()`.
2. Si no existe log previo, usa el **primer día del mes anterior** como cursor inicial (evita descargar años de historial).
3. Después de cada sync exitosa, escribe un registro en `sync_logs` con `cursor_from`, `cursor_to` y estado `"success"`.
4. Se agregó endpoint `POST /api/v1/sync/backfill` para sincronización manual de rangos históricos.
**Impacto:** Primera sync limitada a ~30-60 días de datos. Backfill disponible bajo demanda.

### CT5: `CalculateCategorySummary` es placeholder
**Ubicación:** `sync/transform.go`  
**Problema:** El cálculo de `order_category_summary` requiere mapeo completo de categorías analíticas, que no existe todavía.  
**Solución temporal:** Función devuelve `nil` siempre.  
**Cuándo corregir:** Fase 4/5 — cuando `products/` tenga sqlc implementado con `category_mappings`.  
**Impacto:** El dashboard no mostrará cobertura de categorías por empleado todavía.

### ~~CT6: Payments sin asociación a Orders/Employees~~ ✅ RESUELTO
**Ubicación:** `sync/payments.go` + `sync/engine.go`  
**Estado:** Resuelto. El sync engine ahora resuelve `order_id` y `employee_id` internos durante la sincronización de pagos:
1. `resolveOrderID(ctx, merchantID, cloverOrderID)` → busca en `orders` por `clover_order_id`.
2. `resolveEmployeeID(ctx, merchantID, cloverEmployeeID)` → busca en `employees` por `clover_employee_id`.
3. Si la orden o empleado aún no existe (ej. sync de pagos ocurrió antes que órdenes), el campo queda `nil` y se resuelve en la próxima sync gracias al `ON CONFLICT DO UPDATE` de la query `UpsertPayment`.
**Impacto:** Métricas de pagos por empleado y por orden ahora son calculables. Dashboard analytics funcionan correctamente.

---

## 5. Estado de Implementación

Para el plan detallado de implementación con checklists por fase, ver:
**[implementacion-mvp.md](./implementacion-mvp.md)**

Este documento se enfoca en el estado técnico actual, decisiones arquitectónicas 
y deuda técnica. El plan de ejecución vive centralizado en implementacion-mvp.md.

Resumen de fases:
- ✅ **Fases 1-5 completadas:** Plataforma, Merchant, Clover Client, Auth, Sync Engine, sqlc de dominio (orders, products, employees, payments, analytics), Dashboard API + Sync Status
- ⏳ **Fases 6-7 pendientes:** Frontend React + Integración / Docker Compose

---

## 6. Instrucciones para el Siguiente Agente

### Para continuar el desarrollo:
1. **Leer este archivo primero** — Es la fuente de verdad del estado actual
2. **Verificar compilación:** `cd backend && go build ./cmd/api/` debe funcionar
3. **Verificar tests:** `go test ./...` debe pasar
4. **Chequear `.env`:** Debe existir `backend/.env` con valores de desarrollo

### Para implementar un nuevo feature con sqlc:
1. Crear `internal/{feature}/sqlc/queries.sql` con comentarios `-- name: MethodName :one/:many/:exec`
2. Agregar sección en `sqlc.yaml` apuntando al nuevo `queries.sql`
3. Ejecutar `sqlc generate` desde `backend/`
4. Crear `internal/{feature}/sqlc/repository.go` que implemente la interfaz del feature
   - Usar `toDomain()` y `fromDomain()` helpers (ver `merchant/sqlc/repository.go` como referencia)
5. Crear `internal/{feature}/service.go` e inyectar la interfaz del repositorio
6. Crear `internal/{feature}/handler.go` con `RegisterRoutes()`
7. Wirear el nuevo repositorio y service en `cmd/api/main.go`
8. Ejecutar `go test ./...`

### Para agregar un nuevo endpoint HTTP:
1. Modificar `feature/handler.go` con el nuevo método
2. Agregar la ruta en `feature/handler.go:RegisterRoutes()`
3. No es necesario tocar `cmd/api/main.go` (las rutas se registran automáticamente)

### Archivos de referencia obligatorios:
- `docs/context/ARCHITECTURE.md` — Arquitectura del sistema
- `implementacion-mvp.md` — Plan detallado con checklist
- `backend/migrations/001_initial_schema.up.sql` — Schema de base de datos
- `backend/.env.example` — Variables de entorno

### Decisiones que NO cambiar sin discutir:
- Layout feature-based (no volver a layered)
- DI manual (no adoptar framework de DI)
- sqlc en lugar de ORM
- API + Worker unificados en MVP
- Un solo binario `cmd/api` (worker es placeholder)

---

## 7. Diagrama de Dependencias Actual

```
cmd/api/main.go
├── config.Load() ───────────────────────────────┐
├── logger.New()                                 │
├── db.NewPostgres()                             │
├── security.NewEncrypter()                      │
├── merchant/sqlc.NewSQLCRepository(db, enc)     │── DI Manual
├── merchant.NewService(repo)                    │
├── auth.NewService(oauth, repo, enc, ...)       │
├── auth.NewHandler(authSvc, ...)                │
├── clover.NewClient(env)                        │
├── sync.NewEngine(clover, cache, repo, ...)     │
├── scheduler.NewScheduler(engine, repo, log)    │
└── echo.New() ──────────────────────────────────┘

sync.Engine
├── clover.Client (REST API)
├── token_cache.Cache
├── merchant.Repository
├── orders.Repository
├── orders.OrderItemRepository
├── orders.CategorySummaryRepository
├── products.ProductRepository
├── products.CategoryRepository
├── employees.Repository
└── payments.Repository
```

---

## 8. Notas de Diseño

### Por qué `merchant` usa `sqlc.NullRawMessage` para JSONB
- sqlc genera `pqtype.NullRawMessage` para columnas JSONB
- El wrapper `merchant/sqlc/repository.go` hace `json.Unmarshal()` al convertir a dominio
- Esto mantiene los modelos de dominio limpios (sin dependencias de sqlc)

### Por qué el config loader usa godotenv + mapeo manual
- Viper's `Unmarshal` con `mapstructure:",squash"` **NO** lee variables de entorno a structs anidados correctamente
- `v.GetString("ENCRYPTION_KEY")` funciona, `v.Unmarshal(&cfg)` deja `cfg.Security.EncryptionKey` vacío
- Solución: `godotenv.Load(".env")` primero (carga en `os.Getenv`), luego `v.GetString()` campo por campo
- `findEnvFile()` busca `.env` recursivamente desde el working directory (permite `go run .` desde `cmd/api/`)

### Rate Limiter: token bucket + semaphore
- `golang.org/x/time/rate` maneja el rate limit (16 req/s)
- Un `chan struct{}` con buffer 5 maneja el concurrent limit
- Ambos se aplican en `clover/client.go` mediante `rateLimiter.Do()`
- Ante 429: lee `retry-after` header si existe, sino backoff exponencial `2^attempt + random`

---

*Documento vivo. Actualizar al final de cada sesión de desarrollo.*
