// backend/categorias/categoria_repository_test.go
package categorias

import (
	"backend/shared/config"
	"backend/shared/database" // Usado para conectar en SetupSuite
	"backend/shared/repository"
	"context"
	//"errors" // Necesario para errors.Is
	"fmt"     // Necesario para Printf y Sprintf
	"os"      // Necesario para os.Getenv y os.Setenv
	"testing" // Paquete de testing estándar
	"time"    // Necesario para time.Now y WithinDuration

	//"github.com/stretchr/testify/require" // Usamos require para fallos críticos
	"github.com/stretchr/testify/suite" // Para organizar tests
	"gorm.io/gorm"                      // Necesario para gorm.DeletedAt y gorm.Session
)

// CategoriaRepositoryIntegrationTestSuite define la suite para tests de integración.
type CategoriaRepositoryIntegrationTestSuite struct {
	suite.Suite                      // Embeber suite para métodos como T(), SetupTest, etc.
	db     *gorm.DB                   // Conexión a la BD de prueba REAL
	repo   CategoriaRepository // Instancia REAL del repositorio bajo test
	cfg    config.Config              // Configuración cargada
	dbName string                     // Nombre de la BD de test
}

// SetupSuite se ejecuta UNA VEZ antes de que todos los tests de la suite corran.
func (s *CategoriaRepositoryIntegrationTestSuite) SetupSuite() {
	s.T().Log("--- Iniciando SetupSuite ---")

	// --- Forzar Entorno de Test ---
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "test") // Forzar carga de config.test.yaml
	s.T().Logf("SetupSuite: Forzado APP_ENV=test (Original: '%s')", originalEnv)
	s.T().Cleanup(func() { // Restaurar al final de la suite
		os.Setenv("APP_ENV", originalEnv)
		s.T().Logf("SetupSuite Cleanup: Restaurado APP_ENV a '%s'", originalEnv)
	})

	// --- Cargar Configuración de Test ---
	cfg, err := config.LoadConfig("../shared/config") // Ruta relativa desde este archivo
	s.Require().NoError(err, "SetupSuite: Falló al cargar config de test desde ../shared/config")
	s.Require().Equal("test", cfg.AppEnv, "SetupSuite: La configuración cargada debe ser del entorno 'test'")
	s.cfg = cfg
	s.dbName = cfg.Database.Name
	s.Require().NotEmpty(s.dbName, "SetupSuite: El nombre de la BD en config.test.yaml no puede estar vacío")

	// --- Conectar a la Base de Datos de Prueba ---
	db, err := database.ConnectDB(cfg.Database)
	s.Require().NoError(err, "SetupSuite: Falló al conectar a la BD de test '%s'", s.dbName)
	s.db = db
	s.T().Logf("SetupSuite: Conexión a BD '%s' exitosa.", s.dbName)

	// --- Migrar/Asegurar Tabla con AutoMigrate usando el Modelo GORM ---
	s.T().Log("SetupSuite: Ejecutando AutoMigrate para mysql.CategoriaModel...")
	err = s.db.AutoMigrate(&CategoriaModel{}) // ¡Usa el Modelo GORM!
	s.Require().NoError(err, "SetupSuite: Falló al ejecutar AutoMigrate para CategoriaModel")
	s.T().Log("SetupSuite: Tabla 'categorias' asegurada/creada vía AutoMigrate.")

	// --- Crear la Instancia del Repositorio ---
	s.repo = NewCategoriaRepository(s.db) // Inyectar la conexión de prueba
	s.T().Log("SetupSuite: Repositorio real creado.")
	s.T().Log("--- Fin SetupSuite ---")
}

// TearDownSuite se ejecuta UNA VEZ después de que todos los tests de la suite hayan corrido.
func (s *CategoriaRepositoryIntegrationTestSuite) TearDownSuite() {
	s.T().Log("--- Iniciando TearDownSuite ---")
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			err = sqlDB.Close()
			s.Require().NoError(err, "TearDownSuite: Error al cerrar conexión BD")
			s.T().Log("TearDownSuite: Conexión a BD cerrada.")
		} else {
			s.T().Logf("TearDownSuite: No se pudo obtener sql.DB para cerrar: %v", err)
		}
	}
	s.T().Log("--- Fin TearDownSuite ---")
}

// SetupTest se ejecuta antes de CADA test individual.
// Limpia la tabla para asegurar aislamiento entre tests.
func (s *CategoriaRepositoryIntegrationTestSuite) SetupTest() {
	s.T().Logf("--- Iniciando SetupTest para [%s] ---", s.T().Name())
	s.Require().NotNil(s.db, "SetupTest: s.db no debe ser nil")

	// --- Usar Nombre de Tabla Explícito (obtenido de TableName()) ---
	// Es la forma más segura ya que AutoMigrate usa TableName() si existe.
	tableName := CategoriaModel{}.TableName() // Llama al método TableName() del modelo
	s.Require().NotEmpty(tableName, "SetupTest: El método TableName() no debe devolver vacío")
	s.T().Logf("SetupTest: Limpiando tabla explícita '%s'...", tableName)

	// Usar TRUNCATE TABLE (más rápido y resetea auto_increment)
	// ¡Asegúrate que no haya FKs apuntando a 'categorias' o fallará!
	err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)).Error
	s.Require().NoError(err, "SetupTest: Falló TRUNCATE TABLE para %s", tableName)

	s.T().Logf("SetupTest: Tabla '%s' truncada.", tableName)
	s.T().Log("--- Fin SetupTest ---")
}

// TestCategoriaRepositoryIntegrationTestSuite es el runner para la suite de testify.
func TestCategoriaRepositoryIntegrationTestSuite(t *testing.T) {
	if os.Getenv("APP_ENV") != "test" {
		t.Skip("Saltando tests de integración: APP_ENV no es 'test'")
	}
	suite.Run(t, new(CategoriaRepositoryIntegrationTestSuite))
}

