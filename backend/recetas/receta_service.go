// backend/recetas/service.go

// Este archivo define la implementación con GORM de RecetaRepository.
// Pertenece al paquete de la característica 'recetas'

package recetas // Paquete de la característica 'recetas'

import (
	"context" // El contexto es necesario para las operaciones asíncronas
	"errors" // Para errors.Is y crear nuevos errores
	"fmt" // Para formateo básico de errores si no usamos helper/middleware
	"log" // Temporal, reemplazar con logger estructurado
	"strings" // Para formatear errores
	"backend/categorias" // Importar para usar la interfaz CategoriaService y errores de dominio de categoría
	"github.com/gosimple/slug" // Para generar slugs
	"backend/shared/repository" // Importar para usar la interfaz RecetaRepository y errores de dominio de receta
)

// RecetaService define el contrato para la lógica de negocio de Recetas.
// Devuelve y acepta objetos de dominio (Receta de este paquete).
type RecetaService interface {
	GetAll(ctx context.Context) ([]Receta, error) // Devuelve todas las recetas
	GetByID(ctx context.Context, id uint) (*Receta, error) // Devuelve una receta por su ID
	GetBySlug(ctx context.Context, slug string) (*Receta, error) // Devuelve una receta por su slug
	Create(ctx context.Context, input RecetaInputDTO) (*Receta, error)           // Devuelve la receta creada
	Update(ctx context.Context, id uint, input RecetaInputDTO) (*Receta, error) // Devuelve la receta actualizada
	Delete(ctx context.Context, id uint) error
	FindByCategoriaID(ctx context.Context, categoriaID uint) ([]Receta, error) // Devuelve recetas por categoría
}

type recetaService struct { // no exportado
	recetaRepo    RecetaRepository    // Dependencia de la interfaz del repo de este paquete
	categoriaSvc  categorias.CategoriaService // Dependencia de la interfaz de CategoriaService del paquete 'categorias'
	// logger      *zap.Logger       // Idealmente inyectar logger
}

// NewRecetaService crea una nueva instancia de RecetaService.
func NewRecetaService(
	recetaRepo RecetaRepository,
	categoriaSvc categorias.CategoriaService,
	/* logger *zap.Logger */
) RecetaService {
	return &recetaService{
		recetaRepo:    recetaRepo,
		categoriaSvc:  categoriaSvc,
		// logger: logger,
	}
}

// --- Implementación de Métodos ---

// GetAll obtiene todas las recetas.
func (s *recetaService) GetAll(ctx context.Context) ([]Receta, error) {
	recs, err := s.recetaRepo.GetAll(ctx)
	if err != nil {
		// s.logger.Error("Error en servicio GetAll Recetas", zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error al obtener todas: %w", err)
	}
	return recs, nil
}

// GetByID obtiene una receta por su ID.
func (s *recetaService) GetByID(ctx context.Context, id uint) (*Receta, error) {
	rec, err := s.recetaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) { // Asume ErrRecordNotFound definido en este paquete (en repository.go)
			return nil, ErrRecetaNotFound // Traducir a error de dominio de este paquete
		}
		// s.logger.Error("Error en servicio GetByID Receta", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error al obtener por id %d: %w", id, err)
	}
	return rec, nil
}

// GetBySlug obtiene una receta por su slug.
func (s *recetaService) GetBySlug(ctx context.Context, slug string) (*Receta, error) {
	rec, err := s.recetaRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrRecetaNotFound
		}
		// s.logger.Error("Error en servicio GetBySlug Receta", zap.String("slug", slug), zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error al obtener por slug %s: %w", slug, err)
	}
	return rec, nil
}

// Create crea una nueva receta.
func (s *recetaService) Create(ctx context.Context, input RecetaInputDTO) (*Receta, error) {
	// 1. Validar entrada básica
	nombreLimpio := strings.TrimSpace(input.Nombre)
	if nombreLimpio == "" {
		return nil, ErrRecetaNombreInvalido // Error de dominio
	}
	if input.CategoriaID == 0 {
		return nil, ErrRecetaSinCategoria // Error de dominio
	}

	// 2. Validar que la CategoriaID exista (usando CategoriaService)
	_, err := s.categoriaSvc.GetByID(ctx, input.CategoriaID)
	if err != nil {
		if errors.Is(err, categorias.ErrCategoriaNotFound) {
			// Envolver el error de dominio de receta Y el error original de categoría
			return nil, fmt.Errorf("%w (causa original: %w): la categoría ID %d no existe",
				ErrRecetaSinCategoria, // El error de este paquete recetas
				err,                   // El err original (categorias.ErrCategoriaNotFound)
				input.CategoriaID)
		}
		return nil, fmt.Errorf("servicio recetas: error validando categoría %d: %w", input.CategoriaID, err)
	}

	// 3. Generar Slug
	slugReceta := slug.Make(nombreLimpio)
	// Opcional: verificar si el slug ya existe y añadir sufijo si es necesario (lógica más compleja)
	// _, errSlug := s.recetaRepo.GetBySlug(ctx, slugReceta)
	// if errSlug == nil { return nil, errors.New("servicio recetas: el slug generado ya existe") }
	// if !errors.Is(errSlug, ErrRecordNotFound) { return nil, fmt.Errorf("error verificando slug: %w", errSlug)}

	// 4. Preparar entidad de dominio Receta
	nuevaReceta := &Receta{ // Tipo de dominio de este paquete
		Nombre:            nombreLimpio,
		Slug:              slugReceta,
		TiempoPreparacion: input.TiempoPreparacion,
		Descripcion:       input.Descripcion,
		Foto:              input.Foto, // El servicio podría procesar/validar la foto aquí
		CategoriaID:       input.CategoriaID,
		// Categoria (el struct) no se asigna aquí, el repo lo carga con Preload si se consulta
	}

	// 5. Llamar al repositorio para crear
	if err := s.recetaRepo.Create(ctx, nuevaReceta); err != nil {
		// s.logger.Error("Error en repo.Create Receta", zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error al crear: %w", err)
	}

	log.Printf("Servicio: Receta '%s' creada con ID: %d para CategoriaID: %d\n", nuevaReceta.Nombre, nuevaReceta.ID, nuevaReceta.CategoriaID)
	// Devolver la receta con ID y timestamps (que el repo debería haber llenado)
	return nuevaReceta, nil
}

