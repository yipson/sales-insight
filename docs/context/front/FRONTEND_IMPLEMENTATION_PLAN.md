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

- [ ] **B.1 Crear `core/api/client.ts`**
  - Instancia de Axios con `baseURL` apuntando a `import.meta.env.VITE_API_URL`
  - Interceptor `request`: inyecta `Authorization: Bearer <token>` desde `authStore`
  - Interceptor `response`: maneja 401 (logout), 500 (toast de error genérico)
  - Exportar `apiClient` tipado

- [ ] **B.2 Crear `core/api/types.ts`**
  - Definir `ApiResponse<T>`, `ApiError`, `PaginationParams`

- [ ] **B.3 Crear `core/store/authStore.ts`**
  - Zustand store con: `token`, `restaurantId`, `isAuthenticated`, `login(token)`, `logout()`
  - Agregar persistencia opcional a `localStorage` (evaluar en sesión)

- [ ] **B.4 Crear `core/store/themeStore.ts`**
  - Zustand store con: `theme: 'light' | 'dark'`, `toggleTheme()`
  - Persistencia a `localStorage`
  - Aplicar clase `dark` a `<html>` cuando theme === 'dark'

- [ ] **B.5 Crear `core/store/uiStore.ts`**
  - Zustand store con: `sidebarOpen: boolean`, `toggleSidebar()`
  - Opcional: `toastQueue: Toast[]`, `addToast()`, `removeToast(id)`

- [ ] **B.6 Crear `core/router/routes.tsx`**
  - Definir array de rutas:
    - `/login` → `LoginPage` (sin layout)
    - `/dashboard` → `DashboardPage`
    - `/employees` → `EmployeesPage`
    - `/products` → `ProductsPage`
    - `/orders` → `OrdersPage`
    - `/settings` → `SettingsPage`
  - Todas las rutas excepto `/login` son protegidas

- [ ] **B.7 Crear `core/router/ProtectedRoute.tsx`**
  - Lee `authStore.isAuthenticated`
  - Si no está autenticado: redirige a `/login`
  - Si está autenticado: renderiza `<Outlet />`

- [ ] **B.8 Crear `core/router/Layout.tsx`**
  - Estructura: `Sidebar` (izquierda, fija) + `Header` (arriba) + `<Outlet>` (contenido)
  - Responsive: en MVP solo desktop, sin breakpoints móviles

- [ ] **B.9 Crear `core/router/AppRouter.tsx`**
  - `BrowserRouter` con `Routes`
  - Ruta `/login` fuera del `Layout`
  - Resto de rutas envueltas en `Layout` + `ProtectedRoute`

- [ ] **B.10 Actualizar `App.tsx`**
  - Importar `AppRouter`
  - Envolver con `QueryClientProvider` (TanStack Query)

- [ ] **B.11 Verificar build**
  - `pnpm run build` sin errores

**Commit sugerido:** `feat(frontend): core infrastructure — api client, router, layout, stores`

---

## Fase C: Shared Components Base
**Objetivo:** Instalar componentes shadcn/ui base y crear utilidades compartidas.

- [ ] **C.1 Instalar componentes shadcn/ui**
  - `npx shadcn add button card input dialog table badge sonner dropdown-menu`
  - Verificar que se crearon en `shared/components/ui/`

- [ ] **C.2 Crear `shared/utils/cn.ts`**
  - Combinar `clsx` + `tailwind-merge` en helper `cn(...inputs)`

- [ ] **C.3 Crear `shared/utils/formatters.ts`**
  - `formatCurrency(amount: number): string`
  - `formatDate(date: string | Date, format?: string): string`
  - `formatPercentage(value: number): string`

- [ ] **C.4 Crear `shared/utils/constants.ts`**
  - `API_BASE_URL` (fallback a `http://localhost:8080`)
  - `APP_NAME`
  - `DEFAULT_DATE_RANGE` (últimos 30 días)

- [ ] **C.5 Crear `shared/hooks/useDebounce.ts`**
  - Hook genérico `useDebounce<T>(value: T, delay: number): T`

- [ ] **C.6 Crear `shared/hooks/useLocalStorage.ts`**
  - Hook genérico `useLocalStorage<T>(key: string, initialValue: T): [T, (v: T) => void]`

- [ ] **C.7 Crear `shared/hooks/useEChart.ts`**
  - Hook con `useRef<HTMLDivElement>` + `useEffect`
  - `echarts.init()` en mount, `chart.setOption()` cuando cambian props, `chart.dispose()` en unmount
  - Listener de `resize`

- [ ] **C.8 Crear `shared/components/charts/EChart.tsx`**
  - Wrapper simple sobre `echarts-for-react` para el 90% de los casos
  - Props: `option`, `style`, `onEvents?`

- [ ] **C.9 Crear `shared/components/Layout/Sidebar.tsx`**
  - Links de navegación: Dashboard, Employees, Products, Orders, Settings
  - Indicador de ruta activa (usa `useLocation` de React Router)
  - Botón de colapsar (usa `uiStore.toggleSidebar`)

