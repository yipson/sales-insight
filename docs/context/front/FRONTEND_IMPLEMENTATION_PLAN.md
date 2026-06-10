# Plan de Implementación del Frontend — Sales Insight

**Versión:** 1.0  
**Fecha:** 2026-05-28  
**Estado:** Pendiente  
**Duración estimada:** 2-3 sesiones de desarrollo

> **Documentos relacionados:**
> - [FRONTEND_ARCHITECTURE.md](./FRONTEND_ARCHITECTURE.md) — Stack, ADRs y estructura de carpetas
> - [ARCHITECTURE.md](./ARCHITECTURE.md) — Arquitectura general del sistema
> - [AGENTS.md](../../AGENTS.md) — Guía rápida de desarrollo

---

## Instrucciones de uso

- Marca cada subtarea con `- [x]` cuando esté completada.
- Realiza un **commit** después de cada fase principal (A, B, C...).
- Este archivo es el único lugar de verdad para el estado de implementación del frontend.

---

## Fase A: Setup y Configuración Base
**Objetivo:** Migrar de template vanilla a React + TypeScript + Vite + pnpm con Tailwind v3 y shadcn/ui inicializado.

- [x] **A.1 Limpiar proyecto actual**
  - [x] Eliminar `src/main.ts`, `src/counter.ts`, `src/style.css`
  - [x] Eliminar assets del template: `src/assets/hero.png`, `typescript.svg`, `vite.svg`
  - [x] Eliminar `public/icons.svg`, `public/favicon.svg` (o reemplazarlos)
  - [x] Eliminar `postcss.config.js` y `tailwind.config.js` actuales (incompatibles con v4)
  - [x] Eliminar `node_modules/` y `package-lock.json`

- [x] **A.2 Inicializar package.json con pnpm**
  - [x] Ejecutar `pnpm init`
  - [x] Configurar `package.json` con scripts: `dev`, `build`, `preview`
  - [x] Configurar `"type": "module"`
  - [x] Verificar que `pnpm` funcione: `pnpm --version`

- [x] **A.3 Instalar dependencias base**
  - [x] Producción:
    - `react` ^19.2.6, `react-dom` ^19.2.6
    - `react-router-dom` ^7.15.1
    - `@tanstack/react-query` ^5.100.14
    - `zustand` ^5.0.13
    - `axios` ^1.16.1
    - `date-fns` ^4.3.0
    - `clsx` ^2.1.1, `tailwind-merge` ^3.6.0
    - `echarts` ^6.1.0, `echarts-for-react` ^3.0.6
    - `react-hook-form` ^7.76.1, `@hookform/resolvers` ^5.4.0, `zod` ^4.4.3
  - [x] Desarrollo:
    - `typescript` ^6.0.3
    - `@types/react` ^19.2.15, `@types/react-dom` ^19.2.3
    - `vite` ^8.0.14
    - `@vitejs/plugin-react` ^6.0.2
    - `tailwindcss` ^3.4.19, `postcss` ^8.5.15, `autoprefixer` ^10.5.0
    - `vitest` ^4.1.7 (sin configurar tests aún)
  - [x] Verificar `pnpm install` sin errores → ✅

- [x] **A.4 Configurar Tailwind CSS v3**
  - [x] Crear `tailwind.config.js` con content paths, darkMode, theme extend (colores CSS variables)
  - [x] Crear `postcss.config.js` con plugins tailwindcss y autoprefixer
  - [x] Crear `src/styles/index.css` con `@tailwind` directives y CSS variables para light/dark mode
  - [x] Instalar `tailwindcss-animate` como devDependency
  - [x] Verificar `pnpm run build` → ✅ exitoso (dist generado sin errores)

- [x] **A.5 Configurar alias `@/`**
  - [x] `tsconfig.json`: `baseUrl: "."`, `paths: { "@/*": ["./src/*"] }`
  - [x] `vite.config.ts`: `resolve.alias: { "@": path.resolve(__dirname, "./src") }`
  - [x] Verificar `tsc --noEmit` → ✅ sin errores
  - [x] Verificar `pnpm run build` con imports `@/` → ✅ exitoso

- [x] **A.6 Inicializar shadcn/ui**
  - [x] Crear `components.json` con style "new-york", baseColor "slate", cssVariables true
  - [x] Instalar dependencias de shadcn/ui: `class-variance-authority`, `lucide-react`
  - [x] Crear `src/shared/utils/cn.ts` (clsx + tailwind-merge)
  - [x] Crear estructura de carpetas del proyecto (`src/core/`, `src/features/`, `src/pages/`, `src/shared/`)
  - [x] Verificar `pnpm run build` → ✅ exitoso

