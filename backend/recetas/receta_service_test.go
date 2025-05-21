// backend/recetas/service_test.go

// Este archivo define los tests unitarios para el RecetaService.
// Utiliza mocks para simular dependencias externas.

package recetas_test // Usar paquete _test para forzar testing de API pública

import (
	"backend/categorias"    // Para Categoria y CategoriaService, ErrCategoriaNotFound
	"backend/recetas"       // El paquete que estamos probando
	"backend/recetas/mocks" // Nuestros mocks
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Test unitarios para RecetaService usando mocks
type RecetaServiceTestSuite struct {
	suite.Suite
	mockRecetaRepo   *mocks.RecetaRepositoryMock
	mockCategoriaSvc *mocks.CategoriaServiceMock
	service          recetas.RecetaService // Interfaz del servicio bajo test
	fixedTime        time.Time
}

// SetupTest se ejecuta antes de cada test
func (s *RecetaServiceTestSuite) SetupTest() {
	s.mockRecetaRepo = new(mocks.RecetaRepositoryMock)
	s.mockCategoriaSvc = new(mocks.CategoriaServiceMock)
	s.service = recetas.NewRecetaService(s.mockRecetaRepo, s.mockCategoriaSvc)
	s.fixedTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Para consistencia en CreatedAt/UpdatedAt
}

// TestRecetaServiceTestSuite se ejecuta antes de todos los tests
func TestRecetaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RecetaServiceTestSuite))
}

// TestCreate_Success
func (s *RecetaServiceTestSuite) TestCreate_Success() {
	ctx := context.Background()
	input := recetas.RecetaInputDTO{
		Nombre:            "  Paella Valenciana  ",
		CategoriaID:       1,
		TiempoPreparacion: "1 hora",
		Descripcion:       "Auténtica paella.",
		Foto:              "paella.jpg",
	}
	nombreLimpio := "Paella Valenciana"
	slugEsperado := "paella-valenciana"

	// Arrange: CategoriaService.GetByID para validar CategoriaID
	s.mockCategoriaSvc.On("GetByID", ctx, input.CategoriaID).Return(&categorias.Categoria{ID: input.CategoriaID, Nombre: "Platos Principales"}, nil).Once()

	// Arrange: RecetaRepository.Create
	// Esperamos que se llame a Create con el objeto Receta correcto.
	// También simulamos que el repo asigna ID y timestamps.
	s.mockRecetaRepo.On("Create", ctx, mock.MatchedBy(func(rec *recetas.Receta) bool {
		rec.ID = 10 // Simular ID asignado por DB
		rec.CreatedAt = s.fixedTime
		rec.UpdatedAt = s.fixedTime
		return rec.Nombre == nombreLimpio &&
			rec.Slug == slugEsperado &&
			rec.CategoriaID == input.CategoriaID &&
			rec.TiempoPreparacion == input.TiempoPreparacion &&
			rec.Descripcion == input.Descripcion &&
			rec.Foto == input.Foto
	})).Return(nil).Once()

	// Act
	nuevaReceta, err := s.service.Create(ctx, input)

	// Assert
	s.NoError(err)
	s.NotNil(nuevaReceta)
	s.Equal(uint(10), nuevaReceta.ID)
	s.Equal(nombreLimpio, nuevaReceta.Nombre)
	s.Equal(slugEsperado, nuevaReceta.Slug)
	s.Equal(input.CategoriaID, nuevaReceta.CategoriaID)
	s.Equal(s.fixedTime, nuevaReceta.CreatedAt)
	s.Equal(s.fixedTime, nuevaReceta.UpdatedAt)

	s.mockCategoriaSvc.AssertExpectations(s.T())
	s.mockRecetaRepo.AssertExpectations(s.T())
}

func (s *RecetaServiceTestSuite) TestCreate_Fail_EmptyNombre() {
	ctx := context.Background()
	input := recetas.RecetaInputDTO{Nombre: " ", CategoriaID: 1} // Nombre vacío

	nuevaReceta, err := s.service.Create(ctx, input)

	s.Error(err)
	s.Nil(nuevaReceta)
	s.ErrorIs(err, recetas.ErrRecetaNombreInvalido)
	s.mockCategoriaSvc.AssertNotCalled(s.T(), "GetByID") // No debería llamar si el nombre falla
	s.mockRecetaRepo.AssertNotCalled(s.T(), "Create")
}

