// backend/internal/handler/routes.go (Enfoque: Solo Categorías)
package handler

import (
	"log" // Necesario para logs
	// Importar net/http si usaras constantes de status aquí (aunque se usan en el handler)
	// "net/http"

	"github.com/gin-gonic/gin"
	// No necesitamos importar otros handlers ahora
)

// RegisterRoutes configura las rutas de la API para categorías.
// [✨ AJUSTE] Solo recibe los handlers que estamos usando AHORA.
func RegisterRoutes(
	router *gin.Engine,
	categoriaHandler *CategoriaHandler, // Handler para categorías
) {
	// Prefijo base para la API
	api := router.Group("/api/v1")

	// --- Rutas de Categorías ---
	// [⚠️ ASUNCIÓN] Por ahora, asumimos que NO requieren autenticación.
	// Si necesitas autenticación, deberías añadir un middleware aquí (incluso el placeholder).
	// Ejemplo si necesitaran auth básica (a futuro):
	// categoriaRoutesAuth := api.Group("/categorias", placeholderAuthMiddleware())
	categoriaRoutes := api.Group("/categorias") // Rutas públicas por ahora
	{
		// Usamos los métodos del categoriaHandler que pasamos
		categoriaRoutes.GET("", categoriaHandler.GetAll)
		categoriaRoutes.GET("/:id", categoriaHandler.GetByID)
		categoriaRoutes.POST("", categoriaHandler.Create)
		categoriaRoutes.PUT("/:id", categoriaHandler.Update)
		categoriaRoutes.DELETE("/:id", categoriaHandler.Delete)
	}

	log.Println("🛣️  Rutas de CATEGORÍAS configuradas.")
}

// [ℹ️ INFO] Dejamos el placeholder aquí por si lo necesitas activar rápidamente luego,
// pero no se está usando en las rutas de arriba actualmente.
func placeholderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("⚠️ ADVERTENCIA: Placeholder Auth Middleware NO ACTIVO en rutas actuales, pero disponible.")
		c.Next()
	}
}