// backend/contactos/repository_gorm.go
// Funcionalidad: Implementación GORM de ContactoRepository.
// Capa: Repositorio (Implementación de Persistencia).
package contactos

import (
	"context"
	"errors"
	"fmt"
	"time"
	"gorm.io/gorm"
)

type gormContactoRepository struct {
	db *gorm.DB
}

// NewContactoRepository crea una instancia de la implementación GORM de ContactoRepository.
func NewContactoRepository(db *gorm.DB) ContactoRepository {
	return &gormContactoRepository{db: db}
}

func (r *gormContactoRepository) Create(ctx context.Context, contacto *ContactoForm) error {
	model := FromContactoFormDomain(contacto)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("repo gorm contactos: error creando contacto: %w", err)
	}
	// Actualizar el objeto de dominio con el ID generado y timestamps
	contacto.ID = model.ID
	contacto.CreatedAt = model.CreatedAt
	contacto.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *gormContactoRepository) GetByID(ctx context.Context, id uint) (*ContactoForm, error) {
	var model ContactoModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContactoRepoNotFound // Usar error definido en este paquete
		}
		return nil, fmt.Errorf("repo gorm contactos: error obteniendo por id %d: %w", id, err)
	}
	return model.ToDomain(), nil
}

func (r *gormContactoRepository) GetAll(ctx context.Context) ([]ContactoForm, error) {
	var models []ContactoModel
	// Ordenar por fecha de contacto descendente para ver los más nuevos primero
	if err := r.db.WithContext(ctx).Order("fecha_contacto desc").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("repo gorm contactos: error obteniendo todos: %w", err)
	}
	return ContactoModelsToDomains(models), nil
}

func (r *gormContactoRepository) MarkAsRead(ctx context.Context, id uint) error {
	// Actualizar solo el campo 'leido' y 'updated_at'
	result := r.db.WithContext(ctx).Model(&ContactoModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"leido":      true,
		"updated_at": time.Now(), // Forzar actualización de UpdatedAt
	})
	if result.Error != nil {
		return fmt.Errorf("repo gorm contactos: error marcando como leído id %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrContactoRepoNotFound // No se encontró para actualizar
	}
	return nil
}