- [ ] **C.10 Crear `shared/components/Layout/Header.tsx`**
  - Título de la página actual
  - Botón de tema (light/dark)
  - Botón de logout

- [ ] **C.11 Verificar build**
  - `pnpm run build` sin errores

**Commit sugerido:** `feat(frontend): shared components, shadcn/ui base, hooks and utils`

---

## Fase D: Feature auth
**Objetivo:** Implementar login y estado de autenticación.

- [ ] **D.1 Crear `features/auth/types.ts`**
  - `LoginDTO`, `AuthStatus`, `BootstrapRequest`

- [ ] **D.2 Crear `features/auth/api.ts`**
  - `useAuthStatus()` — TanStack Query, GET `/api/v1/auth/status`
  - `useBootstrap()` — TanStack Query mutation, POST `/api/v1/auth/bootstrap`
  - `useRevoke()` — TanStack Query mutation, POST `/api/v1/auth/revoke`

- [ ] **D.3 Crear `features/auth/components/LoginForm.tsx`**
  - Formulario con React Hook Form + Zod
  - Campos: `name`, `cloverMerchantId`, `accessToken`, `refreshToken` (opcional)
  - Botón "Connect" que llama a `useBootstrap().mutate()`
  - Mensaje de éxito/error

- [ ] **D.4 Crear `features/auth/hooks/useAuthStatus.ts`**
  - Wrapper sobre `features/auth/api.ts` para consumo en componentes

- [ ] **D.5 Crear `pages/Login/LoginPage.tsx`**
  - Layout centrado, sin sidebar
  - Renderiza `LoginForm`

- [ ] **D.6 Verificar flujo de login**
  - Probar que `pnpm run dev` carga `/login`
  - Probar submit del formulario (aunque backend no esté corriendo, verificar request en Network tab)

**Commit sugerido:** `feat(frontend): auth feature — login, bootstrap, auth store`

---

## Fase E: Feature dashboard (MVP crítico)
**Objetivo:** Construir el dashboard analítico con gráficos y métricas.

- [ ] **E.1 Crear `features/dashboard/types.ts`**
  - `SummaryResponse`, `SalesByEmployeeResponse`, `TopProductsResponse`, `CategoryCoverageResponse`, `TicketIdealResponse`, `DateRange`

- [ ] **E.2 Crear `features/dashboard/api.ts`**
  - `useDashboardSummary(restaurantId, from, to)` — GET `/api/v1/dashboard/summary`
  - `useSalesByEmployee(restaurantId, from, to)` — GET `/api/v1/dashboard/sales-by-employee`
  - `useTopProducts(restaurantId, from, to, limit?)` — GET `/api/v1/dashboard/top-products`
  - `useCategoryCoverage(restaurantId, from, to)` — GET `/api/v1/dashboard/category-coverage`
  - `useTicketIdeal(restaurantId, from, to)` — GET `/api/v1/dashboard/ticket-ideal`

- [ ] **E.3 Crear `features/dashboard/store.ts`**
  - Zustand store: `dateRange: { from: Date, to: Date }`
  - Default: últimos 30 días

- [ ] **E.4 Crear `features/dashboard/components/MetricCard.tsx`**
  - Props: `title`, `value`, `suffix?`, `trend?` (porcentaje de cambio)
  - Usa shadcn `Card`

- [ ] **E.5 Crear `features/dashboard/components/DateRangePicker.tsx`**
  - Inputs tipo `date` o librería ligera
  - Actualiza `dashboardStore.dateRange`
  - Botón "Apply"

- [ ] **E.6 Crear `features/dashboard/components/SalesTrendChart.tsx`**
  - Gráfico de línea temporal (ECharts)
  - Consume `useDashboardSummary` (o endpoint de ventas por día si existe)

- [ ] **E.7 Crear `features/dashboard/components/SalesByEmployeeChart.tsx`**
  - Gráfico de barras horizontales
  - Consume `useSalesByEmployee`

- [ ] **E.8 Crear `features/dashboard/components/TopProductsTable.tsx`**
  - Tabla con shadcn `Table`
  - Consume `useTopProducts`

- [ ] **E.9 Crear `features/dashboard/components/CategoryCoverageChart.tsx`**
  - Gráfico de barras o radar
  - Consume `useCategoryCoverage`
  - Nota: puede estar vacío si backend aún devuelve placeholder

- [ ] **E.10 Crear `features/dashboard/components/TicketIdealCard.tsx`**
  - Card con métricas del ticket ideal
  - Consume `useTicketIdeal`

- [ ] **E.11 Crear `pages/Dashboard/DashboardPage.tsx`**
  - Layout de grid (Tailwind `grid`)
  - Fila 1: `DateRangePicker` + 4 `MetricCard`s
  - Fila 2: `SalesTrendChart` (ancho completo)
  - Fila 3: `SalesByEmployeeChart` + `TopProductsTable`
  - Fila 4: `CategoryCoverageChart` + `TicketIdealCard`

