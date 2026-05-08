# ARCHITECTURE.md

## Sales Insight — Arquitectura del Sistema

**Versión:** 1.0-MVP  
**Fecha:** 2026-05-06  
**Ambiente:** PC Local (Docker Compose)  
**Infraestructura destino (post-MVP):** DigitalOcean Droplet  
**Patrón:** Monolito desplegado en contenedores Docker  
**Locales (MVP):** 1 (preparado para 15+)  

---

## 1. Principios Arquitectónicos

1. **Simplicidad operativa:** Un solo comando (`docker compose up`), mínimos componentes.
2. **Sin estado en la aplicación:** Toda la persistencia vive en PostgreSQL. El binario Go puede reiniciarse sin pérdida de información.
3. **Concurrencia eficiente:** Uso nativo de goroutines de Go para manejar polling concurrente sin bloqueos.
4. **Preparación multi-local:** El esquema de datos y la estructura de código permiten la incorporación futura de múltiples `merchant_id` sin refactorización profunda.
5. **Sin dependencias externas innecesarias:** No se incluye Redis, message broker, Nginx ni orquestadores en el MVP local. PostgreSQL y el scheduler interno de Go cubren todos los casos de uso actuales.

---

## 2. Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           USUARIO (Browser Desktop)                         │
└───────────────────────────────┬─────────────────────────────────────────────┘
                                │ HTTP
                                ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         PC LOCAL (Docker Compose)                           │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    GO APPLICATION (Binario Único)                     │    │
│  │                                                                     │    │
│  │  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────┐ │    │
│  │  │   HTTP API  │  │   ETL/       │  │   Token      │  │ Scheduler│ │    │
│  │  │   (Echo)    │  │   Sync       │  │   Cache      │  │ (Cron)   │ │    │
│  │  │             │  │   Engine     │  │   (Memoria)  │  │          │ │    │
│  │  │ • Dashboard │  │ • Polling    │  │ • Access     │  │ • Polling│ │    │
│  │  │   endpoints │  │ • Clover API │  │   Token      │  │ • Token  │ │    │
│  │  │ • Auth      │  │ • Transform  │  │ • Refresh    │  │   refresh│ │    │
│  │  │ • Admin     │  │ • Persist    │  │   Token      │  │ • Retry  │ │    │
│  │  └─────────────┘  └──────────────┘  └──────────────┘  └──────────┘ │    │
│  │                                                                     │    │
│  │  ┌───────────────────────────────────────────────────────────────┐  │    │
│  │  │              SHARED LAYER (Internal Packages)                 │  │    │
│  │  │  • Config  • Database (sqlc/PostgreSQL)  • Clover Client    │  │    │
│  │  │  • Models  • Logger  • Encryption        • Rate Limiter     │  │    │
│  │  └───────────────────────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    │ SQL/TCP (Docker Network)                 │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      POSTGRESQL 15                                  │    │
│  │                                                                     │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │    │
│  │  │restaurant│ │employees │ │products  │ │orders    │ │payments  │  │    │
│  │  │tokens    │ │          │ │          │ │          │ │          │  │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐               │    │
│  │  │order_items│ │categories│ │sync_logs │ │sync_errors│               │    │
│  │  │          │ │(mapping) │ │          │ │          │               │    │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘               │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │              REACT SPA (Vite + Apache ECharts)                      │    │
│  │  ┌───────────────────────────────────────────────────────────────┐  │    │
│  │  │  • Dashboard Desktop   • DateRangePicker   • MetricCard     │  │    │
│  │  │  • SalesByEmployee     • TopProducts        • TicketAnalysis│  │    │
│  │  │  • SalesTrend          • TanStack Query    • Tailwind CSS  │  │    │
│  │  └───────────────────────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        DOCKER NETWORK                               │    │
│  │              (Go ↔ PostgreSQL ↔ React dev server)                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                │
                                │ OAuth 2.0 / REST / Export API
                                ▼
                        ┌──────────────┐
                        │  CLOVER API  │
                        │   CLOUD      │
                        └──────────────┘
