// Archivo: backend/contactos/contactos_model.go
// Funcionalidad: Modelo de dominio para un Formulario de Contacto.
// Capa: Dominio / Lógica de negocio.

// Descripción:
// Define la estructura de la información que un usuario envía a través del formulario de contacto.
// Representa la entidad de negocio pura, sin acoplamientos a frameworks o bases de datos.
//
// Uso:
// - Usado por el ContactoService para procesar la información del formulario.
// - Sirve como base para la creación del ContactoModel (GORM) y los DTOs.
//
// Ciclo de Vida:
// [✔] Definido según la necesidad de recopilar información de contacto.
// [✔] Usado internamente en la lógica de validación y notificación del servicio.
// [✘] No usado directamente por la capa de base de datos (se usa ContactoModel) ni exposición HTTP (se usan DTOs).
//
// Responsabilidades:
// - Representar la información de un mensaje de contacto.
//
// Reglas de Negocio (Ejemplos que el servicio podría aplicar):
// - Nombre, Email y Mensaje son obligatorios.
// - Email debe tener un formato válido.
//
// Posibles extensiones futuras:
// - Campo para adjuntar archivos.
// - Integración con un sistema de ticketing.

package contactos

import (
	"errors"
	"time"
	// "backend/users" // Para futura integración con 'users.User'
)

// ContactoForm representa la información enviada a través del formulario de contacto.
// Es la entidad pura de negocio.
type ContactoForm struct {
	ID                uint      // Identificador único del registro de contacto
	UserID            *uint     // ID del usuario registrado que envía (opcional, puede ser nil)
	NombreRemitente   string    // Nombre proporcionado por el remitente
	EmailRemitente    string    // Email proporcionado por el remitente
	TelefonoRemitente string    // Teléfono proporcionado por el remitente (opcional)
	Asunto            string    // Asunto del mensaje (opcional)
	Mensaje           string    // Cuerpo del mensaje
	Leido             bool      // Indica si el mensaje ha sido leído por un admin
	FechaContacto     time.Time // Fecha y hora en que el usuario envió el formulario
	IPOrigen          string    // IP del remitente (para auditoría, considerar privacidad)
	UserAgent         string    // User-Agent del navegador del remitente (para auditoría)
	CreatedAt         time.Time // Timestamp de creación del registro en BD
	UpdatedAt         time.Time // Timestamp de última actualización del registro en BD
	// User           *users.User // Para cargar el objeto User completo (cuando exista el paquete 'users')
}

// Errores específicos del dominio de Contacto.
var (
	ErrContactoInvalido       = errors.New("los datos del formulario de contacto son inválidos")
	ErrNombreRemitenteVacio   = errors.New("el nombre del remitente es requerido")
	ErrEmailRemitenteInvalido = errors.New("el email del remitente es inválido o requerido")
	ErrMensajeVacio           = errors.New("el mensaje es requerido")
	ErrAsuntoDemasiadoLargo   = errors.New("el asunto no debe exceder los 255 caracteres") // Ejemplo
)
