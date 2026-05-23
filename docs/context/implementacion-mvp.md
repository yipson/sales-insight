# Plan de Implementación — Sales Insight MVP

**Versión:** 2.0 — Feature-Based Architecture  
**Fecha:** 2026-05-13  
**Estado:** En progreso  
**Repositorio:** `sales-insight`  

> **Documentos relacionados:**
> - [ESTADO_ACTUAL.md](./ESTADO_ACTUAL.md) — Estado actual de implementación y deuda técnica (fuente de verdad del estado)
> - [AGENTS.md](../../AGENTS.md) — Guía rápida de desarrollo y setup local
> - [ARCHITECTURE.md](./ARCHITECTURE.md) — Decisiones arquitectónicas y stack tecnológico
> - [DATABASE_SCHEMA.md](./DATABASE_SCHEMA.md) — Esquema de base de datos

---

## Instrucciones de uso

- Marca cada subtarea con `- [x]` cuando esté completada.
- Realiza un **commit** después de cada tarea principal (1.1, 1.2, etc.).
- Este archivo es el único lugar de verdad para el estado de implementación.

---

## Principios de esta implementación

1. **Feature-first:** Cada fase entrega un componente vertical (modelo → repositorio → servicio → handler) en lugar de capas horizontales (toda la DB, luego toda la API).
2. **Patrón por feature:** Todo feature sigue la estructura: `model.go` → `repository.go` (interfaz) → `sqlc/` (impl) → `service.go` → `handler.go`.
3. **DI manual:** Cada struct recibe sus dependencias por constructor. El wiring vive exclusivamente en `cmd/api/main.go`.
4. **Tests sin DB real:** Los `service.go` se testean con mocks de repositorios. Los `sqlc/` se testean con integración contra PostgreSQL real.

---

## Fase 1: Plataforma + Primer Feature (Merchant)
**Objetivo:** Arquitectura base funcional de punta a punta. Un solo feature completo valida el patrón.

- [x] **1.1 Reestructurar backend (tabula rasa)**
  - Eliminar `internal/config/`, `internal/database/` (estructura obsoleta layered)
  - Crear layout feature-based: `internal/platform/*`, `internal/merchant/`, `internal/auth/`, `internal/clover/`, `internal/sync/`, `internal/orders/`, `internal/products/`, `internal/employees/`, `internal/analytics/`, `internal/dashboard/`, `cmd/api/`, `cmd/worker/`, `scripts/`
  - Mover migraciones a `backend/migrations/`

- [x] **1.2 Instalar dependencias Go**
  - `go get github.com/labstack/echo/v4`
  - `go get github.com/spf13/viper`
  - `go get github.com/sqlc-dev/sqlc`
  - `go get golang.org/x/time/rate`
  - `go get github.com/robfig/cron/v3`
  - `go get github.com/golang-jwt/jwt/v5`
  - `go mod tidy`

- [x] **1.3 Implementar `internal/platform/config/config.go`**
  - Leer variables de entorno con Viper
  - Estructura `Config` con validación de campos obligatorios
  - Defaults seguros para desarrollo local
  - Actualizar `backend/.env.example`

- [x] **1.4 Implementar `internal/platform/db/postgres.go`**
  - Pool de conexiones PostgreSQL con `pgx`/`database/sql`
  - Configuración: max open/idle conns, conn lifetime/idle time
  - Health check con context timeout
  - Test de conexión `postgres_test.go`

- [x] **1.5 Implementar `internal/platform/logger/logger.go`**
  - Inicializar `slog` (JSON en prod, texto en dev)
  - Helper para redactar campos sensibles
  - Pasar por constructor a todos los servicios

- [x] **1.6 Implementar `internal/platform/security/encryption.go`**
  - AES-256-GCM para tokens en reposo
  - Funciones: `Encrypt(plaintext string) ([]byte, error)`, `Decrypt(ciphertext []byte) (string, error)`
  - Sin loguear ni exponer tokens

- [x] **1.7 Crear `cmd/api/main.go` con DI manual**
  - Cargar config → logger → DB
  - Inicializar `security.Encrypter`
  - Wiring: `merchant.Repository` (sqlc impl) → `merchant.Service` → `merchant.Handler`
  - Echo server con health check `GET /health`
  - Graceful shutdown: SIGINT/SIGTERM, cerrar DB y esperar requests activos

- [x] **1.8 Configurar `sqlc.yaml` global**
  - Motor PostgreSQL
  - Package por feature: `merchant`, `orders`, `products`, `employees`
  - Output dentro de cada `internal/{feature}/sqlc/`