```

> **Nota:** En el despliegue local, React corre en modo desarrollo (`npm run dev`, puerto 5173) o build estático servido por Go. Nginx se agregará únicamente cuando se migre a DigitalOcean.

---

## 3. Stack Tecnológico Detallado

### 3.1 Infraestructura
| Componente | Tecnología | Versión | Justificación |
|---|---|---|---|
| **Entorno** | Docker + Docker Compose | 24.x / 2.x | Orquestación simple. Un solo `docker-compose.yml` levanta todo el stack en la PC local. |
| **OS (host)** | Windows / macOS / Linux | — | El host puede ser cualquiera; los contenedores normalizan el ambiente. |
| **Proxy / SSL (post-MVP)** | Nginx | 1.24 | **Pendiente.** Se usará en DigitalOcean para SSL (Let's Encrypt), compresión gzip, y serve de archivos estáticos. No se incluye en el stack local. |

### 3.2 Backend
| Componente | Tecnología | Versión | Justificación |
|---|---|---|---|
| **Lenguaje** | Go | 1.22+ | Binario nativo, bajo consumo de RAM, concurrencia con goroutines ideal para polling eficiente. |
| **Framework HTTP** | Echo | v4 | Ligero, rendimiento cercano a stdlib, middlewares built-in (CORS, recovery, logging), validación de requests simple. |
| **Base de datos** | PostgreSQL | 15 | Robusta para datos relacionales + analíticos temporales. Soporte nativo de JSONB si se requiere almacenar respuestas crudas de Clover. |
| **Acceso a BDD** | sqlc | v1.25 | Genera código Go tipado a partir de queries SQL. Elimina reflexión de ORMs, más rápido en runtime y con type-safety en compile time. |
| **Migrations** | golang-migrate | v4 | CLI estándar para migraciones versionadas de PostgreSQL. |
| **Scheduler** | robfig/cron | v3 | Librería Go madura para tareas periódicas dentro del mismo proceso. Sin procesos externos. |
| **HTTP Client** | net/http + retryablehttp | stdlib / HashiCorp | Cliente HTTP robusto con reintentos automáticos y backoff exponencial para la API de Clover. |
| **Rate Limiting** | golang.org/x/time/rate | — | Rate limiter client-side (token bucket) para respetar los límites de Clover: 16 req/s por token, 50 req/s por app, 5 concurrentes por token. |
| **Configuración** | Viper | v1.18 | Lectura de variables de entorno y archivos de config (YAML/JSON). Estándar en el ecosistema Go. |
| **Logging** | slog (stdlib) | Go 1.21+ | Logger estructurado incluido en stdlib. Sin dependencias externas. |
| **Encriptación** | AES-GCM (stdlib crypto) | — | Encriptación de tokens OAuth en reposo. Clave derivada de variable de entorno `ENCRYPTION_KEY`. |

### 3.3 Frontend
| Componente | Tecnología | Versión | Justificación |
|---|---|---|---|
| **Framework** | React | 18 | Estándar de la industria, ecosistema maduro. |
| **Build tool** | Vite | 5 | Builds rápidos, HMR instantáneo, output optimizado para producción. |
| **Estilos** | Tailwind CSS | 3.4 | Utility-first. En MVP se configura para desktop-first (sin breakpoints móviles). |
| **Charts** | Apache ECharts | 5.5+ | Superior a Recharts para datasets grandes, zoom temporal, tooltips ricos y múltiples tipos de gráficos. Bundle más grande pero justificado para un dashboard analítico profesional. No impacta recursos del servidor. |
| **HTTP Client** | TanStack Query (React Query) | 5 | Cacheo de datos del API, revalidación automática, manejo de estados de carga/error. |
| **Routing** | React Router | 6 | SPA routing simple. |

---

## 4. Estructura de Directorios del Proyecto

```
sales-insight/
├── docker-compose.yml
├── Makefile
├── README.md
├── SPEC.md
├── ARCHITECTURE.md
├── DATABASE_SCHEMA.md
│
├── backend/                          # Aplicación Go (monolito)
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── server/
│   │       └── main.go               # Entrypoint: inicializa API, scheduler, DB, token cache
│   │
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go             # Viper: lectura de env vars
│   │   ├── database/
│   │   │   ├── db.go                 # Conexión PostgreSQL + pool
│   │   │   └── queries/              # Archivos .sql para sqlc
│   │   ├── models/                   # Structs generados por sqlc
│   │   ├── clover/
│   │   │   ├── client.go             # HTTP client para Clover API + rate limiting
│   │   │   ├── auth.go               # OAuth 2.0 flow + token refresh
│   │   │   ├── types.go              # Structs de respuesta de Clover
│   │   │   ├── export.go             # Cliente para Clover Export API
│   │   │   └── rate_limiter.go       # Token bucket rate limiter
│   │   ├── sync/
│   │   │   ├── engine.go             # Orquestador de sincronización
│   │   │   ├── orders.go             # Lógica de extracción de órdenes
│   │   │   ├── items.go              # Lógica de extracción de productos
│   │   │   ├── employees.go          # Lógica de extracción de empleados
│   │   │   ├── payments.go           # Lógica de extracción de pagos
│   │   │   ├── export_sync.go        # Lógica para descargar archivos Export API
│   │   │   └── transform.go          # Normalización de datos crudos → modelo propio
│   │   ├── scheduler/
│   │   │   └── scheduler.go          # Tareas cron: polling periódico + refresh de tokens
│   │   ├── api/
│   │   │   ├── server.go             # Setup de Echo, middlewares, routes
│   │   │   ├── handlers/
│   │   │   │   ├── auth.go           # Callback OAuth, estado de conexión, bootstrap manual
│   │   │   │   ├── dashboard.go      # Endpoints de métricas
│   │   │   │   ├── employees.go      # CRUD empleados + métricas
│   │   │   │   ├── products.go       # Listado + mapeo de categorías
│   │   │   │   ├── orders.go         # Listado y detalle de órdenes
│   │   │   │   └── sync_status.go    # Estado de sincronización y logs
│   │   │   └── middleware/
│   │   │       └── auth.go           # Validación JWT de sesiones del dashboard
│   │   ├── encryption/
│   │   │   └── encryption.go         # AES-GCM para tokens en reposo
│   │   ├── token_cache/
│   │   │   └── cache.go              # Caché en memoria de tokens (reconstruida desde DB al iniciar)
│   │   └── logger/
│   │       └── logger.go             # Wrapper de slog con redacción de PII
│   │
│   └── migrations/
│       ├── 001_initial_schema.up.sql
│       ├── 001_initial_schema.down.sql
│       └── ...
│
├── frontend/                         # Aplicación React
│   ├── Dockerfile
│   ├── package.json
│   ├── vite.config.ts
│   ├── index.html
│   ├── tailwind.config.js            # Desktop-first, sin breakpoints móviles en MVP
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       ├── api/
│       │   └── client.ts             # Axios/fetch configurado con base URL
│       ├── hooks/
│       │   └── useDashboard.ts       # TanStack Query wrappers
│       ├── components/
│       │   ├── Layout.tsx
│       │   ├── Sidebar.tsx
│       │   ├── DateRangePicker.tsx
│       │   ├── MetricCard.tsx
│       │   └── charts/
│       │       ├── SalesByEmployee.tsx
│       │       ├── TopProducts.tsx
│       │       ├── TicketAnalysis.tsx
│       │       └── SalesTrend.tsx
│       ├── pages/
│       │   ├── Dashboard.tsx
│       │   ├── Employees.tsx
│       │   ├── Products.tsx
│       │   ├── Orders.tsx
│       │   └── Settings.tsx
│       └── types/
│           └── index.ts
│
└── nginx/                            # Pendiente: solo para despliegue en DO
    └── nginx.conf                    # Configuración de proxy reverso + SSL
