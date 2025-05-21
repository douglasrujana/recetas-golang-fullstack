// backend/contactos/repository_gorm_test.go
package contactos_test // Usar paquete _test

import (
	"context"
	//"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"backend/contactos" // El paquete bajo test
	"backend/shared/config"
	"backend/shared/database"
	// "backend/users" // Si tuvieras que crear un UserModel para la FK en ContactoModel

//"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ContactoRepositoryIntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	repo   contactos.ContactoRepository
	cfg    config.Config
	dbName string
	// testUser users.UserModel // Si necesitaras un usuario para la FK
}

func (s *ContactoRepositoryIntegrationTestSuite) SetupSuite() {
	s.T().Log("--- Iniciando SetupSuite para ContactoRepository ---")
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "test")
	s.T().Logf("SetupSuite: Forzado APP_ENV=test (Original: '%s')", originalEnv)
	s.T().Cleanup(func() { os.Setenv("APP_ENV", originalEnv); s.T().Logf("SetupSuite Cleanup: Restaurado APP_ENV a '%s'", originalEnv) })

	cfg, err := config.LoadConfig("../config") // Ruta a la carpeta 'config'
	s.Require().NoError(err, "SetupSuite: Falló al cargar config de test")
	s.Require().Equal("test", cfg.AppEnv)
	s.cfg = cfg
	s.dbName = cfg.Database.Name
	s.Require().NotEmpty(s.dbName)

	db, err := database.ConnectDB(cfg.Database)
	s.Require().NoError(err, "SetupSuite: Falló al conectar a BD de test '%s'", s.dbName)
	s.db = db
	s.T().Logf("SetupSuite: Conexión a BD '%s' exitosa.", s.dbName)

	s.T().Log("SetupSuite: Ejecutando AutoMigrate para ContactoModel...")
	// Si ContactoModel tuviera FK a UserModel, también necesitarías migrar UserModel:
	// err = s.db.AutoMigrate(&users.UserModel{}, &contactos.ContactoModel{})
	err = s.db.AutoMigrate(&contactos.ContactoModel{}) // Usa el Modelo GORM del paquete 'contactos'
	s.Require().NoError(err, "SetupSuite: Falló AutoMigrate para ContactoModel")
	s.T().Log("SetupSuite: Tabla 'contactos' asegurada/creada vía AutoMigrate.")

	s.repo = contactos.NewContactoRepository(s.db)
	s.T().Log("SetupSuite: Repositorio real creado.")
	s.T().Log("--- Fin SetupSuite para ContactoRepository ---")
}

func (s *ContactoRepositoryIntegrationTestSuite) TearDownSuite() {
	s.T().Log("--- Iniciando TearDownSuite para ContactoRepository ---")
	if s.db != nil {
		sqlDB, _ := s.db.DB(); if sqlDB != nil { sqlDB.Close() }
		s.T().Log("TearDownSuite: Conexión a BD cerrada.")
	}
	s.T().Log("--- Fin TearDownSuite para ContactoRepository ---")
}

func (s *ContactoRepositoryIntegrationTestSuite) SetupTest() {
	s.T().Logf("--- Iniciando SetupTest para [%s] ---", s.T().Name())
	s.Require().NotNil(s.db, "SetupTest: s.db no debe ser nil")

	tableName := contactos.ContactoModel{}.TableName()
	s.T().Logf("SetupTest: Limpiando tabla '%s'...", tableName)
	err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)).Error
	s.Require().NoError(err, "SetupTest: Falló TRUNCATE TABLE para %s", tableName)
	s.T().Logf("SetupTest: Tabla '%s' truncada.", tableName)
}

func TestContactoRepositoryIntegrationTestSuite(t *testing.T) {
	if os.Getenv("APP_ENV") != "test" {
		t.Skip("Saltando tests de integración: APP_ENV no es 'test'")
	}
	suite.Run(t, new(ContactoRepositoryIntegrationTestSuite))
}

// --- Tests de Integración ---

