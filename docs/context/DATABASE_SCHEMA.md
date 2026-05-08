# DATABASE_SCHEMA.md

## Sales Insight — Esquema de Base de Datos

**Versión:** 1.0-MVP  
**Fecha:** 2026-05-06  
**Motor:** PostgreSQL 15  
**Migrations:** golang-migrate  
**Acceso:** sqlc  

---

## 1. Principios del Esquema

1. **Multi-local preparado:** Todas las tablas relevantes incluyen `restaurant_id` como foreign key, aunque el MVP opere con un único registro.
2. **Trazabilidad Clover:** Cada entidad extraída de Clover conserva su identificador original (`clover_*_id`) como unique constraint, permitiendo deduplicación y sincronización incremental.
3. **Categorización analítica:** Cada restaurante define sus propias categorías analíticas (ej. `Tacos`, `Bebidas`, `Postres`) y mapea las categorías nativas de Clover hacia ellas mediante una tabla de configuración.
4. **Temporalidad:** Todos los registros llevan `created_at`, `updated_at` y, cuando aplica, `synced_at` para auditar la última sincronización exitosa.
5. **Índices estratégicos:** Las consultas del dashboard operan sobre rangos de fecha y agregaciones por `restaurant_id`, `employee_id` y `product_id`. Los índices cubren estos patrones de acceso.

---

## 2. Diagrama Entidad-Relación (Simplificado)

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   restaurants   │       │  employees      │       │   categories    │
│   (tokens)      │       │  (meseros)      │       │  (clover nativas)│
└────────┬────────┘       └────────┬────────┘       └────────┬────────┘
         │ 1:N                     │ 1:N                     │ 1:1
         │                         │                         │
┌────────▼────────┐       ┌────────▼────────┐       ┌────────▼────────┐
│     orders      │       │   order_items   │       │category_mappings│
│   (ventas)      │       │  (líneas)       │       │  (mapeo)        │
└────────┬────────┘       └────────┬────────┘       └────────┬────────┘
         │ 1:N                     │ N:1                     │ N:1
         │                         │                         │
┌────────▼────────┐       ┌────────▼────────┐       ┌─────────────────┐
│    payments     │       │    products     │       │  analytic_types │
│   (pagos)       │       │  (platos)       │       │  (enum virtual) │
└─────────────────┘       └─────────────────┘       └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│   sync_logs     │       │  sync_errors    │
│ (auditoría)     │       │  (errores)      │
└─────────────────┘       └─────────────────┘
```

---

## 3. Tablas

### 3.1 `restaurants`

Información del merchant conectado y tokens OAuth encriptados.

```sql
CREATE TABLE restaurants (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255), -- se obtiene de Clover tras el OAuth
    clover_merchant_id  VARCHAR(32) UNIQUE, -- se completa tras el OAuth

    -- Tokens encriptados (AES-256-GCM). NULL hasta completar el primer OAuth.
    access_token        BYTEA,
    refresh_token       BYTEA,

    -- Metadatos del token (no encriptados, útiles para diagnóstico)
    access_token_expires_at TIMESTAMPTZ,
    refresh_token_expires_at TIMESTAMPTZ,

    -- Configuración
    clover_env          VARCHAR(16) NOT NULL DEFAULT 'sandbox', -- 'sandbox' | 'production'
    is_connected        BOOLEAN NOT NULL DEFAULT false, -- true tras OAuth exitoso
    is_active           BOOLEAN NOT NULL DEFAULT true,

    -- Reglas configurables de "ticket ideal" (por restaurante). NULL = no configurado.
    -- Ejemplo: {"tacos": 2, "bebidas": 1} donde la key es el slug de analytic_categories.
    -- Permite que cada restaurante defina qué constituye un ticket completo para su menú.
    ticket_completo_rules JSONB,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_restaurants_clover_merchant_id ON restaurants(clover_merchant_id);
