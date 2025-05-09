package recetas // Paquete de la característica

import "backend/categorias"

// Importar el DTO de categoría si lo anidamos y está en otro paquete (NO es el caso ahora si todo está en su paquete)
// Si CategoriaResponseDTO estuviera en 'backend/categorias', necesitaríamos:
// import "backend/categorias"
// Pero si 'RecetaResponseDTO' anida 'categorias.CategoriaResponseDTO' (el tipo del paquete categorias),
// entonces el import se manejará donde se use este DTO.
// Por ahora, para el DTO de respuesta de Receta, SÍ es común anidar el DTO de Categoría.

// RecetaRequestDTO para crear/actualizar recetas desde la API.
type RecetaRequestDTO struct {
	Nombre            string `json:"nombre" binding:"required"`
	CategoriaID       uint   `json:"categoria_id" binding:"required,gt=0"` // Validar que el ID sea positivo
	TiempoPreparacion string `json:"tiempo_preparacion" binding:"required"`
	Descripcion       string `json:"descripcion" binding:"required"`
	// Foto se manejará probablemente por un endpoint de subida separado
	// Ingredientes se manejarán por separado (ej: array de IDs o nombres)
}

// RecetaResponseDTO para enviar datos de receta al cliente desde la API.
type RecetaResponseDTO struct {
	ID                uint                   `json:"id"`
	Nombre            string                 `json:"nombre"`
	Slug              string                 `json:"slug"`
	TiempoPreparacion string                 `json:"tiempo_preparacion"`
	Descripcion       string                 `json:"descripcion"`
	Foto              string                 `json:"foto"`      // URL completa o path relativo
	FechaCreacion     string                 `json:"fecha_creacion"` // Formato consistente
	// Para anidar la información de la categoría, necesitamos el DTO de Categoría
	// Asumiendo que CategoriaResponseDTO está definido en el paquete 'categorias'
	// y ese paquete será importado donde se construya esta respuesta.
	// Si 'categorias' es un paquete, sería algo como:
	// Categoria         categorias.CategoriaResponseDTO `json:"categoria"`
	// Si está en el mismo paquete (NO recomendado para DTOs de diferentes dominios):
	Categoria categorias.CategoriaResponseDTO `json:"categoria"` // Solo si CategoriaResponseDTO está en este paquete 'recetas'
}

// Para la anidación correcta, si `CategoriaResponseDTO` vive en `backend/categorias/categoria_api_dto.go`
// y este archivo está en `backend/recetas/receta_api_dto.go`, entonces en el Handler de Recetas
// al construir la respuesta, importarías ambos paquetes:
//
// import (
//    "backend/categorias" // Para categorias.CategoriaResponseDTO
//    "backend/recetas"    // Para recetas.RecetaResponseDTO
// )
//
// func construirRespuesta() recetas.RecetaResponseDTO {
//     catDTO := categorias.CategoriaResponseDTO{...}
//     return recetas.RecetaResponseDTO{..., Categoria: catDTO}
// }