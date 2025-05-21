// backend/contactos/api_dto.go
package contactos

// ContactoRequestDTO es el DTO para la entrada de la API al crear un mensaje de contacto.
type ContactoRequestDTO struct {
	Nombre   string `json:"nombre" binding:"required,min=2,max=150"`
	Email    string `json:"email" binding:"required,email"`
	Telefono string `json:"telefono,omitempty" binding:"omitempty,min=7,max=30"` // Opcional
	Asunto   string `json:"asunto,omitempty" binding:"omitempty,max=255"`        // Opcional
	Mensaje  string `json:"mensaje" binding:"required,min=10"`
}

// ContactoResponseDTO es el DTO para la respuesta de la API al crear un mensaje.
// Podría ser más elaborado o simplemente un mensaje de éxito.
type ContactoResponseDTO struct {
	ID        uint   `json:"id,omitempty"`     // Opcional: devolver el ID del mensaje guardado
	Status    string `json:"status" example:"Mensaje enviado y guardado correctamente."`
	Timestamp string `json:"timestamp" example:"2025-05-18T10:20:30Z"`
}

// ContactoListItemDTO es el DTO para listar mensajes en un panel de admin (ejemplo).
type ContactoListItemDTO struct {
	ID            uint   `json:"id"`
	NombreRemitente string `json:"nombre_remitente"`
	EmailRemitente string `json:"email_remitente"`
	Asunto        string `json:"asunto"`
	FechaContacto string `json:"fecha_contacto"` // Formateada
	Leido         bool   `json:"leido"`
}