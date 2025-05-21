// backend/contactos/mocks/contacto_repository_mock.go
package mocks

import (
	"backend/contactos" // Para los tipos de dominio y la interfaz
	"context"
	"github.com/stretchr/testify/mock"
)

type ContactoRepositoryMock struct {
	mock.Mock
}

// Asegurar que implementa la interfaz
var _ contactos.ContactoRepository = (*ContactoRepositoryMock)(nil)

func (m *ContactoRepositoryMock) Create(ctx context.Context, contacto *contactos.ContactoForm) error {
	args := m.Called(ctx, contacto)
	// Simular que el repo asigna ID y timestamps si el Create es exitoso
	if args.Error(0) == nil && contacto != nil {
		contacto.ID = 1 // ID de ejemplo
		// contacto.CreatedAt = time.Now() // O un tiempo fijo para tests
		// contacto.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *ContactoRepositoryMock) GetByID(ctx context.Context, id uint) (*contactos.ContactoForm, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contactos.ContactoForm), args.Error(1)
}

func (m *ContactoRepositoryMock) GetAll(ctx context.Context) ([]contactos.ContactoForm, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]contactos.ContactoForm), args.Error(1)
}

func (m *ContactoRepositoryMock) MarkAsRead(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}