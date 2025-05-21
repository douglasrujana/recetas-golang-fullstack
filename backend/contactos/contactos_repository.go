// backend/contactos/repository.go
// Funcionalidad: Interfaz para la persistencia de ContactoForm.
// Capa: Repositorio (Abstracción).
package contactos

import (
	"context"
	"errors" // Para errores comunes de repositorio
)

// Errores específicos o comunes para el repositorio de contactos.
// Podrían importarse de un paquete shared/repositoryerrors si fueran más genéricos.
var (
	ErrContactoRepoNotFound = errors.New("repositorio contactos: registro no encontrado")
)

// ContactoRepository define el contrato para las operaciones de datos de ContactoForm.
type ContactoRepository interface {
	// Create guarda un nuevo mensaje de contacto en la base de datos.
	// Modifica el puntero 'contacto' para incluir el ID generado.
	Create(ctx context.Context, contacto *ContactoForm) error

	// GetByID recupera un mensaje de contacto por su ID.
	GetByID(ctx context.Context, id uint) (*ContactoForm, error)

	// GetAll recupera todos los mensajes de contacto (podría tener paginación/filtros a futuro).
	GetAll(ctx context.Context) ([]ContactoForm, error)

	// MarkAsRead marca un mensaje como leído.
	MarkAsRead(ctx context.Context, id uint) error

	// (Opcional) Delete para eliminar mensajes de contacto (considerar soft delete).
	// Delete(ctx context.Context, id uint) error
}