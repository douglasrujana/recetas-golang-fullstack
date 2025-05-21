// backend/recetas/repository_gorm_test.go

// Este archivo define los tests de integración para el RecetaRepository.
// Utiliza GORM para la definición de la tabla y el mapeo de campos.
// No contiene tags de Gorm. Representa el negocio, no la base de datos.

package recetas_test // Usar paquete _test para probar la API pública del paquete 'recetas'

import (
	"context"
	// "errors" // Descomentar si usas errors.Is con errores de repo locales
	"fmt"
	"os"
	"testing"
	"time" // Necesario para time.Now() y WithinDuration

	"backend/categorias" // Necesitaremos crear categorías para las recetas y sus tipos
	"backend/recetas"    // El paquete bajo test y sus tipos
	"backend/shared/config"
	"backend/shared/database"

	//"github.com/stretchr/testify/require" // Usamos require para fallos críticos
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// --- Tests de Integración ---
type RecetaRepositoryIntegrationTestSuite struct {
	suite.Suite
	db            *gorm.DB                   // Conexión a la BD de prueba REAL
	recetaRepo    recetas.RecetaRepository   // Repositorio de recetas bajo test
	categoriaRepo categorias.CategoriaRepository // Para crear/gestionar categorías de prueba
	cfg           config.Config
	dbName        string
	testCategoria categorias.Categoria // Una categoría de prueba reutilizable para esta suite
}

// SetupSuite: Conectar a BD, AutoMigrate AMBOS modelos, crear repos.
func (s *RecetaRepositoryIntegrationTestSuite) SetupSuite() {
	s.T().Log("--- Iniciando SetupSuite para RecetaRepository ---")
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "test")
	s.T().Logf("SetupSuite: Forzado APP_ENV=test (Original: '%s')", originalEnv)
	s.T().Cleanup(func() { os.Setenv("APP_ENV", originalEnv); s.T().Logf("SetupSuite Cleanup: Restaurado APP_ENV a '%s'", originalEnv) })

	cfg, err := config.LoadConfig("../config") // Ruta a la carpeta 'config' que contiene config.test.yaml
	s.Require().NoError(err, "SetupSuite: Falló al cargar config de test")
	s.Require().Equal("test", cfg.AppEnv)
	s.cfg = cfg
	s.dbName = cfg.Database.Name
	s.Require().NotEmpty(s.dbName)

	db, err := database.ConnectDB(cfg.Database)
	s.Require().NoError(err, "SetupSuite: Falló al conectar a BD de test '%s'", s.dbName)
	s.db = db
	s.T().Logf("SetupSuite: Conexión a BD '%s' exitosa.", s.dbName)

	s.T().Log("SetupSuite: Ejecutando AutoMigrate para CategoriaModel y RecetaModel...")
	// ¡IMPORTANTE! Migrar AMBOS modelos para que GORM cree la FK correctamente.
	err = s.db.AutoMigrate(&categorias.CategoriaModel{}, &recetas.RecetaModel{})
	s.Require().NoError(err, "SetupSuite: Falló AutoMigrate para CategoriaModel y/o RecetaModel")
	s.T().Log("SetupSuite: Tablas 'categorias' y 'recetas' aseguradas/creadas.")

	s.recetaRepo = recetas.NewRecetaRepository(s.db)
	s.categoriaRepo = categorias.NewCategoriaRepository(s.db)
	s.T().Log("SetupSuite: Repositorios creados.")

	// --- NUEVA SECCIÓN: LIMPIAR TABLAS EN ORDEN ANTES DE CREAR DATOS PARA LA SUITE ---
	s.T().Log("SetupSuite: Limpiando tablas (recetas primero, luego categorias) antes de crear datos de prueba para la suite...")

	recetaTableName := recetas.RecetaModel{}.TableName()
	categoriaTableName := categorias.CategoriaModel{}.TableName()

	// 1. Limpiar tabla hija (recetas) PRIMERO
	// Usamos DELETE FROM y ALTER TABLE para resetear el auto_increment.
	err = s.db.Exec(fmt.Sprintf("DELETE FROM `%s`", recetaTableName)).Error
	s.Require().NoError(err, "SetupSuite: Falló DELETE FROM para '%s'", recetaTableName)
	err = s.db.Exec(fmt.Sprintf("ALTER TABLE `%s` AUTO_INCREMENT = 1", recetaTableName)).Error
	s.Require().NoError(err, "SetupSuite: Falló ALTER TABLE para '%s'", recetaTableName)
	s.T().Logf("SetupSuite: Tabla hija '%s' limpiada.", recetaTableName)

	// 2. Limpiar tabla padre (categorias) DESPUÉS
	err = s.db.Exec(fmt.Sprintf("DELETE FROM `%s`", categoriaTableName)).Error
	s.Require().NoError(err, "SetupSuite: Falló DELETE FROM para '%s'", categoriaTableName)
	err = s.db.Exec(fmt.Sprintf("ALTER TABLE `%s` AUTO_INCREMENT = 1", categoriaTableName)).Error
	s.Require().NoError(err, "SetupSuite: Falló ALTER TABLE para '%s'", categoriaTableName)
	s.T().Logf("SetupSuite: Tabla padre '%s' limpiada.", categoriaTableName)
	// --- FIN NUEVA SECCIÓN ---

	// Crear una categoría de prueba base para usar en los tests de recetas
	// Ahora no debería haber conflictos de duplicados si se ejecuta la suite varias veces
	// sin `docker-compose down -v` entre medio.
	catDePrueba := &categorias.Categoria{Nombre: "Categoría de Test para Recetas", Slug: "categoria-test-recetas"}
	err = s.categoriaRepo.Create(context.Background(), catDePrueba)
	s.Require().NoError(err, "SetupSuite: Falló al crear categoría de prueba base")
	s.testCategoria = *catDePrueba
	s.T().Logf("SetupSuite: Categoría de prueba base creada con ID: %d", s.testCategoria.ID)

	s.T().Log("--- Fin SetupSuite para RecetaRepository ---")
}

