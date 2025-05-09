// backend/categorias/dto/dto.go
package categorias

// --- DTOs de Entrada (Request) ---

// CategoriaRequestDTO para crear/actualizar categorías
// [✨ MEJORA] Renombrar para claridad (Request)
type CategoriaRequestDTO struct {
	Nombre string `json:"nombre" binding:"required"`
}

// --- DTOs de Salida (Response) ---

// CategoriaResponseDTO para enviar datos de categoría al cliente
// [✨ NUEVO Y RECOMENDADO]
type CategoriaResponseDTO struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Slug   string `json:"slug"`
}

// Para listas de recetas en la respuesta
// type RecetaResponses []RecetaResponseDTO // Si prefieres un alias
