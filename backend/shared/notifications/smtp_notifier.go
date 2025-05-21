// backend/shared/notifications/smtp_notifier.go
package notifications

import (
	"backend/shared/config" // Para obtener las credenciales y configuración SMTP
	"context"
	"errors" // Para crear nuevos errores
	"fmt"
	"net/smtp" // Paquete estándar de Go para SMTP
	"strings"  // Para strings.Join y strings.Builder
)

// smtpNotifier implementa la interfaz EmailNotifier usando el protocolo SMTP.
type smtpNotifier struct {
	cfg config.SMTPConfig // Configuración SMTP inyectada (host, puerto, usuario, pass, from)
}

// NewSMTPNotifier es la factory function para crear una instancia de smtpNotifier.
// Recibe la configuración SMTP y devuelve la interfaz EmailNotifier.
func NewSMTPNotifier(cfg config.SMTPConfig) (EmailNotifier, error) {
	// Validar configuración SMTP esencial
	if cfg.Host == "" || cfg.Port == 0 {
		return nil, errors.New("configuración SMTP inválida: host o puerto no definidos")
	}
	// Para Mailtrap y muchos otros, username y password son necesarios para la autenticación.
	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("configuración SMTP inválida: username o password no definidos")
	}
	if cfg.From == "" { // El remitente por defecto
		return nil, errors.New("configuración SMTP inválida: email 'From' no definido")
	}

	return &smtpNotifier{cfg: cfg}, nil
}

// SendEmail envía un email usando las credenciales y servidor SMTP configurados.
func (n *smtpNotifier) SendEmail(ctx context.Context, data EmailData) error {
	// Construir la dirección del servidor SMTP
	addr := fmt.Sprintf("%s:%d", n.cfg.Host, n.cfg.Port)

	// Autenticación SMTP
	// smtp.PlainAuth se usa comúnmente. El primer argumento es la identidad (usualmente vacío).
	auth := smtp.PlainAuth("", n.cfg.Username, n.cfg.Password, n.cfg.Host)

	// Construir el mensaje del email
	var msgBuilder strings.Builder

	// From (usar el 'From' de EmailData si se provee, sino el de config)
	from := data.From
	if from == "" {
		from = n.cfg.From // Usar el remitente por defecto de la configuración
	}
	msgBuilder.WriteString(fmt.Sprintf("From: %s\r\n", from))

	// To
	if len(data.To) == 0 {
		return errors.New("email no puede ser enviado: no hay destinatarios (To)")
	}
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(data.To, ",")))

	// Subject
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))

	// Headers de Contenido
	if data.IsHTML {
		msgBuilder.WriteString("MIME-version: 1.0;\r\n")
		msgBuilder.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	} else {
		msgBuilder.WriteString("Content-Type: text/plain; charset=\"UTF-8\";\r\n")
	}
	msgBuilder.WriteString("Content-Transfer-Encoding: 8bit\r\n") // Común para UTF-8

	// Línea en blanco requerida entre headers y cuerpo
	msgBuilder.WriteString("\r\n")

	// Cuerpo del mensaje
	msgBuilder.WriteString(data.Body)

	// Enviar el email
	// El contexto (ctx) podría usarse aquí para implementar timeouts si smtp.SendMail lo soportara directamente
	// o si se envolviera en una goroutine con select y ctx.Done(). Por ahora, es para consistencia de interfaz.
	err := smtp.SendMail(addr, auth, from, data.To, []byte(msgBuilder.String()))
	if err != nil {
		return fmt.Errorf("smtpNotifier: error al enviar email: %w", err)
	}

	// log.Printf("Email enviado exitosamente a: %s (Subject: %s)", strings.Join(data.To, ", "), data.Subject)
	return nil
}