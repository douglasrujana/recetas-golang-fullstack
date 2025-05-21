// backend/contactos/service_test.go
package contactos_test // Usar paquete _test

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/contactos" // El paquete bajo test
	contactosMocks "backend/contactos/mocks" // Mocks del paquete contactos
	notificationsMocks "backend/shared/notifications/mocks" // Mocks del paquete shared/notifications
	"backend/shared/notifications" // Para el tipo notifications.EmailData

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ContactoServiceTestSuite struct {
	suite.Suite
	mockContactoRepo *contactosMocks.ContactoRepositoryMock
	mockEmailNotifier *notificationsMocks.EmailNotifierMock
	service           contactos.ContactoService
	adminEmail        string
	fromEmail         string
}

func (s *ContactoServiceTestSuite) SetupTest() {
	s.mockContactoRepo = new(contactosMocks.ContactoRepositoryMock)
	s.mockEmailNotifier = new(notificationsMocks.EmailNotifierMock)
	s.adminEmail = "admin@test.com"
	s.fromEmail = "noreply@test.com"
	s.service = contactos.NewContactoService(s.mockContactoRepo, s.mockEmailNotifier, s.adminEmail, s.fromEmail)
}

func TestContactoServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContactoServiceTestSuite))
}

func (s *ContactoServiceTestSuite) TestProcesarNuevoContacto_Success() {
	ctx := context.Background()
	input := contactos.EnviarContactoInput{
		Nombre:    "Usuario de Prueba",
		Email:     "prueba@example.com",
		Mensaje:   "Este es un mensaje de prueba.",
		Asunto:    "Test Asunto",
		IPOrigen:  "127.0.0.1",
		UserAgent: "Go Test",
	}

	// Arrange: Mockear ContactoRepository.Create
	s.mockContactoRepo.On("Create", ctx, mock.AnythingOfType("*contactos.ContactoForm")).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*contactos.ContactoForm)
		arg.ID = 1 // Simular que el repo asigna ID
		arg.CreatedAt = time.Now()
		arg.UpdatedAt = time.Now()
		arg.FechaContacto = arg.CreatedAt // Asumir que se setea igual que CreatedAt en el servicio
	}).Return(nil).Once()

	// Arrange: Mockear EmailNotifier.SendEmail
	s.mockEmailNotifier.On("SendEmail", ctx, mock.MatchedBy(func(data notifications.EmailData) bool {
		s.Contains(data.To, s.adminEmail)
		s.Equal(s.fromEmail, data.From)
		s.Contains(data.Subject, input.Asunto)
		s.Contains(data.Body, input.Nombre)
		s.Contains(data.Body, input.Email)
		s.Contains(data.Body, input.Mensaje)
		return true
	})).Return(nil).Once()

	// Act
	contactoGuardado, err := s.service.ProcesarNuevoContacto(ctx, input)

	// Assert
	s.NoError(err)
	s.NotNil(contactoGuardado)
	s.Equal(uint(1), contactoGuardado.ID)
	s.Equal(input.Nombre, contactoGuardado.NombreRemitente)
	s.Equal(input.Email, contactoGuardado.EmailRemitente)
	s.WithinDuration(time.Now(), contactoGuardado.FechaContacto, 5*time.Second) // Verificar que FechaContacto se seteó

	s.mockContactoRepo.AssertExpectations(s.T())
	s.mockEmailNotifier.AssertExpectations(s.T())
}

func (s *ContactoServiceTestSuite) TestProcesarNuevoContacto_Fail_NombreVacio() {
	ctx := context.Background()
	input := contactos.EnviarContactoInput{Email: "test@example.com", Mensaje: "Mensaje"}

	_, err := s.service.ProcesarNuevoContacto(ctx, input)

	s.Error(err)
	s.ErrorIs(err, contactos.ErrNombreRemitenteVacio)
	s.mockContactoRepo.AssertNotCalled(s.T(), "Create")
	s.mockEmailNotifier.AssertNotCalled(s.T(), "SendEmail")
}

func (s *ContactoServiceTestSuite) TestProcesarNuevoContacto_Fail_EmailInvalido() {
	ctx := context.Background()
	input := contactos.EnviarContactoInput{Nombre: "Test", Email: "email-invalido", Mensaje: "Mensaje"}

	_, err := s.service.ProcesarNuevoContacto(ctx, input)

	s.Error(err)
	s.ErrorIs(err, contactos.ErrEmailRemitenteInvalido)
}

func (s *ContactoServiceTestSuite) TestProcesarNuevoContacto_Fail_MensajeVacio() {
	ctx := context.Background()
	input := contactos.EnviarContactoInput{Nombre: "Test", Email: "test@example.com"}

	_, err := s.service.ProcesarNuevoContacto(ctx, input)

	s.Error(err)
	s.ErrorIs(err, contactos.ErrMensajeVacio)
}

func (s *ContactoServiceTestSuite) TestProcesarNuevoContacto_Fail_RepoCreateError() {
	ctx := context.Background()
	input := contactos.EnviarContactoInput{Nombre: "Test", Email: "test@example.com", Mensaje: "Msg"}
	repoError := errors.New("error de BD al crear")

	s.mockContactoRepo.On("Create", ctx, mock.AnythingOfType("*contactos.ContactoForm")).Return(repoError).Once()

	_, err := s.service.ProcesarNuevoContacto(ctx, input)

	s.Error(err)
	s.ErrorIs(err, repoError) // Esperamos el error del repo envuelto
	s.Contains(err.Error(), "error guardando contacto")
	s.mockEmailNotifier.AssertNotCalled(s.T(), "SendEmail") // No se debe enviar email si falla el guardado
}

func (s *ContactoServiceTestSuite) TestProcesarNuevoContacto_Success_EmailSendFails() {
	ctx := context.Background()
	input := contactos.EnviarContactoInput{Nombre: "Test", Email: "test@example.com", Mensaje: "Msg"}
	emailError := errors.New("smtp falló")

	s.mockContactoRepo.On("Create", ctx, mock.AnythingOfType("*contactos.ContactoForm")).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*contactos.ContactoForm); arg.ID = 1
	}).Return(nil).Once()
	s.mockEmailNotifier.On("SendEmail", ctx, mock.AnythingOfType("notifications.EmailData")).Return(emailError).Once()

	// Act: El servicio DEBERÍA guardar el contacto aunque el email falle, y loguear el error del email.
	// El error devuelto por ProcesarNuevoContacto debería ser nil en este caso,
	// ya que el guardado en BD fue exitoso.
	contactoGuardado, err := s.service.ProcesarNuevoContacto(ctx, input)

	// Assert
	s.NoError(err, "ProcesarNuevoContacto no debería fallar si solo el envío de email falla (pero sí loguear)")
	s.NotNil(contactoGuardado)
	s.Equal(uint(1), contactoGuardado.ID)
	// Aquí podrías verificar que se logueó una advertencia (requiere mockear el logger global o inyectar logger).
	// Por ahora, confiamos en el log.Printf dentro del servicio.

	s.mockContactoRepo.AssertExpectations(s.T())
	s.mockEmailNotifier.AssertExpectations(s.T())
}

// TODO: Añadir tests para ObtenerTodosLosContactos, ObtenerContactoPorID, MarcarContactoComoLeido
// cubriendo casos de éxito, "no encontrado" (repo devuelve error), y error genérico del repo.