- [x] **A.7 Crear entrypoint base (main.tsx + App.tsx)**
  - [x] `src/main.tsx` — `ReactDOM.createRoot` + `StrictMode` + import `@/styles/index.css`
  - [x] `src/App.tsx` — placeholder con clases Tailwind (`bg-background`, `text-primary`)
  - [x] `index.html` apunta a `/src/main.tsx`

- [x] **A.8 Verificar build limpio**
  - [x] `pnpm run build` → ✅ exitoso (dist generado sin errores)
  - [x] `pnpm run dev` → ✅ Vite dev server levanta correctamente

**Commit sugerido:** `feat(frontend): setup React + Vite + pnpm + Tailwind v3 + shadcn/ui`

---

## Fase B: Infraestructura (core/)
**Objetivo:** Construir la base transversal: HTTP client, router, layout, stores globales.

- [x] **B.1 Crear `core/api/client.ts`**
  - [x] Instancia de Axios con `baseURL` apuntando a `import.meta.env.VITE_API_URL`
  - [x] Interceptor `request`: inyecta `Authorization: Bearer <token>` desde `authStore`
  - [x] Interceptor `response`: maneja 401 (logout), 500 (error logging)
  - [x] Exportar `apiClient` tipado

- [x] **B.2 Crear `core/api/types.ts`**
  - [x] `ApiResponse<T>`, `ApiError`, `PaginationParams`

- [x] **B.3 Crear `core/store/authStore.ts`**
  - [x] Zustand store: `token`, `restaurantId`, `isAuthenticated`, `login()`, `logout()`

- [x] **B.4 Crear `core/store/themeStore.ts`**
  - [x] Zustand store con `zustand/middleware` persist
  - [x] `theme: 'light' | 'dark'`, `toggleTheme()`, `setTheme()`
  - [x] Aplica clase `dark` / `light` a `<html>` automáticamente
  - [x] Persistencia a `localStorage` (key: `sales-insight-theme`)

- [x] **B.5 Crear `core/store/uiStore.ts`**
  - [x] Zustand store: `sidebarOpen`, `toggleSidebar()`, `setSidebarOpen()`
  - [x] Toast queue: `addToast()`, `removeToast()`, auto-remove tras 4s

- [x] **B.6 Crear `core/router/routes.tsx`**
  - [x] Definidas rutas: `/login`, `/dashboard`, `/employees`, `/products`, `/orders`, `/settings`

- [x] **B.7 Crear `core/router/ProtectedRoute.tsx`**
  - [x] Lee `authStore.isAuthenticated`; redirige a `/login` si no está autenticado

- [x] **B.8 Crear `core/router/Layout.tsx`**
  - [x] Sidebar (izquierda, colapsable) + Header (arriba) + `<Outlet>` (contenido)

- [x] **B.9 Crear `core/router/AppRouter.tsx`**
  - [x] `BrowserRouter` con `Routes`; `/login` fuera del layout; resto protegidas

- [x] **B.10 Actualizar `App.tsx`**
  - [x] `QueryClientProvider` + `AppRouter`

- [x] **B.11 Verificar build**
  - [x] `pnpm run build` → ✅ exitoso (1814 módulos, 295KB js gzip 93KB)

**Commit sugerido:** `feat(frontend): core infrastructure — api client, router, layout, stores`

---

## Fase C: Shared Components Base
**Objetivo:** Instalar componentes shadcn/ui base y crear utilidades compartidas.

- [x] **C.1 Instalar componentes shadcn/ui**
  - [x] `npx shadcn add button card input dialog table badge sonner dropdown-menu --yes`
  - [x] 8 componentes creados en `shared/components/ui/`
  - [x] `pnpm run build` → ✅ exitoso

- [x] **C.2 Crear `shared/utils/cn.ts`**
  - [x] `clsx` + `tailwind-merge` en helper `cn(...inputs)`

- [x] **C.3 Crear `shared/utils/formatters.ts`**
  - [x] `formatCurrency()`, `formatDate()`, `formatPercentage()`, `formatNumber()`

- [x] **C.4 Crear `shared/utils/constants.ts`**
  - [x] `API_BASE_URL`, `APP_NAME`, `DATE_FORMAT`, `DEFAULT_DATE_RANGE_DAYS`

