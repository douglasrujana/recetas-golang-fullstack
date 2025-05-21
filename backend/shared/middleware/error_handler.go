// backend/shared/middleware/error_handler.go

// Middleware para manejar errores de forma centralizada.
//
// Este middleware se encarga de capturar y manejar errores que ocurran en los handlers
// de la aplicación, proporcionando una respuesta consistente y estructurada.
//
// El middleware:
// - Captura errores de forma centralizada
// - Proporciona respuestas JSON con formato estándar
// - Maneja diferentes tipos de errores (4xx, 5xx)
// - Incluye detalles útiles para el cliente
// - Puede ser personalizado para adaptarse a diferentes entornos
//
// Ejemplo de uso:
//
// router.Use(shared.ErrorHandler())
//
// Nota: Se recomienda usar este middleware en la mayoría de las aplicaciones
// para una mejor experiencia de usuario y manejo de errores.

package middleware

import (
	"errors"  // Para errors.Is y errors.As
	"fmt"     // Para formatear mensajes de error
	"log"     // Para logging. A futuro, reemplazar con un logger estructurado.
	"net/http" // Para los códigos de estado HTTP
	//"strings"  // Para operaciones con strings, si fueran necesarias

	// --- Paquetes de Características ---
	// Importamos los paquetes de características para acceder a sus errores de dominio definidos.
	"backend/categorias"
	"backend/recetas"

	// --- Paquetes Compartidos ---
	"backend/shared/apitypes" // Para nuestro DTO estándar de respuesta de error

	// --- Paquetes de Terceros ---
	"github.com/gin-gonic/gin"                // El framework web
	"github.com/go-playground/validator/v10" // Para manejar errores de validación del binding de Gin
)

// formatValidationErrors es una función helper para convertir errores de validación de `validator/v10`
// en un mapa más legible para la respuesta de la API.
func formatValidationErrors(verrs validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, f := range verrs {
		// f.StructNamespace() // Proporciona el path completo al campo, ej: "User.Email"
		// f.Field()           // Proporciona solo el nombre del campo, ej: "Email"
		// f.Tag()             // Proporciona la regla de validación que falló, ej: "required", "email"
		// f.Param()           // Proporciona el parámetro de la regla, ej: "6" para "min=6"
		// f.Value()           // Proporciona el valor que falló la validación

		// Construir un mensaje de error útil basado en el tag de validación
		var msg string
		switch f.Tag() {
		case "required":
			msg = fmt.Sprintf("El campo '%s' es requerido.", f.Field())
		case "email":
			msg = fmt.Sprintf("El campo '%s' debe ser una dirección de email válida.", f.Field())
		case "min":
			msg = fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres.", f.Field(), f.Param())
		case "max":
			msg = fmt.Sprintf("El campo '%s' no debe exceder los %s caracteres.", f.Field(), f.Param())
		case "gt":
			msg = fmt.Sprintf("El campo '%s' debe ser mayor que %s.", f.Field(), f.Param())
		// Puedes añadir más 'case' para otros tags de validación comunes que uses.
		default:
			msg = fmt.Sprintf("El campo '%s' es inválido (regla de validación: '%s').", f.Field(), f.Tag())
		}
		errs[f.Field()] = msg // Usar el nombre del campo (sin el nombre del struct) como clave
	}
	return errs
}