// --- Tests de Integración Específicos ---
// (Los métodos de test individuales que ya tenías y funcionaban)

func (s *CategoriaRepositoryIntegrationTestSuite) TestCreateAndGetByID_Success() {
	ctx := context.Background()
	categoriaACrear := &Categoria{
		Nombre: " Postres Fríos ",
		Slug:   "postres-frios",
	}
	err := s.repo.Create(ctx, categoriaACrear)
	s.Require().NoError(err, "Create debe ser exitoso")
	s.Require().NotZero(categoriaACrear.ID, "Create debe asignar ID")
	idCreado := categoriaACrear.ID

	categoriaObtenida, err := s.repo.GetByID(ctx, idCreado)
	s.Require().NoError(err, "GetByID debe ser exitoso post-create")
	s.Require().NotNil(categoriaObtenida)
	s.Equal(idCreado, categoriaObtenida.ID)
	s.Equal(" Postres Fríos ", categoriaObtenida.Nombre)
	s.Equal("postres-frios", categoriaObtenida.Slug)
	s.WithinDuration(time.Now(), categoriaObtenida.CreatedAt, 10*time.Second)
	s.WithinDuration(time.Now(), categoriaObtenida.UpdatedAt, 10*time.Second)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	categoriaObtenida, err := s.repo.GetByID(ctx, 99999)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repository.ErrRecordNotFound)
	s.Nil(categoriaObtenida)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetAll_EmptyAndWithData() {
	ctx := context.Background()
	categoriasVacias, err := s.repo.GetAll(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(categoriasVacias)
	s.Require().Len(categoriasVacias, 0)

	cat1 := &Categoria{Nombre: "Ensaladas", Slug: "ensaladas"}
	cat2 := &Categoria{Nombre: "Sopas", Slug: "sopas"}
	s.Require().NoError(s.repo.Create(ctx, cat1))
	s.Require().NoError(s.repo.Create(ctx, cat2))

	categoriasConDatos, err := s.repo.GetAll(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(categoriasConDatos)
	s.Require().Len(categoriasConDatos, 2)
	s.Equal("Sopas", categoriasConDatos[0].Nombre)
	s.Equal("Ensaladas", categoriasConDatos[1].Nombre)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetBySlug_Success() {
	ctx := context.Background()
	slugTarget := "bebidas-sin-alcohol"
	cat := &Categoria{Nombre: "Bebidas Sin Alcohol", Slug: slugTarget}
	s.Require().NoError(s.repo.Create(ctx, cat))
	categoriaObtenida, err := s.repo.GetBySlug(ctx, slugTarget)
	s.Require().NoError(err)
	s.Require().NotNil(categoriaObtenida)
	s.Equal(cat.ID, categoriaObtenida.ID)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetBySlug_NotFound() {
	ctx := context.Background()
	categoriaObtenida, err := s.repo.GetBySlug(ctx, "slug-que-no-existe")
	s.Require().Error(err)
	s.Require().ErrorIs(err, repository.ErrRecordNotFound)
	s.Nil(categoriaObtenida)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetByNombre_Success() {
	ctx := context.Background()
	nombreTarget := " Comida Italiana "
	cat := &Categoria{Nombre: nombreTarget, Slug: "comida-italiana"}
	s.Require().NoError(s.repo.Create(ctx, cat))
	categoriaObtenida, err := s.repo.GetByNombre(ctx, nombreTarget)
	s.Require().NoError(err)
	s.Require().NotNil(categoriaObtenida)
	s.Equal(cat.ID, categoriaObtenida.ID)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetByNombre_NotFound() {
	ctx := context.Background()
	categoriaObtenida, err := s.repo.GetByNombre(ctx, "Comida Que No Existe")
	s.Require().Error(err)
	s.Require().ErrorIs(err, repository.ErrRecordNotFound)
	s.Nil(categoriaObtenida)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestUpdate_Success() {
	ctx := context.Background()
	catInicial := &Categoria{Nombre: "Vegana", Slug: "vegana"}
	s.Require().NoError(s.repo.Create(ctx, catInicial))
	id := catInicial.ID
	tiempoCreacion := catInicial.CreatedAt

	catActualizada := &Categoria{
		ID:     id,
		Nombre: " Vegetariana ",
		Slug:   "vegetariana-slug",
	}
	err := s.repo.Update(ctx, catActualizada)
	s.Require().NoError(err)

	catVerificada, err := s.repo.GetByID(ctx, id)
	s.Require().NoError(err)
	s.Require().NotNil(catVerificada)
	s.Equal(" Vegetariana ", catVerificada.Nombre)
	s.Equal("vegetariana-slug", catVerificada.Slug)
	s.True(catVerificada.UpdatedAt.After(tiempoCreacion))
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestUpdate_NotFound() {
	ctx := context.Background()
	catInexistente := &Categoria{ID: 9999, Nombre: "No Existe", Slug: "no-existe"}
	err := s.repo.Update(ctx, catInexistente)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repository.ErrRecordNotFound)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestDelete_Success() {
	ctx := context.Background()
	cat := &Categoria{Nombre: "Para Borrar", Slug: "para-borrar"}
	s.Require().NoError(s.repo.Create(ctx, cat))
	id := cat.ID
	err := s.repo.Delete(ctx, id)
	s.Require().NoError(err)
	_, err = s.repo.GetByID(ctx, id)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repository.ErrRecordNotFound)
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestDelete_NotFound() {
	ctx := context.Background()
	err := s.repo.Delete(ctx, 8888)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repository.ErrRecordNotFound)
}