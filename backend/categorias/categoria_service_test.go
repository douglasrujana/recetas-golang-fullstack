// backend/categorias/categoria_service_test.go
// Test unitarios para CategoriaService usando mocks
package categorias

import (
	"context"
	"errors"
	"testing" // Paquete estándar de testing
	//"github.com/stretchr/testify/assert" // Para aserciones
	"github.com/stretchr/testify/mock"   // Para configurar el mock
	"github.com/stretchr/testify/suite"  // Para organizar tests
	"backend/shared/repository"
)

// Definimos una suite de tests para CategoriaService
type CategoriaServiceTestSuite struct {
	suite.Suite // Embeber suite.Suite
	mockRepo    *CategoriaRepositoryMock // Instancia de nuestro mock
	service     CategoriaService               // La interfaz del servicio a probar
}

// SetupTest se ejecuta antes de CADA test en la suite
func (s *CategoriaServiceTestSuite) SetupTest() {
	s.mockRepo = new(CategoriaRepositoryMock) // Crear nueva instancia del mock
	s.service = NewCategoriaService(s.mockRepo)     // Crear nueva instancia del servicio con el mock inyectado
}

// TestCategoriaServiceTestSuite corre la suite completa
func TestCategoriaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoriaServiceTestSuite))
}

// --- Tests para cada método del servicio ---

func (s *CategoriaServiceTestSuite) TestGetAll_Success() {
	// Arrange (Preparación)
	ctx := context.Background()
	mockCategorias := []Categoria{
		{ID: 1, Nombre: "Postres", Slug: "postres"},
		{ID: 2, Nombre: "Ensaladas", Slug: "ensaladas"},
	}
	// Configurar el mock: Cuando se llame a GetAll con cualquier contexto,
	// debe devolver nuestro slice mockCategorias y un error nil.
	s.mockRepo.On("GetAll", ctx).Return(mockCategorias, nil).Once()

	// Act (Ejecución)
	categorias, err := s.service.GetAll(ctx)

	// Assert (Verificación)
	s.NoError(err)              // No debería haber error
	s.NotNil(categorias)        // El resultado no debería ser nil
	s.Len(categorias, 2)        // Debería haber 2 categorías
	s.Equal("Postres", categorias[0].Nombre)
	s.mockRepo.AssertExpectations(s.T()) // Verifica que GetAll fue llamado como se esperaba
}

func (s *CategoriaServiceTestSuite) TestGetAll_RepositoryError() {
	// Arrange
	ctx := context.Background()
	mockError := errors.New("database connection failed")
	// Configurar mock: Cuando se llame a GetAll, devolver nil y nuestro error mock.
	s.mockRepo.On("GetAll", ctx).Return(nil, mockError).Once()

	// Act
	categorias, err := s.service.GetAll(ctx)

	// Assert
	s.Error(err)                 // Debería haber un error
	s.Nil(categorias)            // El resultado debería ser nil
	s.ErrorContains(err, "servicio: error al obtener todas las categorias") // Verificar el error envuelto
	s.ErrorIs(err, mockError)    // Verificar que el error original está presente
	s.mockRepo.AssertExpectations(s.T())
}

func (s *CategoriaServiceTestSuite) TestGetByID_Success() {
	// Arrange
	ctx := context.Background()
	mockID := uint(5)
	mockCategoria := &Categoria{ID: mockID, Nombre: "Sopas", Slug: "sopas"}
	s.mockRepo.On("GetByID", ctx, mockID).Return(mockCategoria, nil).Once()

	// Act
	categoria, err := s.service.GetByID(ctx, mockID)

	// Assert
	s.NoError(err)
	s.NotNil(categoria)
	s.Equal(mockID, categoria.ID)
	s.Equal("Sopas", categoria.Nombre)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *CategoriaServiceTestSuite) TestGetByID_NotFound() {
	// Arrange
	ctx := context.Background()
	mockID := uint(99)
	// Configurar mock: GetByID debe devolver nil y el error específico del repo
	s.mockRepo.On("GetByID", ctx, mockID).Return(nil, repository.ErrRecordNotFound).Once()

	// Act
	categoria, err := s.service.GetByID(ctx, mockID)

	// Assert
	s.Error(err)                         // Esperamos un error
	s.Nil(categoria)                     // No debe devolver categoría
	s.ErrorIs(err, ErrCategoriaNotFound) // Esperamos el error de DOMINIO traducido
	s.mockRepo.AssertExpectations(s.T())
}

