// backend/recetas/receta_repository_gorm.go
// Implementación con GORM de RecetaRepository.

// Este archivo define la implementación con GORM de RecetaRepository.
// Pertenece al paquete de la característica 'recetas'

package recetas // Paquete de la característica 'recetas'

import (
	// No necesita importar "backend/recetas" para RecetaModel o Receta (dominio)
	// ya que están en el mismo paquete (en diferentes archivos).
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	//"gorm.io/gorm/clause" // Para Preload anidado si es necesario
	"backend/shared/repository"
)

type gormRecetaRepository struct { // no exportado
	db *gorm.DB
}

// NewRecetaRepository crea una instancia de RecetaRepository (implementación GORM).
func NewRecetaRepository(db *gorm.DB) RecetaRepository { // Devuelve la interfaz RecetaRepository
	return &gormRecetaRepository{db: db}
}

// --- Implementación de Métodos ---

func (r *gormRecetaRepository) GetAll(ctx context.Context) ([]Receta, error) {
	var models []RecetaModel
	// ¡IMPORTANTE! Preload("Categoria") para cargar la relación.
	// GORM usará el struct CategoriaModel (del paquete 'categorias') definido en RecetaModel.
	if err := r.db.WithContext(ctx).Preload("Categoria").Order("id desc").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("repo gorm recetas: getall: %w", err)
	}
	return RecetaModelsToDomains(models), nil // Usar el mapeador de RecetaModel
}

// GetByID recupera una receta por su ID.
func (r *gormRecetaRepository) GetByID(ctx context.Context, id uint) (*Receta, error) {
	var model RecetaModel
	// ¡IMPORTANTE! Preload("Categoria")
	if err := r.db.WithContext(ctx).Preload("Categoria").First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Podríamos devolver nuestro propio ErrRecordNotFound del paquete 'recetas' si lo definimos
			return nil, repository.ErrRecordNotFound // Asumiendo que ErrRecordNotFound está definido en este paquete (en repository.go)
		}
		return nil, fmt.Errorf("repo gorm recetas: getbyid %d: %w", id, err)
	}
	return model.ToDomain(), nil
}

// GetBySlug recupera una receta por su slug.
func (r *gormRecetaRepository) GetBySlug(ctx context.Context, slug string) (*Receta, error) {
	var model RecetaModel
	// ¡IMPORTANTE! Preload("Categoria")
	if err := r.db.WithContext(ctx).Preload("Categoria").Where("slug = ?", slug).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("repo gorm recetas: getbyslug %s: %w", slug, err)
	}
	return model.ToDomain(), nil
}

// Create crea una receta en la base de datos.
func (r *gormRecetaRepository) Create(ctx context.Context, receta *Receta) error {
	model := FromRecetaDomain(receta) // Mapear dominio a modelo GORM
	// GORM se encargará de la CategoriaID si está presente en el modelo.
	// Si quisiéramos asociar un CategoriaModel completo, la lógica sería diferente.
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("repo gorm recetas: create: %w", err)
	}
	// Actualizar el ID en el objeto de dominio original
	receta.ID = model.ID
	// GORM también actualiza CreatedAt y UpdatedAt en el modelo.
	// Si el dominio los necesita actualizados inmediatamente:
	receta.CreatedAt = model.CreatedAt
	receta.UpdatedAt = model.UpdatedAt
	return nil
}

// Update actualiza una receta existente en la base de datos.
func (r *gormRecetaRepository) Update(ctx context.Context, receta *Receta) error {
	model := FromRecetaDomain(receta)
	// Para Update, es crucial que el modelo tenga el ID correcto.
	// GORM .Updates solo actualiza campos no cero, o usa .Select para especificar.
	// Si quieres actualizar CategoriaID, asegúrate que esté en el modelo.
	result := r.db.WithContext(ctx).Model(&RecetaModel{}).Where("id = ?", model.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("repo gorm recetas: update %d: %w", model.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrRecordNotFound // ID no encontrado para actualizar
	}
	// Podrías querer recargar el modelo para obtener el UpdatedAt actualizado por la BD
	// y luego actualizar el objeto de dominio 'receta' si es necesario.
	// Por ahora, asumimos que la actualización fue exitosa.
	return nil
}

// Delete elimina una receta de la base de datos.
func (r *gormRecetaRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&RecetaModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("repo gorm recetas: delete %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrRecordNotFound // ID no encontrado para borrar
	}
	return nil
}

// FindByCategoriaID encuentra todas las recetas de una categoría específica.
func (r *gormRecetaRepository) FindByCategoriaID(ctx context.Context, categoriaID uint) ([]Receta, error) {
	var models []RecetaModel
	// Preload Categoria también aquí para consistencia
	if err := r.db.WithContext(ctx).Preload("Categoria").Where("categoria_id = ?", categoriaID).Order("id desc").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("repo gorm recetas: findbycategoriaid %d: %w", categoriaID, err)
	}
	return RecetaModelsToDomains(models), nil
}