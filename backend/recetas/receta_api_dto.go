// backend/recetas/receta_api_dto.go

// Este archivo define los DTOs (Data Transfer Objects) para la característica 'recetas'.
// DTOs son estructuras que definen cómo se serializan los datos para la API.
// Utilizan tags 'json' para el binding del cuerpo de la petición y 'binding' para validaciones de Gin.

//
// @Tags Recetas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /recetas/:id [get]
// @Param id path uint true "ID de la receta"
// @Success 200 {object} recetas.RecetaResponseDTO "Receta"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Receta no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
package recetas // Paquete de la característica 'recetas'

// Importar el DTO de respuesta de categoría del paquete 'categorias'
// para poder anidarlo en nuestra RecetaResponseDTO.
import "backend/categorias"

// RecetaRequestDTO define la estructura para crear/actualizar recetas desde la API.
// Utiliza tags 'json' para el binding del cuerpo de la petición y 'binding' para validaciones de Gin.
type RecetaRequestDTO struct {
	Nombre            string `json:"nombre" binding:"required,min=3,max=150" example:"Paella de Mariscos"` // @description Nombre de la receta
	CategoriaID       uint   `json:"categoria_id" binding:"required,gt=0" example:"1"`                   // @description ID de la categoría a la que pertenece
	TiempoPreparacion string `json:"tiempo_preparacion" binding:"required,max=50" example:"1 hora 30 mins"` // @description Tiempo estimado de preparación
	Descripcion       string `json:"descripcion" binding:"required" example:"Una deliciosa paella tradicional..."` // @description Pasos o descripción de la receta
	Foto              string `json:"foto,omitempty" example:"paella.jpg"`                               // @description Nombre del archivo de imagen o URL (opcional en request)
	// Ingredientes se manejarán probablemente a través de otro mecanismo o un array de IDs/structs.
	// Ej: IngredientesID []uint `json:"ingredientes_id,omitempty"`
}

// RecetaResponseDTO para la documentación con Swagger
// RecetaResponseDTO define la estructura para enviar datos de receta al cliente desde la API.
// Utiliza tags 'json' para la serialización a JSON.
// @description Estructura para enviar datos de receta al cliente desde la API.
type RecetaResponseDTO struct {
	ID                uint                            `json:"id" example:"1"`
	Nombre            string                          `json:"nombre" example:"Paella de Mariscos"`
	Slug              string                          `json:"slug" example:"paella-de-mariscos"`
	TiempoPreparacion string                          `json:"tiempo_preparacion" example:"1 hora 30 mins"`
	Descripcion       string                          `json:"descripcion" example:"Una deliciosa paella tradicional..."`
	Foto              string                          `json:"foto,omitempty" example:"uploads/recetas/paella.jpg"` // URL completa o path relativo accesible
	CreatedAt         string                          `json:"created_at" example:"2025-05-17T10:00:00Z"` // Formato consistente (ej: RFC3339)
	UpdatedAt         string                          `json:"updated_at" example:"2025-05-17T10:00:00Z"` // Formato consistente
	Categoria         categorias.CategoriaResponseDTO `json:"categoria"` // Objeto de Categoría anidado (usando el DTO de 'categorias')  
}