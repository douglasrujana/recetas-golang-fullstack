// backend/recetas/mocks/categoria_service_mock.go
package mocks

import (
	"backend/categorias" // Para la interfaz y tipos de dominio de Categoria
	"context"
	"github.com/stretchr/testify/mock"
)

// CategoriaServiceMock es una implementación mock de CategoriaService.
type CategoriaServiceMock struct {
	mock.Mock
}

// Verifica que CategoriaServiceMock implementa la interfaz CategoriaService.
var _ categorias.CategoriaService = (*CategoriaServiceMock)(nil)

// GetAll es un mock de la función GetAll de la interfaz CategoriaService.
func (m *CategoriaServiceMock) GetAll(ctx context.Context) ([]categorias.Categoria, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]categorias.Categoria), args.Error(1)
}

// GetByID es un mock de la función GetByID de la interfaz CategoriaService.
func (m *CategoriaServiceMock) GetByID(ctx context.Context, id uint) (*categorias.Categoria, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*categorias.Categoria), args.Error(1)
}

// Create es un mock de la función Create de la interfaz CategoriaService.
func (m *CategoriaServiceMock) Create(ctx context.Context, input categorias.CategoriaInputDTO) (*categorias.Categoria, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*categorias.Categoria), args.Error(1)
}

// Update es un mock de la función Update de la interfaz CategoriaService.
func (m *CategoriaServiceMock) Update(ctx context.Context, id uint, input categorias.CategoriaInputDTO) (*categorias.Categoria, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*categorias.Categoria), args.Error(1)
}

// Delete es un mock de la función Delete de la interfaz CategoriaService.
func (m *CategoriaServiceMock) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id); return args.Error(0)
}