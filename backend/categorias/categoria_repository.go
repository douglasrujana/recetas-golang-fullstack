// backend/categorias/categoria_repository.go

// Este archivo define la interface CategoriaRepository.
// Depende SOLO del dominio y no contiene implementaciones.

package categorias

import (
	//"backend/internal/domain" // Depende SOLO del dominio
	"context"
)

// CategoriaRepository define los métodos para interactuar con el almacenamiento de categorías.
type CategoriaRepository interface {
	GetAll(ctx context.Context) ([]Categoria, error)
	GetByID(ctx context.Context, id uint) (*Categoria, error)
	GetBySlug(ctx context.Context, slug string) (*Categoria, error) // Útil para buscar por slug
	GetByNombre(ctx context.Context, nombre string) (*Categoria, error) // Necesario para verificar duplicados
	Create(ctx context.Context, categoria *Categoria) error // Recibe y potencialmente modifica el puntero (ej: asignando ID)
	Update(ctx context.Context, categoria *Categoria) error
	Delete(ctx context.Context, id uint) error
}