- [x] **1.9 Implementar feature `merchant/` completo**
  - `merchant/model.go` — structs de dominio (`Merchant`, `TokenData`)
  - `merchant/repository.go` — interfaz `Repository` (boundaries reales)
  - `merchant/sqlc/queries.sql` — CRUD + token management (sqlc generate)
  - `merchant/sqlc/repository.go` — wrapper que implementa `merchant.Repository` usando código sqlc generado
  - `merchant/service.go` — lógica de negocio, usa `Repository` (interfaz)
  - `merchant/handler.go` — endpoints: `GET /api/v1/merchant/status`, `POST /api/v1/merchant/bootstrap`
  - `merchant/service_test.go` — tests con mock de `Repository`

**Commit sugerido:** `feat: plataforma base + feature merchant con arquitectura feature-based`

---

## Fase 2: Cliente Clover + Autenticación
**Objetivo:** Comunicación segura con Clover API y flujo OAuth para conectar merchants.

- [x] **2.1 Implementar `internal/clover/client.go` + `oauth_client.go` + `dto.go`**
  - `client.go` — HTTP client genérico para API REST v3 (`Authorization: Bearer`, `Content-Type`)
  - `oauth_client.go` — Cliente OAuth separado con URLs distintas (`apisandbox.dev.clover.com` vs `sandbox.dev.clover.com`)
  - `dto.go` — Structs de respuesta: `TokenResponse`, `OrderResponse`, `EmployeeResponse`, etc.
  - Métodos OAuth: `ExchangeCode()`, `RefreshTokens()`, `AuthorizeURL()`
  - ✅ Archivos creados y compilando

- [x] **2.2 Implementar `internal/clover/rate_limiter.go`**
  - Token bucket con `golang.org/x/time/rate`: 16 req/s, burst 5
  - Semáforo para max 5 requests concurrentes
  - Respeto del header `retry-after` ante 429
  - Backoff exponencial con jitter (max 10 retries)
  - Integrado en `clover.Client` para todas las requests REST
  - ✅ Archivo creado y compilando

- [x] **2.3 Implementar token cache + auth service**
  - `token_cache/cache.go` — mapa en memoria [merchantID]TokenData con RWMutex
  - `auth/service.go` — orquesta flujo OAuth completo:
    - `HandleCallback()` — exchange code → guarda tokens encriptados en DB + cache
    - `RefreshAccessToken()` — detecta expiración, llama `/oauth/v2/refresh`, actualiza DB + cache
    - `RebuildCache()` — reconstruye cache desde PostgreSQL al iniciar
  - ✅ Archivos creados y compilando

- [x] **2.4 DTOs y OAuth client completados**
  - `clover/dto.go` — structs de respuesta de Clover API
  - `clover/oauth_client.go` — ExchangeCode, RefreshTokens, AuthorizeURL

- [x] **2.5 Implementar feature `auth/` completo**
  - `auth/service.go` — inicia OAuth, maneja callback, genera JWT de sesión
  - `auth/handler.go` — endpoints: `GET /api/v1/auth/clover`, `GET /api/v1/auth/clover/callback`, `POST /api/v1/auth/bootstrap`, `GET /api/v1/auth/status`, `POST /api/v1/auth/revoke`
  - `auth/middleware.go` — validación JWT en rutas protegidas
  - ✅ Archivos creados y compilando

- [x] **2.6 Actualizar `cmd/api/main.go` + fix config loader**
  - Wire `clover.OAuthClient`, `token_cache.Cache`, `auth.Service`, `auth.Handler`
  - Agregar rutas de auth al router Echo
  - Rebuild token cache desde PostgreSQL al iniciar
  - Fix: cambiar `config.go` a lectura manual con `godotenv` (Viper `Unmarshal` + `squash` no lee env vars correctamente)
  - ✅ Servidor arranca correctamente leyendo `backend/.env`

**Commit:** `feat: auth flow, token cache, OAuth client y fix config loader` ✅

---

## Fase 3: Sync Engine (Extracción + Transformación)
**Objetivo:** Traer datos de Clover y persistirlos en nuestro modelo. Diseñado como componente separable.

- [x] **3.1 Implementar features de dominio (solo modelos + interfaces)**
  - `orders/model.go` — structs: `Order`, `OrderItem`, `OrderCategorySummary`
  - `orders/repository.go` — interfaces: `Repository`, `OrderItemRepository`, `CategorySummaryRepository`
  - `products/model.go` — structs: `Product`, `Category`, `AnalyticCategory`, `CategoryMapping`
  - `products/repository.go` — interfaces: `ProductRepository`, `CategoryRepository`, `AnalyticCategoryRepository`, `CategoryMappingRepository`
  - `employees/model.go` — struct: `Employee`
  - `employees/repository.go` — interface: `Repository`
  - ✅ Compilando (implementaciones sqlc en Fase 4)