// ErrorHandler es un middleware de Gin para manejar errores de forma centralizada en la aplicación.
// Captura errores adjuntados al contexto por los handlers y los traduce a respuestas HTTP JSON estándar.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Next() permite que los siguientes handlers en la cadena (incluyendo el handler del endpoint)
		// se ejecuten primero. Este middleware actuará DESPUÉS de ellos, al "desenrollar" la pila de llamadas.
		c.Next()

		// Si no hay errores adjuntos al contexto por los handlers, no hay nada que hacer.
		if len(c.Errors) == 0 {
			return
		}

		// Gin puede acumular múltiples errores. Usualmente, el último es el más relevante
		// o el que causó que el handler detuviera su procesamiento.
		lastErrorRegistered := c.Errors.Last()
		if lastErrorRegistered == nil { // Doble chequeo por si acaso
			return
		}

		// err es el error real que el handler (o un middleware anterior) adjuntó.
		err := lastErrorRegistered.Err

		// --- Logging del Error ---
		// Siempre es buena práctica loguear el error en el servidor para depuración,
		// especialmente los errores 5xx o cualquier error inesperado.
		// TODO: Reemplazar 'log.Printf' con un logger estructurado (slog, zap, zerolog)
		//       que incluya más contexto (RequestID, UserID si hay auth, etc.).
		log.Printf("[ErrorHandler] Error detectado: [%T] %v\n", err, err)
		// Para errores envueltos, %+v puede dar un stack trace si el error lo soporta:
		// log.Printf("[ErrorHandler] Detalle completo del error: %+v\n", err)


		// --- Mapeo del Error a Respuesta HTTP ---
		var statusCode int
		var responseBody apitypes.ErrorResponse // Usar nuestro DTO estándar para errores

		switch {
		// --- Errores de Dominio de Categorias ---
		case errors.Is(err, categorias.ErrCategoriaNotFound):
			statusCode = http.StatusNotFound // 404
			responseBody = apitypes.ErrorResponse{Error: err.Error()}
		case errors.Is(err, categorias.ErrCategoriaNombreYaExiste):
			statusCode = http.StatusConflict // 409
			responseBody = apitypes.ErrorResponse{Error: err.Error()}

		// --- Errores de Dominio de Recetas ---
		case errors.Is(err, recetas.ErrRecetaNotFound):
			statusCode = http.StatusNotFound // 404
			responseBody = apitypes.ErrorResponse{Error: err.Error()}
		case errors.Is(err, recetas.ErrRecetaNombreInvalido):
			statusCode = http.StatusBadRequest // 400
			responseBody = apitypes.ErrorResponse{Error: err.Error()}
		case errors.Is(err, recetas.ErrRecetaSinCategoria):
			// Este error es devuelto por el servicio de recetas si, por ejemplo,
			// se intenta crear una receta con un CategoriaID que no existe.
			// Un 400 (Bad Request) o 422 (Unprocessable Entity) es apropiado.
			statusCode = http.StatusBadRequest
			responseBody = apitypes.ErrorResponse{Error: err.Error()}
		// Puedes añadir más 'case errors.Is(err, recetas.OtroErrorDeReceta)' aquí.

		// --- Errores de Validación del Binding de Gin ---
		case errors.As(err, &validator.ValidationErrors{}):
			statusCode = http.StatusBadRequest // 400
			validationErrors := err.(validator.ValidationErrors)
			// Formatear los errores de validación para incluirlos en 'details'.
			responseBody = apitypes.ErrorResponse{
				Error:   "Uno o más campos fallaron la validación.",
				Details: formatValidationErrors(validationErrors),
			}
		
		// --- Otros Errores Comunes que podrías querer manejar explícitamente ---
		// Ejemplo: si un servicio devuelve un error genérico de "acción no autorizada"
		// var ErrAccionNoAutorizada = errors.New("acción no autorizada") // Definido en algún paquete
		// case errors.Is(err, ErrAccionNoAutorizada):
		//  statusCode = http.StatusForbidden // 403
		//  responseBody = apitypes.ErrorResponse{Error: err.Error()}

		// --- Caso por Defecto: Error Interno del Servidor ---
		default:
			statusCode = http.StatusInternalServerError // 500
			// En producción, NUNCA exponer err.Error() directamente para errores 500
			// si podría contener información sensible de la infraestructura.
			responseBody = apitypes.ErrorResponse{Error: "Ocurrió un error interno inesperado en el servidor."}
			// Loguear el error completo con stack trace es VITAL aquí para el equipo de desarrollo.
			log.Printf("--- DETALLE ERROR INTERNO NO MANEJADO (DEFAULT CASE) --- \nError: %v\nStack (si disponible): %+v\n---------------------------\n", err, err)
		}

		// Enviar la respuesta JSON al cliente y detener cualquier procesamiento adicional.
		// Solo enviar respuesta si no se ha enviado ya (c.Writer.Written() es false).
		// Esto previene errores de "http: superfluous response.WriteHeader call".
		if !c.Writer.Written() {
			c.AbortWithStatusJSON(statusCode, responseBody)
		}
	}
}