```

---

## 5. Flujos de Datos

### 5.1 Autenticación OAuth 2.0

```
Admin Browser ──► Go API (GET /auth/clover)
     │                              │
     │                              ▼
     │                    Redirect a Clover OAuth
     │                    (client_id, redirect_uri, scope)
     │
     │◄──────────────────────────────────────────────
     │    Clover redirige a /auth/clover/callback?code=...
     │    (Requiere endpoint público accesible)
     │
     ▼
Go API ──► Intercambia code por tokens (POST Clover)
     │
     ▼
Encripta tokens ──► Guarda en PostgreSQL (tabla restaurant_tokens)
     │
     ▼
Carga en memoria (token_cache) para requests subsiguientes
     │
     ▼
Inicia sincronización completa inicial (polling o Export API)
```

**Estrategia para entorno local sin IP pública:**
- **Opción A (recomendada):** Desplegar temporalmente en un VPS o usar ngrok para obtener los tokens una sola vez. Luego copiar los tokens encriptados de PostgreSQL al entorno local.
- **Opción B:** Implementar un endpoint manual de bootstrap (`POST /auth/bootstrap`) que reciba `access_token` y `refresh_token` obtenidos externamente (ej. desde Postman) y los almacene encriptados.

**Scopes de Clover requeridos:**
- `ORDERS_R` — lectura de órdenes
- `INVENTORY_R` — lectura de productos/ítems
- `EMPLOYEES_R` — lectura de empleados
- `MERCHANT_R` — lectura de información del merchant
- `PAYMENTS_R` — lectura de pagos

### 5.2 Sincronización por Polling

El scheduler interno ejecuta tareas periódicas:

| Tarea | Frecuencia | Descripción |
|---|---|---|
| `sync_orders` | Cada 5 minutos | Extrae órdenes modificadas desde la última sync. Usa `modifiedTime` como cursor. |
| `sync_items` | Cada 30 minutos | Extrae productos/ítems modificados. |
| `sync_employees` | Cada 1 hora | Extrae empleados (baja frecuencia de cambio). |
| `sync_payments` | Cada 15 minutos | Extrae pagos nuevos. |
| `refresh_tokens` | Dinámico | Refresca `access_token` antes de que expire (basado en `access_token_expiration`). |
| `export_backfill` | Bajo demanda | Solicita exportación masiva vía Export API para períodos >2 meses. |

**Estrategia de cursor:**
Cada tabla de sync guarda el `modifiedTime` más alto procesado. La siguiente extracción consulta `?filter=modifiedTime>=ULTIMO_CURSOR`, minimizando transferencia de datos y respetando rate limits.

**Rate limiting client-side:**
El `clover.Client` usa un token bucket (`golang.org/x/time/rate`) configurado con:
- `Limit: 16` (requests por segundo, por token)
- `Burst: 5` (requests concurrentes máximos, por token)

Ante un `429`:
1. Leer header `retry-after`.
2. Si no existe, pausar 1 segundo.
3. Reintentar con backoff exponencial (`sleep((2^attempt) + random)`).
4. Loggear el evento en `sync_errors`.

### 5.3 Carga Inicial Masiva (Backfill)

Para restaurantes con historial extenso:

```
Admin solicita backfill desde UI
        │
        ▼