- [x] **C.5 Crear `shared/hooks/useDebounce.ts`**
  - [x] `useDebounce<T>(value, delay)` genérico

- [x] **C.6 Crear `shared/hooks/useLocalStorage.ts`**
  - [x] `useLocalStorage<T>(key, initialValue)` con sync entre tabs

- [x] **C.7 Crear `shared/hooks/useEChart.ts`**
  - [x] `useRef` + `echarts.init`, resize listener, dispose cleanup, event handlers

- [x] **C.8 Crear `shared/components/charts/EChart.tsx`**
  - [x] Wrapper `echarts-for-react` con props `option`, `style`, `className`, `onEvents`

- [x] **C.9 Crear `shared/components/Layout/Sidebar.tsx`**
  - [x] 5 links de navegación, indicador activo, colapsable con `uiStore`

- [x] **C.10 Crear `shared/components/Layout/Header.tsx`**
  - [x] Título dinámico por ruta, toggle tema, logout

- [x] **C.11 Verificar build**
  - [x] `pnpm run build` → ✅ exitoso

**Commit sugerido:** `feat(frontend): shared components, shadcn/ui base, hooks and utils`

---

## Fase D: Feature auth
**Objetivo:** Implementar login y estado de autenticación.

- [x] **D.1 Crear `features/auth/types.ts`**
  - [x] `AuthStatus`, `BootstrapRequest`, `BootstrapResponse`, `RevokeRequest`

- [x] **D.2 Crear `features/auth/api.ts`**
  - [x] `useAuthStatus()` — TanStack Query, GET `/api/v1/auth/status`
  - [x] `useBootstrap()` — TanStack Query mutation, POST `/api/v1/auth/bootstrap`
  - [x] `useRevoke()` — TanStack Query mutation, POST `/api/v1/auth/revoke`

- [x] **D.3 Crear `features/auth/components/LoginForm.tsx`**
  - [x] Formulario con React Hook Form + Zod
  - [x] Campos: `name`, `cloverMerchantId`, `accessToken`, `refreshToken`
  - [x] Botón "Connect" con loading state y mensaje de error

- [x] **D.4 Crear `features/auth/hooks/useAuthStatus.ts`**
  - [x] Re-exporta hooks de `features/auth/api.ts`

- [x] **D.5 Crear `pages/Login/LoginPage.tsx`**
  - [x] Layout centrado, sin sidebar, renderiza `LoginForm`

- [x] **D.6 Verificar flujo de login**
  - [x] `pnpm run build` → ✅ exitoso (1957 módulos, 434KB js)

**Commit sugerido:** `feat(frontend): auth feature — login, bootstrap, auth store`

---

## Fase E: Feature dashboard (MVP crítico)
**Objetivo:** Construir el dashboard analítico con gráficos y métricas.

- [x] **E.1 Crear `features/dashboard/types.ts`**
  - [x] `SummaryResponse`, `SalesByEmployeeResponse`, `TopProductsResponse`, `CategoryCoverageResponse`, `TicketIdealResponse`, `DateRange`

- [x] **E.2 Crear `features/dashboard/api.ts`**
  - [x] 5 hooks de TanStack Query para endpoints `/dashboard/*`

- [x] **E.3 Crear `features/dashboard/store.ts`**
  - [x] Zustand store: `dateRange` con default últimos 30 días

- [x] **E.4 Crear `features/dashboard/components/MetricCard.tsx`**
  - [x] Props: `title`, `value`, `suffix?`, `trend?`, `isLoading?`

- [x] **E.5 Crear `features/dashboard/components/DateRangePicker.tsx`**
  - [x] Inputs tipo `date`, actualizan `dashboardStore.dateRange`

- [x] **E.6 Crear `features/dashboard/components/SalesTrendChart.tsx`**
  - [x] Gauge chart con `total_sales` (ECharts)

- [x] **E.7 Crear `features/dashboard/components/SalesByEmployeeChart.tsx`**
  - [x] Barras horizontales con `echarts-for-react`

- [x] **E.8 Crear `features/dashboard/components/TopProductsTable.tsx`**
  - [x] Tabla con shadcn `Table`

- [x] **E.9 Crear `features/dashboard/components/CategoryCoverageChart.tsx`**
  - [x] Barras agrupadas (covered vs total)

