
```markdown
# Caso de Uso: API REST para Formulario de Contacto

**ID del Caso de Uso:** `UC-API-001` (Puedes usar un sistema de numeración)
**Versión:** `1.0`
**Fecha de Creación:** `YYYY-MM-DD`
**Autor:** `[Tu Nombre/Nombre del Curso]`

## 1. Descripción General

Este caso de uso describe la funcionalidad del backend para recibir, validar y almacenar los datos enviados desde un formulario de contacto del sitio web. El frontend (Vue.js) consumirá esta API para enviar la información.

## 2. Actores

*   **Usuario Anónimo/Visitante del Sitio (Primario):** La persona que llena y envía el formulario de contacto.
*   **Sistema (Backend Go):** El componente que procesa la solicitud.
*   **Base de Datos (MySQL):** Donde se almacenan los mensajes.

## 3. Precondiciones

*   El servicio backend (Go) está en ejecución y accesible.
*   La base de datos MySQL está configurada y accesible por el servicio backend.
*   El usuario ha navegado a la página que contiene el formulario de contacto en el frontend (Vue.js).

## 4. Flujo Principal (Happy Path)

1.  **Usuario:** Completa los campos del formulario de contacto en la interfaz Vue.js (ej: Nombre, Email, Mensaje).
2.  **Usuario:** Hace clic en el botón "Enviar Mensaje" (o similar).
3.  **Frontend (Vue.js):** Realiza una validación básica de los campos (ej: campos requeridos, formato de email).
4.  **Frontend (Vue.js):** Si la validación es exitosa, empaqueta los datos del formulario en formato JSON.
5.  **Frontend (Vue.js):** Envía una petición HTTP `POST` al endpoint `/api/contact` del backend Go, incluyendo el JSON en el cuerpo de la solicitud.
6.  **Backend (Go):** Recibe la petición HTTP `POST` en el endpoint `/api/contact`.
7.  **Backend (Go):** Deserializa (unmarshal) el cuerpo JSON de la petición a una estructura de datos de Go.
8.  **Backend (Go):** Realiza una validación exhaustiva de los datos recibidos:
    *   Campos obligatorios presentes (Nombre, Email, Mensaje).
    *   Formato de email válido.
    *   Longitud máxima/mínima de los campos.
    *   (Opcional) Sanitización de datos para prevenir XSS si se van a mostrar en algún panel de admin.
9.  **Backend (Go):** Si la validación es exitosa, prepara los datos para la inserción en la base de datos.
10. **Backend (Go):** Establece una conexión (si no existe ya en un pool) con la base de datos MySQL.
11. **Backend (Go):** Ejecuta una sentencia SQL `INSERT` para guardar el nuevo mensaje en la tabla `contacts` (o similar), incluyendo campos como `name`, `email`, `message`, `created_at`.
12. **Backend (Go):** Verifica que la inserción en la base de datos haya sido exitosa.
13. **Backend (Go):** Responde al frontend con un código de estado HTTP `201 Created` y, opcionalmente, un mensaje de éxito en formato JSON (ej: `{"message": "Mensaje enviado correctamente"}`).
14. **Frontend (Vue.js):** Recibe la respuesta `201 Created`.
15. **Frontend (Vue.js):** Muestra un mensaje de confirmación al usuario (ej: "¡Gracias! Tu mensaje ha sido enviado.").
16. **Frontend (Vue.js):** (Opcional) Limpia los campos del formulario.

## 5. Flujos Alternativos y Excepciones

*   **5.1. Validación Fallida en Frontend:**
    1.  **Frontend (Vue.js):** La validación básica falla (ej: email inválido).
    2.  **Frontend (Vue.js):** Muestra mensajes de error específicos junto a los campos correspondientes.
    3.  **Fin del caso de uso.**

*   **5.2. Validación Fallida en Backend:**
    1.  **Backend (Go):** La validación de datos (paso 8 del Flujo Principal) falla.
    2.  **Backend (Go):** Responde al frontend con un código de estado HTTP `400 Bad Request` y un cuerpo JSON detallando los errores de validación (ej: `{"errors": {"email": "Formato de email inválido", "message": "El mensaje es demasiado corto"}}`).
    3.  **Frontend (Vue.js):** Recibe la respuesta `400 Bad Request`.
    4.  **Frontend (Vue.js):** Muestra los mensajes de error al usuario.
    5.  **Fin del caso de uso.**

*   **5.3. Error al Guardar en Base de Datos:**
    1.  **Backend (Go):** Ocurre un error durante la inserción en MySQL (paso 11 del Flujo Principal) (ej: la base de datos no está disponible, error de constraint).
    2.  **Backend (Go):** Registra el error internamente (logging).
    3.  **Backend (Go):** Responde al frontend con un código de estado HTTP `500 Internal Server Error` y un mensaje genérico de error (ej: `{"error": "Ocurrió un error al procesar tu solicitud. Inténtalo más tarde."}`).
    4.  **Frontend (Vue.js):** Recibe la respuesta `500 Internal Server Error`.
    5.  **Frontend (Vue.js):** Muestra un mensaje de error genérico al usuario.
    6.  **Fin del caso de uso.**

*   **5.4. Endpoint No Encontrado:**
    1.  **Frontend (Vue.js):** Intenta enviar la petición a un endpoint incorrecto (ej: `/api/conctact`).
    2.  **Backend (Go Router):** No encuentra una ruta que coincida.
    3.  **Backend (Go Router):** Responde con un `404 Not Found`.
    4.  **Frontend (Vue.js):** Maneja el error 404 apropiadamente (puede ser un error genérico al usuario si no es esperado).
    5.  **Fin del caso de uso.**

## 6. Postcondiciones

*   **Éxito:**
    *   Un nuevo registro de contacto se ha almacenado en la tabla `contacts` de la base de datos MySQL.
    *   El usuario ha recibido una notificación de éxito en la interfaz.
*   **Fallo:**
    *   No se ha creado ningún nuevo registro en la base de datos.
    *   El usuario ha recibido una notificación de error en la interfaz.

## 7. Requisitos Especiales / No Funcionales

*   **Seguridad:**
    *   Validación de entrada tanto en frontend como en backend.
    *   Sanitización de datos para prevenir ataques XSS si los mensajes se muestran en otro lugar.
    *   Considerar protección contra CSRF si se manejan sesiones de usuario (aunque para un contacto anónimo es menos crítico).
    *   Rate limiting para prevenir spam/abusos.
*   **Rendimiento:** La API debe responder en menos de 500ms bajo carga normal.
*   **Escalabilidad:** (Consideración futura, para este caso de uso simple puede no ser crítico inicialmente).
*   **Logging:** El backend debe registrar las solicitudes recibidas (con datos anonimizados si es necesario) y cualquier error que ocurra.

## 8. Diseño de la API (Endpoint)

*   **Método HTTP:** `POST`
*   **Ruta:** `/api/contact` (o `/api/v1/contact` si versionas tu API)
*   **Cuerpo de la Solicitud (Request Body - JSON):**

    ```json
    {
      "name": "Nombre Apellido",
      "email": "usuario@ejemplo.com",
      "message": "Este es el contenido del mensaje.",
      "subject": "Asunto del Mensaje (Opcional)"
    }
    ```

*   **Respuesta Exitosa (201 Created):**

    ```json
    {
      "message": "Mensaje enviado correctamente.",
      "contact_id": 123 // Opcional: ID del registro creado
    }
    ```

*   **Respuesta de Error de Validación (400 Bad Request):**

    ```json
    {
      "errors": {
        "email": "El formato del email es inválido.",
        "message": "El mensaje no puede estar vacío."
      }
    }
    ```

*   **Respuesta de Error del Servidor (500 Internal Server Error):**

    ```json
    {
      "error": "Ocurrió un error interno al procesar la solicitud."
    }
    ```

## 9. Estructura de la Tabla MySQL (Ejemplo)

Nombre de la tabla: `contacts`

| Columna     | Tipo         | Restricciones                    | Descripción                         |
| :---------- | :----------- | :------------------------------- | :---------------------------------- |
| `id`        | `INT`        | `PRIMARY KEY AUTO_INCREMENT`     | Identificador único del mensaje     |
| `name`      | `VARCHAR(255)`| `NOT NULL`                       | Nombre del remitente                |
| `email`     | `VARCHAR(255)`| `NOT NULL`                       | Email del remitente                 |
| `subject`   | `VARCHAR(255)`| `NULL`                           | Asunto del mensaje (opcional)       |
| `message`   | `TEXT`       | `NOT NULL`                       | Contenido del mensaje               |
| `created_at`| `TIMESTAMP`  | `DEFAULT CURRENT_TIMESTAMP`      | Fecha y hora de creación            |
| `ip_address`| `VARCHAR(45)`| `NULL`                           | (Opcional) IP del remitente         |
| `user_agent`| `TEXT`       | `NULL`                           | (Opcional) User agent del navegador |

## 10. Tecnologías Implicadas

*   **Frontend:** Vue.js (con Axios o Fetch API para las peticiones HTTP)
*   **Backend:** Go (con un router como `gorilla/mux`, `chi`, o el `http.ServeMux` estándar; y un driver de DB como `database/sql` con `go-sql-driver/mysql`)
*   **Base de Datos:** MySQL

## 11. Criterios de Aceptación (Resumen para Pruebas)

*   [ ] El usuario puede enviar un formulario con datos válidos y recibir un mensaje de éxito.
*   [ ] Los datos enviados se almacenan correctamente en la base de datos MySQL.
*   [ ] El usuario recibe mensajes de error claros si los datos del formulario son inválidos (frontend y backend).
*   [ ] El sistema maneja correctamente los errores de conexión o escritura en la base de datos.
*   [ ] El endpoint `/api/contact` responde con los códigos HTTP y cuerpos JSON esperados para cada escenario.

## 12. Notas Adicionales / Pendientes

*   Considerar la implementación de reCAPTCHA o similar para prevenir bots.
*   Definir la estrategia de logging con más detalle.
*   (Cualquier otra nota o tarea pendiente relacionada con este caso de uso)

```

**Cómo usar esta plantilla:**

1.  **Copia y Pega:** Guarda este contenido en un archivo con extensión `.md` (ej: `uc_api_contacto.md`).
2.  **Adapta:** Modifica los campos como el ID, fecha, autor, y detalles específicos de tu implementación.
3.  **Reutiliza:** Para cada nueva funcionalidad o endpoint de tu API, crea un nuevo archivo `.md` basado en esta plantilla y adáptalo.

Esta estructura te ayudará a ti y a tus estudiantes a entender claramente qué se espera de cada parte de la aplicación y cómo interactúan los componentes. ¡Mucha suerte con tu curso!