func (s *CategoriaServiceTestSuite) TestGetByID_RepositoryError() {
	// Arrange
	ctx := context.Background()
	mockID := uint(10)
	mockError := errors.New("random DB error")
	s.mockRepo.On("GetByID", ctx, mockID).Return(nil, mockError).Once()

	// Act
	categoria, err := s.service.GetByID(ctx, mockID)

	// Assert
	s.Error(err)
	s.Nil(categoria)
	s.ErrorContains(err, "servicio: error al obtener categoria id 10") // Check wrapped error
	s.ErrorIs(err, mockError)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *CategoriaServiceTestSuite) TestCreate_Success() {
    // Arrange
    ctx := context.Background()
    input := CategoriaInputDTO{Nombre: " Bebidas  "} // Con espacios extra
    nombreLimpio := "Bebidas"
    slugEsperado := "bebidas"

    // 1. Mock para GetByNombre: Esperamos que NO encuentre nada
    s.mockRepo.On("GetByNombre", ctx, nombreLimpio).Return(nil, repository.ErrRecordNotFound).Once()

    // 2. Mock para Create: Esperamos que se llame con el objeto correcto y devuelva nil (sin error)
    // Usamos mock.MatchedBy para verificar el contenido del argumento *domain.Categoria
    s.mockRepo.On("Create", ctx, mock.MatchedBy(func(cat *Categoria) bool {
        // Asignar ID simulado por el repo mockeado (GORM lo haría en la realidad)
        cat.ID = 100 
        // Verificar que los datos pasados al repo son correctos ANTES de asignar ID
        return cat.Nombre == nombreLimpio && cat.Slug == slugEsperado && cat.ID == 100 // El ID lo asignamos aquí para simular GORM
    })).Return(nil).Once()

    // Act
    nuevaCategoria, err := s.service.Create(ctx, input)

    // Assert
    s.NoError(err)
    s.NotNil(nuevaCategoria)
    s.Equal(nombreLimpio, nuevaCategoria.Nombre)
    s.Equal(slugEsperado, nuevaCategoria.Slug)
    s.Equal(uint(100), nuevaCategoria.ID) // Verificar el ID asignado por el "repo"
    s.mockRepo.AssertExpectations(s.T())
}

func (s *CategoriaServiceTestSuite) TestCreate_NombreVacio() {
	// Arrange
	ctx := context.Background()
	input := CategoriaInputDTO{Nombre: "   "} // Nombre vacío o solo espacios

	// Act
	nuevaCategoria, err := s.service.Create(ctx, input)

	// Assert
	s.Error(err)
	s.Nil(nuevaCategoria)
	s.EqualError(err, "servicio: el nombre de la categoría no puede estar vacío")
	// Verificar que NINGÚN método del repo fue llamado
	s.mockRepo.AssertNotCalled(s.T(), "GetByNombre")
	s.mockRepo.AssertNotCalled(s.T(), "Create")
}

func (s *CategoriaServiceTestSuite) TestCreate_NombreYaExiste() {
	// Arrange
	ctx := context.Background()
	input := CategoriaInputDTO{Nombre: "Postres"}
	nombreLimpio := "Postres"
	categoriaExistente := &Categoria{ID: 1, Nombre: nombreLimpio, Slug: "postres"}

	// Mock para GetByNombre: Devuelve la categoría existente (sin error)
	s.mockRepo.On("GetByNombre", ctx, nombreLimpio).Return(categoriaExistente, nil).Once()

	// Act
	nuevaCategoria, err := s.service.Create(ctx, input)

	// Assert
	s.Error(err)
	s.Nil(nuevaCategoria)
	s.ErrorIs(err, ErrCategoriaNombreYaExiste) // Esperamos el error de dominio específico
	s.mockRepo.AssertExpectations(s.T())
	s.mockRepo.AssertNotCalled(s.T(), "Create") // Create no debe llamarse si el nombre existe
}