// TODO: Escribir tests para Create_Fail_ZeroCategoriaID
func (s *RecetaServiceTestSuite) TestCreate_Fail_ZeroCategoriaID() {
	ctx := context.Background()
	input := recetas.RecetaInputDTO{Nombre: "Test Receta", CategoriaID: 0} // CategoriaID 0

	nuevaReceta, err := s.service.Create(ctx, input)

	s.Error(err)
	s.Nil(nuevaReceta)
	s.ErrorIs(err, recetas.ErrRecetaSinCategoria)
	s.mockCategoriaSvc.AssertNotCalled(s.T(), "GetByID")
	s.mockRecetaRepo.AssertNotCalled(s.T(), "Create")
}

// TODO: Escribir tests para Create_Fail_CategoriaNotFound
func (s *RecetaServiceTestSuite) TestCreate_Fail_CategoriaNotFound() {
	ctx := context.Background()
	input := recetas.RecetaInputDTO{Nombre: "Test Receta", CategoriaID: 99} // ID de categoría inexistente

	// Arrange: CategoriaService.GetByID devuelve el error específico de dominio del paquete 'categorias'
	s.mockCategoriaSvc.On("GetByID", ctx, input.CategoriaID).Return(nil, categorias.ErrCategoriaNotFound).Once()

	// Act
	nuevaReceta, err := s.service.Create(ctx, input)

	// Assert
	s.Require().Error(err, "Se esperaba un error al crear receta con categoría no existente")
	s.Nil(nuevaReceta, "No se debería devolver una receta")

	// Verificar que el error devuelto por RecetaService es (o envuelve) recetas.ErrRecetaSinCategoria
	s.True(errors.Is(err, recetas.ErrRecetaSinCategoria), "El error debe ser o envolver recetas.ErrRecetaSinCategoria")

	// Verificar que el error original categorias.ErrCategoriaNotFound está en la cadena de errores
	// Esto es importante porque el servicio DEBE haber detectado este error específico.
	s.True(errors.Is(err, categorias.ErrCategoriaNotFound), "El error original categorias.ErrCategoriaNotFound debe estar en la cadena")

	// Opcional: Verificar parte del mensaje si quieres ser específico sobre el texto
	s.Contains(err.Error(), fmt.Sprintf("la categoría ID %d no existe", input.CategoriaID), "El mensaje debe indicar qué categoría ID falló")

	s.mockCategoriaSvc.AssertExpectations(s.T())
	s.mockRecetaRepo.AssertNotCalled(s.T(), "Create") // Create del repo no debe llamarse
}

// TODO: Escribir tests para Create_Fail_CategoriaServiceError
func (s *RecetaServiceTestSuite) TestCreate_Fail_CategoriaServiceError() {
	ctx := context.Background()
	input := recetas.RecetaInputDTO{Nombre: "Test Receta", CategoriaID: 1}
	someOtherError := errors.New("error inesperado del servicio de categorías")

	s.mockCategoriaSvc.On("GetByID", ctx, input.CategoriaID).Return(nil, someOtherError).Once()

	nuevaReceta, err := s.service.Create(ctx, input)

	s.Error(err)
	s.Nil(nuevaReceta)
	s.ErrorIs(err, someOtherError) // Esperamos el error original del categoriaSvc envuelto
	s.Contains(err.Error(), "error validando categoría")

	s.mockCategoriaSvc.AssertExpectations(s.T())
	s.mockRecetaRepo.AssertNotCalled(s.T(), "Create")
}

// TODO: Escribir tests para Create_Fail_RecetaRepoError
func (s *RecetaServiceTestSuite) TestCreate_Fail_RecetaRepoError() {
	ctx := context.Background()
	input := recetas.RecetaInputDTO{Nombre: "Receta que falla al guardar", CategoriaID: 1}
	repoError := errors.New("error de base de datos al crear")

	s.mockCategoriaSvc.On("GetByID", ctx, input.CategoriaID).Return(&categorias.Categoria{ID: input.CategoriaID}, nil).Once()
	s.mockRecetaRepo.On("Create", ctx, mock.AnythingOfType("*recetas.Receta")).Return(repoError).Once()

	nuevaReceta, err := s.service.Create(ctx, input)

	s.Error(err)
	s.Nil(nuevaReceta)
	s.ErrorIs(err, repoError) // Esperamos el error original del recetaRepo envuelto
	s.Contains(err.Error(), "error al crear")

	s.mockCategoriaSvc.AssertExpectations(s.T())
	s.mockRecetaRepo.AssertExpectations(s.T())
}

// TODO: Escribir tests para GetAll, GetByID, GetBySlug, Update, Delete, FindByCategoriaID
// Siguiendo patrones similares: caso de éxito, caso de "no encontrado", caso de error del repositorio.
// Para Update, también casos de validación de CategoriaID y receta no encontrada.
