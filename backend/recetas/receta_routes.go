// backend/recetas/routes.go

// Este archivo define las rutas específicas para la entidad Receta.
// Recibe el grupo de router BASE de la API (ej: el que se crea con router.Group("/api/v1"))
// y el handler específico para recetas.

package recetas // Paquete de la característica 'recetas'

import (
	"log" // Para loguear que las rutas se configuraron
	"github.com/gin-gonic/gin"
)

// RegisterRecetaRoutes registra las rutas específicas para la entidad Receta.
// Recibe el grupo de router BASE de la API (ej: el que se crea con router.Group("/api/v1"))
// y el handler específico para recetas.
func RegisterRecetaRoutes(apiBaseGroup *gin.RouterGroup, h *RecetaHandler) {
	// Crear un subgrupo específico para recetas a partir del grupo base.
	// Esto resultará en rutas como /api/v1/recetas
	recetaRoutes := apiBaseGroup.Group("/recetas")
	{
		recetaRoutes.GET("", h.GetAll)                 // GET /api/v1/recetas
		recetaRoutes.POST("", h.Create)                // POST /api/v1/recetas
		recetaRoutes.GET("/:id", h.GetByID)            // GET /api/v1/recetas/:id
		recetaRoutes.PUT("/:id", h.Update)             // PUT /api/v1/recetas/:id
		recetaRoutes.DELETE("/:id", h.Delete)          // DELETE /api/v1/recetas/:id
		// Podríamos añadir GET /slug/:slug si implementamos GetBySlug en el handler
		// recetaRoutes.GET("/slug/:slug", h.GetBySlug)
	}

	// (Opcional) Rutas para obtener recetas por categoría.
	// Esto crea un endpoint como /api/v1/categorias/:categoria_id/recetas
	// Es una forma común de estructurar APIs anidadas o relacionadas.
	// Nota: Esto asume que el grupo base 'apiBaseGroup' es el mismo que se usa para registrar las rutas de categorías.
	// Si las rutas de categorías están en su propio subgrupo (ej: apiBaseGroup.Group("/categorias")),
	// entonces esta ruta tendría que definirse de forma diferente o el RecetaHandler
	// necesitaría ser registrado también bajo el grupo de categorías.
	// Por ahora, para simplicidad, lo ponemos directamente bajo el apiBaseGroup.
	// Si queremos que sea /api/v1/categorias/:categoria_id/recetas, el grupo base debe ser el de /api/v1
	// y el handler de categorias se encarga de /categorias, y el de recetas de /recetas y esta ruta anidada.
	// Esto podría requerir que el handler de recetas se registre también en el contexto de las rutas de categorías,
	// o que tengamos un handler específico para esta ruta anidada.
	// Por ahora, mantengamoslo simple, o considera si esta ruta es mejor que viva en el paquete 'categorias'.

	// Una forma más directa si RecetaHandler maneja esto:
	// Se registraría como: GET /api/v1/categorias/:categoria_id/recetas
	// El apiBaseGroup que recibe RegisterRecetaRoutes sería el de /api/v1
	// Y el handler de recetas debe tener un método que sepa extraer :categoria_id.
	// categoriaSpecificRoutes := apiBaseGroup.Group("/categorias/:categoria_id") // Param anidado
	// {
	//     categoriaSpecificRoutes.GET("/recetas", h.FindByCategoria) // Ruta: /api/v1/categorias/:categoria_id/recetas
	// }
	// ¡OJO! Lo anterior crea un conflicto si el paquete 'categorias' ya define rutas para /categorias/:categoria_id
	// Es mejor que la ruta de "recetas por categoría" sea parte del recurso "recetas" o "categorias" de forma consistente.
	// Por ahora, lo más simple es una ruta como:
	// GET /api/v1/recetas/categoria/:categoria_id
	recetasPorCategoria := recetaRoutes.Group("/categoria/:categoria_id") // Ruta: /api/v1/recetas/categoria/:categoria_id
	{
		recetasPorCategoria.GET("", h.FindByCategoria)
	}


	log.Println("🛣️  Rutas de Recetas configuradas.")
}