- [ ] **E.12 Verificar visualización**
  - `pnpm run dev`, navegar a `/dashboard`
  - Verificar que los charts se renderizan (aunque los datos sean mocks o vacíos)

**Commit sugerido:** `feat(frontend): dashboard feature — metrics, charts, date range filter`

---

## Fase F: Páginas restantes (esqueleto)
**Objetivo:** Crear las páginas faltantes con estructura mínima para completar la navegación.

- [ ] **F.1 Crear `features/employees/api.ts`**
  - `useEmployees()`, `useEmployee(id)`, `useEmployeeOrders(id)`

- [ ] **F.2 Crear `pages/Employees/EmployeesPage.tsx`**
  - Tabla vacía con columnas definidas (nombre, rol, órdenes, ventas)
  - Integrar `features/employees/api.ts`

- [ ] **F.3 Crear `features/products/api.ts`**
  - `useProducts()`, `useProduct(id)`, `useCategories()`

- [ ] **F.4 Crear `pages/Products/ProductsPage.tsx`**
  - Tabla vacía con columnas (nombre, categoría, precio, cantidad vendida)
  - Sección de categorías analíticas (placeholder)

- [ ] **F.5 Crear `features/orders/api.ts`**
  - `useOrders()`, `useOrder(id)`

- [ ] **F.6 Crear `pages/Orders/OrdersPage.tsx`**
  - Tabla vacía con filtros por fecha (placeholder)
  - Drawer o modal para detalle de orden

- [ ] **F.7 Crear `features/sync/api.ts`**
  - `useSyncStatus()`, `useSyncLogs()`, `useSyncTrigger()`, `useSyncBackfill()`

- [ ] **F.8 Crear `features/sync/components/SyncStatusCard.tsx`**
  - Muestra última sync por entidad (órders, items, employees, payments)
  - Indicador visual: éxito (verde), fallo (rojo), en progreso (amarillo)

- [ ] **F.9 Crear `features/sync/components/SyncTriggerButton.tsx`**
  - Botón "Sync Now" que llama a `useSyncTrigger().mutate()`
  - Dropdown para seleccionar entidad

- [ ] **F.10 Crear `features/sync/components/SyncLogTable.tsx`**
  - Tabla con últimos logs de sincronización
  - Consume `useSyncLogs()`

- [ ] **F.11 Crear `pages/Settings/SettingsPage.tsx`**
  - Estado de conexión con Clover (`features/sync/components/SyncStatusCard`)
  - Botón "Sync manual" (`SyncTriggerButton`)
  - Logs de sincronización (`SyncLogTable`)
  - Botón "Revoke connection" (usa `features/auth/api.ts`)

- [ ] **F.12 Actualizar `Sidebar.tsx`**
  - Asegurar que todos los links naveguen correctamente a las 5 rutas protegidas

- [ ] **F.13 Verificar navegación completa**
  - Clic en cada item del sidebar → cambio de ruta correcto
  - Refresh en `/employees` → no 404 (configurar Vite fallback)

**Commit sugerido:** `feat(frontend): employees, products, orders, settings pages with sync integration`

---

## Fase G: Integración y Validación
**Objetivo:** Todo funcionando junto contra el backend local.

- [ ] **G.1 Configurar `.env` y `.env.example`**
  - `VITE_API_URL=http://localhost:8080/api/v1`
  - `VITE_APP_NAME=Sales Insight`

- [ ] **G.2 Verificar proxy CORS**
  - Levantar backend local: `cd backend && go run ./cmd/api`
  - Frontend en `pnpm run dev` (puerto 5173)
  - Verificar que CORS del backend permite `http://localhost:5173`

- [ ] **G.3 Probar flujo end-to-end**
  - Login → Dashboard → cambio de fecha → charts se actualizan
  - Navegación a Settings → ver estado de sync
  - Logout → redirige a /login

- [ ] **G.4 Verificar build de producción**
  - `pnpm run build` → output en `dist/`
  - Servir `dist/` con `vite preview` o backend Go
  - Verificar que rutas SPA funcionan (fallback a `index.html`)

- [ ] **G.5 Verificar tipos**
  - `pnpm exec tsc --noEmit` → sin errores de TypeScript

- [ ] **G.6 Actualizar `AGENTS.md` y `ESTADO_ACTUAL.md`**
  - Reflejar que frontend ya no es "starter vacío"
  - Documentar comandos de pnpm y estructura de carpetas

**Commit sugerido:** `feat(frontend): integration, validation, and documentation update`

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

---

## Notas

- **No tests por ahora:** Vitest está instalado pero no se configuran tests en este plan. Se agregarán en una fase posterior.
- **Manejo de errores:** Implementar `ErrorBoundary` de React en una mejora futura (wrap de `AppRouter`).
- **Loading states:** Usar los estados `isLoading`/`isFetching` de TanStack Query en cada componente que consume datos.
- **Dark mode:** Implementar toggle en Header desde Fase C; aplicar clase `dark` en `<html>` y usar `dark:` prefixes de Tailwind.

---

*Documento vivo. Actualizar al final de cada sesión de desarrollo.*