- [x] **3.2 Implementar `internal/sync/engine.go`**
  - Orquestador que coordina todas las extracciones
  - Constructor con DI manual de todos los repositories (interfaces)
  - `getToken()` y `getCloverMerchantID()` helpers privados
  - Batch size configurable (default: 100)
  - Diseñado como componente separable para futuro worker
  - ✅ Compilando

- [x] **3.3 Implementar `internal/sync/orders.go`**
  - `SyncOrders(ctx, merchantID, cursor)` — extracción incremental desde Clover
  - Endpoint: `/v3/merchants/{id}/orders?filter=modifiedTime>={cursor}&expand=lineItems`
  - Parseo JSON → modelos de dominio (`orders.Order`, `orders.OrderItem`)
  - Cálculo de agregados: `item_count`, `unique_category_count`, `total_quantity`
  - Upsert batch de órdenes + líneas (respetando `batchSize` del engine)
  - Retorna nuevo cursor (último `modifiedTime` procesado)
  - ✅ Compilando

- [x] **3.4 Implementar `internal/sync/items.go`**
  - `SyncItems(ctx, merchantID)` — extracción completa de productos y categorías
  - Endpoints: `/v3/merchants/{id}/categories` y `/v3/merchants/{id}/items?expand=categories`
  - Transforma categorías Clover → `products.Category` + productos → `products.Product`
  - Mapeo de `CategoryID` en productos buscando por `CloverCategoryID`
  - Upsert batch de categorías + productos
  - Retorna cantidad de productos y categorías procesadas
  - ✅ Compilando

- [x] **3.5 Implementar `internal/sync/employees.go`**
  - `SyncEmployees(ctx, merchantID)` — extracción completa de empleados
  - Endpoint: `/v3/merchants/{id}/employees`
  - Transforma respuesta Clover → `employees.Employee` con `IsActive` desde Clover
  - Upsert batch de empleados
  - Retorna cantidad de empleados procesados
  - ✅ Compilando

- [x] **3.6 Implementar `internal/sync/payments.go`**
  - `SyncPayments(ctx, merchantID, cursor)` — extracción incremental de pagos
  - Endpoint: `/v3/merchants/{id}/payments?filter=modifiedTime>={cursor}`
  - Transforma respuesta Clover → `payments.Payment`
  - Upsert batch de pagos
  - Retorna nuevo cursor + cantidad de pagos procesados
  - ✅ Compilando

- [x] **3.7 Implementar `internal/sync/transform.go`**
  - `PriceToCents(price)` — validación de precios no negativos
  - `TimestampToUTC(ms)` — conversión Unix ms → UTC
  - `RoundToMicrosecond(t)` — precisión PostgreSQL
  - `CalculateCategorySummary()` — placeholder para Phase 4 (requiere mapeo completo de categorías)
  - ✅ Compilando

- [x] **3.8 Implementar `internal/platform/scheduler/cron.go`**
  - `Scheduler` struct wrapper de `robfig/cron/v3`
  - Tareas registradas:
    - `sync_orders` — cada 5 minutos
    - `sync_items` — cada 30 minutos
    - `sync_employees` — cada 1 hora
    - `sync_payments` — cada 15 minutos
  - `runForAllMerchants()` — ejecuta sync para cada merchant conectado
  - `Start()` / `Stop()` — control del cron con graceful shutdown
  - Logging estructurado con slog
  - ✅ Compilando

- [x] **3.9 Actualizar `cmd/api/main.go`**
  - Wire `sync.Engine` con todos los repositories (merchant real + stubs para orders/products/employees/payments)
  - Wire `scheduler.Scheduler` con el engine y merchantRepo
  - `cronScheduler.Start()` al iniciar la aplicación
  - `cronScheduler.Stop()` en graceful shutdown (antes de cerrar HTTP server)
  - ✅ Compilando y tests pasando

**Commit sugerido:** `feat: sync engine con polling incremental, scheduler y transform`

---

## Fase 4: Features de API REST (Dominio)
**Objetivo:** Endpoints CRUD y listados para employees, products/categories y orders.

- [x] **4.1 Implementar feature `employees/` completo**
  - `employees/sqlc/queries.sql` + `employees/sqlc/repository.go`
  - `employees/service.go` + `employees/handler.go`
  - Endpoints: `GET /api/v1/employees`, `GET /api/v1/employees/:id`, `GET /api/v1/employees/:id/orders` (delega a orders.Service)
  - `employees/service_test.go` con mocks

