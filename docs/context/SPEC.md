# SPEC.md

## Sales Insight — Especificación del Producto

**Versión:** 1.0-MVP  
**Fecha:** 2026-05-06  
**Stack propuesto:** Go + PostgreSQL + React (Vite) + Apache ECharts  
**Ambiente inicial:** PC Local (Docker Compose)  
**Infraestructura destino (post-MVP):** DigitalOcean Droplet  
**Locales objetivo (MVP):** 1  
**Locales objetivo (futuro):** 15+

---

## 1. Visión

Construir un sistema automatizado —denominado **Sales Insight**— que extraiga datos operativos de la plataforma Clover (órdenes, productos, pagos, empleados), los estructure en un modelo analítico propio, y los presente en un dashboard web de escritorio que permita evaluar el desempeño de los meseros, identificar productos más vendidos, y detectar oportunidades de upselling.

El sistema debe operar sin intervención manual constante, iniciar su ejecución de forma local mediante Docker, y estar preparado para migrar a un VPS en la nube en el futuro.

---

## 2. Objetivos del MVP

1. Conectar de forma segura con la API de Clover de **un (1) restaurante** mediante OAuth 2.0.
2. Extraer automáticamente órdenes, productos, pagos y empleados mediante **polling periódico**.
3. Procesar y estructurar los datos para soportar análisis por: empleado, producto, fecha/hora.
4. Proveer un **dashboard web de escritorio** con métricas clave de ventas y desempeño.
5. Permitir que cada restaurante defina sus **propias categorías analíticas** (ej. `Tacos`, `Bebidas`, `Postres`) y mapee las categorías nativas de Clover hacia ellas.
6. Operar en un entorno local mediante **Docker Compose**, con el mínimo consumo de recursos posible.

---

## 3. Requerimientos Funcionales

### 3.1 Autenticación y Conexión
- **RF-1.1:** El sistema debe implementar el flujo OAuth 2.0 Authorization Code de Clover para obtener tokens de acceso.
- **RF-1.2:** Los tokens (`access_token`, `refresh_token`, `merchant_id`) deben almacenarse de forma **encriptada en PostgreSQL** como fuente de verdad.
- **RF-1.2b:** El backend debe mantener una **caché en memoria** de los tokens activos para minimizar lecturas a base de datos en cada request a Clover. La caché se reconstruye desde PostgreSQL al iniciar el proceso.
- **RF-1.3:** El sistema debe refrescar automáticamente el `access_token` cuando expire, sin intervención manual.
- **RF-1.4:** Debe existir una pantalla de administración donde se pueda iniciar o revocar la conexión con Clover.
- **RF-1.5:** Para el entorno local sin IP pública, se debe documentar y soportar una **estrategia de bootstrap OAuth** (ej. ngrok temporal o migración de tokens desde una instancia con endpoint público) para obtener los tokens iniciales. El sistema productivo final requerirá un dominio fijo (DO/VPS) para el callback OAuth.

### 3.2 Extracción de Datos
- **RF-2.1:** Extraer datos de los endpoints relevantes de Clover:
  - Órdenes (`/v3/merchants/{id}/orders`)
  - Productos/ítems (`/v3/merchants/{id}/items`)
  - Pagos (`/v3/merchants/{id}/payments`)
  - Empleados (`/v3/merchants/{id}/employees`)
- **RF-2.2:** Un **scheduler interno (polling)** debe ejecutarse periódicamente para la sincronización continua.
- **RF-2.3:** El sistema debe detectar y evitar duplicados utilizando identificadores únicos de Clover (`clover_order_id`, `clover_item_id`, etc.).
- **RF-2.4:** Debe existir un mecanismo de recuperación ante fallos (reintentos con **backoff exponencial**).
- **RF-2.5:** Los datos extraídos deben registrarse con timestamp de sincronización para trazabilidad.
- **RF-2.6:** El cliente HTTP de Clover debe implementar **rate limiting client-side** respetando los límites de Clover: máximo 16 req/s por token y 50 req/s por app, con máximo 5 requests concurrentes por token. Ante un error `429 Too Many Requests`, debe pausar y reintentar respetando el header `retry-after`.
- **RF-2.7:** Para cargas iniciales masivas (backfill) mayores a 2 meses, el sistema debe soportar la **Clover Export API** como alternativa al polling REST, minimizando el riesgo de rate limits. *(Nota: requiere solicitar permisos a developer-relations@clover.com).*
- **RF-2.8:** El polling debe usar `modifiedTime` como cursor para evitar re-descargar datos no modificados.