// TearDownSuite: Cerrar conexión a BD.
func (s *RecetaRepositoryIntegrationTestSuite) TearDownSuite() {
	s.T().Log("--- Iniciando TearDownSuite para RecetaRepository ---")
	if s.db != nil {
		sqlDB, _ := s.db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		s.T().Log("TearDownSuite: Conexión a BD cerrada.")
	}
	s.T().Log("--- Fin TearDownSuite para RecetaRepository ---")
}

// SetupTest: Limpiar tabla de recetas ANTES de cada test individual.
// La categoría de prueba base creada en SetupSuite se mantiene para todos los tests de esta suite.
func (s *RecetaRepositoryIntegrationTestSuite) SetupTest() {
	s.T().Logf("--- Iniciando SetupTest para [%s] ---", s.T().Name())
	s.Require().NotNil(s.db, "SetupTest: s.db no debe ser nil")

	recetaTableName := recetas.RecetaModel{}.TableName()

	// Limpiar tabla HIJA (recetas)
	s.T().Logf("SetupTest: Limpiando tabla (DELETE FROM) '%s'...", recetaTableName)
	err := s.db.Exec(fmt.Sprintf("DELETE FROM `%s`", recetaTableName)).Error
	s.Require().NoError(err, "SetupTest: Falló DELETE FROM para %s", recetaTableName)
	err = s.db.Exec(fmt.Sprintf("ALTER TABLE `%s` AUTO_INCREMENT = 1", recetaTableName)).Error
	s.Require().NoError(err, "SetupTest: Falló ALTER TABLE para '%s'", recetaTableName)
	s.T().Logf("SetupTest: Tabla '%s' limpiada.", recetaTableName)
}

// TestRecetaRepositoryIntegrationTestSuite ejecuta la suite de tests de integración.
func TestRecetaRepositoryIntegrationTestSuite(t *testing.T) {
	if os.Getenv("APP_ENV") != "test" {
		t.Skip("Saltando tests de integración: APP_ENV no es 'test'")
	}
	suite.Run(t, new(RecetaRepositoryIntegrationTestSuite))
}

// --- Tests de Integración ---
// (El resto de los métodos de test TestCreateRecetaAndGetByID_Success, TestGetAll_WithPreloadCategoria, etc.,
// permanecen igual que en la versión anterior que te pasé que SÍ pasaba los tests individuales)

