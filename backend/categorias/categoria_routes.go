// backend/categorias/categoria_routes.go

// Este archivo define las rutas específicas para la entidad Categoria.
// Recibe el grupo de router BASE de la API (ej: el que se crea con router.Group("/api/v1"))
// y el handler específico para categorías.

package categorias // ¡El paquete debe ser 'categorias'!

import (
	"log"
	"github.com/gin-gonic/gin"
)

// RegisterCategoriaRoutes registra las rutas específicas para la entidad Categoria.
// Recibe el grupo de router BASE de la API (ej: el que se crea con router.Group("/api/v1"))
// y el handler específico para categorías.
func RegisterCategoriaRoutes(apiBaseGroup *gin.RouterGroup, h *CategoriaHandler) {
	// Crear un subgrupo específico para categorías a partir del grupo base
	// Esto resultará en rutas como /api/v1/categorias
	categoriaRoutes := apiBaseGroup.Group("/categorias")
	{
		categoriaRoutes.GET("", h.GetAll)
		categoriaRoutes.GET("/:id", h.GetByID)
		categoriaRoutes.POST("", h.Create)
		categoriaRoutes.PUT("/:id", h.Update)
		categoriaRoutes.DELETE("/:id", h.Delete)
	}
	log.Println("🛣️  Rutas de Categorías configuradas bajo el grupo API base.")
}

// El placeholderAuthMiddleware no es usado por RegisterCategoriaRoutes directamente.
// Si lo necesitas, debería estar en un paquete de middleware compartido
// o dentro de un paquete 'auth' cuando creemos esa característica.
// Por ahora, lo puedes dejar aquí comentado o moverlo si lo usas en otro lado.
/*
func placeholderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("⚠️ ADVERTENCIA: Placeholder Auth Middleware NO ACTIVO en rutas actuales, pero disponible.")
		c.Next()
	}
}
*/