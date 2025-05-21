// backend/internal/repository/categoria_repository.go
package repository

import (
	"backend/internal/domain" // Depende SOLO del dominio
	"context"
)

// CategoriaRepository define los métodos para interactuar con el almacenamiento de categorías.
type CategoriaRepository interface {
	GetAll(ctx context.Context) ([]domain.Categoria, error)
	GetByID(ctx context.Context, id uint) (*domain.Categoria, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Categoria, error) // Útil para buscar por slug
	GetByNombre(ctx context.Context, nombre string) (*domain.Categoria, error) // Necesario para verificar duplicados
	Create(ctx context.Context, categoria *domain.Categoria) error // Recibe y potencialmente modifica el puntero (ej: asignando ID)
	Update(ctx context.Context, categoria *domain.Categoria) error
	Delete(ctx context.Context, id uint) error
}