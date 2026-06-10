# Sales Insight — Frontend

Aplicación web del dashboard de Sales Insight. Consumida por el backend Go y diseñada para desktop-first.

## Stack

- **React 18** + **TypeScript**
- **Vite** (build tool)
- **pnpm** (package manager)
- **Tailwind CSS v3** (estilos)
- **shadcn/ui** (componentes base)
- **TanStack Query v5** (server state / cache)
- **Zustand** (estado global)
- **React Router v6** (routing)
- **Apache ECharts** (gráficos)
- **React Hook Form + Zod** (formularios)
- **Axios** (HTTP client)

## Estructura de carpetas

```
src/
├── core/           # Infraestructura transversal (api, router, stores globales)
├── pages/          # Thin wrappers de rutas (1 carpeta = 1 ruta)
├── features/       # Lógica encapsulada por dominio (dashboard, auth, sync...)
├── shared/         # Reutilizable entre cualquier feature (ui, hooks, utils, types)
└── styles/         # Tailwind directives
```

Ver [FRONTEND_ARCHITECTURE.md](../docs/context/FRONTEND_ARCHITECTURE.md) para el detalle completo de decisiones arquitectónicas.

## Comandos

```bash
# Instalar dependencias
pnpm install

# Servidor de desarrollo
pnpm run dev

# Build de producción
pnpm run build

# Previsualizar build
pnpm run preview
```

## Variables de entorno

Copiar `.env.example` a `.env` y ajustar:

```bash
VITE_API_URL=http://localhost:8080/api/v1
```

## Requisitos previos

- Node.js 18+
- pnpm instalado globalmente: `npm install -g pnpm`
- Backend corriendo localmente (ver `../backend/`)

## Conexión con backend

El frontend espera el backend en `http://localhost:8080` por defecto. Asegurarse de que el backend tenga configurado `FRONTEND_URL=http://localhost:5173` para que CORS funcione correctamente.

## Convenciones

- **Alias `@/`**: apunta a `src/`. Ejemplo: `import { Button } from "@/shared/components/ui/button"`
- **Componentes**: PascalCase (`MetricCard.tsx`)
- **Hooks**: camelCase con prefijo `use` (`useDashboardSummary.ts`)
- **Stores**: camelCase con sufijo `Store` (`authStore.ts`)

## Páginas

| Ruta | Descripción |
|---|---|
| `/login` | Login / Bootstrap manual de tokens Clover |
| `/dashboard` | Métricas principales, gráficos de ventas, top productos |
| `/employees` | Lista de empleados y métricas |
| `/products` | Productos, categorías y mapeo analítico |
| `/orders` | Órdenes con detalle de líneas |
| `/settings` | Estado de sincronización, logs, trigger manual |

---

*Para el plan de implementación detallado, ver [FRONTEND_IMPLEMENTATION_PLAN.md](../docs/context/FRONTEND_IMPLEMENTATION_PLAN.md)*