### 3.3 Procesamiento y Modelado de Datos
- **RF-3.1:** El sistema debe importar y respetar la **categoría nativa** que cada producto tiene en Clover.
- **RF-3.2:** El administrador podrá crear **categorías analíticas propias** del restaurante (nombre, slug, color) y configurar un **mapeo** desde las categorías nativas de Clover hacia ellas.
- **RF-3.3:** Las órdenes deben desagregarse en líneas de ítem (`order_items`) para permitir análisis a nivel de producto vendido.
- **RF-3.4:** Se debe calcular y almacenar el desglose de cada orden por categoría analítica (ej. cuántas entradas, platos principales, bebidas y postres contiene).

### 3.4 Dashboard y Visualización
- **RF-4.1:** El dashboard debe ser accesible vía navegador web en **escritorio** (desktop-first). El diseño responsive se pospone para una fase posterior.
- **RF-4.2:** Vista de **Ventas por Empleado**: total vendido, ticket promedio, ranking.
- **RF-4.3:** Vista de **Productos Más Vendidos**: ranking por cantidad y por revenue, filtrable por categoría.
- **RF-4.4:** Vista de **Cobertura de Categorías**: porcentaje de órdenes que incluyen cada categoría analítica definida por el restaurante. Métricas de upselling: órdenes con 2+ categorías distintas, promedio de categorías por orden.
- **RF-4.5:** Vista de **Ventas por Período**: filtrado por fecha, con agrupación por día, semana o mes.
- **RF-4.6:** Todas las vistas deben permitir filtrar por rango de fechas.
- **RF-4.7:** Las visualizaciones deben implementarse con **Apache ECharts** para soportar datasets grandes, zoom temporal y tooltips ricos.

### 3.5 Administración
- **RF-5.1:** Pantalla de administración de categorías analíticas (crear, editar, eliminar) y mapeo de categorías nativas de Clover hacia ellas.
- **RF-5.2:** Pantalla de estado de sincronización (última extracción exitosa, errores recientes, rate limits encontrados).

---

## 4. Requerimientos No Funcionales

- **RNF-1 (Performance):** El dashboard debe cargar las vistas principales en menos de 2 segundos. Las consultas deben usar índices de base de datos.
- **RNF-2 (Disponibilidad):** El sistema debe operar continuamente mientras el contenedor Docker esté activo. El scheduler debe ser resiliente a reinicios del proceso.
- **RNF-3 (Seguridad):**
  - Tokens OAuth encriptados en reposo (AES-256-GCM).
  - HTTPS obligatorio para producción; para desarrollo local se permite HTTP.
  - Sin exposición de credenciales en logs (redacción de campos sensibles).
  - Validación y saneamiento de todos los parámetros de query (OWASP A03).
  - API keys de Clover en variables de entorno, nunca commiteadas al repositorio (OWASP A05).
  - Endpoints de debug deshabilitados en build de producción.
  - Validación del header `X-Clover-Auth` cuando se habiliten webhooks post-MVP.
- **RNF-4 (Escalabilidad horizontal - preparación):** El esquema de base de datos debe incluir `restaurant_id` en todas las tablas relevantes, aunque el MVP opere con un solo valor fijo.
- **RNF-5 (Mantenibilidad):** Código modular. Separación clara entre capa de extracción, procesamiento, API y frontend.
- **RNF-6 (Observabilidad básica):** Logs estructurados de las extracciones. Registro de errores en base de datos o archivo para diagnóstico. Alertas ante llamadas fallidas consecutivas a la API de Clover.
- **RNF-7 (Recursos):** Debe operar establemente en un entorno local con Docker Compose, sin componentes adicionales de cache o colas (Redis omitido en MVP). PostgreSQL y la memoria del proceso Go son suficientes.