Go API crea export job vía /v3/merchants/{id}/exports
        │
        ▼
Polea estado del export hasta DONE
        │
        ▼
Descarga archivos JSON (máx. 1,000 objetos por archivo)
        │
        ▼
Procesa e inserta en PostgreSQL por batches
```

**Consideraciones Export API:**
- Disponible solo en horarios específicos (US: lun-vie 10:00-15:00 UTC).
- Sandbox no tiene restricción horaria.
- Archivos se eliminan después de 24h.
- Rango máximo: 30 días por export.
- Requiere permisos previos de `developer-relations@clover.com`.

### 5.4 Procesamiento de Datos (ETL Ligero)

```
Datos crudos Clover API
        │
        ▼
┌───────────────┐
│ Normalización │  • Convertir centavos a unidad monetaria
│               │  • Parsear timestamps a UTC
│               │  • Mapear IDs de Clover a UUIDs internos
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ Categorización│  • Lee `category_id` nativo de Clover
│ Analítica     │  • Aplica mapeo configurado (tabla category_mappings)
│               │  • Determina: entrada / plato_principal / postre / bebida / otro
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ Persistencia  │  • INSERT/UPDATE idempotente (ON CONFLICT)
│ PostgreSQL    │  • Transacciones por batch (100-500 registros)
└───────────────┘
        │
        ▼
