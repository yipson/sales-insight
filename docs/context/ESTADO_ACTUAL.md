# Estado Actual del Proyecto — Sales Insight Backend

**Última actualización:** 2026-05-19  
**Sesión:** Refactorización arquitectónica + Fases 1-3 implementadas  
**Branch:** `development`  

---

## 1. Resumen Ejecutivo

El backend ha sido completamente reconstruido con arquitectura **feature-based**, **DI manual** y **repository pattern sobre sqlc**. Las Fases 1, 2 y 3 del MVP están implementadas y compilando. La Fase 4 (implementaciones sqlc de dominio) y Fase 5 (Dashboard API) están pendientes.

**Estado de compilación:** ✅ `go build ./cmd/api/` funciona  
**Estado de tests:** ✅ `go test ./...` pasa (tests unitarios con mocks)  
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
│   │   ├── db/        # PostgreSQL pool
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
│   ├── orders/        # Feature: solo modelos + interfaz + stub
│   ├── products/      # Feature: solo modelos + interfaz + stub
│   ├── employees/     # Feature: solo modelos + interfaz + stub
│   ├── payments/      # Feature: solo modelos + interfaz + stub
│   ├── auth/          # Feature completo (OAuth + JWT middleware)
│   ├── sync/          # Engine de sincronización
│   │   ├── engine.go    # Orquestador con DI
│   │   ├── orders.go    # Extracción incremental
│   │   ├── items.go     # Extracción productos/categorías
│   │   ├── employees.go # Extracción empleados
│   │   ├── payments.go  # Extracción pagos
│   │   └── transform.go # Utilidades de normalización
│   ├── analytics/     # Feature: vacío (Fase 5)
│   ├── dashboard/     # Feature: vacío (Fase 5)
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
| D10 | Stub repositories para features no implementadas | Permite compilar el sync engine sin sqlc completo | `orders/stub.go`, `products/stub.go`, etc. |

---

## 4. Correcciones Temporales (Deuda Técnica)

### CT1: Stub repositories en dominio
**Ubicación:** `orders/stub.go`, `products/stub.go`, `employees/stub.go`, `payments/stub.go`  
**Problema:** El sync engine necesita los repositories para compilar, pero las implementaciones sqlc de estos features no existen todavía.  
**Solución temporal:** Stub repositories que devuelven `ErrNotImplemented`.  
**Cuándo corregir:** Fase 4 — crear `queries.sql` + `sqlc generate` + `sqlc/repository.go` para cada feature, luego eliminar stubs.  
**Impacto:** El scheduler intentará ejecutar sync pero fallará silenciosamente (los stubs devuelven error, el scheduler loguea y continúa).

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

### CT4: Cursor inicial en scheduler usa `time.Time{}`
**Ubicación:** `platform/scheduler/cron.go`  
**Problema:** `syncOrders` y `syncPayments` usan cursor en zero time (`time.Time{}`), lo que significa que la primera ejecución traerá TODO el historial de Clover.  
**Solución temporal:** Comentario `TODO` indica que en producción debería cargarse desde `sync_logs`.  
**Cuándo corregir:** Fase 6 (integración) — implementar lectura/escritura de cursores en `sync_logs`.  
**Impacto:** Primera sync muy pesada; riesgo de rate limiting agresivo.

### CT5: `CalculateCategorySummary` es placeholder
**Ubicación:** `sync/transform.go`  
**Problema:** El cálculo de `order_category_summary` requiere mapeo completo de categorías analíticas, que no existe todavía.  
**Solución temporal:** Función devuelve `nil` siempre.  
**Cuándo corregir:** Fase 4/5 — cuando `products/` tenga sqlc implementado con `category_mappings`.  
**Impacto:** El dashboard no mostrará cobertura de categorías por empleado todavía.

### CT6: Payments sin asociación a Orders/Employees
**Ubicación:** `sync/payments.go`  
**Problema:** Los pagos se extraen pero no se resuelven los `order_id` y `employee_id` internos (se guardan como nil).  
**Solución temporal:** Se extraen los Clover IDs pero no se hace lookup.  
**Cuándo corregir:** Fase 4 — cuando `orders` y `employees` tengan repositorios sqlc funcionando.  
**Impacto:** Métricas de pagos por empleado/orden no disponibles.

---

## 5. Estado de Implementación por Fase