func (s *CategoriaServiceTestSuite) TestCreate_RepoGetError() {
    // Arrange
    ctx := context.Background()
    input := CategoriaInputDTO{Nombre: "Nueva Cat"}
    nombreLimpio := "Nueva Cat"
    mockError := errors.New("error verificando nombre")

    // Mock para GetByNombre: Devuelve un error inesperado
    s.mockRepo.On("GetByNombre", ctx, nombreLimpio).Return(nil, mockError).Once()

    // Act
    nuevaCategoria, err := s.service.Create(ctx, input)

    // Assert
    s.Error(err)
    s.Nil(nuevaCategoria)
    s.ErrorContains(err, "servicio: error inesperado al verificar nombre")
    s.ErrorIs(err, mockError)
    s.mockRepo.AssertExpectations(s.T())
    s.mockRepo.AssertNotCalled(s.T(), "Create")
}

func (s *CategoriaServiceTestSuite) TestCreate_RepoCreateError() {
	// Arrange
	ctx := context.Background()
	input := CategoriaInputDTO{Nombre: "Desayunos"}
	nombreLimpio := "Desayunos"
	slugEsperado := "desayunos"
	mockError := errors.New("failed to insert")

	// Mock GetByNombre: No encuentra nada
	s.mockRepo.On("GetByNombre", ctx, nombreLimpio).Return(nil, repository.ErrRecordNotFound).Once()
	// Mock Create: Devuelve un error
	s.mockRepo.On("Create", ctx, mock.MatchedBy(func(cat *Categoria) bool {
		return cat.Nombre == nombreLimpio && cat.Slug == slugEsperado
	})).Return(mockError).Once()

	// Act
	nuevaCategoria, err := s.service.Create(ctx, input)

	// Assert
	s.Error(err)
	s.Nil(nuevaCategoria)
	s.ErrorContains(err, "servicio: error al crear categoria en repositorio")
	s.ErrorIs(err, mockError)
	s.mockRepo.AssertExpectations(s.T())
}

// --- Tests para Update ---
// (Similar estructura: Success, NotFound, NombreVacio, NombreYaExisteEnOtro, RepoGetError, RepoUpdateError)
func (s *CategoriaServiceTestSuite) TestUpdate_Success() {
    // Arrange
    ctx := context.Background()
    id := uint(1)
    input := CategoriaInputDTO{Nombre: " Comida Asiática "}
    nombreLimpio := "Comida Asiática"
    slugEsperado := "comida-asiatica"
    categoriaExistente := &Categoria{ID: id, Nombre: "Comida China", Slug: "comida-china"} // Nombre anterior

    // 1. Mock GetByID: Encuentra la categoría a actualizar
    s.mockRepo.On("GetByID", ctx, id).Return(categoriaExistente, nil).Once()
    // 2. Mock GetByNombre (porque el nombre cambió): No encuentra conflicto
    s.mockRepo.On("GetByNombre", ctx, nombreLimpio).Return(nil, repository.ErrRecordNotFound).Once()
    // 3. Mock Update: Se llama con el objeto actualizado y no devuelve error
    s.mockRepo.On("Update", ctx, mock.MatchedBy(func(cat *Categoria) bool {
        return cat.ID == id && cat.Nombre == nombreLimpio && cat.Slug == slugEsperado
    })).Return(nil).Once()

    // Act
    categoriaActualizada, err := s.service.Update(ctx, id, input)

    // Assert
    s.NoError(err)
    s.NotNil(categoriaActualizada)
    s.Equal(id, categoriaActualizada.ID)
    s.Equal(nombreLimpio, categoriaActualizada.Nombre)
    s.Equal(slugEsperado, categoriaActualizada.Slug)
    s.mockRepo.AssertExpectations(s.T())
}

// TODO: Añadir más tests para Update (NotFound, NombreVacio, NombreYaExisteEnOtro, RepoGetError, RepoUpdateError)

// --- Tests para Delete ---
// (Similar estructura: Success, NotFound, RepoGetError, RepoDeleteError)
func (s *CategoriaServiceTestSuite) TestDelete_Success() {
    // Arrange
    ctx := context.Background()
    id := uint(1)
     categoriaExistente := &Categoria{ID: id, Nombre: "A Borrar", Slug: "a-borrar"}

    // 1. Mock GetByID: Encuentra la categoría
    s.mockRepo.On("GetByID", ctx, id).Return(categoriaExistente, nil).Once()
    // 2. Mock Delete: No devuelve error
    s.mockRepo.On("Delete", ctx, id).Return(nil).Once()

    // Act
    err := s.service.Delete(ctx, id)

    // Assert
    s.NoError(err)
    s.mockRepo.AssertExpectations(s.T())
}

// TODO: Añadir más tests para Delete (NotFound, RepoGetError, RepoDeleteError)