//backend/recetas/domain/receta.go
package recetas

import (
	"time"
	"backend/categorias"
)

// Receta representa la entidad de negocio pura para una receta.
type Receta struct {
	ID          uint      // Identificador único
	CategoriaId uint      // ID de la categoría a la que pertenece
	Categoria   categorias.Categoria // [❓ OPCIONAL] Puedes incluir el objeto Categoria completo si es útil en el dominio.
	                      // Alternativamente, puedes quitarlo y cargarlo solo cuando sea necesario en la capa de servicio/handler.
	Nombre      string    // Nombre de la receta
	Slug        string    // Slug para URLs
	Tiempo      string    // Tiempo de preparación/cocción
	Foto        string    // Nombre/ruta del archivo de foto
	Descripcion string    // Descripción o pasos
	Fecha       time.Time // Fecha de creación/publicación
}

// type Recetas []Receta