CREATE INDEX idx_restaurants_is_connected ON restaurants(is_connected) WHERE is_connected = true;
CREATE INDEX idx_restaurants_is_active ON restaurants(is_active) WHERE is_active = true;
CREATE INDEX idx_restaurants_ticket_rules ON restaurants USING GIN(ticket_completo_rules) WHERE ticket_completo_rules IS NOT NULL;
```

**Notas:**
- `name` y `clover_merchant_id` pueden ser NULL inicialmente. Se completan tras el primer OAuth exitoso.
- `access_token` y `refresh_token` son NULL hasta la conexión. Se almacenan como `BYTEA` (binario encriptado).
- La clave de encriptación (`ENCRYPTION_KEY`) vive en variable de entorno, nunca en la base de datos.
- El backend reconstruye su caché en memoria desde esta tabla al iniciar, filtrando `is_connected = true`.
- **Estrategia OAuth:** La fila se crea *vacía* al iniciar la app (o durante el setup inicial) con `is_connected = false`. Tras el callback OAuth exitoso, se actualizan `clover_merchant_id`, tokens encriptados y `is_connected = true`.
- PostgreSQL crea automáticamente un índice único en la columna `id` por ser PRIMARY KEY. No es necesario (ni recomendable) crear un `CREATE INDEX` adicional.

---

### 3.2 `employees`

Empleados (meseros, cajeros, etc.) extraídos de Clover.

```sql
CREATE TABLE employees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- Identificador Clover
    clover_employee_id  VARCHAR(32) NOT NULL,

    -- Datos básicos
    name                VARCHAR(255) NOT NULL,
    email               VARCHAR(255),
    phone               VARCHAR(50),
    role                VARCHAR(50), -- 'admin', 'manager', 'employee', etc.
    is_active           BOOLEAN NOT NULL DEFAULT true,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Constraints
    UNIQUE(restaurant_id, clover_employee_id)
);

CREATE INDEX idx_employees_restaurant_id ON employees(restaurant_id);
CREATE INDEX idx_employees_clover_employee_id ON employees(clover_employee_id);
CREATE INDEX idx_employees_is_active ON employees(is_active) WHERE is_active = true;
```

---

### 3.3 `categories`

Categorías de productos tal como vienen de Clover (nativas).

```sql
CREATE TABLE categories (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- Identificador Clover
    clover_category_id  VARCHAR(32) NOT NULL,

    -- Datos nativos
    name                VARCHAR(255) NOT NULL,
    sort_order          INTEGER,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(restaurant_id, clover_category_id)
);

CREATE INDEX idx_categories_restaurant_id ON categories(restaurant_id);
```

---

### 3.4 `category_mappings`

Mapeo de categorías nativas de Clover a categorías analíticas del sistema.

```sql
CREATE TABLE analytic_categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id   UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name            VARCHAR(50) NOT NULL,        -- ej. "Tacos", "Bebidas", "Postres"
    slug            VARCHAR(50) NOT NULL,        -- ej. "tacos", "bebidas" (para queries y JSONB keys)
    display_order   INTEGER NOT NULL DEFAULT 0,  -- orden en dashboard
    color           VARCHAR(7) DEFAULT '#3B82F6', -- color en gráficos
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, slug)
);

CREATE INDEX idx_analytic_categories_restaurant_id ON analytic_categories(restaurant_id);
CREATE INDEX idx_analytic_categories_is_active ON analytic_categories(is_active) WHERE is_active = true;

CREATE TABLE category_mappings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    category_id         UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    analytic_category_id UUID NOT NULL REFERENCES analytic_categories(id) ON DELETE CASCADE,

    -- Auditoría del mapeo
    mapped_by           VARCHAR(255), -- usuario que configuró el mapeo
    mapped_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(restaurant_id, category_id)
);