┌───────────────┐
│ Agregaciones  │  • Las métricas del dashboard se calculan en tiempo real
│ (Runtime)     │    mediante queries SQL agregadas + índices.
│               │  • No hay tablas de pre-agregación en el MVP.
└───────────────┘
```

---

## 6. Decisiones Arquitectónicas Clave (ADRs)

### ADR-01: Monolito vs. Microservicios
**Decisión:** Monolito en un solo binario Go.  
**Justificación:** Un solo restaurante, un solo equipo, ambiente local. La complejidad de red, deploys y monitoreo de microservicios no se justifica. La separación interna por packages (`internal/api`, `internal/sync`, `internal/scheduler`) permite extracción futura si escala.

### ADR-02: Go vs. Python
**Decisión:** Go como lenguaje principal del backend.  
**Justificación:** Menor consumo de RAM, binario único sin runtime, concurrencia nativa para polling eficiente. Python se reserva para futuros microservicios de análisis estadístico (ML/forecasting) si se requiere, comunicados vía HTTP/gRPC desde Go.

### ADR-03: Sin Redis / Sin Message Broker
**Decisión:** No se incluye Redis, RabbitMQ, Kafka ni similar en el MVP.  
**Justificación:**
- Volumen de datos: ~500 órdenes/día = insignificante para PostgreSQL.
- Concurrencia de lectura: 1-5 usuarios del dashboard.
- Colas de trabajo: Un solo restaurante no genera suficiente carga como para necesitar un broker. El scheduler de Go procesa secuencialmente en goroutines ligeras.
- Recursos: En PC local, se optimizan recursos para PostgreSQL + Go.
**Revisión:** Reevaluar cuando se conecten 5+ locales simultáneamente.

### ADR-04: sqlc vs. ORM (GORM)
**Decisión:** sqlc en lugar de GORM u ORM completo.  
**Justificación:**
- Genera código Go tipado en compile time a partir de queries SQL explícitas.
- Sin reflexión en runtime = menor consumo de CPU y RAM.
- Queries analíticas complejas (GROUP BY, window functions) se escriben en SQL nativo, no en DSL de ORM.
- Migraciones con `golang-migrate` mantienen el esquema versionado.

### ADR-05: React + Vite vs. Next.js
**Decisión:** React con Vite (SPA), sin Next.js.  
**Justificación:**
- No se requiere SSR (Server-Side Rendering): el dashboard es una aplicación interna, no necesita SEO.
- Vite genera archivos estáticos que se pueden servir directamente (por Go o Nginx en DO), sin necesidad de un servidor Node.js corriendo.
- Menor consumo de recursos en build y en runtime.
- TanStack Query maneja el fetching y cacheo de datos del API eficientemente.

### ADR-06: Solo Polling (Webhooks Post-MVP)
**Decisión:** No se usan webhooks en el MVP. Solo polling periódico.  
**Justificación:**
- Webhooks requieren un endpoint HTTPS público. El MVP corre en PC local sin IP pública.
- El polling cada 5 minutos es suficiente para un dashboard de análisis; no se necesita latencia de milisegundos.
- Reduce complejidad operativa: no hay que validar firmas de webhooks ni manejar reintentos de Clover.
- Los webhooks se habilitarán como mejora cuando se migre a DigitalOcean con dominio fijo.

### ADR-07: Sin Tablas de Pre-Agregación
**Decisión:** Las métricas del dashboard se calculan con queries SQL agregadas en tiempo real.  
**Justificación:**
- Con el volumen de datos de 1 restaurante, PostgreSQL responde en <100ms con índices adecuados.
- Evita complejidad de mantener materialized views o jobs de pre-computo.
- Si en el futuro el volumen o la lentitud lo justifican, se pueden agregar índices compuestos, particiones o materialized views sin cambiar la arquitectura.

### ADR-08: Token Cache en Memoria + PostgreSQL
**Decisión:** Los tokens se almacenan en PostgreSQL (encriptados) como fuente de verdad, y se cachean en memoria del proceso Go.  
**Justificación:**
- PostgreSQL garantiza persistencia ante reinicios.
- La caché en memoria evita una query a la BDD en cada request a Clover (decenas de requests por minuto).
- Al iniciar, el proceso reconstruye la caché desde PostgreSQL automáticamente.
- Si se escala a múltiples réplicas del backend en el futuro, la fuente de verdad sigue siendo PostgreSQL.

### ADR-09: Apache ECharts vs. Recharts
**Decisión:** Apache ECharts para visualizaciones.  
**Justificación:**
- Recharts es más simple pero se degrada con datasets grandes.
- Sales Insight necesitará analizar meses de datos, comparar múltiples empleados, y hacer zoom en series temporales.
- ECharts soporta zoom, brush, múltiples ejes, y decenas de miles de puntos sin problemas.
- El bundle extra corre en el navegador del usuario; no impacta el servidor.

### ADR-10: Entorno Local (Docker) vs. DO Droplet
**Decisión:** El MVP corre localmente en Docker Compose. DO Droplet se deja para migración post-MVP.  
**Justificación:**
- Permite desarrollo y validación con datos productivos sin costo de infraestructura.
- Facilita iteraciones rápidas sin deploys remotos.
- Nginx y SSL se posponen hasta la migración a nube, reduciendo componentes en local.
- El código está preparado para migrar: basta agregar Nginx y un dominio.

---

## 7. Seguridad

### 7.1 Medidas por Capa

| Capa | Medida |
|---|---|
| **Transporte (post-MVP)** | HTTPS obligatorio (Let's Encrypt vía Nginx). Redirect 80→443. |
| **Transporte (local)** | HTTP permitido. Sin SSL requerido en desarrollo local. |
| **Tokens OAuth** | Encriptación AES-256-GCM en reposo. Clave maestra en variable de entorno `ENCRYPTION_KEY`. |
| **Token Cache** | En memoria del proceso Go (no accesible externamente). Se reconstruye desde DB encriptada al iniciar. |
| **API Dashboard** | JWT para sesiones de usuarios del frontend. Expiración corta (24h). Refresh token rotativo. |
| **Base de datos** | PostgreSQL accesible solo desde la red interna de Docker. Sin puerto expuesto al exterior del host. |
| **Logs** | Sin loguear tokens, contraseñas ni PII sensible. `slog` con redacción de campos sensibles. |
| **CORS** | Configurado en Echo para permitir solo el origen del frontend (localhost:5173 en dev). |
| **Rate Limiting** | Cliente HTTP con token bucket para Clover. Echo limita requests a endpoints de auth para prevenir fuerza bruta. |
| **Inyección (OWASP A03)** | sqlc genera queries parametrizadas. Sin concatenación de strings SQL. |
| **Configuración (OWASP A05)** | API keys en `.env`, nunca en repositorio. `.env.example` con valores dummy. Endpoints de debug desactivados en producción. |
| **Integridad (OWASP A08)** | Validación de tokens OAuth por firma (cuando se habiliten webhooks post-MVP, validar `X-Clover-Auth`). |

### 7.2 OWASP Top 10 — Mapeo a Sales Insight

| Riesgo | Mitigación en Sales Insight |
|---|---|
| **A01: Control de acceso roto** | Endpoints del dashboard protegidos por JWT. Sin exposición de `merchant_id` en URLs públicas. |
| **A02: Fallos criptográficos** | Tokens OAuth encriptados (AES-GCM). HTTPS en producción. |
| **A03: Inyección** | sqlc con queries parametrizadas. Validación de inputs con Echo validator. |
| **A04: Diseño inseguro** | Segregación de datos por `restaurant_id` preparada desde el esquema. |
| **A05: Configuración incorrecta** | Config vía env vars. Sin defaults inseguros. Dependencias auditadas con `govulncheck`. |
| **A06: Componentes vulnerables** | Go modules auditados. Actualizaciones regulares. |
| **A07: Fallos de autenticación** | OAuth 2.0 con Clover. Refresh automático. Sin almacenamiento de contraseñas locales. |
| **A08: Integridad de datos** | Validación de checksums en descargas Export API. Firmar releases. |
| **A09: Fallos en logging** | Logs estructurados de todas las llamadas a Clover. Alertas ante 3+ errores consecutivos. |

---

## 8. Escalabilidad y Multi-Local (Preparación)

Aunque el MVP opera con 1 solo `restaurant_id`, la arquitectura está preparada para escalar a 15+ locales con cambios mínimos:

1. **Esquema de base de datos:** Todas las tablas incluyen `restaurant_id` como FK. Los queries del dashboard ya filtran por este campo.
2. **Tokens OAuth:** La tabla `restaurant_tokens` almacena múltiples filas, una por merchant conectado.
3. **Token Cache:** Se extiende a un mapa en memoria indexado por `restaurant_id`.
4. **Scheduler:** El cron puede iterar sobre una lista de `restaurant_id` activos, ejecutando la sync de cada uno secuencialmente (o en goroutines paralelas controladas con semáforo).
5. **Rate Limiting:** El token bucket se instancia por `restaurant_id`, respetando los límites per-token de Clover.
6. **Infraestructura:** DigitalOcean permite escalar el droplet verticalmente o migrar a múltiples droplets con un load balancer cuando el monolito ya no alcance.
7. **Webhooks:** Al migrar a DO con dominio fijo, se habilitan webhooks para reducir la frecuencia de polling y la carga en la API de Clover.

---

## 9. Plan de Despliegue

### 9.1 Despliegue Local (MVP)

**Requisitos previos:**
- Docker Desktop instalado.
- App registrada en Clover Developer Dashboard con `client_id` y `client_secret`.
- *(Opcional)* Tokens OAuth obtenidos previamente (si no se usará ngrok para el callback).

**Pasos:**

```bash
# 1. Clonar el repositorio
$ git clone <repo> && cd sales-insight

