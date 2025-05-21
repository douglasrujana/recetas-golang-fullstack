// backend/categorias/dto/dto.go
//
// @Tags Categorias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /categorias/:id [get]
// @Param id path uint true "ID de la categoría"
// @Success 200 {object} categorias.CategoriaResponseDTO "Categoría"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Categoría no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
package categorias // Paquete de la característica 'categorias'
// --- DTOs de Entrada (Request) ---

// CategoriaRequestDTO para crear/actualizar categorías
// [✨ MEJORA] Renombrar para claridad (Request)
type CategoriaRequestDTO struct {
	Nombre string `json:"nombre" binding:"required, min=3, example:"Postres"` // @description Nombre de la categoría (mínimo 3 caracteres)
}

// --- DTOs de Salida (Response) ---

// CategoriaResponseDTO para enviar datos de categoría al cliente
// [✨ NUEVO Y RECOMENDADO]
type CategoriaResponseDTO struct {
	ID     uint   `json:"id" example:"1"` // @description ID de la categoría
	Nombre string `json:"nombre" example:"Postres"`
	Slug   string `json:"slug" example:"postres"`
}
// Para listas de recetas en la respuesta
// type RecetaResponses []RecetaResponseDTO // Si prefieres un alias
