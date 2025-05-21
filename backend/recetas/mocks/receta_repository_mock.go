// backend/recetas/mocks/receta_repository_mock.go
package mocks

import (
	"backend/recetas" // Para los tipos de dominio y la interfaz
	"context"
	"github.com/stretchr/testify/mock"
)

type RecetaRepositoryMock struct {
	mock.Mock
}

var _ recetas.RecetaRepository = (*RecetaRepositoryMock)(nil) // Verifica interfaz

func (m *RecetaRepositoryMock) GetAll(ctx context.Context) ([]recetas.Receta, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]recetas.Receta), args.Error(1)
}
func (m *RecetaRepositoryMock) GetByID(ctx context.Context, id uint) (*recetas.Receta, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*recetas.Receta), args.Error(1)
}
func (m *RecetaRepositoryMock) GetBySlug(ctx context.Context, slug string) (*recetas.Receta, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*recetas.Receta), args.Error(1)
}
func (m *RecetaRepositoryMock) Create(ctx context.Context, receta *recetas.Receta) error {
	args := m.Called(ctx, receta); return args.Error(0)
}
func (m *RecetaRepositoryMock) Update(ctx context.Context, receta *recetas.Receta) error {
	args := m.Called(ctx, receta); return args.Error(0)
}
func (m *RecetaRepositoryMock) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id); return args.Error(0)
}
func (m *RecetaRepositoryMock) FindByCategoriaID(ctx context.Context, categoriaID uint) ([]recetas.Receta, error) {
	args := m.Called(ctx, categoriaID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]recetas.Receta), args.Error(1)
}