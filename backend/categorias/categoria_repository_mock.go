// backend/categorias/mocks/categoria_repository_mock.go
// O podría ser backend/categorias/mocks/categoria_repository_mock.go
package categorias

import (
	"context"
	"github.com/stretchr/testify/mock" // Importar el paquete de mock
)

// CategoriaRepositoryMock es una implementación mock de CategoriaRepository.
type CategoriaRepositoryMock struct {
	mock.Mock // Embeber mock.Mock
}

// Aseguramos en tiempo de compilación que implementa la interfaz.
var _ CategoriaRepository = (*CategoriaRepositoryMock)(nil)

// Implementar CADA método de la interfaz repository.CategoriaRepository

func (m *CategoriaRepositoryMock) GetAll(ctx context.Context) ([]Categoria, error) {
	// Le decimos a testify/mock qué argumentos esperamos y qué debemos devolver
	args := m.Called(ctx)
	// Args.Get(0) es el primer valor devuelto (el slice), Args.Error(1) es el segundo (el error)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Categoria), args.Error(1)
}

func (m *CategoriaRepositoryMock) GetByID(ctx context.Context, id uint) (*Categoria, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Categoria), args.Error(1)
}

func (m *CategoriaRepositoryMock) GetBySlug(ctx context.Context, slug string) (*Categoria, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Categoria), args.Error(1)
}

func (m *CategoriaRepositoryMock) GetByNombre(ctx context.Context, nombre string) (*Categoria, error) {
	args := m.Called(ctx, nombre)
	if args.Get(0) == nil {
		// Caso especial: si devolvemos nil y error, es probable que no encontrara nada
		return nil, args.Error(1)
	}
	return args.Get(0).(*Categoria), args.Error(1) // Si encuentra, devuelve categoría y error nil
}

func (m *CategoriaRepositoryMock) Create(ctx context.Context, categoria *Categoria) error {
	args := m.Called(ctx, categoria)
	// Create solo devuelve error
	return args.Error(0)
}

func (m *CategoriaRepositoryMock) Update(ctx context.Context, categoria *Categoria) error {
	args := m.Called(ctx, categoria)
	return args.Error(0)
}

func (m *CategoriaRepositoryMock) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}