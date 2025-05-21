//backend/recetas/repository.go

//Este archivo definirá el contrato para todas las operaciones 
//de persistencia de datos relacionadas con las recetas. Vivirá en nuestro paquete de característica recetas.

package recetas // Pertenece al paquete de la característica 'recetas'

import (
	"context"
	// "backend/shared/repositoryerrors" // Si tuvieras errores comunes de repo en shared
	// O definir errores específicos aquí si es necesario, aunque ErrRecordNotFound podría venir de shared
	"errors" // Por ahora, para definir un error base si es necesario
)

// Errores que pueden ser devueltos por implementaciones de RecetaRepository.
// Podríamos usar los errores genéricos de un paquete 'shared/repository' o
// definir/reutilizar aquí si hay especificidad.
// Por simplicidad, asumimos que Errores de Dominio se usarán más arriba (en Servicio).
var (
	// ErrRecetaRepoNotFound es un ejemplo si quisiéramos un error de "no encontrado"
	// específico de este repositorio, pero usualmente se usa uno más genérico.
	// Por ahora, confiaremos en que el servicio traduzca un error genérico del repo.
	ErrRepoGeneneral = errors.New("repositorio: ocurrió un error inesperado")
)

// RecetaRepository define el contrato para las operaciones de datos de Recetas.
// Nota: Devuelve y acepta *Receta (el tipo de dominio de este paquete).
type RecetaRepository interface {
	// GetAll recupera todas las recetas, potencialmente con su categoría precargada.
	GetAll(ctx context.Context) ([]Receta, error)

	// GetByID recupera una receta por su ID, con su categoría precargada.
	GetByID(ctx context.Context, id uint) (*Receta, error)

	// GetBySlug recupera una receta por su slug, con su categoría precargada.
	GetBySlug(ctx context.Context, slug string) (*Receta, error)

	// Create inserta una nueva receta en la base de datos.
	// El *Receta de entrada se modifica para incluir el ID generado.
	Create(ctx context.Context, receta *Receta) error

	// Update actualiza una receta existente en la base de datos.
	Update(ctx context.Context, receta *Receta) error

	// Delete elimina una receta por su ID.
	Delete(ctx context.Context, id uint) error

	// FindByCategoriaID recupera todas las recetas pertenecientes a una categoría específica.
	FindByCategoriaID(ctx context.Context, categoriaID uint) ([]Receta, error)

	// (A futuro, cuando integremos ingredientes)
	// AddIngredienteToReceta(ctx context.Context, recetaID uint, ingredienteID uint, cantidad string) error
	// RemoveIngredienteFromReceta(ctx context.Context, recetaID uint, ingredienteID uint) error
	// GetIngredientesByRecetaID(ctx context.Context, recetaID uint) ([]Ingrediente, error) // Asumiendo que Ingrediente es un tipo de dominio
}