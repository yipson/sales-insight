-- Migration: 001_initial_schema
-- Crea el esquema completo de base de datos para Sales Insight MVP

-- ============================================================
-- TIPOS ENUM
-- ============================================================

CREATE TYPE sync_entity AS ENUM ('orders', 'items', 'employees', 'payments', 'export');

-- ============================================================
-- TABLAS
-- ============================================================

-- 1. restaurants
CREATE TABLE restaurants (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255),
    clover_merchant_id  VARCHAR(32) UNIQUE,

    access_token        BYTEA,
    refresh_token       BYTEA,

    access_token_expires_at TIMESTAMPTZ,
    refresh_token_expires_at TIMESTAMPTZ,

    clover_env          VARCHAR(16) NOT NULL DEFAULT 'sandbox',
    is_connected        BOOLEAN NOT NULL DEFAULT false,
    is_active           BOOLEAN NOT NULL DEFAULT true,

    -- Reglas configurables de "ticket ideal" (por restaurante). NULL = no configurado.
    -- Ejemplo: {"tacos": 2, "bebidas": 1} donde la key es el slug de analytic_categories.
    -- Permite que cada restaurante defina qué constituye un ticket completo para su menú.
    ticket_completo_rules JSONB,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. employees
CREATE TABLE employees (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    clover_employee_id  VARCHAR(32) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    email               VARCHAR(255),
    phone               VARCHAR(50),
    role                VARCHAR(50),
    is_active           BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, clover_employee_id)
);

-- 3. categories (nativas de Clover)
CREATE TABLE categories (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    clover_category_id  VARCHAR(32) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    sort_order          INTEGER,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, clover_category_id)
);

-- 4. analytic_categories (personalizables por restaurante)
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

-- 5. category_mappings
CREATE TABLE category_mappings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    category_id         UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    analytic_category_id UUID NOT NULL REFERENCES analytic_categories(id) ON DELETE CASCADE,
    mapped_by           VARCHAR(255),
    mapped_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, category_id)
);

-- 6. products
CREATE TABLE products (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    clover_item_id      VARCHAR(32) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,
    price               BIGINT NOT NULL,
    cost                BIGINT,
    sku                 VARCHAR(100),
    category_id         UUID REFERENCES categories(id) ON DELETE SET NULL,
    is_available        BOOLEAN NOT NULL DEFAULT true,
    is_deleted          BOOLEAN NOT NULL DEFAULT false,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, clover_item_id)
);

-- 7. orders
CREATE TABLE orders (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    clover_order_id     VARCHAR(32) NOT NULL,
    employee_id         UUID REFERENCES employees(id) ON DELETE SET NULL,
    order_number        VARCHAR(50),
    order_type          VARCHAR(50),
    state               VARCHAR(50) NOT NULL DEFAULT 'open',
    total_amount        BIGINT NOT NULL DEFAULT 0,
    tax_amount          BIGINT NOT NULL DEFAULT 0,
    discount_amount     BIGINT NOT NULL DEFAULT 0,
    tip_amount          BIGINT NOT NULL DEFAULT 0,
    created_time        TIMESTAMPTZ NOT NULL,
    modified_time       TIMESTAMPTZ NOT NULL,
    pay_type            VARCHAR(50),
    item_count          INTEGER NOT NULL DEFAULT 0,
    unique_category_count INTEGER NOT NULL DEFAULT 0,
    total_quantity      INTEGER NOT NULL DEFAULT 0,
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, clover_order_id)
);

-- 8. order_items
CREATE TABLE order_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    order_id            UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id          UUID REFERENCES products(id) ON DELETE SET NULL,
    clover_line_item_id VARCHAR(32),
    name                VARCHAR(255) NOT NULL,
    quantity            INTEGER NOT NULL DEFAULT 1,
    unit_price          BIGINT NOT NULL,
    total_price         BIGINT NOT NULL,
    analytic_category_id UUID REFERENCES analytic_categories(id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 9. order_category_summary
CREATE TABLE order_category_summary (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    order_id            UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    analytic_category_id UUID NOT NULL REFERENCES analytic_categories(id) ON DELETE CASCADE,
    item_count          INTEGER NOT NULL DEFAULT 0,
    total_quantity      INTEGER NOT NULL DEFAULT 0,
    total_amount        BIGINT NOT NULL DEFAULT 0,
    UNIQUE(restaurant_id, order_id, analytic_category_id)
);

-- 10. payments
CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    clover_payment_id   VARCHAR(32) NOT NULL,
    order_id            UUID REFERENCES orders(id) ON DELETE SET NULL,
    employee_id         UUID REFERENCES employees(id) ON DELETE SET NULL,
    amount              BIGINT NOT NULL,
    tip_amount          BIGINT NOT NULL DEFAULT 0,
    tax_amount          BIGINT NOT NULL DEFAULT 0,
    payment_type        VARCHAR(50),
    card_type           VARCHAR(50),
    result              VARCHAR(50) NOT NULL DEFAULT 'success',
    external_payment_id VARCHAR(255),
    created_time        TIMESTAMPTZ NOT NULL,
    modified_time       TIMESTAMPTZ NOT NULL,
    synced_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(restaurant_id, clover_payment_id)
);

