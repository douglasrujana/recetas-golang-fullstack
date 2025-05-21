// backend/shared/notifications/mocks/email_notifier_mock.go
package mocks // O el nombre de paquete que uses para mocks compartidos

import (
	"backend/shared/notifications" // La interfaz que estamos mockeando
	"context"
	"github.com/stretchr/testify/mock"
)

type EmailNotifierMock struct {
	mock.Mock
}

// Asegurar que implementa la interfaz
var _ notifications.EmailNotifier = (*EmailNotifierMock)(nil)

func (m *EmailNotifierMock) SendEmail(ctx context.Context, data notifications.EmailData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}