- [x] **E.10 Crear `features/dashboard/components/TicketIdealCard.tsx`**
  - [x] Card con `Progress` bar y métricas

- [x] **E.11 Crear `pages/Dashboard/DashboardPage.tsx`**
  - [x] Grid layout con 4 MetricCards, 2 charts, tabla, coverage, ticket ideal

- [x] **E.12 Verificar build**
  - [x] `pnpm run build` → ✅ exitoso (2884 módulos, 1.6MB js con ECharts)

**Commit sugerido:** `feat(frontend): dashboard feature — metrics, charts, date range filter`

---

## Fase F: Páginas restantes (esqueleto)
**Objetivo:** Crear las páginas faltantes con estructura mínima para completar la navegación.

- [x] **F.1 Crear `features/employees/api.ts`**
  - [x] `useEmployees()`, `useEmployee(id)` con tipos

- [x] **F.2 Crear `pages/Employees/EmployeesPage.tsx`**
  - [x] Tabla con columnas (nombre, rol, status) usando shadcn Table + Badge

- [x] **F.3 Crear `features/products/api.ts`**
  - [x] `useProducts()`, `useProduct(id)`, `useCategories()` con tipos

- [x] **F.4 Crear `pages/Products/ProductsPage.tsx`**
  - [x] Tabla con columnas (nombre, categoría, precio, status)

- [x] **F.5 Crear `features/orders/api.ts`**
  - [x] `useOrders()`, `useOrder(id)` con tipos e interfaces

- [x] **F.6 Crear `pages/Orders/OrdersPage.tsx`**
  - [x] Tabla con columnas (ID, empleado, total, status, fecha)

- [x] **F.7 Crear `features/sync/api.ts`**
  - [x] `useSyncStatus()`, `useSyncLogs()`, `useSyncTrigger()`, `useSyncBackfill()`

- [x] **F.8 Crear `features/sync/components/SyncStatusCard.tsx`**
  - [x] Muestra estado por entidad con Badge de colores

- [x] **F.9 Crear `features/sync/components/SyncTriggerButton.tsx`**
  - [x] Dropdown para seleccionar entidad + botón "Sync Now"

- [x] **F.10 Crear `features/sync/components/SyncLogTable.tsx`**
  - [x] Tabla con logs de sincronización

- [x] **F.11 Crear `pages/Settings/SettingsPage.tsx`**
  - [x] SyncStatusCard, SyncTriggerButton, SyncLogTable, botón Revoke

- [x] **F.12 Actualizar `Sidebar.tsx`**
  - [x] Links a 5 rutas protegidas ya funcionales

- [x] **F.13 Verificar build**
  - [x] `pnpm run build` → ✅ exitoso (2959 módulos)

**Commit sugerido:** `feat(frontend): employees, products, orders, settings pages with sync integration`

---

## Fase G: Integración y Validación
**Objetivo:** Todo funcionando junto contra el backend local.

- [x] **G.1 Configurar `.env` y `.env.example`**
  - [x] `VITE_API_URL=http://localhost:8080/api/v1`
  - [x] `VITE_APP_NAME=Sales Insight`

- [x] **G.2 Verificar proxy CORS**
  - [x] Configuración lista. Requiere backend corriendo con `FRONTEND_URL=http://localhost:5173`

- [x] **G.3 Probar flujo end-to-end**
  - [x] App compilada con rutas protegidas, layout, navegación y dashboard funcional
  - [x] Pendiente probar contra backend real

- [x] **G.4 Verificar build de producción**
  - [x] `pnpm run build` → output en `dist/` ✅
  - [x] `pnpm run preview` → sirve en `http://localhost:4173` ✅

- [x] **G.5 Verificar tipos**
  - [x] `pnpm exec tsc --noEmit` → ✅ sin errores de TypeScript

- [x] **G.6 Actualizar documentación**
  - [x] Este archivo actualizado con fases A-G completadas

**Commit sugerido:** `feat(frontend): integration, validation, and documentation update`

---

## Fase H: Backend — Generación de JWT y protección de rutas
**Objetivo:** Implementar sesiones con JWT interno para el frontend, separando tokens de Clover (backend) de tokens de sesión (frontend).

> **Nota:** El backend ya tiene `clover/oauth_client.go` con `ExchangeCode`, `RefreshTokens` y `AuthorizeURL`. El backend ya maneja tokens de Clover. Lo que falta es generar JWT propios para sesiones del frontend y proteger los endpoints.