- [x] **4.2 Implementar feature `products/` completo**
  - `products/sqlc/queries.sql` (CRUD productos + categorías nativas/analíticas + mappings)
  - `products/sqlc/repository.go`
  - `products/service.go` + `products/handler.go`
  - Endpoints: `GET /api/v1/products`, `GET /api/v1/products/:id`, `GET /api/v1/products/categories`
  - Endpoints admin: `PUT /api/v1/products/categories/:id/map`, `POST /api/v1/analytic-categories`, `PUT /api/v1/analytic-categories/:id`, `DELETE /api/v1/analytic-categories/:id`
  - `products/service_test.go`
  - Nota: `GET /api/v1/products/top` queda pendiente (requiere order_items analytics, ver 4.4)

- [x] **4.3 Implementar feature `orders/` completo**
  - `orders/sqlc/queries.sql` (CRUD + upsert con ON CONFLICT)
  - `orders/sqlc/repository.go` (3 repos separados para evitar colisión de nombres de métodos)
  - `orders/service.go` + `orders/handler.go`
  - Endpoints: `GET /api/v1/orders`, `GET /api/v1/orders/:id` (devuelve OrderDetail con items + category_summary)
  - `orders/service_test.go`

- [x] **4.4 Implementar feature `analytics/` completo**
  - `analytics/model.go` — structs de KPIs
  - `analytics/repository.go` — interfaz para queries agregadas
  - `analytics/sqlc/queries.sql` — 6 queries agregadas (summary, sales-by-employee, top-products, category-coverage, count-categories, ticket-ideal-breakdown)
  - `analytics/sqlc/repository.go` — implementación sqlc
  - `analytics/service.go` — orquesta queries + cálculo de cobertura % y ticket ideal
  - `analytics/service_test.go` — mocks + tests
  - Sin handlers propios; consumido por `dashboard.Service` (Fase 5)

- [x] **4.5 Implementar feature `payments/` completo**
  - `payments/sqlc/queries.sql` — upsert + lookups
  - `payments/sqlc/repository.go` — implementación sqlc
  - `payments/service.go` — lógica de negocio (para sync y API)
  - `payments/service_test.go` — mocks + tests
  - Reemplaza stub en sync engine; el pipeline de persistencia de pagos queda 100% funcional

- [x] **4.6 Actualizar `cmd/api/main.go` y eliminar stubs**
  - Wire de repos reales: employees, products (4 structs), orders (3 structs), analytics, payments
  - Wire de services: orderSvc → employeeSvc (delegación), productSvc, analyticsSvc
  - Registro de handlers en Echo: employees, products, orders
  - Sync engine usa repos reales de todos los dominios
  - **Limpieza:** Eliminados archivos `stub.go` de employees, products, orders, payments

**Commit sugerido:** `feat: API REST employees, products, orders, analytics y payments; wire real repos + eliminar stubs`

---

## Fase 5: Dashboard + Sync Status API
**Objetivo:** Endpoints de métricas y estado de sincronización.

- [ ] **5.1 Implementar feature `dashboard/` completo**
  - `dashboard/model.go` — DTOs de respuesta (`DashboardSummary`, `SalesByEmployee`, etc.)
  - `dashboard/service.go` — orquesta `analytics.Service`, `orders.Service`, `employees.Service`
  - `dashboard/handler.go` — endpoints:
    - `GET /api/v1/dashboard/summary`
    - `GET /api/v1/dashboard/sales-by-employee`
    - `GET /api/v1/dashboard/top-products`
    - `GET /api/v1/dashboard/category-coverage`
    - `GET /api/v1/dashboard/ticket-ideal`
  - Todos con parámetro `?from=YYYY-MM-DD&to=YYYY-MM-DD`

- [ ] **5.2 Implementar sync status handlers**
  - Agregar a `sync/` (o crear `sync/handler.go`):
    - `GET /api/v1/sync/status` — última sync por entidad
    - `GET /api/v1/sync/logs` — historial de sync
    - `POST /api/v1/sync/trigger` — disparar sync manual
    - `GET /api/v1/sync/errors` — errores no resueltos

- [ ] **5.3 Implementar `internal/api/server.go`**
  - Setup centralizado de Echo con middlewares: CORS, logger, recovery, request ID
  - Montar todos los handlers en grupos de rutas
  - `GET /health` con chequeo de DB

- [ ] **5.4 Actualizar `cmd/api/main.go`**
  - Usar `api.NewServer(...)` en lugar de setup inline
  - Wire `dashboard.Service` y `sync.Handler`