---

## 5. Flujos de Usuario Principales

### 5.1 Primer Uso: Conexión con Clover
1. Administrador accede al panel y hace clic en "Conectar con Clover".
2. Es redirigido a la página de autorización de Clover.
3. Autoriza el acceso a los datos del merchant.
4. Clover redirige al callback del sistema con un `code`.
5. El backend intercambia el `code` por tokens, encripta y almacena en PostgreSQL, y carga en caché de memoria.
6. Se ejecuta la primera sincronización completa de datos (usando Export API para histórico si está disponible, o polling).

> **Nota:** El callback OAuth requiere un endpoint público. Para el MVP local, se puede usar ngrok temporalmente o importar tokens manualmente desde una instancia con acceso público.

### 5.2 Uso Diario: Revisión de Dashboard
1. Administrador/gerente accede al dashboard web desde su navegador (desktop).
2. Visualiza automáticamente los datos del día actual.
3. Aplica filtros de fecha para comparar períodos.
4. Revisa el ranking de meseros y los productos más vendidos.
5. Identifica meseros con bajo índice de tickets completos.

### 5.3 Mantenimiento: Categorías Analíticas y Mapeo
1. Administrador accede a la sección "Categorías".
2. Crea sus categorías analíticas propias: "Tacos", "Bebidas", "Postres" (nombre, color, orden).
3. Visualiza la lista de categorías nativas traídas de Clover.
4. Asigna o corrige el mapeo hacia la categoría analítica correspondiente (ej. categoria Clover "Bebestibles" → `Bebidas`).
5. El sistema recalcula las métricas de ticket afectadas por este cambio.

---

## 6. Lógica de Negocio Clave

### 6.1 Cobertura de Categorías (Análisis de Ticket)
En lugar de un concepto rígido de "ticket completo" que asume una estructura de menú universal, el sistema calcula métricas relativas al catálogo real de cada restaurante:

- **Cobertura por categoría:** % de órdenes que incluyen al menos un ítem de cada categoría analítica configurada por el restaurante (ej. `Tacos`, `Bebidas`, `Postres`).
- **Órdenes multi-categoría:** % de órdenes que incluyen ítems de **2 o más categorías analíticas distintas**. Mide cross-selling sin depender de una fórmula fija.
- **Promedio de categorías por orden:** Indicador general de diversificación del ticket.
- **Cobertura por empleado:** Comparación de % de cobertura por categoría entre meseros.

> **Ejemplo contextual:** Un café tendrá alta cobertura de `Bebidas` y baja de `Alimentos`. Una taquería tendrá alta cobertura de `Tacos`. El dashboard muestra solo las categorías definidas por ese restaurante.

### 6.2 Ticket Ideal (Opcional y Configurable)
Además del análisis flexible, cada restaurante puede optar por configurar su propia definición de "ticket ideal" (o "ticket completo"):

- El administrador define un JSON de reglas usando los slugs de sus categorías analíticas: `{"tacos": 2, "bebidas": 1}`.
- El sistema calcula dinámicamente, por cada orden, si cumple con las reglas configuradas.
- El dashboard muestra esta métrica **solo si el restaurante tiene reglas configuradas**. Si no, se omite sin afectar el resto del análisis.
- Esto permite que una taquería mida "2 tacos + 1 bebida", mientras un café mide "1 bebida + 1 alimento".

> **Nota:** Esta funcionalidad es **opcional**. El MVP funciona completamente sin configurar reglas, mostrando únicamente cobertura de categorías.

### 6.2 Producto Más Vendido
Se mide por:
- Cantidad total de unidades vendidas.
- Revenue total generado (`quantity * unit_price` acumulado).

### 6.3 Ventas por Empleado
Se mide por:
- Suma del `total_amount` de todas las órdenes asignadas al empleado.
- Ticket promedio = `total ventas / cantidad de órdenes`.
- Horas trabajadas (aproximación): diferencia en horas entre la primera y la última orden del empleado en un día dado. Si no registra órdenes, se considera 0 horas.
- Ventas por hora trabajada = `total ventas del día / horas trabajadas aproximadas`.

