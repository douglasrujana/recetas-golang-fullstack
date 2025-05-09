// backend/categorias/dto/dto.go
package categorias

// --- DTOs de Entrada (Request) ---

// EjemploDTO se mantiene si lo usas en algún handler de ejemplo
type EjemploDTO struct {
	Correo   string `json:"correo" binding:"required,email"`   // Añadir validación básica
	Password string `json:"password" binding:"required,min=6"` // Añadir validación básica
}

// CategoriaRequestDTO para crear/actualizar categorías
// [✨ MEJORA] Renombrar para claridad (Request)
type CategoriaRequestDTO struct {
	Nombre string `json:"nombre" binding:"required"`
}

// RecetaRequestDTO para crear/actualizar recetas
// [✨ MEJORA] Renombrar para claridad (Request)
type RecetaRequestDTO struct {
	Nombre      string `json:"nombre" binding:"required"`
	CategoriaID uint   `json:"categoria_id" binding:"required,gt=0"` // Validar que el ID sea positivo
	Tiempo      string `json:"tiempo" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	// La Foto y el Slug se suelen manejar de forma diferente (Slug generado, Foto subida)
	// Foto string `json:"foto"` // Probablemente se maneje por endpoint de upload
}

// --- DTOs de Salida (Response) ---

// CategoriaResponseDTO para enviar datos de categoría al cliente
// [✨ NUEVO Y RECOMENDADO]
type CategoriaResponseDTO struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Slug   string `json:"slug"`
}

// RecetaResponseDTO para enviar datos de receta al cliente
// [✨ MEJORA] Renombrar para claridad y asegurar consistencia
type RecetaResponseDTO struct {
	ID          uint                 `json:"id"`
	Nombre      string               `json:"nombre"`
	Slug        string               `json:"slug"`
	Tiempo      string               `json:"tiempo"`
	Descripcion string               `json:"descripcion"`
	Foto        string               `json:"foto"`      // URL completa de la foto
	Fecha       string               `json:"fecha"`     // Formato de fecha consistente (e.g., RFC3339)
	Categoria   CategoriaResponseDTO `json:"categoria"` // Incluir DTO de categoría anidado
}

// Para listas de recetas en la respuesta
// type RecetaResponses []RecetaResponseDTO // Si prefieres un alias