CREATE INDEX idx_category_mappings_restaurant_id ON category_mappings(restaurant_id);
CREATE INDEX idx_category_mappings_analytic ON category_mappings(analytic_category_id);
```

**Notas:**
- Si una categoría Clover no tiene mapeo, se le asigna una categoría analítica genérica (ej. "Otros") creada automáticamente por el sistema.
- El administrador configura este mapeo desde el dashboard. Un cambio aquí dispara el recálculo de métricas de ticket.

---

### 3.5 `products`

Productos / ítems del menú extraídos de Clover.

```sql
CREATE TABLE products (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- Identificador Clover
    clover_item_id      VARCHAR(32) NOT NULL,

    -- Datos básicos
    name                VARCHAR(255) NOT NULL,
    description         TEXT,
    price               BIGINT NOT NULL, -- en centavos (ej. 1250 = $12.50)
    cost                BIGINT,          -- en centavos, si está disponible
    sku                 VARCHAR(100), -- Stock Keeping Unit: código interno de inventario del producto

    -- Relaciones
    category_id         UUID REFERENCES categories(id) ON DELETE SET NULL,

    -- Estado
    is_available        BOOLEAN NOT NULL DEFAULT true,
    is_deleted          BOOLEAN NOT NULL DEFAULT false, -- soft delete si Clover marca como eliminado

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(restaurant_id, clover_item_id)
);

