// backend/recetas/receta_domain.go
package recetas // <--- PAQUETE 'recetas'

import (
	"time"
	// "backend/categorias" // <-- Necesitaremos importar el paquete 'categorias' si Categoria (struct) vive allí
                           // y es un paquete diferente. Si 'Categoria' también se define dentro de
                           // este paquete 'recetas' (lo cual NO sería lo ideal para Clean Architecture),
                           // entonces no se necesitaría este import.
                           // ASUMIREMOS QUE 'Categoria' está en el paquete 'categorias'.
	"backend/categorias" // Importamos el paquete donde está definido domain.Categoria
	"errors"             // Para definir errores específicos del dominio
)

// Receta representa la entidad de negocio pura para una receta.
type Receta struct {
	ID                uint      // Identificador único
	CategoriaID       uint      // ID de la categoría a la que pertenece (clave foránea)
	Categoria         *categorias.Categoria // Objeto Categoria anidado (del paquete 'categorias')
	Nombre            string    // Nombre de la receta
	Slug              string    // Slug para URLs
	TiempoPreparacion string    // Tiempo de preparación/cocción (ej: "30 minutos")
	Foto              string    // Nombre/ruta del archivo de foto o URL
	Descripcion       string    // Descripción o pasos
	// Fecha       time.Time // Fecha de creación/publicación - GORM puede manejar CreatedAt/UpdatedAt
	CreatedAt         time.Time // Manejado por GORM o explícitamente
	UpdatedAt         time.Time // Manejado por GORM o explícitamente
}

// Errores específicos del dominio Receta
var (
	ErrRecetaNotFound          = errors.New("receta no encontrada")
	ErrRecetaNombreInvalido    = errors.New("el nombre de la receta no es válido o está vacío")
	ErrRecetaSinCategoria      = errors.New("la receta debe pertenecer a una categoría válida")
	ErrRecetaIngredientesInvalidos = errors.New("los ingredientes proporcionados para la receta no son válidos")
	// ... otros errores que puedan surgir ...
)

// type Recetas []Receta // No es estrictamente necesario, un slice []Receta funciona igual.