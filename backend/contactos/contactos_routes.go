// backend/contactos/routes.go
package contactos

import (
	"log"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas para la funcionalidad de Contactos.
func RegisterContactoRoutes(apiBaseGroup *gin.RouterGroup, h *ContactoHandler /*, authMiddleware gin.HandlerFunc */) {
	// Rutas públicas para enviar mensajes de contacto
	contactosPublicRoutes := apiBaseGroup.Group("/contactos")
	{
		contactosPublicRoutes.POST("", h.EnviarMensaje)
	}

	// Rutas para administración de contactos (requerirían autenticación de admin)
	// Ejemplo: Asumiendo que el middleware de autenticación de admin se aplica a 'adminApiGroup' en main.go
	// o se pasa aquí como 'authMiddleware'.
	contactosAdminRoutes := apiBaseGroup.Group("/admin/contactos")
	// Ejemplo de cómo aplicar un middleware específico aquí si se pasara:
	// if authMiddleware != nil {
	//  contactosAdminRoutes.Use(authMiddleware)
	// }
	{
		contactosAdminRoutes.GET("", h.GetAllContactos)          // Necesita Auth Admin
		contactosAdminRoutes.PATCH("/:id/leido", h.MarcarComoLeido) // Necesita Auth Admin
	}

	log.Println("🛣️  Rutas de Contactos configuradas.")
}