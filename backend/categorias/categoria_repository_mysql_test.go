// backend/categorias/categoria_repository_mysql_test.go (Adaptado para Modelos Separados)
// Test de integración para CategoriaRepository
package categorias

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"backend/shared/repository"
)

// Struct y New (sin cambios)
type categoriaRepository struct { db *gorm.DB }
func NewCategoriaRepository(db *gorm.DB) CategoriaRepository { return &categoriaRepository{db: db} }

func (r *categoriaRepository) GetAll(ctx context.Context) ([]Categoria, error) {
	var models []CategoriaModel // Slice del modelo GORM
	// Usar el modelo GORM en la consulta. TableName() se usará automáticamente.
	if err := r.db.WithContext(ctx).Order("id desc").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("repositorio mysql: error al obtener todas las categorias: %w", err)
	}
	// Mapear el resultado a un slice de dominio usando el helper
	return ModelsToDomains(models), nil
}

func (r *categoriaRepository) GetByID(ctx context.Context, id uint) (*Categoria, error) {
	var model CategoriaModel // Modelo GORM
	// Usar el modelo GORM en First. TableName() se usará.
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("repositorio mysql: error al obtener categoria por id %d: %w", id, err)
	}
	return model.ToDomain(), nil // Mapear a dominio antes de devolver
}

 func (r *categoriaRepository) GetBySlug(ctx context.Context, slug string) (*Categoria, error) {
	var model CategoriaModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("repositorio mysql: error al obtener categoria por slug %s: %w", slug, err)
	}
	return model.ToDomain(), nil
}

func (r *categoriaRepository) GetByNombre(ctx context.Context, nombre string) (*Categoria, error) {
	var model CategoriaModel
	if err := r.db.WithContext(ctx).Where("nombre = ?", nombre).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("repositorio mysql: error al obtener categoria por nombre %s: %w", nombre, err)
	}
	return model.ToDomain(), nil
}

func (r *categoriaRepository) Create(ctx context.Context, categoria *Categoria) error {
	model := FromDomain(categoria) // Mapear Dominio -> Modelo GORM
	// Crear usando el Modelo GORM
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		// Podríamos añadir chequeo de error UNIQUE aquí basado en el string del error MySQL
		return fmt.Errorf("repositorio mysql: error al crear categoria: %w", err)
	}
	// Actualizar el ID en el objeto de dominio original que nos pasaron
	categoria.ID = model.ID
	return nil
}

func (r *categoriaRepository) Update(ctx context.Context, categoria *Categoria) error {
	model := FromDomain(categoria) // Mapear Dominio -> Modelo GORM
	// Usar Updates con el Modelo GORM es más seguro que Save
	// Asegúrate que el modelo tenga el ID correcto
	result := r.db.WithContext(ctx).Model(&CategoriaModel{}).Where("id = ?", model.ID).Updates(model)
	if result.Error != nil {
		// Podríamos chequear error UNIQUE aquí también
		return fmt.Errorf("repositorio mysql: error al actualizar categoria id %d: %w", model.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrRecordNotFound // No encontró el ID para actualizar
	}
	// No es necesario actualizar el objeto dominio original aquí, GORM actualizó la BD
	return nil
}

func (r *categoriaRepository) Delete(ctx context.Context, id uint) error {
	// Borrar usando el Modelo GORM como referencia y el ID
	result := r.db.WithContext(ctx).Delete(&CategoriaModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("repositorio mysql: error al eliminar categoria id %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrRecordNotFound // No encontró el ID para borrar
	}
	return nil
}