func (s *RecetaRepositoryIntegrationTestSuite) TestCreateRecetaAndGetByID_Success() {
	ctx := context.Background()
	s.Require().NotZero(s.testCategoria.ID, "La categoría de prueba base debe tener un ID")

	recetaACrear := &recetas.Receta{
		Nombre:            "Sopa de Tomate Test",
		Slug:              "sopa-tomate-test",
		TiempoPreparacion: "25 minutos",
		Descripcion:       "Una sopa de tomate simple y deliciosa.",
		Foto:              "sopa_tomate.jpg",
		CategoriaID:       s.testCategoria.ID,
	}

	err := s.recetaRepo.Create(ctx, recetaACrear)
	s.Require().NoError(err, "Create receta debe ser exitoso")
	s.Require().NotZero(recetaACrear.ID, "Create debe asignar ID a la receta")
	idCreado := recetaACrear.ID

	recetaObtenida, err := s.recetaRepo.GetByID(ctx, idCreado)
	s.Require().NoError(err, "GetByID receta debe ser exitoso")
	s.Require().NotNil(recetaObtenida)
	s.Equal(idCreado, recetaObtenida.ID)
	s.Equal(recetaACrear.Nombre, recetaObtenida.Nombre)
	s.Equal(recetaACrear.CategoriaID, recetaObtenida.CategoriaID)
	s.Require().NotNil(recetaObtenida.Categoria, "La categoría precargada no debe ser nil")
	s.Equal(s.testCategoria.ID, recetaObtenida.Categoria.ID)
	s.Equal(s.testCategoria.Nombre, recetaObtenida.Categoria.Nombre)
	// Añadir verificación de timestamps si son importantes
	s.WithinDuration(time.Now(), recetaObtenida.CreatedAt, 10*time.Second)
	s.WithinDuration(time.Now(), recetaObtenida.UpdatedAt, 10*time.Second)
}

func (s *RecetaRepositoryIntegrationTestSuite) TestGetAll_WithPreloadCategoria() {
	ctx := context.Background()
	s.Require().NotZero(s.testCategoria.ID)

	rec1 := &recetas.Receta{Nombre: "Receta A", Slug: "receta-a", CategoriaID: s.testCategoria.ID, TiempoPreparacion: "10m"}
	rec2 := &recetas.Receta{Nombre: "Receta B", Slug: "receta-b", CategoriaID: s.testCategoria.ID, TiempoPreparacion: "20m"}
	s.Require().NoError(s.recetaRepo.Create(ctx, rec1))
	s.Require().NoError(s.recetaRepo.Create(ctx, rec2))

	listaRecetas, err := s.recetaRepo.GetAll(ctx)
	s.Require().NoError(err)
	s.Require().Len(listaRecetas, 2)

	for _, rec := range listaRecetas {
		s.Require().NotNil(rec.Categoria, "Cada receta en GetAll debe tener su categoría precargada")
		s.Equal(s.testCategoria.ID, rec.Categoria.ID)
	}
}

func (s *RecetaRepositoryIntegrationTestSuite) TestFindByCategoriaID_Success() {
	ctx := context.Background()
	s.Require().NotZero(s.testCategoria.ID)

	catOtra, err := s.createTestCategoria(ctx, "Otra Categoría Recetas", "otra-cat-recetas")
	s.Require().NoError(err)
	s.Require().NotZero(catOtra.ID)

	rec1 := &recetas.Receta{Nombre: "R1 CatTest", Slug: "r1-cattest-rec", CategoriaID: s.testCategoria.ID}
	rec2 := &recetas.Receta{Nombre: "R2 CatOtra", Slug: "r2-catotra-rec", CategoriaID: catOtra.ID}
	rec3 := &recetas.Receta{Nombre: "R3 CatTest", Slug: "r3-cattest-rec", CategoriaID: s.testCategoria.ID}
	s.Require().NoError(s.recetaRepo.Create(ctx, rec1))
	s.Require().NoError(s.recetaRepo.Create(ctx, rec2))
	s.Require().NoError(s.recetaRepo.Create(ctx, rec3))

	recetasFiltradas, err := s.recetaRepo.FindByCategoriaID(ctx, s.testCategoria.ID)
	s.Require().NoError(err)
	s.Require().Len(recetasFiltradas, 2, "Debe encontrar 2 recetas para s.testCategoria.ID")
	for _, rec := range recetasFiltradas {
		s.Equal(s.testCategoria.ID, rec.CategoriaID)
		s.Require().NotNil(rec.Categoria)
		s.Equal(s.testCategoria.Nombre, rec.Categoria.Nombre)
	}
}

// Helper para crear categorías de test adicionales si es necesario
func (s *RecetaRepositoryIntegrationTestSuite) createTestCategoria(ctx context.Context, nombre, slug string) (*categorias.Categoria, error) {
	// Asegurar que el slug sea único para este helper
	cat := &categorias.Categoria{Nombre: nombre, Slug: slug + fmt.Sprintf("-%d", time.Now().UnixNano())}
	err := s.categoriaRepo.Create(ctx, cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

// TODO: Implementar los tests que faltan:
// - GetByID_NotFound
// - GetBySlug_Success y GetBySlug_NotFound
// - Update_Success, Update_NotFound, Update_ChangeCategoria
// - Delete_Success, Delete_NotFound