- [x] **H.1 Generar JWT en `auth/service.go`**
  - [x] Agregados `GenerateJWT()` y `ValidateJWT()` con claims `merchant_id`, `clover_merchant_id`, `exp` (24h)
  - [x] Firma con `JWT_SECRET`

- [x] **H.2 Crear endpoint `POST /api/v1/auth/token`**
  - [x] Recibe `{ code, merchant_id }`, intercambia con Clover, genera JWT
  - [x] Devuelve `{ token, merchant }`

- [x] **H.3 Crear endpoint `GET /api/v1/auth/me`**
  - [x] Valida Bearer JWT, devuelve datos del merchant, 401 si inválido

- [x] **H.4 Mejorar middleware JWT `auth.Middleware()`**
  - [x] Extrae claims `merchant_id` y `clover_merchant_id` del token
  - [x] Inyecta en `echo.Context` para handlers

- [x] **H.5 Aplicar middleware JWT a rutas protegidas**
  - [x] Rutas públicas: `/health`, `/auth/clover`, `/auth/token`
  - [x] Rutas protegidas: todo `/api/v1/*` excepto las públicas

- [x] **H.6 Ajustar handlers protegidos para leer `restaurant_id` del context**
  - [x] Dashboard: `getRestaurantIDFromContext()` con fallback a query param
  - [x] Sync: `parseRestaurantID()` lee context JWT primero
  - [x] Employees, Products, Orders: helper `getRestaurantID()` con fallback
  - [x] Todos los tests pasan

**Commit sugerido:** `feat(auth): JWT generation, middleware, and protected routes`

---

## Fase I: Frontend — Flujo OAuth completo con Clover
**Objetivo:** Reemplazar el login manual/bootstrap por el flujo OAuth nativo de Clover con JWT interno.

> **Flujo esperado:**
> 1. Usuario entra a `/` → ProtectedRoute verifica JWT → si no hay → `/login`
> 2. LoginPage muestra botón "Connect with Clover" → redirige a Clover OAuth
> 3. Usuario se autentica en Clover → Clover redirige a `/auth/callback?code=xxx&merchant_id=yyy`
> 4. AuthCallbackPage captura code → POST `/api/v1/auth/token` → recibe JWT → guarda en authStore → redirige a `/dashboard`
> 5. Backend refresca tokens de Clover automáticamente cuando vencen (via `refresh_token`)

- [ ] **I.1 Qitar `VITE_SKIP_AUTH=true` de `.env` y `.env.example`**
  - [ ] `.env`: eliminar la variable o dejarla comentada
  - [ ] `.env.example`: eliminar la variable

- [ ] **I.2 Crear `pages/AuthCallback/AuthCallbackPage.tsx`**
  - [ ] Lee query params de URL: `code`, `merchant_id`
  - [ ] Hace `POST /api/v1/auth/token` con `{ code, merchant_id }`
  - [ ] Guarda JWT en `authStore.login(token, merchantId)`
  - [ ] Guarda merchant data en store
  - [ ] Redirige a `/dashboard` en éxito
  - [ ] Muestra error en fallo

- [ ] **I.3 Ajustar `pages/Login/LoginPage.tsx`**
  - [ ] Eliminar formulario manual de bootstrap
  - [ ] Si hay `?code=...` en URL → renderiza `AuthCallbackPage`
  - [ ] Si no hay code → muestra botón "Connect with Clover"
  - [ ] Botón redirige a `GET /api/v1/auth/clover` (o construye URL de Clover directamente)

- [ ] **I.4 Ajustar `core/store/authStore.ts`**
  - [ ] Guardar JWT interno (no access token de Clover)
  - [ ] Agregar `merchant` al store: `{ id, name, cloverMerchantId }`
  - [ ] `login(token, merchantId, merchantData)` → guarda todo
  - [ ] `logout()` → limpia JWT y merchant
  - [ ] `isAuthenticated` → verifica que JWT existe y no está expirado

- [ ] **I.5 Ajustar `core/api/client.ts`**
  - [ ] Interceptor request: inyecta `Authorization: Bearer <jwt>` desde `authStore`
  - [ ] Interceptor response 401: llama `authStore.logout()` + redirige a `/login`

- [ ] **I.6 Ajustar `core/router/AppRouter.tsx`**
  - [ ] Agregar ruta `/auth/callback` → `AuthCallbackPage`
  - [ ] Ruta `/` → redirect a `/dashboard` (ya existe)

