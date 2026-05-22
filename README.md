# Sales Insight

Sistema de análisis de ventas y desempeño para restaurantes con Clover POS.

---

## 📚 Documentación del Proyecto

### Para continuar el desarrollo (empezar aquí)

1. **[ESTADO_ACTUAL.md](./docs/context/ESTADO_ACTUAL.md)** — Estado actual de implementación y deuda técnica
2. **[AGENTS.md](./AGENTS.md)** — Guía rápida para desarrolladores: comandos, setup local, reglas de arquitectura
3. **[implementacion-mvp.md](./docs/context/implementacion-mvp.md)** — Plan de implementación detallado con checklists por fase
4. **[ARCHITECTURE.md](./docs/context/ARCHITECTURE.md)** — Decisiones arquitectónicas, stack tecnológico y ADRs
5. **[DATABASE_SCHEMA.md](./docs/context/DATABASE_SCHEMA.md)** — Esquema completo de PostgreSQL, índices y decisiones de diseño

### Documentación Clover (referencia API)

- [OAuth Authorization Flow](./docs/context/clover/oauth/00-diagram.md)
- [Generate OAuth Tokens](./docs/context/clover/oauth/01-clover-generate-oauth-expiring-tokens.md)
- [Refresh Access Tokens](./docs/context/clover/oauth/02-clover-refresh-access-tokens.md)
- [API Rate Limits](./docs/context/clover/clover-api-rate-limits.md)
- [Webhooks](./docs/context/clover/clover-webhooks.md)
- [Export Merchant Data](./docs/context/clover/clover-export-merchant-data.md)

### Documento para cliente

- [SALES_INSIGHT.md](./docs/context/SALES_INSIGHT.md) — Propuesta de MVP presentada a Ricardo Caicedo

---

## 🚀 Inicio Rápido

```bash
# 1. Base de datos (PostgreSQL en Docker)
make db-up

# 2. Backend (requiere backend/.env configurado)
cd backend && go run ./cmd/api

# 3. Frontend
cd frontend && npm install && npm run dev
```

Para instrucciones detalladas de setup, ver [AGENTS.md](./AGENTS.md).

---

## 🗺️ Mapa de Lectura por Rol

| Si sos... | Empezá por... |
|---|---|
| **Desarrollador continuando el proyecto** | [ESTADO_ACTUAL.md](./docs/context/ESTADO_ACTUAL.md) → [AGENTS.md](./AGENTS.md) |
| **Desarrollador nuevo en el equipo** | [AGENTS.md](./AGENTS.md) → [ARCHITECTURE.md](./docs/context/ARCHITECTURE.md) |
| **Revisor de código / Tech lead** | [ESTADO_ACTUAL.md](./docs/context/ESTADO_ACTUAL.md) → [ARCHITECTURE.md](./docs/context/ARCHITECTURE.md) |
| **DevOps / Infraestructura** | [ARCHITECTURE.md](./docs/context/ARCHITECTURE.md) → [AGENTS.md#development-setup-local-non-docker](./AGENTS.md) |
| **Cliente / Stakeholder** | [SALES_INSIGHT.md](./docs/context/SALES_INSIGHT.md) |

---

## 📁 Estructura del Repositorio

```
sales-insight/
├── backend/                 # Aplicación Go (API + Sync Engine)
│   ├── cmd/api/             # Entrypoint único
│   ├── internal/            # Código de dominio (feature-based)
│   ├── migrations/          # SQL migrations (golang-migrate)
│   └── .env.example         # Variables de entorno
├── frontend/                # React + Vite + Tailwind
├── docs/
│   └── context/             # Documentación del proyecto
│       ├── ESTADO_ACTUAL.md
│       ├── AGENTS.md
│       ├── ARCHITECTURE.md
│       ├── DATABASE_SCHEMA.md
│       ├── implementacion-mvp.md
│       ├── SALES_INSIGHT.md
│       └── clover/          # Documentación de referencia Clover
└── docker-compose.yml       # PostgreSQL + servicios futuros
```

---

*Última actualización: 2026-05-22*

---

## 🔍 Análisis de Documentación

Para un análisis detallado de la precisión y calidad de esta documentación, ver:
[Análisis de Precisión de Documentación](../documents/analisis-precision-documentacion.md)