func (s *ContactoRepositoryIntegrationTestSuite) TestCreateAndGetByID_Success() {
	ctx := context.Background()
	now := time.Now().Truncate(time.Second) // Truncar para facilitar comparación
	// var testUserID uint = 1 // Descomentar y asignar si tuvieras un usuario de prueba

	contactoACrear := &contactos.ContactoForm{
		// UserID:            &testUserID, // Descomentar si usas UserID
		NombreRemitente:   "Juan Test",
		EmailRemitente:    "juan.test@example.com",
		TelefonoRemitente: "555-1234",
		Asunto:            "Consulta de Test",
		Mensaje:           "Este es el contenido del mensaje de prueba.",
		Leido:             false,
		FechaContacto:     now, // Usar un tiempo conocido
		IPOrigen:          "192.168.1.100",
		UserAgent:         "Test Agent/1.0",
	}

	// ACT: Create
	err := s.repo.Create(ctx, contactoACrear)
	// ASSERT: Create
	s.Require().NoError(err, "Create contacto debe ser exitoso")
	s.Require().NotZero(contactoACrear.ID, "Create debe asignar ID al objeto dominio original")
	idCreado := contactoACrear.ID
	// Verificar que CreatedAt y UpdatedAt (manejados por GORM en el modelo) se populen
	s.WithinDuration(now, contactoACrear.CreatedAt, 5*time.Second, "Create: CreatedAt debe ser reciente")
	s.WithinDuration(now, contactoACrear.UpdatedAt, 5*time.Second, "Create: UpdatedAt debe ser reciente")


	// ACT: GetByID
	contactoObtenido, err := s.repo.GetByID(ctx, idCreado)
	// ASSERT: GetByID
	s.Require().NoError(err, "GetByID contacto debe ser exitoso para ID creado")
	s.Require().NotNil(contactoObtenido, "GetByID debe devolver un contacto")

	s.Equal(idCreado, contactoObtenido.ID)
	s.Equal(contactoACrear.NombreRemitente, contactoObtenido.NombreRemitente)
	s.Equal(contactoACrear.EmailRemitente, contactoObtenido.EmailRemitente)
	s.Equal(contactoACrear.TelefonoRemitente, contactoObtenido.TelefonoRemitente)
	s.Equal(contactoACrear.Asunto, contactoObtenido.Asunto)
	s.Equal(contactoACrear.Mensaje, contactoObtenido.Mensaje)
	s.Equal(contactoACrear.Leido, contactoObtenido.Leido)
	// Comparar tiempo con una pequeña tolerancia o después de truncar
	s.Equal(contactoACrear.FechaContacto.Unix(), contactoObtenido.FechaContacto.Unix())
	s.Equal(contactoACrear.IPOrigen, contactoObtenido.IPOrigen)
	s.Equal(contactoACrear.UserAgent, contactoObtenido.UserAgent)
	// if contactoACrear.UserID != nil { // Descomentar si usas UserID
	//  s.Require().NotNil(contactoObtenido.UserID)
	//  s.Equal(*contactoACrear.UserID, *contactoObtenido.UserID)
	// } else {
	//  s.Nil(contactoObtenido.UserID)
	// }
}

func (s *ContactoRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	idInexistente := uint(77777)

	contactoObtenido, err := s.repo.GetByID(ctx, idInexistente)

	s.Require().Error(err, "GetByID debe devolver error para ID inexistente")
	s.Require().ErrorIs(err, contactos.ErrContactoRepoNotFound, "El error devuelto debe ser ErrContactoRepoNotFound")
	s.Nil(contactoObtenido, "GetByID no debe devolver contacto si no lo encuentra")
}

func (s *ContactoRepositoryIntegrationTestSuite) TestGetAll_EmptyAndWithData() {
	ctx := context.Background()

	// --- Test GetAll con tabla vacía ---
	contactosVacios, err := s.repo.GetAll(ctx)
	s.Require().NoError(err, "GetAll no debe fallar en tabla vacía")
	s.Require().NotNil(contactosVacios, "GetAll debe devolver un slice vacío, no nil")
	s.Require().Len(contactosVacios, 0, "El slice debe estar vacío inicialmente")

	// --- Insertar datos ---
	form1 := &contactos.ContactoForm{NombreRemitente: "Ana", EmailRemitente: "ana@test.com", Mensaje: "Msg1", FechaContacto: time.Now().Add(-1 * time.Hour)}
	form2 := &contactos.ContactoForm{NombreRemitente: "Luis", EmailRemitente: "luis@test.com", Mensaje: "Msg2", FechaContacto: time.Now()}
	s.Require().NoError(s.repo.Create(ctx, form1))
	s.Require().NoError(s.repo.Create(ctx, form2))

	// --- Test GetAll con datos ---
	contactosConDatos, err := s.repo.GetAll(ctx)
	s.Require().NoError(err, "GetAll no debe fallar con datos")
	s.Require().NotNil(contactosConDatos)
	s.Require().Len(contactosConDatos, 2, "Debe haber 2 contactos")

	// Verificar orden (ORDER BY fecha_contacto desc en la implementación del repo)
	s.Equal("Luis", contactosConDatos[0].NombreRemitente) // El más nuevo primero
	s.Equal("Ana", contactosConDatos[1].NombreRemitente)
}

func (s *ContactoRepositoryIntegrationTestSuite) TestMarkAsRead_SuccessAndNotFound() {
	ctx := context.Background()
	form := &contactos.ContactoForm{NombreRemitente: "Pedro", EmailRemitente: "pedro@test.com", Mensaje: "Para leer", FechaContacto: time.Now(), Leido: false}
	s.Require().NoError(s.repo.Create(ctx, form))
	idCreado := form.ID

	// --- Test MarkAsRead Success ---
	err := s.repo.MarkAsRead(ctx, idCreado)
	s.Require().NoError(err, "MarkAsRead debe ser exitoso para ID existente")

	// Verificar que se marcó como leído
	contactoLeido, err := s.repo.GetByID(ctx, idCreado)
	s.Require().NoError(err)
	s.Require().NotNil(contactoLeido)
	s.True(contactoLeido.Leido, "El contacto debe estar marcado como leído")
	s.True(contactoLeido.UpdatedAt.After(contactoLeido.CreatedAt), "UpdatedAt debe ser posterior a CreatedAt después de marcar como leído")

	// --- Test MarkAsRead NotFound ---
	idInexistente := uint(6666)
	err = s.repo.MarkAsRead(ctx, idInexistente)
	s.Require().Error(err, "MarkAsRead debe devolver error para ID inexistente")
	s.Require().ErrorIs(err, contactos.ErrContactoRepoNotFound)
}

// TODO:
// - Test Create con campos obligatorios faltantes (aunque esto es más del servicio o validación del handler).
//   El repositorio asumirá que los datos que llegan son válidos para inserción.
// - Test Update (si se implementa un método Update en el repositorio).
// - Test Delete (si se implementa).