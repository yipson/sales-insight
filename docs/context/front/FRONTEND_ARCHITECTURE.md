# Arquitectura del Frontend — Sales Insight

**Versión:** 1.0  
**Fecha:** 2026-05-28  
**Stack:** React 18 + TypeScript + Vite + pnpm  
**Autor:** Agent (sesión de arquitectura)

> **Documentos relacionados:**
> - [FRONTEND_IMPLEMENTATION_PLAN.md](./FRONTEND_IMPLEMENTATION_PLAN.md) — Paso a paso de implementación
> - [ARCHITECTURE.md](./ARCHITECTURE.md) — Arquitectura general del sistema (backend + infra)
> - [AGENTS.md](../../AGENTS.md) — Guía rápida de desarrollo

---

## 1. Stack Tecnológico Completo

| Capa | Tecnología | Versión | Justificación |
|---|---|---|---|
| **Framework** | React | 18 | Estándar de la industria, ecosistema maduro. |
| **Lenguaje** | TypeScript | 5.x | Type safety en compile time, mejor DX, refactorización segura. |
| **Build Tool** | Vite | ^6 | Builds rápidos, HMR instantáneo, output optimizado para producción. |
| **Package Manager** | pnpm | — | Rápido, eficiente en disco, compatible con monorepo workspaces. |
| **Routing** | React Router DOM | 6 | SPA routing declarativo, nested routes, loaders. |
| **Estado Global** | Zustand | 4 | Ligero, sin boilerplate, type-safe, slices independientes. |
| **Server State / Cache** | TanStack Query (React Query) | 5 | Cacheo, revalidación, deduping, manejo de errores y loading states automáticos. |
| **Charts** | Apache ECharts + echarts-for-react | 5 | Superior para datasets grandes, zoom temporal, tooltips ricos. |
| **Estilos** | Tailwind CSS | 3 | Utility-first, desktop-first (sin breakpoints móviles en MVP). |
| **UI Base** | shadcn/ui | latest | Componentes accesibles, headless, estilizados con Tailwind, copiados al codebase (no dependencia opaca). |
| **Formularios** | React Hook Form + Zod | latest | Performance (no re-renders innecesarios) + validación type-safe con schemas. |
| **HTTP Client** | Axios | 1.x | Interceptores globales para JWT injection, error handling, cancelación de requests. |
| **Fechas** | date-fns | 3.x | Funcional, tree-shakeable, inmutable. |
| **Iconos** | lucide-react | latest | Iconos limpios, consistentes, tree-shakeable. |
| **Testing** | Vitest | 2.x | Testing framework rápido integrado con Vite. *(Por ahora no se implementan tests)* |

---

## 2. Estructura de Carpetas

```
frontend/
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   │
│   ├── core/                    ← Infraestructura transversal
│   │   ├── api/
│   │   │   ├── client.ts        # Axios instance: baseURL, interceptores (JWT), error handling
│   │   │   └── types.ts         # Tipos genéricos de API (ApiResponse, ApiError)
│   │   ├── router/
│   │   │   ├── routes.tsx         # Definición de todas las rutas
│   │   │   ├── AppRouter.tsx      # BrowserRouter + rutas + layout
│   │   │   ├── ProtectedRoute.tsx # Guard: redirige a /login si no hay auth
│   │   │   └── Layout.tsx         # Layout principal (Sidebar + Header + <Outlet>)
│   │   └── store/
│   │       ├── authStore.ts       # Zustand: token, user, isAuthenticated, login, logout
│   │       ├── themeStore.ts      # Zustand: dark/light mode
│   │       └── uiStore.ts         # Zustand: sidebarOpen, toastQueue, modalStack
│   │
│   ├── pages/                   ← Thin wrappers: 1 carpeta = 1 ruta
│   │   ├── Login/
│   │   │   └── LoginPage.tsx
│   │   ├── Dashboard/
│   │   │   └── DashboardPage.tsx
│   │   ├── Employees/
│   │   │   └── EmployeesPage.tsx
│   │   ├── Products/
│   │   │   └── ProductsPage.tsx
│   │   ├── Orders/
│   │   │   └── OrdersPage.tsx
│   │   └── Settings/
│   │       └── SettingsPage.tsx
│   │
│   ├── features/                ← Lógica encapsulada por dominio
│   │   ├── auth/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── api.ts           # TanStack Query: login, logout, status, bootstrap
│   │   │   └── types.ts
│   │   ├── dashboard/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── api.ts           # TanStack Query: 5 endpoints GET /dashboard/*
│   │   │   ├── store.ts         # Zustand: dateRange compartido entre widgets
│   │   │   └── types.ts
│   │   ├── employees/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── api.ts
│   │   │   └── types.ts
│   │   ├── products/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── api.ts
│   │   │   └── types.ts
│   │   ├── orders/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── api.ts
│   │   │   └── types.ts
│   │   └── sync/
│   │       ├── components/
│   │       ├── hooks/
│   │       ├── api.ts           # TanStack Query: status, logs, trigger
│   │       └── store.ts         # Zustand: último estado de sync
│   │
│   ├── shared/                  ← Reutilizable entre cualquier feature
│   │   ├── components/
│   │   │   ├── ui/              # ← shadcn/ui components (Button, Input, Card, Dialog, Table...)
│   │   │   ├── Layout/          # Sidebar.tsx, Header.tsx, MainContent.tsx
│   │   │   └── charts/          # EChart.tsx (wrapper), ChartContainer
│   │   ├── hooks/
│   │   │   ├── useDebounce.ts
│   │   │   ├── useLocalStorage.ts
│   │   │   └── useEChart.ts     # useRef + echarts.init para datasets complejos
│   │   ├── utils/
│   │   │   ├── cn.ts            # clsx + tailwind-merge
│   │   │   ├── formatters.ts    # formatCurrency, formatDate, formatPercentage
│   │   │   └── constants.ts     # API_BASE_URL, APP_NAME, DATE_FORMAT
│   │   └── types/
│   │       └── index.ts         # DTOs globales: Restaurant, Pagination, DateRange
│   │
│   └── styles/
│       └── index.css            # Tailwind directives (@tailwind base/components/utilities)
│
├── components.json              # shadcn/ui config
├── package.json
├── pnpm-lock.yaml
├── vite.config.ts
├── tsconfig.json
├── tsconfig.app.json
├── tailwind.config.js
├── postcss.config.js
└── index.html
```