CREATE INDEX idx_products_restaurant_id ON products(restaurant_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_clover_item_id ON products(clover_item_id);
CREATE INDEX idx_products_is_available ON products(is_available) WHERE is_available = true;
```

**Notas:**
- `price` en `BIGINT` de centavos evita errores de punto flotante en cálculos monetarios.
- `is_deleted` permite detectar ítems eliminados en Clover sin borrar histórico de ventas.

---

### 3.6 `orders`

Órdenes / tickets de venta extraídos de Clover.

```sql
CREATE TABLE orders (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- Identificador Clover
    clover_order_id     VARCHAR(32) NOT NULL,

    -- Relaciones
    employee_id         UUID REFERENCES employees(id) ON DELETE SET NULL,

    -- Datos del ticket
    order_number        VARCHAR(50),
    order_type          VARCHAR(50), -- 'online', 'in_house', 'pickup', 'delivery'
    state               VARCHAR(50) NOT NULL DEFAULT 'open', -- 'open', 'locked', 'paid', 'refunded'

    -- Totales (en centavos)
    total_amount        BIGINT NOT NULL DEFAULT 0,
    tax_amount          BIGINT NOT NULL DEFAULT 0,
    discount_amount     BIGINT NOT NULL DEFAULT 0,
    tip_amount          BIGINT NOT NULL DEFAULT 0,

    -- Fechas del ticket (horario del restaurante)
    created_time        TIMESTAMPTZ NOT NULL,
    modified_time       TIMESTAMPTZ NOT NULL,
    pay_type            VARCHAR(50), -- 'full', 'split'

    -- Campos calculados para análisis (se actualizan después de insertar order_items)
    -- NOTA: Eliminados los campos fijos 'has_*' e 'is_ticket_completo' porque asumen
    -- una estructura de menú (entrada+principal+postre+bebida) que no todos los
    -- restaurantes siguen. Ver 'order_category_summary' para análisis flexible.
    item_count              INTEGER NOT NULL DEFAULT 0,           -- total de líneas
    unique_category_count   INTEGER NOT NULL DEFAULT 0,           -- categorías analíticas distintas
    total_quantity          INTEGER NOT NULL DEFAULT 0,           -- suma de cantidades

    -- Timestamps del sistema
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(restaurant_id, clover_order_id)
);

CREATE INDEX idx_orders_restaurant_id ON orders(restaurant_id);
CREATE INDEX idx_orders_employee_id ON orders(employee_id);
CREATE INDEX idx_orders_created_time ON orders(created_time);
CREATE INDEX idx_orders_modified_time ON orders(modified_time);
CREATE INDEX idx_orders_state ON orders(state);

-- Índice compuesto clave para queries de dashboard: ventas por empleado en rango de fechas
CREATE INDEX idx_orders_employee_created ON orders(employee_id, created_time);

-- Índice para análisis de cobertura de categorías (upselling) por empleado
CREATE INDEX idx_orders_employee_categories ON orders(employee_id, unique_category_count, created_time);
```

**Notas:**
- Los campos `has_*` y `is_ticket_completo` se calculan durante el ETL (después de insertar los `order_items`). Permiten filtrar rápidamente sin hacer JOINs complejos.
- `created_time` es el timestamp original de Clover (no el `created_at` del sistema). Se usa para análisis de ventas por hora/día.

---

### 3.7 `order_items`

Líneas de cada orden (desagregación del ticket por producto).

```sql
CREATE TABLE order_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- Relaciones
    order_id            UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id          UUID REFERENCES products(id) ON DELETE SET NULL,

    -- Identificador Clover (si existe)
    clover_line_item_id VARCHAR(32),

    -- Datos de la línea
    name                VARCHAR(255) NOT NULL, -- nombre del producto en el momento de la venta
    quantity            INTEGER NOT NULL DEFAULT 1,
    unit_price          BIGINT NOT NULL, -- en centavos
    total_price         BIGINT NOT NULL, -- quantity * unit_price

    -- Categoría analítica calculada en el momento (histórica)
    -- Si el mapeo cambia en el futuro, estos registros conservan la categoría original
    analytic_category_id UUID REFERENCES analytic_categories(id) ON DELETE SET NULL,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_items_analytic_category ON order_items(analytic_category_id);

-- Índice para "productos más vendidos" por rango de fechas
CREATE INDEX idx_order_items_product_created ON order_items(product_id, created_at);
```

**Notas:**
- `name` se replica aquí para preservar el nombre histórico, aunque el producto cambie de nombre en Clover.
- `analytic_category_id` se denormaliza para facilitar agregaciones sin JOIN adicional a `products` → `categories` → `category_mappings` → `analytic_categories`.

---

### 3.8 `order_category_summary`

Desglose de cada orden por categoría analítica. Permite calcular métricas de cobertura sin depender de columnas fijas en `orders`.

```sql
CREATE TABLE order_category_summary (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    order_id            UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

    analytic_category_id UUID NOT NULL REFERENCES analytic_categories(id) ON DELETE CASCADE,
    item_count          INTEGER NOT NULL DEFAULT 0,  -- cuántas líneas de esa categoría
    total_quantity      INTEGER NOT NULL DEFAULT 0,  -- suma de cantidades
    total_amount        BIGINT NOT NULL DEFAULT 0,   -- suma de total_price en centavos

    UNIQUE(restaurant_id, order_id, analytic_category_id)
);

CREATE INDEX idx_order_category_summary_order_id ON order_category_summary(order_id);
CREATE INDEX idx_order_category_summary_category ON order_category_summary(analytic_category_id);
CREATE INDEX idx_order_category_summary_restaurant ON order_category_summary(restaurant_id, analytic_category_id);
```

**Notas:**
- Esta tabla se recalcula/actualiza cada vez que se sincroniza una orden.
- Permite responder: "¿Qué % de órdenes incluye bebidas?" sin importar si el restaurante tiene entradas o postres.
- Facilita comparar cobertura entre empleados: "El mesero A incluye postre en 60% de sus órdenes".

---

---

### 3.8 `payments`

Pagos asociados a órdenes.

```sql
CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    -- Identificador Clover
    clover_payment_id   VARCHAR(32) NOT NULL,

    -- Relaciones
    order_id            UUID REFERENCES orders(id) ON DELETE SET NULL,
    employee_id         UUID REFERENCES employees(id) ON DELETE SET NULL,

    -- Datos del pago
    amount              BIGINT NOT NULL, -- en centavos
    tip_amount          BIGINT NOT NULL DEFAULT 0,
    tax_amount          BIGINT NOT NULL DEFAULT 0,
    payment_type        VARCHAR(50), -- 'cash', 'credit_card', 'debit_card', 'gift_card', 'other'
    card_type           VARCHAR(50), -- 'VISA', 'MASTERCARD', etc.
    result              VARCHAR(50) NOT NULL DEFAULT 'success', -- 'success', 'fail', 'refund'
    external_payment_id VARCHAR(255), -- ID de pasarela externa si aplica

    -- Fechas
    created_time        TIMESTAMPTZ NOT NULL,
    modified_time       TIMESTAMPTZ NOT NULL,

    -- Timestamps
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(restaurant_id, clover_payment_id)
);

