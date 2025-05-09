// backend/internal/repository/mysql/categoria_repository_test.go
package mysql

import (
	"backend/internal/config"
	"backend/internal/database" // Necesitaremos conectar a la BD de test
	"backend/internal/domain"
	"backend/internal/repository" // Importar interfaz y errores de repo
	"context"
	"os"      // Para leer APP_ENV
	"testing" // Paquete estándar

	"github.com/stretchr/testify/require" // Usaremos require para fallar rápido si el setup falla
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm" // Importar GORM
)

// CategoriaRepositoryIntegrationTestSuite define la suite para tests de integración.
type CategoriaRepositoryIntegrationTestSuite struct {
	suite.Suite
	db   *gorm.DB                       // Conexión a la BD de prueba REAL
	repo repository.CategoriaRepository // Instancia REAL del repositorio
	cfg  config.Config                  // Configuración cargada (para DSN si es necesario)
}

// SetupSuite se ejecuta UNA VEZ antes de que todos los tests de la suite corran.
// Perfecto para inicializar la conexión a la base de datos.
func (s *CategoriaRepositoryIntegrationTestSuite) SetupSuite() {
	// Asegurar que estamos usando la config de test
	// Nota: Idealmente, forzar APP_ENV aquí o verificar que ya está seteada.
	// Por simplicidad, asumimos que APP_ENV=test se setea al correr `go test`.
	if os.Getenv("APP_ENV") != "test" {
		s.T().Fatalf("ERROR: Los tests de integración deben ejecutarse con APP_ENV=test")
	}

	// --- Añadir log de diagnóstico (OPCIONAL PERO ÚTIL SI SIGUE FALLANDO) ---
	wd, _ := os.Getwd() // Obtener Working Directory
	s.T().Logf("Directorio de trabajo actual para test: %s", wd)
	s.T().Logf("Intentando cargar config desde ruta relativa: %s", "../../../config")
	// --------------------------------------------------------------------

	// Cargar la configuración (leerá config/config.test.yaml)
	// Usamos ".." porque el test corre desde el directorio del paquete (mysql),
	// y necesitamos subir un nivel para encontrar la carpeta config.
	cfg, err := config.LoadConfig("../../../config") // Ajusta la ruta si es necesario
	require.NoError(s.T(), err, "SetupSuite: Falló al cargar config de test")
	s.cfg = cfg // Guardar config si se necesita

	// Conectar a la base de datos de prueba
	db, err := database.ConnectDB(cfg.Database)
	require.NoError(s.T(), err, "SetupSuite: Falló al conectar a la BD de test")
	s.db = db

	// Crear la instancia REAL del repositorio con la conexión de test
	s.repo = NewCategoriaRepository(s.db)

	// (Opcional pero RECOMENDADO) Ejecutar migraciones si las tienes
	// err = s.db.AutoMigrate(&domain.Categoria{}) // ¡CUIDADO CON AUTOMIGRATE EN TESTS!
	// require.NoError(s.T(), err, "SetupSuite: Falló al ejecutar migraciones")
	// Es MEJOR tener un script SQL o una función que cree la tabla limpiamente.
	// Por ahora, asumimos que la tabla 'categorias' existe en la BD de test (MySQL la crea vacía).
	s.T().Log("SetupSuite: Conexión a BD de prueba y repositorio real listos.")
}

// TearDownSuite se ejecuta UNA VEZ después de que todos los tests de la suite hayan corrido.
// Perfecto para cerrar la conexión a la base de datos.
func (s *CategoriaRepositoryIntegrationTestSuite) TearDownSuite() {
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			err = sqlDB.Close()
			require.NoError(s.T(), err, "TearDownSuite: Error al cerrar conexión BD")
			s.T().Log("TearDownSuite: Conexión a BD de prueba cerrada.")
		}
	}
}

// BeforeTest / SetupTest se ejecuta antes de CADA test.
// Ideal para limpiar tablas y asegurar un estado inicial conocido.
func (s *CategoriaRepositoryIntegrationTestSuite) SetupTest() {
	// ¡LIMPIAR LA TABLA ANTES DE CADA TEST ES CRUCIAL!
	// Esto asegura que los tests no interfieran entre sí.
	err := s.db.Exec("DELETE FROM categorias").Error // Borrado rápido
	// err := s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.Categoria{}).Error // Alternativa GORM (más segura)
	require.NoError(s.T(), err, "SetupTest: Falló al limpiar la tabla categorias")
	// Resetear auto_increment (depende de MySQL, puede necesitar TRUNCATE o ALTER)
	// Usar TRUNCATE es más completo si no hay FKs apuntando a esta tabla.
	err = s.db.Exec("ALTER TABLE categorias AUTO_INCREMENT = 1").Error // Reiniciar contador
	// err = s.db.Exec("TRUNCATE TABLE categorias").Error // Alternativa más limpia si es posible
	require.NoError(s.T(), err, "SetupTest: Falló al resetear AUTO_INCREMENT")
	s.T().Log("SetupTest: Tabla 'categorias' limpiada.")
}