# 2. Configurar variables de entorno
$ cp backend/.env.example backend/.env
$ nano backend/.env
# CLOVER_CLIENT_ID=xxx
# CLOVER_CLIENT_SECRET=xxx
# CLOVER_ENV=production    # o sandbox
# ENCRYPTION_KEY=una-clave-segura-de-32-caracteres-exactos
# DATABASE_URL=postgres://user:pass@db:5432/sales_insight?sslmode=disable
# JWT_SECRET=otra-clave-segura
# API_PORT=8080
# FRONTEND_URL=http://localhost:5173

# 3. Levantar stack
$ docker compose up -d --build

# 4. Ejecutar migraciones
$ docker compose exec backend migrate -path /migrations -database "$DATABASE_URL" up

# 5. Verificar estado
$ docker compose ps
$ docker compose logs -f backend

# 6. Acceder al dashboard
# Frontend: http://localhost:5173
# API:      http://localhost:8080
```

**Sobre OAuth en local:**
- Si tienes tokens válidos, usa el endpoint de bootstrap manual para insertarlos.
- Si necesitas hacer el flujo completo, ejecuta temporalmente ngrok: `ngrok http 8080`, configura la URL en Clover Developer Dashboard, realiza el OAuth, y luego detén ngrok.

### 9.2 Migración a DigitalOcean (Post-MVP)

**Componentes a agregar:**
1. **Nginx:** Proxy reverso, compresión gzip, serve de estáticos.
2. **SSL:** Certbot + Let's Encrypt para dominio propio.
3. **Dominio:** Configurar DNS A record apuntando al droplet.
4. **Webhooks:** Habilitar endpoint público `/webhooks/clover` y suscribir eventos en Clover Developer Dashboard.
5. **Firewall:** UFW en Ubuntu, permitir solo 80/443.

**Pasos:**

```bash
# 1. Crear droplet en DO (Ubuntu 22.04, 2GB RAM mínimo)
# 2. Configurar DNS y esperar propagación
# 3. Copiar .env de local (ajustar URLs y DATABASE_URL)
# 4. docker compose up -d --build
# 5. Configurar SSL
$ docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
# 6. Configurar webhooks en Clover Dashboard apuntando a https://tu-dominio.com/webhooks/clover
```

---

## 10. Interfaces y Contratos

### 10.1 Entrada/Salida con Clover
| Dirección | Protocolo | Formato | Descripción |
|---|---|---|---|
| Salida | HTTPS | JSON | Peticiones OAuth y REST API a Clover. |
| Entrada | HTTPS | JSON | *(Post-MVP)* Webhooks enviados por Clover a nuestro endpoint. |
| Salida | HTTPS | JSON | Export API (solicitud de archivos bulk). |

### 10.2 Entrada/Salida con Frontend
| Dirección | Protocolo | Formato | Descripción |
|---|---|---|---|
| Entrada/Salida | HTTP/HTTPS | JSON | API REST (`/api/v1/*`) consumida por React SPA. |
| Salida | HTTP | Static | Go sirve archivos estáticos del build de React (en local) o Nginx (en DO). |

---

*Documento vivo. Las decisiones arquitectónicas pueden ser revisadas cuando el sistema escale o los requerimientos cambien.*