### ✅ Fase 1: Plataforma + Merchant (COMPLETADA)
- [x] Layout feature-based
- [x] `platform/config` — Viper + godotenv + lectura manual
- [x] `platform/db` — PostgreSQL pool
- [x] `platform/logger` — slog
- [x] `platform/security` — AES-256-GCM
- [x] `cmd/api/main.go` — DI manual + graceful shutdown
- [x] `sqlc.yaml` — Configuración por feature
- [x] `merchant/` — Feature completo con sqlc generado
  - [x] Modelos, interfaz, service, handler, tests con mocks
  - [x] `merchant/sqlc/queries.sql` + código generado
  - [x] `merchant/sqlc/repository.go` — wrapper que implementa interfaz

### ✅ Fase 2: Cliente Clover + Autenticación (COMPLETADA)
- [x] `clover/client.go` — REST API client con rate limiter
- [x] `clover/oauth_client.go` — OAuth endpoints separados
- [x] `clover/dto.go` — Structs de respuesta
- [x] `clover/rate_limiter.go` — Token bucket + semaphore + retry 429
- [x] `token_cache/cache.go` — Cache en memoria con RWMutex
- [x] `auth/service.go` — HandleCallback, RefreshAccessToken, RebuildCache
- [x] `auth/handler.go` — OAuth, bootstrap, status, revoke
- [x] `auth/middleware.go` — JWT validation
- [x] `cmd/api/main.go` — Wire auth + rebuild cache al iniciar

### ✅ Fase 3: Sync Engine (COMPLETADA)
- [x] Domain models + interfaces: `orders/`, `products/`, `employees/`, `payments/`
- [x] `sync/engine.go` — Orquestador con DI
- [x] `sync/orders.go` — Extracción incremental con cursor
- [x] `sync/items.go` — Extracción productos/categorías
- [x] `sync/employees.go` — Extracción empleados
- [x] `sync/payments.go` — Extracción pagos incremental
- [x] `sync/transform.go` — Normalización de datos
- [x] `platform/scheduler/cron.go` — Cron con frecuencias
- [x] Stub repositories para features no implementadas
- [x] `cmd/api/main.go` — Wire engine + scheduler

### ⏳ Fase 4: Implementaciones sqlc de Dominio (PENDIENTE)
- [ ] `orders/sqlc/queries.sql` + `sqlc/repository.go`
- [ ] `products/sqlc/queries.sql` + `sqlc/repository.go`
- [ ] `employees/sqlc/queries.sql` + `sqlc/repository.go`
- [ ] `payments/sqlc/queries.sql` + `sqlc/repository.go`
- [ ] Eliminar stub repositories
- [ ] Implementar `sync_logs` y `sync_errors` repositories
- [ ] Actualizar `sqlc.yaml` con todas las queries

### ⏳ Fase 5: Dashboard API (PENDIENTE)
- [ ] `analytics/` — Queries agregadas (KPIs, rankings)
- [ ] `dashboard/` — Service que orquesta analytics
- [ ] `dashboard/handler.go` — Endpoints del dashboard

### ⏳ Fase 6: Integración y Validación (PENDIENTE)
- [ ] `docker-compose.yml` actualizado
- [ ] Flujo OAuth end-to-end
- [ ] Primera sincronización real
- [ ] Validación de métricas

### ⏳ Fase 7: Frontend (PENDIENTE)
- [ ] React + Vite + Tailwind
- [ ] Apache ECharts
- [ ] TanStack Query

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
5. Reemplazar stub en `cmd/api/main.go` con la implementación real
6. Ejecutar `go test ./...`

### Para agregar un nuevo endpoint HTTP:
1. Modificar `feature/handler.go` con el nuevo método
2. Agregar la ruta en `feature/handler.go:RegisterRoutes()`
3. No es necesario tocar `cmd/api/main.go` (las rutas se registran automáticamente)

### Archivos de referencia obligatorios:
- `docs/context/ARCHITECTURE.md` — Arquitectura del sistema
- `implementacion/implementacion-mvp.md` — Plan detallado con checklist
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
├── orders.Repository (STUB)
├── orders.OrderItemRepository (STUB)
├── orders.CategorySummaryRepository (STUB)
├── products.ProductRepository (STUB)
├── products.CategoryRepository (STUB)
├── employees.Repository (STUB)
└── payments.Repository (STUB)
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
