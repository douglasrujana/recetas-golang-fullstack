// backend/categorias/categoria_service.go
package categorias

import (
	// --- Imports necesarios ---
	"context"
	"errors"
	"fmt"
	"log" // Temporal, reemplazar con logger estructurado
	"strings"
	"github.com/gosimple/slug" // Para generar slugs
	"backend/shared/repository"
)

// --- Interfaz del Servicio ---
// Define el contrato que este servicio ofrece a las capas externas (Handlers).

// CategoriaService define la lógica de negocio para las categorías.
type CategoriaService interface {
	GetAll(ctx context.Context) ([]Categoria, error)
	GetByID(ctx context.Context, id uint) (*Categoria, error)
	Create(ctx context.Context, input CategoriaInputDTO) (*Categoria, error)          // Devuelve la categoría creada
	Update(ctx context.Context, id uint, input CategoriaInputDTO) (*Categoria, error) // Devuelve la categoría actualizada
	Delete(ctx context.Context, id uint) error
}
// --- Implementación Concreta del Servicio ---
// Proporciona la lógica real para la interfaz CategoriaService.

// categoriaService implementa la interfaz CategoriaService.
type categoriaService struct {
	repo CategoriaRepository // Dependencia de la interfaz del repo (Inyección)
	// logger *zap.Logger               // Ejemplo: Inyectar un logger sería ideal
}

// NewCategoriaService es la Factory Function para crear instancias del servicio.
// Recibe las dependencias (el repo) y devuelve la interfaz del servicio.
func NewCategoriaService(repo CategoriaRepository /*, logger *zap.Logger*/) CategoriaService {
	return &categoriaService{
		repo: repo,
		// logger: logger,
	}
}

// --- Implementación de los Métodos de la Interfaz ---

func (s *categoriaService) GetAll(ctx context.Context) ([]Categoria, error) {
	categorias, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("servicio: error al obtener todas las categorias: %w", err)
	}
	return categorias, nil
}

func (s *categoriaService) GetByID(ctx context.Context, id uint) (*Categoria, error) {
	categoria, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrCategoriaNotFound // Traducir error
		}
		return nil, fmt.Errorf("servicio: error al obtener categoria id %d: %w", id, err)
	}
	return categoria, nil
}

func (s *categoriaService) Create(ctx context.Context, input CategoriaInputDTO) (*Categoria, error) {
	nombreLimpio := strings.TrimSpace(input.Nombre)
	if nombreLimpio == "" {
		return nil, errors.New("servicio: el nombre de la categoría no puede estar vacío")
	}

	_, err := s.repo.GetByNombre(ctx, nombreLimpio) // Verificar duplicados
	if err == nil {
		return nil, ErrCategoriaNombreYaExiste
	}
	if !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, fmt.Errorf("servicio: error inesperado al verificar nombre '%s': %w", nombreLimpio, err)
	}

	nuevaCategoria := &Categoria{
		Nombre: nombreLimpio,
		Slug:   slug.Make(nombreLimpio), // Generar Slug
	}

	err = s.repo.Create(ctx, nuevaCategoria) // Llamar al repo
	if err != nil {
		return nil, fmt.Errorf("servicio: error al crear categoria en repositorio: %w", err)
	}

	log.Printf("Servicio: Categoría '%s' creada con ID: %d\n", nuevaCategoria.Nombre, nuevaCategoria.ID)
	return nuevaCategoria, nil // Devolver la categoría con ID
}

func (s *categoriaService) Update(ctx context.Context, id uint, input CategoriaInputDTO) (*Categoria, error) {
	nombreLimpio := strings.TrimSpace(input.Nombre)
	if nombreLimpio == "" {
		return nil, errors.New("servicio: el nombre de la categoría no puede estar vacío")
	}

	categoriaAActualizar, err := s.repo.GetByID(ctx, id) // Buscar primero
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrCategoriaNotFound
		}
		return nil, fmt.Errorf("servicio: error al buscar categoría %d para actualizar: %w", id, err)
	}

	if nombreLimpio != categoriaAActualizar.Nombre { // Verificar duplicados solo si cambia nombre
		existenteConNuevoNombre, err := s.repo.GetByNombre(ctx, nombreLimpio)
		if err == nil && existenteConNuevoNombre.ID != id {
			return nil, ErrCategoriaNombreYaExiste
		}
		if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("servicio: error inesperado al verificar nuevo nombre '%s': %w", nombreLimpio, err)
		}
	}

	categoriaAActualizar.Nombre = nombreLimpio // Actualizar datos
	categoriaAActualizar.Slug = slug.Make(nombreLimpio)

	err = s.repo.Update(ctx, categoriaAActualizar) // Llamar al repo
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("Advertencia: La categoría %d desapareció justo antes de actualizarla.\n", id)
			return nil, ErrCategoriaNotFound
		}
		return nil, fmt.Errorf("servicio: error al actualizar categoría en repositorio: %w", err)
	}

	log.Printf("Servicio: Categoría ID %d actualizada a nombre '%s'\n", categoriaAActualizar.ID, categoriaAActualizar.Nombre)
	return categoriaAActualizar, nil // Devolver categoría actualizada
}

func (s *categoriaService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id) // Verificar si existe
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrCategoriaNotFound
		}
		return fmt.Errorf("servicio: error al buscar categoría %d para eliminar: %w", id, err)
	}

	err = s.repo.Delete(ctx, id) // Llamar al repo
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Printf("Advertencia: Se intentó borrar la categoría %d pero ya no existía.\n", id)
			return ErrCategoriaNotFound
		}
		return fmt.Errorf("servicio: error al eliminar categoría en repositorio: %w", err)
	}

	log.Printf("Servicio: Categoría ID %d eliminada.\n", id)
	return nil
}