// TestCategoriaRepositoryIntegrationTestSuite corre la suite.
func TestCategoriaRepositoryIntegrationTestSuite(t *testing.T) {
	// Solo correr estos tests si APP_ENV está seteado a "test"
	if os.Getenv("APP_ENV") != "test" {
		t.Skip("Saltando tests de integración: APP_ENV no es 'test'")
	}
	suite.Run(t, new(CategoriaRepositoryIntegrationTestSuite))
}

// --- Tests de Integración ---

func (s *CategoriaRepositoryIntegrationTestSuite) TestCreateAndGetByID_Success() {
	ctx := context.Background()
	categoriaACrear := &domain.Categoria{
		Nombre: " Carnes Rojas ", // Con espacios para probar trim/limpieza si lo hiciera el repo (aunque lo hace el servicio)
		Slug:   "carnes-rojas",   // El repo recibe el slug ya generado por el servicio
	}

	// --- Test Create ---
	err := s.repo.Create(ctx, categoriaACrear)
	s.Require().NoError(err, "Create debe ser exitoso")
	s.Require().NotZero(categoriaACrear.ID, "Create debe asignar un ID")
	idCreado := categoriaACrear.ID

	// --- Test GetByID ---
	categoriaObtenida, err := s.repo.GetByID(ctx, idCreado)
	s.Require().NoError(err, "GetByID debe ser exitoso para ID creado")
	s.Require().NotNil(categoriaObtenida, "GetByID debe devolver una categoría")

	// Verificar los datos
	s.Equal(idCreado, categoriaObtenida.ID)
	// El repo guarda lo que recibe, no limpia espacios (eso es tarea del servicio)
	s.Equal(" Carnes Rojas ", categoriaObtenida.Nombre)
	s.Equal("carnes-rojas", categoriaObtenida.Slug)

	// (Opcional) Verificar directamente en la BD si quieres doble seguridad
	var count int64
	s.db.Model(&domain.Categoria{}).Where("id = ?", idCreado).Count(&count)
	s.Equal(int64(1), count, "Debe existir 1 registro en la BD con el ID creado")
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	idInexistente := uint(99999)

	categoriaObtenida, err := s.repo.GetByID(ctx, idInexistente)

	s.Require().Error(err, "GetByID debe devolver error para ID inexistente")
	// Verificar que el error es el específico del repositorio
	s.Require().ErrorIs(err, repository.ErrRecordNotFound, "El error debe ser ErrRecordNotFound")
	s.Nil(categoriaObtenida, "GetByID no debe devolver categoría si no la encuentra")
}

func (s *CategoriaRepositoryIntegrationTestSuite) TestGetAll_EmptyAndWithData() {
	ctx := context.Background()

	// --- Test GetAll con tabla vacía ---
	categoriasVacias, err := s.repo.GetAll(ctx)
	s.Require().NoError(err, "GetAll no debe fallar en tabla vacía")
	s.NotNil(categoriasVacias, "GetAll debe devolver un slice vacío, no nil")
	s.Len(categoriasVacias, 0, "El slice debe estar vacío inicialmente")

	// --- Insertar datos ---
	cat1 := &domain.Categoria{Nombre: "Verduras", Slug: "verduras"}
	cat2 := &domain.Categoria{Nombre: "Frutas", Slug: "frutas"}
	s.Require().NoError(s.repo.Create(ctx, cat1))
	s.Require().NoError(s.repo.Create(ctx, cat2))

	// --- Test GetAll con datos ---
	categoriasConDatos, err := s.repo.GetAll(ctx)
	s.Require().NoError(err, "GetAll no debe fallar con datos")
	s.NotNil(categoriasConDatos)
	s.Len(categoriasConDatos, 2, "Debe haber 2 categorías")

	// Verificar orden (asumimos ORDER BY id desc como en la implementación)
	// GORM puede devolverlos en orden de ID si no hay Order o si es ASC. Ajustar si es necesario.
	// El repo los pide desc, así que el último insertado (frutas) debería ser el primero.
	s.Equal("Frutas", categoriasConDatos[0].Nombre)
	s.Equal("Verduras", categoriasConDatos[1].Nombre)
}

// TODO: Escribir tests similares para:
// - GetBySlug (Success, NotFound)
// - GetByNombre (Success, NotFound)
// - Update (Success, NotFound al intentar actualizar)
// - Delete (Success, NotFound al intentar borrar, NotFound después de borrar)
// - Create (Considerar test para violación de constraint UNIQUE si aplica en BD)