---

## 3. Decisiones Arquitectónicas (ADRs Frontend)

### ADR-F01: Feature-Based Layout
**Decisión:** Organizar el código por dominio (`features/`) en lugar de por capa técnica (`components/`, `hooks/`, `services/`).  
**Justificación:**
- Escalabilidad cognitiva: un desarrollador trabaja en un dominio modificando archivos contiguos.
- Cada feature encapsula su propia API (TanStack Query), estado local (Zustand), componentes y types.
- Reduce merge conflicts cuando múltiples desarrolladores trabajan en features distintas.
- `pages/` son thin wrappers que orquestan features; no contienen lógica de negocio.

### ADR-F02: Zustand en Dos Niveles
**Decisión:** Dividir el estado global en `core/store/` (transversal) y `features/{x}/store.ts` (local al dominio).  
**Justificación:**
- `core/store/`: Solo lo que múltiples features necesitan leer/escribir simultáneamente (auth, theme, UI shell).
- `features/{x}/store.ts`: Estado que solo componentes de esa feature comparten (ej. `dateRange` del dashboard).
- Evita que el store global se convierta en un "god object" difícil de mantener.

### ADR-F03: TanStack Query por Feature
**Decisión:** Cada feature expone sus hooks de query/mutation en `features/{x}/api.ts`. No existe una "API layer" global con funciones sueltas.  
**Justificación:**
- Encapsulamiento: los hooks de `dashboard` no se filtran a `products`.
- Type safety: cada `api.ts` importa los DTOs de su propio `types.ts`.
- Facilita mocking en tests futuros: cada feature es auto-contenida.

### ADR-F04: shadcn/ui en shared/components/ui/
**Decisión:** Los componentes de shadcn/ui viven en `shared/components/ui/` (el path por defecto de la CLI).  
**Justificación:**
- shadcn/ui no es una librería de npm opaca; son archivos copiados al codebase.
- El CLI de shadcn espera encontrar `components/ui/` para poder hacer `npx shadcn add button` sin reconfigurar paths.
- Los componentes de shadcn son primitives reutilizables; encajan perfectamente en `shared/`.

### ADR-F05: ECharts Dual Mode
**Decisión:** Usar `echarts-for-react` como wrapper por defecto, y un hook `useEChart.ts` (useRef + echarts.init) para casos complejos.  
**Justificación:**
- El 90% de los gráficos del dashboard son configuraciones declarativas simples que funcionan perfectamente con el wrapper.
- El 10% restante (datasets masivos, eventos personalizados de zoom/brush) requieren control imperativo que el wrapper no expone bien.
- Centralizar el control imperativo en un hook reutilizable evita duplicar lógica de `resize`, `dispose` y `init`.

### ADR-F06: Axios sobre Fetch
**Decisión:** Usar Axios en lugar de `fetch` nativo.  
**Justificación:**
- Interceptores globales: inyección automática de header `Authorization`, logging de requests, manejo centralizado de 401.
- Cancelación de requests: útil cuando un componente se desmonta antes de que termine una query.
- Timeouts configurables y transformación automática de JSON.