// Update actualiza una receta existente.
func (s *recetaService) Update(ctx context.Context, id uint, input RecetaInputDTO) (*Receta, error) {
	// 1. Validar entrada básica
	nombreLimpio := strings.TrimSpace(input.Nombre)
	if nombreLimpio == "" {
		return nil, ErrRecetaNombreInvalido
	}
	if input.CategoriaID == 0 {
		return nil, ErrRecetaSinCategoria
	}

	// 2. Validar que la CategoriaID exista
	_, err := s.categoriaSvc.GetByID(ctx, input.CategoriaID)
	if err != nil {
		if errors.Is(err, categorias.ErrCategoriaNotFound) {
			return nil, fmt.Errorf("%w: la categoría ID %d no existe", ErrRecetaSinCategoria, input.CategoriaID)
		}
		// s.logger.Error("Error validando CategoriaID en Update Receta", zap.Uint("categoriaID", input.CategoriaID), zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error validando categoría %d: %w", input.CategoriaID, err)
	}

	// 3. Obtener receta existente para actualizar
	recetaAActualizar, err := s.recetaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrRecetaNotFound
		}
		// s.logger.Error("Error buscando receta para Update", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error buscando para update %d: %w", id, err)
	}

	// 4. Actualizar campos
	recetaAActualizar.Nombre = nombreLimpio
	recetaAActualizar.Slug = slug.Make(nombreLimpio) // Regenerar slug
	recetaAActualizar.TiempoPreparacion = input.TiempoPreparacion
	recetaAActualizar.Descripcion = input.Descripcion
	recetaAActualizar.Foto = input.Foto
	recetaAActualizar.CategoriaID = input.CategoriaID
	// Categoria (el struct) se actualizará en la BD a través de CategoriaID
	// y se cargará con Preload si se consulta de nuevo.

	// 5. Llamar al repositorio para actualizar
	if err := s.recetaRepo.Update(ctx, recetaAActualizar); err != nil {
		// s.logger.Error("Error en repo.Update Receta", zap.Uint("id", id), zap.Error(err))
		if errors.Is(err, repository.ErrRecordNotFound) { // Si se borró justo antes
            return nil, ErrRecetaNotFound
        }
		return nil, fmt.Errorf("servicio recetas: error al actualizar: %w", err)
	}

	log.Printf("Servicio: Receta ID %d actualizada a nombre '%s'\n", recetaAActualizar.ID, recetaAActualizar.Nombre)
	// Devolver la receta actualizada (podría ser el mismo puntero o uno recargado)
	return recetaAActualizar, nil
}

// Delete elimina una receta por su ID.
func (s *recetaService) Delete(ctx context.Context, id uint) error {
	// 1. Verificar si existe antes de borrar (buena práctica)
	_, err := s.recetaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrRecetaNotFound
		}
		// s.logger.Error("Error buscando receta para Delete", zap.Uint("id", id), zap.Error(err))
		return fmt.Errorf("servicio recetas: error buscando para delete %d: %w", id, err)
	}

	// 2. Llamar al repositorio para eliminar
	if err := s.recetaRepo.Delete(ctx, id); err != nil {
		// s.logger.Error("Error en repo.Delete Receta", zap.Uint("id", id), zap.Error(err))
        if errors.Is(err, repository.ErrRecordNotFound) { // Si se borró justo antes
            return ErrRecetaNotFound
        }
		return fmt.Errorf("servicio recetas: error al eliminar: %w", err)
	}

	log.Printf("Servicio: Receta ID %d eliminada.\n", id)
	return nil
}

// FindByCategoriaID encuentra todas las recetas de una categoría específica.
func (s *recetaService) FindByCategoriaID(ctx context.Context, categoriaID uint) ([]Receta, error) {
	// 1. Validar que la CategoriaID exista (opcional, pero bueno para consistencia)
	_, err := s.categoriaSvc.GetByID(ctx, categoriaID)
	if err != nil {
		if errors.Is(err, categorias.ErrCategoriaNotFound) {
			// Devolver un slice vacío si la categoría no existe, en lugar de un error,
			// podría ser una decisión de diseño (el cliente pidió recetas de una categoría que no existe).
			// O devolver el error para ser más estricto. Por ahora, devolvemos error.
			return nil, fmt.Errorf("%w: la categoría ID %d no existe", ErrRecetaSinCategoria, categoriaID)
		}
		// s.logger.Error("Error validando CategoriaID en FindByCategoriaID", zap.Uint("categoriaID", categoriaID), zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error validando categoría %d: %w", categoriaID, err)
	}

	recs, err := s.recetaRepo.FindByCategoriaID(ctx, categoriaID)
	if err != nil {
		// s.logger.Error("Error en servicio FindByCategoriaID", zap.Uint("categoriaID", categoriaID), zap.Error(err))
		return nil, fmt.Errorf("servicio recetas: error buscando por categoriaID %d: %w", categoriaID, err)
	}
	return recs, nil
}