---

## 7. Preguntas Pendientes / Decisiones por Confirmar

| ID | Tema | Impacto | Estado |
|---|---|---|---|
| **P-01** | Fuente de horas trabajadas por empleado. | Resuelto: se usará aproximación por rango de órdenes (primera vs. última orden del día). | Resuelto |
| **P-02** | Uso de Webhooks de Clover vs. solo polling. | Resuelto: solo polling para MVP. Webhooks se habilitarán como mejora cuando se migre a VPS/DO con dominio fijo. | Resuelto |
| **P-03** | Permisos para Clover Export API. | Pendiente: solicitar a developer-relations@clover.com. Sin esto, el backfill inicial se hará por polling con rate limiting cuidadoso. | Pendiente |

---

## 8. Criterios de Aceptación del MVP

- [ ] CA-1: Un administrador puede conectar exitosamente el restaurante a Clover mediante OAuth (o importar tokens manualmente en entorno local).
- [ ] CA-2: Los datos de órdenes, productos, pagos y empleados se sincronizan automáticamente mediante polling sin errores.
- [ ] CA-3: No existen registros duplicados en la base de datos después de múltiples sincronizaciones.
- [ ] CA-4: El administrador puede crear categorías analíticas propias y mapear las categorías nativas de Clover hacia ellas.
- [ ] CA-5: El dashboard muestra correctamente las métricas de ventas por empleado y productos más vendidos.
- [ ] CA-6: El dashboard calcula y muestra la cobertura de categorías por orden y el porcentaje de órdenes multi-categoría (2+ categorías distintas).
- [ ] CA-7: El dashboard es funcional y legible en un navegador de escritorio (Chrome/Firefox/Edge).
- [ ] CA-8: El sistema ejecuta exitosamente en un entorno local con Docker Compose.
- [ ] CA-9: El sistema opera sin intervención manual durante al menos 7 días consecutivos.
- [ ] CA-10: El cliente HTTP respeta los rate limits de Clover y maneja correctamente errores 429 con backoff exponencial.

---

## 9. Futuras Mejoras Post-MVP (No incluidas)

1. Soporte multi-local (15+ restaurantes).
2. Webhooks de Clover para sincronización en tiempo real (requiere VPS con dominio fijo).
3. Integración con fuente real de horas trabajadas (reloj checador) para reemplazar la aproximación por órdenes.
4. UI de configuración de "ticket ideal" por restaurante (el campo ya existe en el esquema, falta la interfaz de administración).
5. Alertas automáticas (ej. mesero sin ventas de postres en X días).
6. Integración con microservicios Python para análisis estadísticos avanzados.
7. Caché distribuida (Redis) para alta concurrencia.
8. Exportación de reportes a PDF/Excel.
9. Autenticación multi-usuario con roles (admin, gerente, solo lectura).
10. Dashboard responsive para navegadores móviles.
11. Despliegue en DigitalOcean Droplet con dominio fijo y SSL.

---

## 10. Glosario

| Término | Definición |
|---|---|
| **Sales Insight** | Nombre del sistema de análisis de ventas y desempeño de restaurantes. |
| **Clover** | Plataforma de punto de venta (POS) utilizada por el restaurante. |
| **Merchant** | Identificador del negocio/restaurante dentro de Clover. |
| **Cobertura de categorías** | Análisis del porcentaje de órdenes que incluyen cada categoría analítica, adaptado al menú real de cada restaurante. |
| **Categoría analítica** | Clasificación personalizable definida por cada restaurante (ej. `Tacos`, `Bebidas`), obtenida mediante mapeo desde las categorías nativas de Clover. |
| **Upselling** | Estrategia de venta que busca que el cliente agregue productos adicionales a su orden. |
| **Bootstrap OAuth** | Estrategia para obtener tokens iniciales de Clover en un entorno sin IP pública (ej. ngrok temporal o importación manual). |

---

*Documento vivo. Cualquier cambio de alcance debe ser acordado y versionado.*
