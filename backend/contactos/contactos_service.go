// backend/contactos/service.go
// Funcionalidad: Lógica de negocio para Contactos.
// Capa: Servicio / Casos de Uso.
package contactos

import (
	"context"
	"fmt"
	"log" // Temporal
	"strings"
	"time"
	"errors"
	"backend/shared/notifications" // Para la interfaz EmailNotifier
)

// ContactoService define el contrato para la lógica de negocio de Contactos.
type ContactoService interface {
	ProcesarNuevoContacto(ctx context.Context, input EnviarContactoInput) (*ContactoForm, error)
	ObtenerTodosLosContactos(ctx context.Context) ([]ContactoForm, error)
	ObtenerContactoPorID(ctx context.Context, id uint) (*ContactoForm, error)
	MarcarContactoComoLeido(ctx context.Context, id uint) error
}

type contactoService struct {
	repo     ContactoRepository          // Repositorio de contactos de este paquete
	notifier notifications.EmailNotifier // Notificador de email compartido
	adminEmail string                    // Email del admin a notificar (de config)
    fromEmail  string                    // Email "De" para notificaciones (de config)
}

// NewContactoService crea una nueva instancia de ContactoService.
// Recibe el repositorio, el notificador y los emails relevantes de la configuración.
func NewContactoService(
	r ContactoRepository,
	n notifications.EmailNotifier,
	adminEmail string, // Email del admin para recibir notificaciones
    fromEmail string,  // Email 'From' para los correos de notificación
) ContactoService {
	return &contactoService{repo: r, notifier: n, adminEmail: adminEmail, fromEmail: fromEmail}
}

func (s *contactoService) ProcesarNuevoContacto(ctx context.Context, input EnviarContactoInput) (*ContactoForm, error) {
	// Validar entrada (lógica de negocio)
	if strings.TrimSpace(input.Nombre) == "" {
		return nil, ErrNombreRemitenteVacio
	}
	if strings.TrimSpace(input.Email) == "" || !strings.Contains(input.Email, "@") { // Validación simple
		return nil, ErrEmailRemitenteInvalido
	}
	if strings.TrimSpace(input.Mensaje) == "" {
		return nil, ErrMensajeVacio
	}
	if len(input.Asunto) > 255 {
		return nil, ErrAsuntoDemasiadoLargo
	}

	// Crear entidad de dominio
	contacto := &ContactoForm{
		UserID:            input.UserID, // Puede ser nil
		NombreRemitente:   input.Nombre,
		EmailRemitente:    input.Email,
		TelefonoRemitente: input.Telefono,
		Asunto:            input.Asunto,
		Mensaje:           input.Mensaje,
		Leido:             false,
		FechaContacto:     time.Now(), // La app establece la fecha de contacto
		IPOrigen:          input.IPOrigen,
		UserAgent:         input.UserAgent,
	}

	// Guardar en la base de datos
	if err := s.repo.Create(ctx, contacto); err != nil {
		return nil, fmt.Errorf("servicio contactos: error guardando contacto: %w", err)
	}
	log.Printf("Servicio: Contacto de '%s' guardado con ID: %d\n", contacto.EmailRemitente, contacto.ID)

	// Enviar notificación por email
	emailSubject := fmt.Sprintf("Nuevo Mensaje de Contacto: %s", contacto.Asunto)
	if contacto.Asunto == "" {
		emailSubject = "Nuevo Mensaje de Contacto Recibido"
	}
	emailBody := fmt.Sprintf(
		"Has recibido un nuevo mensaje de contacto:\n\n"+
			"Nombre: %s\n"+
			"Email: %s\n"+
			"Teléfono: %s\n"+
			"Asunto: %s\n"+
			"Fecha: %s\n"+
			"IP: %s\n"+
			"User Agent: %s\n\n"+
			"Mensaje:\n%s\n\n"+
			"ID del Mensaje: %d",
		contacto.NombreRemitente, contacto.EmailRemitente, contacto.TelefonoRemitente,
		contacto.Asunto, contacto.FechaContacto.Format("2006-01-02 15:04:05"),
		contacto.IPOrigen, contacto.UserAgent, contacto.Mensaje, contacto.ID,
	)

	emailData := notifications.EmailData{
		To:      []string{s.adminEmail}, // Email del admin desde config
		From:    s.fromEmail,          // Email "De" desde config
		Subject: emailSubject,
		Body:    emailBody,
		IsHTML:  false, // Por ahora, texto plano
	}

	if err := s.notifier.SendEmail(ctx, emailData); err != nil {
		// Loguear el error de email pero no hacer fallar la operación principal de guardar contacto.
		// El mensaje ya se guardó en la BD.
		log.Printf("ALERTA: Mensaje de contacto (ID: %d) guardado, PERO falló el envío de email de notificación: %v\n", contacto.ID, err)
		// Podrías tener un sistema de reintentos para emails o una cola aquí.
	} else {
		log.Printf("Servicio: Email de notificación para contacto ID %d enviado a %s\n", contacto.ID, s.adminEmail)
	}

	return contacto, nil // Devolver el contacto guardado (con ID y timestamps)
}

func (s *contactoService) ObtenerTodosLosContactos(ctx context.Context) ([]ContactoForm, error) {
	contactos, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("servicio contactos: error obteniendo todos: %w", err)
	}
	return contactos, nil
}

func (s *contactoService) ObtenerContactoPorID(ctx context.Context, id uint) (*ContactoForm, error) {
	contacto, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrContactoRepoNotFound) {
			return nil, ErrContactoInvalido // O un ErrContactoNoEncontrado específico
		}
		return nil, fmt.Errorf("servicio contactos: error obteniendo por id %d: %w", id, err)
	}
	return contacto, nil
}

func (s *contactoService) MarcarContactoComoLeido(ctx context.Context, id uint) error {
	err := s.repo.MarkAsRead(ctx, id)
	if err != nil {
		if errors.Is(err, ErrContactoRepoNotFound) {
			return ErrContactoInvalido // O un ErrContactoNoEncontrado específico
		}
		return fmt.Errorf("servicio contactos: error marcando como leído id %d: %w", id, err)
	}
	log.Printf("Servicio: Contacto ID %d marcado como leído.\n", id)
	return nil
}