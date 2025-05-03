// backend/internal/repository/mysql/categoria_repository.go
package mysql // Pertenece al adaptador MySQL

import (
	"backend/internal/domain"
	"backend/internal/repository" // Importa la interfaz y errores del repo
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm" // Dependencia de infraestructura (GORM)
)

// categoriaRepository implementa repository.CategoriaRepository usando GORM.
type categoriaRepository struct {
	db *gorm.DB // Dependencia inyectada
}

// NewCategoriaRepository es la Factory Function para crear una instancia.
// Devuelve la INTERFAZ, no la implementación concreta.
func NewCategoriaRepository(db *gorm.DB) repository.CategoriaRepository {
	return &categoriaRepository{db: db}
}

// --- Implementación de los métodos de la interfaz ---

// GetAll implementa repository.CategoriaRepository.GetAll.
func (r *categoriaRepository) GetAll(ctx context.Context) ([]domain.Categoria, error) {
	var categorias []domain.Categoria
	// GORM mapeará directamente si los campos coinciden o usas tags `gorm:"column:..."`
	// El contexto se pasa para posible cancelación/timeout
	if err := r.db.WithContext(ctx).Order("id desc").Find(&categorias).Error; err != nil {
		// Envolver el error para contexto
		return nil, fmt.Errorf("repositorio mysql: error al obtener todas las categorias: %w", err)
	}
	return categorias, nil
}

// GetByID implementa repository.CategoriaRepository.GetByID.
func (r *categoriaRepository) GetByID(ctx context.Context, id uint) (*domain.Categoria, error) {
	var categoria domain.Categoria
	if err := r.db.WithContext(ctx).First(&categoria, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Traducir error de GORM a error de nuestro paquete repository
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("repositorio mysql: error al obtener categoria por id %d: %w", id, err)
	}
	return &categoria, nil
}

// GetBySlug implementa repository.CategoriaRepository.GetBySlug.
func (r *categoriaRepository) GetBySlug(ctx context.Context, slug string) (*domain.Categoria, error) {
	var categoria domain.Categoria
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&categoria).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("repositorio mysql: error al obtener categoria por slug %s: %w", slug, err)
	}
	return &categoria, nil
}
 
// GetByNombre implementa repository.CategoriaRepository.GetByNombre.
func (r *categoriaRepository) GetByNombre(ctx context.Context, nombre string) (*domain.Categoria, error) {
	var categoria domain.Categoria
	if err := r.db.WithContext(ctx).Where("nombre = ?", nombre).First(&categoria).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Es importante devolver ErrRecordNotFound aquí para que el servicio sepa que no existe
			return nil, repository.ErrRecordNotFound 
		}
		return nil, fmt.Errorf("repositorio mysql: error al obtener categoria por nombre %s: %w", nombre, err)
	}
	// Si encuentra una, la devuelve (sin error)
	return &categoria, nil
}

func (r *categoriaRepository) Create(ctx context.Context, categoria *domain.Categoria) error {
	// GORM automáticamente maneja el ID y timestamps si están definidos en el struct y DB
	if err := r.db.WithContext(ctx).Create(categoria).Error; err != nil {
		// Podríamos intentar detectar errores de duplicados aquí si la BD lo permite
		// if strings.Contains(err.Error(), "Duplicate entry") { // Ejemplo muy básico
		// 	 return repository.ErrDuplicateRecord
		// }
		return fmt.Errorf("repositorio mysql: error al crear categoria: %w", err)
	}
	// El ID ahora debería estar seteado en el objeto 'categoria' que pasamos por puntero
	return nil
}

// Update implementa repository.CategoriaRepository.Update.
func (r *categoriaRepository) Update(ctx context.Context, categoria *domain.Categoria) error {
	// Asume que 'categoria' ya tiene el ID correcto.
	// GORM `Save` actualiza todos los campos o solo los cambiados si se usa `Updates`.
	// Usar `Save` es más simple si el objeto viene completo.
	result := r.db.WithContext(ctx).Save(categoria) // Save actualiza si existe PK, o inserta si no (¡cuidado!)
	// Es mejor asegurar que actualizamos usando Updates o seleccionando primero.
	// Alternativa más segura:
	// result := r.db.WithContext(ctx).Model(&domain.Categoria{}).Where("id = ?", categoria.ID).Updates(categoria)

	if result.Error != nil {
		return fmt.Errorf("repositorio mysql: error al actualizar categoria id %d: %w", categoria.ID, result.Error)
	}
	if result.RowsAffected == 0 {
        // Esto podría significar que el ID no existía para actualizar
        return repository.ErrRecordNotFound // O un error específico "update failed"
    }
	return nil
}

// Delete
func (r *categoriaRepository) Delete(ctx context.Context, id uint) error {
	// GORM permite borrar pasando un objeto con ID o directamente el ID.
	// Borrar directamente por ID es más eficiente si no necesitas el objeto.
	result := r.db.WithContext(ctx).Delete(&domain.Categoria{}, id)
	if result.Error != nil {
		// Manejar posibles errores (ej: foreign key constraint)
		// if errors.Is(result.Error, gormForeignKeyConstraintError) { // Pseudo-código
        //     return repository.ErrForeignKeyViolation
        // }
		return fmt.Errorf("repositorio mysql: error al eliminar categoria id %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
        // Si no se afectaron filas, significa que el ID no existía.
        return repository.ErrRecordNotFound
    }
	return nil
}