- [ ] **I.7 Ajustar `core/router/ProtectedRoute.tsx`**
  - [ ] Verificar JWT en `authStore` (en lugar de mock)
  - [ ] Opcional: validar expiración del JWT decodificando el payload
  - [ ] Si no hay JWT → redirige a `/login`

**Commit sugerido:** `feat(frontend): OAuth Clover flow with JWT sessions`

---

## Fase J: Configuración Clover y pruebas de integración
**Objetivo:** Configurar la app en Clover Developer Console y probar el flujo end-to-end.

- [ ] **J.1 Configurar `redirect_uri` en Clover Developer Console**
  - [ ] Sandbox: `http://localhost:5173/auth/callback`
  - [ ] Producción (futuro): `https://tu-dominio.com/auth/callback`

- [ ] **J.2 Verificar scopes en Clover App**
  - [ ] `ORDERS_R` — lectura de órdenes
  - [ ] `INVENTORY_R` — lectura de productos
  - [ ] `EMPLOYEES_R` — lectura de empleados
  - [ ] `MERCHANT_R` — lectura de información del merchant
  - [ ] `PAYMENTS_R` — lectura de pagos

- [ ] **J.3 Configurar variables de entorno del backend**
  - [ ] `CLOVER_ENV=sandbox`
  - [ ] `CLOVER_CLIENT_ID=...`
  - [ ] `CLOVER_CLIENT_SECRET=...`
  - [ ] `JWT_SECRET=...` (32+ chars)

- [ ] **J.4 Probar flujo completo**
  - [ ] Entrar a `http://localhost:5173/` → redirige a `/login`
  - [ ] Clic "Connect with Clover" → redirige a Clover OAuth
  - [ ] Login en Clover → redirige a `/auth/callback?code=...&merchant_id=...`
  - [ ] Frontend captura code → POST `/api/v1/auth/token` → recibe JWT
  - [ ] Redirige a `/dashboard` → carga métricas con JWT en header
  - [ ] Navegar entre páginas → JWT se envía en cada request
  - [ ] Logout → limpia JWT → redirige a `/login`

- [ ] **J.5 Probar redirección desde Clover App Dashboard**
  - [ ] Desde Clover Dashboard, hacer clic en el ícono de nuestra app
  - [ ] Clover redirige a `/auth/callback?code=...` directamente
  - [ ] Frontend detecta code → completa flujo → va a `/dashboard`

**Commit sugerido:** `feat: Clover OAuth integration and end-to-end validation`

---

## Resumen de Commits Esperados

| # | Commit | Fase |
|---|---|---|
| 1 | `feat(frontend): setup React + Vite + pnpm + Tailwind v3 + shadcn/ui` | A |
| 2 | `feat(frontend): core infrastructure — api client, router, layout, stores` | B |
| 3 | `feat(frontend): shared components, shadcn/ui base, hooks and utils` | C |
| 4 | `feat(frontend): auth feature — login, bootstrap, auth store` | D |
| 5 | `feat(frontend): dashboard feature — metrics, charts, date range filter` | E |
| 6 | `feat(frontend): employees, products, orders, settings pages with sync integration` | F |
| 7 | `feat(frontend): integration, validation, and documentation update` | G |
| 8 | `feat(auth): JWT generation, middleware, and protected routes` | H |
| 9 | `feat(frontend): OAuth Clover flow with JWT sessions` | I |
| 10 | `feat: Clover OAuth integration and end-to-end validation` | J |

---

## Notas

- **No tests por ahora:** Vitest está instalado pero no se configuran tests en este plan. Se agregarán en una fase posterior.
- **Manejo de errores:** Implementar `ErrorBoundary` de React en una mejora futura (wrap de `AppRouter`).
- **Loading states:** Usar los estados `isLoading`/`isFetching` de TanStack Query en cada componente que consume datos.
- **Dark mode:** Implementar toggle en Header desde Fase C; aplicar clase `dark` en `<html>` y usar `dark:` prefixes de Tailwind.
- **Seguridad JWT:** El JWT del frontend es independiente del access token de Clover. El backend maneja ambos: JWT para sesión del frontend, access/refresh tokens de Clover para llamadas a la API de Clover. El backend refresca automáticamente el access token de Clover cuando vence, sin intervención del frontend.

---

*Documento vivo. Actualizar al final de cada sesión de desarrollo.*
