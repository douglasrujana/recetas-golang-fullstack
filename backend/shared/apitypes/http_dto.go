// backend/shared/apitypes/http_dto.go
// dto para respuestas de error JSON de la API.

// ErrorResponse es la estructura estándar para respuestas de error JSON de la API.
//
// @Description Estructura estándar para respuestas de error JSON de la API.
// @Field Error string Mensaje principal del error
// @Field Details interface{} Opcional, para detalles adicionales (ej: errores de validación)
// @Example {"error":"Mensaje descriptivo del error","details":{"field":"problema"}}
//
// @Tags API
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /error [get]
// @Success 200 {object} apitypes.ErrorResponse "Respuesta de error"
package apitypes

// ErrorResponse es la estructura estándar para respuestas de error JSON de la API.
type ErrorResponse struct {
	Error   string      `json:"error" example:"Mensaje descriptivo del error"`           // Mensaje principal del error
	Details map[string]string `json:"details,omitempty" example:"field_name:problema de validación"` // Ejemplo más simple para map
}

// (Podrías tener aquí otros DTOs comunes de API, como para paginación)