**Commit sugerido:** `feat: dashboard API y sync status endpoints`

---

## Fase 6: Frontend (Dashboard)
**Objetivo:** Interfaz visual desktop-first.

- [ ] **6.1 Configurar Tailwind CSS**
  - Desktop-first (sin breakpoints móviles en MVP)
  - Paleta de colores base en `tailwind.config.js`

- [ ] **6.2 Instalar Apache ECharts**
  - `npm install echarts echarts-for-react`

- [ ] **6.3 Crear layout base**
  - `src/components/Layout.tsx`
  - `src/components/Sidebar.tsx`
  - `src/components/Header.tsx`

- [ ] **6.4 Página `Dashboard.tsx`**
  - KPIs cards: ventas totales, órdenes, ticket promedio, categorías promedio
  - Gráfico: ventas por período (línea temporal)
  - Gráfico: top 5 empleados (barras horizontales)
  - Filtro de rango de fechas

- [ ] **6.5 Página `Employees.tsx`**
  - Tabla con métricas
  - Gráfico: cobertura de categorías por empleado

- [ ] **6.6 Página `Products.tsx`**
  - Tabla: productos más vendidos
  - Gráfico: distribución por categoría
  - Sección admin: crear/editar categorías analíticas + mapeo

- [ ] **6.7 Página `Orders.tsx`**
  - Tabla paginada con filtros
  - Detalle de orden al hacer clic

- [ ] **6.8 Página `Settings.tsx`**
  - Estado de conexión con Clover
  - Botón "Conectar / Revocar"
  - Botón "Sync manual"
  - Logs de sincronización recientes

- [ ] **6.9 Integrar TanStack Query**
  - `src/hooks/useDashboard.ts`
  - `src/hooks/useEmployees.ts`
  - `src/hooks/useProducts.ts`
  - `src/hooks/useOrders.ts`
  - Invalidación de caché después de sync

**Commit sugerido:** `feat: frontend dashboard con Apache ECharts y TanStack Query`

---

## Fase 7: Integración y Validación
**Objetivo:** Todo funcionando junto.

- [ ] **7.1 Actualizar `docker-compose.yml`**
  - Servicios: db, backend, frontend
  - Variables de entorno compartidas
  - Healthchecks y depends_on
  - Volumes persistentes

- [ ] **7.2 Probar flujo OAuth / Bootstrap**
  - Obtener tokens de Clover (sandbox o producción)
  - Insertar tokens vía bootstrap manual
  - Verificar reconstrucción de caché en memoria

- [ ] **7.3 Probar primera sincronización**
  - Trigger manual desde Settings
  - Verificar datos en PostgreSQL
  - Confirmar que `sync_logs` registra la ejecución
  - Validar que no hay duplicados tras múltiples syncs

- [ ] **7.4 Validar métricas del dashboard**
  - Comparar totales de ventas con Clover
  - Verificar cálculo de cobertura de categorías
  - Confirmar que `ticket_completo_rules` funciona

- [ ] **7.5 Verificar persistencia**
  - Reiniciar contenedor de DB: datos deben permanecer
  - Reiniciar backend: caché de tokens se reconstruye
  - Verificar que scheduler retoma desde el cursor guardado

- [ ] **7.6 Documentar instrucciones de uso**
  - Actualizar `README.md` con pasos de setup
  - Incluir troubleshooting común
  - Notas sobre migración a DO

**Commit sugerido:** `feat: integración completa y validación del MVP`

---

## Resumen de Commits Esperados

| # | Commit | Fase |
|---|---|---|
| 1 | `feat: plataforma base + feature merchant con arquitectura feature-based` | 1 |
| 2 | `feat: cliente Clover con OAuth, rate limiting, encriptación y auth` | 2 |
| 3 | `feat: sync engine con polling incremental, scheduler y transform` | 3 |
| 4 | `feat: API REST features employees, products, orders y analytics` | 4 |
| 5 | `feat: dashboard API y sync status endpoints` | 5 |
| 6 | `feat: frontend dashboard con Apache ECharts y TanStack Query` | 6 |
| 7 | `feat: integración completa y validación del MVP` | 7 |

---

## Notas

- **Este documento se actualiza en cada sesión.** Marca las tareas completadas y añade notas si algo cambia.
- **Cada commit es responsabilidad del usuario.** El agente ejecuta el código, el usuario hace `git add . && git commit -m "..."`.
- **Si una tarea requiere más de una sesión,** divídela en subtareas adicionales.
- **El backend está diseñado para separarse:** `cmd/worker/` puede extraerse a contenedor independiente sin modificar `internal/sync/`.

---

*Última actualización: 2026-05-19*