CREATE INDEX idx_payments_restaurant_id ON payments(restaurant_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_employee_id ON payments(employee_id);
CREATE INDEX idx_payments_created_time ON payments(created_time);
CREATE INDEX idx_payments_payment_type ON payments(payment_type);
```

---

### 3.9 `sync_logs`

Auditoría de cada ejecución de sincronización.

```sql
CREATE TYPE sync_entity AS ENUM ('orders', 'items', 'employees', 'payments', 'export');

CREATE TABLE sync_logs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    entity              sync_entity NOT NULL,
    status              VARCHAR(20) NOT NULL, -- 'success', 'partial', 'error'

    -- Métricas
    records_processed   INTEGER NOT NULL DEFAULT 0,
    records_inserted    INTEGER NOT NULL DEFAULT 0,
    records_updated     INTEGER NOT NULL DEFAULT 0,
    records_skipped     INTEGER NOT NULL DEFAULT 0,
    records_failed      INTEGER NOT NULL DEFAULT 0,

    -- Cursor de sincronización (para polling incremental)
    cursor_from         TIMESTAMPTZ, -- rango de modifiedTime consultado
    cursor_to           TIMESTAMPTZ,

    -- Duración
    started_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at         TIMESTAMPTZ,

    -- Metadatos
    triggered_by        VARCHAR(50) NOT NULL DEFAULT 'scheduler', -- 'scheduler', 'manual', 'webhook'
    details             JSONB       -- datos adicionales (export_id, etc.)
);

