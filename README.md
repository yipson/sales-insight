# Sales Insight

Sistema de análisis de ventas y desempeño para restaurantes con Clover POS.

## Estructura

- `backend/` — Aplicación Go (API + Sync Engine)
- `frontend/` — Aplicación React + Vite + Tailwind CSS
- `docs/context/` — Documentación de planificación (SPEC, ARCHITECTURE, DATABASE_SCHEMA)
- `docker-compose.yml` — Orquestación local

## Inicio rápido

```bash
# 1. Configurar variables de entorno
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
# Editar .env con credenciales de Clover

# 2. Levantar servicios
docker compose up -d db

# 3. Ejecutar migraciones
cd backend && migrate -path migrations -database "$DATABASE_URL" up

# 4. Levantar backend y frontend
docker compose up -d

# 5. Acceder
# Frontend: http://localhost:5173
# API:      http://localhost:8080
```

## Documentación

Ver `docs/context/` para especificación completa.