### ADR-F07: Tailwind CSS v3 (no v4)
**Decisión:** Mantener Tailwind CSS v3 en lugar de migrar a v4.  
**Justificación:**
- Tailwind v4 cambia el modelo de configuración (CSS-first en lugar de `tailwind.config.js`).
- shadcn/ui está optimizado para Tailwind v3 en la mayoría de sus templates.
- Migrar a v4 requiere reescribir la configuración y posiblemente ajustar clases; el beneficio no justifica el esfuerzo en el MVP.

### ADR-F08: pnpm como Package Manager
**Decisión:** Usar pnpm en lugar de npm.  
**Justificación:**
- Instalaciones más rápidas y uso eficiente de disco (hard links).
- `pnpm-lock.yaml` más determinista y legible que `package-lock.json`.
- Compatible con monorepo workspaces si el proyecto escala a un workspace con shared packages.

---

## 4. Reglas de Dependencias entre Carpetas

```
pages/         → puede importar de features/, shared/, core/
features/      → puede importar de shared/, core/
                → NUNCA importa de otra feature/
shared/        → puede importar de core/
                → NUNCA importa de features/ ni pages/
core/          → NUNCA importa de pages/, features/, shared/
```

**Ejemplos válidos:**
- `pages/Dashboard/DashboardPage.tsx` importa `features/dashboard/components/MetricCard` ✅
- `features/dashboard/api.ts` importa `core/api/client` ✅
- `features/dashboard/store.ts` importa `shared/utils/constants` ✅

**Ejemplos inválidos:**
- `features/products/api.ts` importa `features/orders/types.ts` ❌ (circular coupling)
- `shared/components/ui/Button.tsx` importa `features/auth/types.ts` ❌ (shared no depende de features)
- `core/store/authStore.ts` importa `features/dashboard/store.ts` ❌ (core es la capa más interna)

---

## 5. Alias `@/`

Configurado en `tsconfig.json` y `vite.config.ts`:

```json
// tsconfig.json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

Uso:
- `@/core/store/authStore`
- `@/features/dashboard/api`
- `@/shared/components/ui/Button`
- `@/shared/utils/formatters`

---

## 6. Convenciones de Nomenclatura

| Tipo | Convención | Ejemplo |
|---|---|---|
| Componentes React | PascalCase | `MetricCard.tsx`, `SalesChart.tsx` |
| Hooks custom | camelCase con prefijo `use` | `useDashboardSummary.ts`, `useDebounce.ts` |
| Stores Zustand | camelCase con sufijo `Store` | `authStore.ts`, `themeStore.ts` |
| Archivos de API (TanStack Query) | `api.ts` | `features/dashboard/api.ts` |
| Archivos de tipos | `types.ts` | `features/dashboard/types.ts` |
| Utilidades | camelCase | `formatCurrency.ts`, `cn.ts` |
| Constantes | UPPER_SNAKE_CASE dentro de archivos | `API_BASE_URL`, `DATE_FORMAT` |

---

## 7. Manejo de Estado — Matriz de Decisión

| Dato | Tecnología | Ubicación | Ejemplo |
|---|---|---|---|
| Sesión de usuario (token) | Zustand | `core/store/authStore.ts` | Token JWT, merchantID |
| Preferencia de tema | Zustand + persist | `core/store/themeStore.ts` | Dark/light mode |
| Estado del UI shell | Zustand | `core/store/uiStore.ts` | Sidebar abierto, toast activo |
| Datos del servidor | TanStack Query | `features/{x}/api.ts` | Métricas del dashboard, lista de órdenes |
| Estado local de feature | Zustand | `features/{x}/store.ts` | DateRange del dashboard |
| Estado de formulario | React Hook Form | Componente local | Inputs de login, filtros de tabla |
| Estado efímero de UI | useState | Componente local | Hover, modal open, loading local |

---

## 8. Notas sobre Autenticación

El backend actual expone endpoints de auth pero no enforce JWT universalmente (algunos endpoints como `GET /auth/status` listan todos los merchants sin verificar token). A pesar de esto, el frontend implementará un flujo de auth completo:

1. **Login / Bootstrap:** El admin puede iniciar sesión vía OAuth Clover o insertar tokens manualmente (`POST /auth/bootstrap`).
2. **JWT Storage:** El token JWT de sesión se almacena en `authStore.ts` (Zustand). En un futuro se puede agregar persistencia a `localStorage` si el backend emite tokens con expiración.
3. **Interceptor Axios:** Cada request a `/api/v1/*` inyecta automáticamente el header `Authorization: Bearer <token>` si existe.
4. **ProtectedRoute:** `core/router/ProtectedRoute.tsx` verifica `authStore.isAuthenticated`. Si no hay sesión, redirige a `/login`.
5. **Logout:** Limpia el store y redirige a `/login`.

---

*Documento vivo. Actualizar cuando se agreguen nuevas decisiones arquitectónicas.*
