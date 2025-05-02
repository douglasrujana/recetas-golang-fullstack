// backend/internal/domain/contacto.go
package domain

import "time"

// Contacto representa la entidad de negocio para un mensaje de contacto.
type Contacto struct {
	ID       uint      // Identificador único
	Nombre   string    // Nombre de quien contacta
	Correo   string    // Correo electrónico
	Telefono string    // Teléfono (opcional)
	Mensaje  string    // Mensaje enviado
	Fecha    time.Time // Fecha de envío
}

// type Contactos []Contacto