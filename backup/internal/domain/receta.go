// backend/internal/domain/receta.go
package domain

import "time"

// Receta representa la entidad de negocio pura para una receta.
type Receta struct {
	ID          uint      // Identificador único
	CategoriaId uint      // ID de la categoría a la que pertenece
	Categoria   Categoria // [❓ OPCIONAL] Puedes incluir el objeto Categoria completo si es útil en el dominio.
	                      // Alternativamente, puedes quitarlo y cargarlo solo cuando sea necesario en la capa de servicio/handler.
	Nombre      string    // Nombre de la receta
	Slug        string    // Slug para URLs
	Tiempo      string    // Tiempo de preparación/cocción
	Foto        string    // Nombre/ruta del archivo de foto
	Descripcion string    // Descripción o pasos
	Fecha       time.Time // Fecha de creación/publicación
}

// type Recetas []Receta