-- 11. sync_logs
CREATE TABLE sync_logs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    entity              sync_entity NOT NULL,
    status              VARCHAR(20) NOT NULL,
    records_processed   INTEGER NOT NULL DEFAULT 0,
    records_inserted    INTEGER NOT NULL DEFAULT 0,
    records_updated     INTEGER NOT NULL DEFAULT 0,
    records_skipped     INTEGER NOT NULL DEFAULT 0,
    records_failed      INTEGER NOT NULL DEFAULT 0,
    cursor_from         TIMESTAMPTZ,
    cursor_to           TIMESTAMPTZ,
    started_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at         TIMESTAMPTZ,
    triggered_by        VARCHAR(50) NOT NULL DEFAULT 'scheduler',
    details             JSONB
);

-- 12. sync_errors
CREATE TABLE sync_errors (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id       UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    sync_log_id         UUID REFERENCES sync_logs(id) ON DELETE SET NULL,
    entity              sync_entity NOT NULL,
    clover_id           VARCHAR(32),
    error_code          VARCHAR(50),
    error_message       TEXT NOT NULL,
    http_status         INTEGER,
    retry_after         INTEGER,
    is_resolved         BOOLEAN NOT NULL DEFAULT false,
    resolved_at         TIMESTAMPTZ,
    retry_count         INTEGER NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- ÍNDICES
-- ============================================================

-- restaurants
CREATE INDEX idx_restaurants_clover_merchant_id ON restaurants(clover_merchant_id);
CREATE INDEX idx_restaurants_is_connected ON restaurants(is_connected) WHERE is_connected = true;
CREATE INDEX idx_restaurants_is_active ON restaurants(is_active) WHERE is_active = true;
CREATE INDEX idx_restaurants_ticket_rules ON restaurants USING GIN(ticket_completo_rules) WHERE ticket_completo_rules IS NOT NULL;

-- employees
CREATE INDEX idx_employees_restaurant_id ON employees(restaurant_id);
CREATE INDEX idx_employees_clover_employee_id ON employees(clover_employee_id);
CREATE INDEX idx_employees_is_active ON employees(is_active) WHERE is_active = true;

-- categories
CREATE INDEX idx_categories_restaurant_id ON categories(restaurant_id);

-- analytic_categories
CREATE INDEX idx_analytic_categories_restaurant_id ON analytic_categories(restaurant_id);
CREATE INDEX idx_analytic_categories_is_active ON analytic_categories(is_active) WHERE is_active = true;

-- category_mappings
CREATE INDEX idx_category_mappings_restaurant_id ON category_mappings(restaurant_id);
CREATE INDEX idx_category_mappings_analytic ON category_mappings(analytic_category_id);

-- products
CREATE INDEX idx_products_restaurant_id ON products(restaurant_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_clover_item_id ON products(clover_item_id);
CREATE INDEX idx_products_is_available ON products(is_available) WHERE is_available = true;

-- orders
CREATE INDEX idx_orders_restaurant_id ON orders(restaurant_id);
CREATE INDEX idx_orders_employee_id ON orders(employee_id);
CREATE INDEX idx_orders_created_time ON orders(created_time);
CREATE INDEX idx_orders_modified_time ON orders(modified_time);
CREATE INDEX idx_orders_state ON orders(state);
CREATE INDEX idx_orders_employee_created ON orders(employee_id, created_time);
CREATE INDEX idx_orders_employee_categories ON orders(employee_id, unique_category_count, created_time);

-- order_items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_items_analytic_category ON order_items(analytic_category_id);
CREATE INDEX idx_order_items_product_created ON order_items(product_id, created_at);

-- order_category_summary
CREATE INDEX idx_order_category_summary_order_id ON order_category_summary(order_id);
CREATE INDEX idx_order_category_summary_category ON order_category_summary(analytic_category_id);
CREATE INDEX idx_order_category_summary_restaurant ON order_category_summary(restaurant_id, analytic_category_id);

-- payments
CREATE INDEX idx_payments_restaurant_id ON payments(restaurant_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_employee_id ON payments(employee_id);
CREATE INDEX idx_payments_created_time ON payments(created_time);
CREATE INDEX idx_payments_payment_type ON payments(payment_type);

-- sync_logs
CREATE INDEX idx_sync_logs_restaurant_id ON sync_logs(restaurant_id);
CREATE INDEX idx_sync_logs_entity ON sync_logs(entity);
CREATE INDEX idx_sync_logs_status ON sync_logs(status);
CREATE INDEX idx_sync_logs_started_at ON sync_logs(started_at);

-- sync_errors
CREATE INDEX idx_sync_errors_restaurant_id ON sync_errors(restaurant_id);
CREATE INDEX idx_sync_errors_is_resolved ON sync_errors(is_resolved) WHERE is_resolved = false;
CREATE INDEX idx_sync_errors_entity ON sync_errors(entity);
