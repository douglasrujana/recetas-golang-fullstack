// backend/shared/notifications/notifier.go
package notifications

import (
	"context" // Es una buena práctica incluir contexto en las interfaces
)

// EmailData contiene la información necesaria para enviar un email.
type EmailData struct {
	To      []string // Lista de emails destinatarios
	From    string   // Email del remitente (ej: "noreply@tuapp.com")
	Subject string   // Asunto del email
	Body    string   // Cuerpo del email (puede ser texto plano o HTML)
	IsHTML  bool     // Indica si el cuerpo es HTML
}

// EmailNotifier define el contrato para cualquier servicio que envíe emails.
// Esto permite cambiar la implementación (ej: de SMTP a un servicio de email como SendGrid)
// sin modificar los servicios que lo utilizan.
type EmailNotifier interface {
	SendEmail(ctx context.Context, data EmailData) error
}