CREATE INDEX idx_sync_logs_restaurant_id ON sync_logs(restaurant_id);
CREATE INDEX idx_sync_logs_entity ON sync_logs(entity);
CREATE INDEX idx_sync_logs_status ON sync_logs(status);
CREATE INDEX idx_sync_logs_started_at ON sync_logs(started_at);
```

---

### 3.10 `sync_errors`

Registro detallado de errores durante la sincronización.

```sql
CREATE TABLE sync_errors (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,

    sync_log_id         UUID REFERENCES sync_logs(id) ON DELETE SET NULL,

    entity              sync_entity NOT NULL,
    clover_id           VARCHAR(32), -- ID del objeto en Clover que falló, si aplica

    -- Error
    error_code          VARCHAR(50),
    error_message       TEXT NOT NULL,
    http_status         INTEGER, -- 429, 500, 503, etc.
    retry_after         INTEGER, -- valor del header retry-after si aplica

    -- Estado de recuperación
    is_resolved         BOOLEAN NOT NULL DEFAULT false,
    resolved_at         TIMESTAMPTZ,
    retry_count         INTEGER NOT NULL DEFAULT 0,

    -- Timestamps
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sync_errors_restaurant_id ON sync_errors(restaurant_id);
CREATE INDEX idx_sync_errors_is_resolved ON sync_errors(is_resolved) WHERE is_resolved = false;
CREATE INDEX idx_sync_errors_entity ON sync_errors(entity);
```

---

## 4. Vistas (Opcional, Post-MVP)

Para el MVP, las métricas se calculan directamente en queries. Sin embargo, se documentan vistas útiles para el futuro:

```sql
-- Vista: métricas diarias por empleado
--CREATE VIEW v_daily_employee_metrics AS
--SELECT
--    o.restaurant_id,
--    o.employee_id,
--    e.name AS employee_name,
--    DATE(o.created_time) AS sale_date,
--    COUNT(*) AS total_orders,
--    SUM(o.total_amount) AS total_sales,
--    AVG(o.total_amount) AS avg_ticket,
--    SUM(CASE WHEN o.is_ticket_completo THEN 1 ELSE 0 END) AS tickets_completos,
--    ROUND(
--        100.0 * SUM(CASE WHEN o.is_ticket_completo THEN 1 ELSE 0 END) / COUNT(*),
--        2
--    ) AS pct_ticket_completo
--FROM orders o
--JOIN employees e ON o.employee_id = e.id
--WHERE o.state IN ('locked', 'paid')
--GROUP BY o.restaurant_id, o.employee_id, e.name, DATE(o.created_time);
```

---

## 5. Estrategia de Migraciones

```
migrations/
├── 001_initial_schema.up.sql      # Crea todas las tablas, tipos e índices
├── 001_initial_schema.down.sql    # Elimina tablas en orden inverso
├── 002_seed_analytic_categories.up.sql  -- Opcional: inserta mapping inicial
```

### Orden de creación (up):
1. `CREATE TYPE sync_entity`
2. `CREATE TABLE restaurants`
3. `CREATE TABLE employees`
4. `CREATE TABLE categories`
5. `CREATE TABLE analytic_categories`
6. `CREATE TABLE category_mappings`
7. `CREATE TABLE products`
8. `CREATE TABLE orders`
9. `CREATE TABLE order_items`
10. `CREATE TABLE order_category_summary`
11. `CREATE TABLE payments`
12. `CREATE TABLE sync_logs`
13. `CREATE TABLE sync_errors`
14. `CREATE INDEX` (todos)

### Orden de eliminación (down):
1. Eliminar índices
2. Eliminar tablas en orden inverso (sync_errors → restaurants)
3. Eliminar tipos ENUM

---

## 6. Decisiones de Diseño

### DD-01: Precios en `BIGINT` de centavos
**Decisión:** Todos los campos monetarios se almacenan como `BIGINT` en centavos.  
**Justificación:** Evita errores de precisión de punto flotante. El frontend divide por 100 para mostrar.

### DD-02: Denormalización de `analytic_category_id` en `order_items`
**Decisión:** `order_items.analytic_category_id` guarda la referencia a la categoría analítica en el momento de la venta.  
**Justificación:** Si el administrador cambia el mapeo de categorías en el futuro, el historial de ventas conserva la categoría original. Esto garantiza reproducibilidad de reportes históricos.

### DD-03: Campos calculados en `orders`
**Decisión:** `item_count`, `unique_category_count`, `total_quantity` (agregados). Eliminados: `has_entrada`, `has_plato_principal`, etc.  
**Justificación:** Los campos fijos `has_*` asumían una estructura de menú universal (entrada+principal+postre+bebida) que no todos los restaurantes siguen. En su lugar, se usan métricas genéricas (`unique_category_count`) más una tabla de desglose (`order_category_summary`) que permite análisis flexible por categoría real del restaurante.

### DD-04: Soft delete en productos
**Decisión:** `products.is_deleted` en lugar de `DELETE CASCADE`.  
**Justificación:** Preserva el historial de ventas. Un producto eliminado en Clover sigue apareciendo en órdenes antiguas.

### DD-05: `created_time` vs. `created_at`
**Decisión:** `created_time` = timestamp original de Clover. `created_at` = timestamp del sistema.  
**Justificación:** Las métricas de negocio (ventas por hora, por día) usan `created_time`. `created_at` sirve para auditoría del sistema.

### DD-06: `clover_*_id` como VARCHAR(32)
**Decisión:** Los IDs de Clover se almacenan como VARCHAR de 32 caracteres, no como UUID.  
**Justificación:** Los IDs de Clover no siguen el formato UUID estándar (son alfanuméricos de longitud variable). Usar VARCHAR evita conversiones fallidas.

### DD-07: Menú y categorías por restaurante
**Decisión:** Tablas `products` y `categories` incluyen `restaurant_id` y sus unique constraints son compuestas `(restaurant_id, clover_*_id)`.  
**Justificación:** Cada restaurante en Clover tiene su propio catálogo de ítems y categorías. El mismo `clover_item_id` puede existir en múltiples merchants, pero representa productos diferentes. Las constraints compuestas garantizan unicidad por local y permiten multi-tenant desde el esquema.

### DD-08: Ticket Ideal Configurable (`ticket_completo_rules`)
**Decisión:** `restaurants.ticket_completo_rules` como `JSONB` opcional (NULL por defecto).  
**Justificación:**
- Permite que cada restaurante defina su propio "ticket ideal" sin forzar una estructura universal.
- Ejemplo café: `{"bebidas": 1, "alimentos": 1}`. Ejemplo taquería: `{"tacos": 2, "bebidas": 1}`.
- Las keys del JSON corresponden al `slug` de `analytic_categories`, permitiendo que cada restaurante use su propio vocabulario.
- Si es NULL, el dashboard omite esta métrica y muestra solo "Cobertura de Categorías" (siempre disponible).
- Se calcula dinámicamente en runtime desde `order_category_summary` + las reglas del restaurante, sin necesidad de columnas fijas en `orders`.

---

*Documento vivo. El esquema puede evolucionar con migraciones versionadas.*
