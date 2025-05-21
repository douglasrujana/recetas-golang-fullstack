// backend/contactos/service_dto.go
package contactos

// EnviarContactoInput es el DTO para la entrada del servicio de creación de contacto.
type EnviarContactoInput struct {
	Nombre            string
	Email             string
	Telefono          string // Opcional
	Asunto            string // Opcional
	Mensaje           string
	// Campos adicionales que el servicio podría necesitar o derivar:
	IPOrigen          string
	UserAgent         string
	UserID            *uint // ID del usuario logueado, si aplica
}