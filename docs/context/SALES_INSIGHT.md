# SALES INSIGHT
PROPUESTA DE PRODUCTO MÍNIMO VIABLE (MVP)

Fecha: Mayo 2026


---
## 1. ¿QUÉ ES SALES INSIGHT?

Sales Insight es una herramienta de inteligencia de negocios diseñada
específicamente para restaurantes que operan con Clover POS.

Su propósito es simple pero poderoso: convertir los datos que ya estás
generando todos los días (ventas, órdenes, productos, empleados) en
información clara y accionable que te permita tomar mejores decisiones
estratégicas para tu negocio.

En lugar de depender de la intuición o de revisar reportes manuales de
Clover, Sales Insight te presenta un dashboard visual donde podrás:

  • Identificar quiénes son tus meseros más efectivos y por qué.
  • Saber exactamente qué productos venden más y cuáles están
    subutilizados.
  • Detectar oportunidades de venta cruzada (upselling) en cada orden.
  • Medir el desempeño de tu equipo con base en datos reales, no en
    percepciones.
  • Comparar períodos de tiempo para entender tendencias de ventas.

Sales Insight no reemplaza a Clover. Lo complementa: Clover registra las
ventas, Sales Insight las analiza y te dice qué decisiones tomar para
aumentar tu rentabilidad.


---
## 2. ¿QUÉ ES UN MVP?


MVP significa Producto Mínimo Viable. Es la primera versión funcional del
software que incluye las funcionalidades esenciales para empezar a generar
valor inmediato en tu restaurante.

Pensalo de esta manera: no construimos el producto completo de una vez,
sino que empezamos con lo más importante, lo ponemos a funcionar en tu
negocio, y luego vamos agregando mejoras basándonos en lo que realmente
te sirve.

El MVP de Sales Insight incluye todo lo necesario para:
  • Conectar con tu Clover.
  • Extraer automáticamente tus datos.
  • Mostrarlos en un dashboard claro y útil.
  • Operar sin que tengas que hacer nada manual después de la
    configuración inicial.


---
## 3. FUNCIONALIDADES INCLUIDAS EN ESTA VERSIÓN


### 3.1 CONEXIÓN AUTOMÁTICA CON CLOVER

Sales Insight se conecta directamente a tu cuenta de Clover mediante
autenticación segura (OAuth 2.0). Una vez conectado, no necesitás hacer
nada manual: el sistema extrae los datos por su cuenta.

### 3.2 EXTRACCIÓN PERIÓDICA DE DATOS

El sistema consulta tu Clover cada 5 minutos para mantener la información
actualizada. No hay necesidad de exportar archivos ni de hacer copias
manuales.

Datos que se extraen:
  • Órdenes de venta
  • Productos e ítems del menú
  • Pagos
  • Empleados

### 3.3 DASHBOARD DE ANÁLISIS

Un panel visual accesible desde tu navegador web (en computadora) donde
podrás ver:

- A) VENTAS POR EMPLEADO
     - Total vendido por cada mesero.
     - Ticket promedio por empleado.
     - Ranking de desempeño.
     - Aproximación de ventas por hora trabajada.

- B) PRODUCTOS MÁS VENDIDOS
     - Ranking de platos e ítems por cantidad vendida.
     - Ranking por ingresos generados.
     - Filtrado por categoría.

- C) COBERTURA DE CATEGORÍAS
     - Qué porcentaje de órdenes incluyen bebidas, alimentos, postres
       u otras categorías que vos definas.
     - Indicador de órdenes multi-categoría (cuántas órdenes incluyen
       productos de 2 o más categorías distintas).
     - Promedio de diversificación por orden.

-  D) ANÁLISIS POR PERÍODO
     - Comparación de ventas por día, semana o mes.
     - Filtros por rango de fechas.

### 3.4 ADMINISTRACIÓN DE CATEGORÍAS PERSONALIZADAS

Cada restaurante tiene su propio menú y su propia forma de organizar
productos. Sales Insight te permite:

  - Crear tus propias categorías analíticas (ej. "Tacos", "Bebidas",
    "Entradas", "Postres").
  - Mapear las categorías que ya tenés en Clover hacia estas categorías
    propias.
  - Cambiar el mapeo cuando quieras sin perder el historial de ventas.

### 3.5 TICKET IDEAL CONFIGURABLE (OPCIONAL)

Si tu restaurante tiene una estructura de venta ideal (por ejemplo: 2
tacos + 1 bebida), podés configurar una regla personalizada y el sistema
te dirá qué porcentaje de órdenes la cumplen. Esto es completamente
opcional y adaptable a tu menú.

### 3.6 ESTADO DE SINCRONIZACIÓN

Una pantalla donde podés ver:
- Cuándo fue la última sincronización exitosa.
- Si hubo errores recientes.
- Opción de disparar una sincronización manual cuando lo necesites.


--- 
## 4. QUÉ NO INCLUYE ESTA VERSIÓN (PERO ESTÁ PLANEADO)


Para ser transparente, estas funcionalidades NO están en el MVP pero se
agregarán en futuras versiones:

- App móvil (el dashboard funciona solo en computadora por ahora).
- Alertas automáticas (ej. "el mesero X no vendió postres esta semana").
- Soporte para múltiples locales (el MVP es para un solo restaurante).
- Actualización en tiempo real milisegundo a milisegundo (el polling
    cada 5 minutos es suficiente para análisis).
- Exportación de reportes a PDF o Excel.


---
### 5. CRITERIOS DE ÉXITO DEL MVP


Estos son los puntos clave que determinan que el MVP está funcionando
correctamente:

- 1. Conexión exitosa con Clover mediante OAuth.
- 2. Sincronización automática cada 5 minutos sin errores durante 7 días
     consecutivos.
- 3. Sin registros duplicados en la base de datos.
- 4. El administrador puede crear categorías analíticas propias y mapear
     las categorías de Clover hacia ellas.
- 5. El dashboard muestra correctamente:
    - a) Ventas por empleado.
    - b) Productos más vendidos.
    - c) Cobertura de categorías por orden.
- 6. El sistema opera sin intervención manual después de la configuración
     inicial.
- 7. El sistema respeta los límites de la API de Clover y maneja
     correctamente cualquier interrupción temporal.


---
## 6. TECNOLOGÍAS UTILIZADAS

Sales Insight está construido con tecnologías modernas, probadas y de
alto rendimiento:

- Go: lenguaje de programación rápido y eficiente en recursos.
- PostgreSQL: base de datos robusta para almacenar y analizar los
    datos.
- React: interfaz de usuario interactiva y fluida.
- Docker: encapsulamiento del sistema para fácil instalación.

Estas tecnologías garantizan que el sistema sea estable, seguro y
escalable.


---
## 7. PRÓXIMOS PASOS

- 1. Configuración de la app en Clover Developer Dashboard (1-2 días).
- 2. Instalación local de Sales Insight y primeras pruebas (1-2 semanas).
- 3. Demo funcional con datos reales del restaurante.
- 4. Ajustes finales y puesta en marcha.


---
Documento preparado por el equipo de desarrollo de Sales Insight.
Mayo 2026.
