// backend/internal/repository/mysql/categoria_model.go
package mysql

import (
	"backend/internal/domain" // Importa el modelo de dominio PURO
	"time"

	"gorm.io/gorm" // Importa GORM
)

// CategoriaModel representa la tabla 'categorias' en la BD y usa GORM.
// Define la estructura de la tabla y las reglas de mapeo del ORM.
type CategoriaModel struct {
	// gorm.Model // Alternativa: Embeber si usas ID, CreatedAt, UpdatedAt, DeletedAt estándar GORM
	ID        uint           `gorm:"primaryKey"`
	Nombre    string         `gorm:"type:varchar(100);uniqueIndex:uk_categorias_nombre,priority:1"` // uniqueIndex con nombre
	Slug      string         `gorm:"type:varchar(120);uniqueIndex:uk_categorias_slug,priority:1"` // uniqueIndex con nombre
	CreatedAt time.Time      // GORM maneja esto automáticamente
	UpdatedAt time.Time      // GORM maneja esto automáticamente
	DeletedAt gorm.DeletedAt `gorm:"index"` // Para soft delete (si lo activas en GORM o manualmente)
}

// TableName le dice explícitamente a GORM cómo se llama la tabla física.
// Esto anula la convención de pluralización automática de GORM.
func (CategoriaModel) TableName() string {
	// ¡ASEGÚRATE que este nombre ("categorias") coincida con el
	// nombre de la tabla que AutoMigrate creará!
	return "categorias"
}

// --- Funciones de Mapeo ---

// ToDomain convierte el modelo de persistencia (GORM) al modelo de dominio (puro).
func (m *CategoriaModel) ToDomain() *domain.Categoria {
	if m == nil {
		return nil
	}
	return &domain.Categoria{
		ID:        m.ID,
		Nombre:    m.Nombre,
		Slug:      m.Slug,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		// No mapeamos DeletedAt al dominio a menos que sea explícitamente necesario
	}
}

// FromDomain convierte un modelo de dominio (puro) al modelo de persistencia (GORM).
func FromDomain(d *domain.Categoria) *CategoriaModel {
	if d == nil {
		return nil
	}
	return &CategoriaModel{
		ID:     d.ID, // Pasar el ID es crucial para Updates/Saves
		Nombre: d.Nombre,
		Slug:   d.Slug,
		// GORM maneja CreatedAt/UpdatedAt, no necesitamos mapearlos desde el dominio
		// a menos que la lógica de negocio los establezca explícitamente.
	}
}

// --- Funciones de Mapeo para Slices (Helpers) ---

// ModelsToDomains convierte un slice de modelos GORM a un slice de modelos de dominio.
func ModelsToDomains(models []CategoriaModel) []domain.Categoria {
	if models == nil {
		return []domain.Categoria{} // Devolver slice vacío es más seguro que nil
	}
	domainCategorias := make([]domain.Categoria, 0, len(models)) // Pre-alocar capacidad
	for _, model := range models {
		// Convertir cada modelo y añadir al slice de dominio
		if domainModel := model.ToDomain(); domainModel != nil { // Chequeo extra por si ToDomain devolviera nil
			domainCategorias = append(domainCategorias, *domainModel)
		}
	}
	return domainCategorias
}

// DomainsToModels convierte un slice de modelos de dominio a un slice de modelos GORM.
// (Menos común necesitar esto, usualmente se opera de a uno para Create/Update)
// func DomainsToModels(domainCategorias []domain.Categoria) []CategoriaModel {
// 	if domainCategorias == nil {
// 		return []CategoriaModel{}
// 	}
// 	models := make([]CategoriaModel, 0, len(domainCategorias))
// 	for _, dom := range domainCategorias {
//      if model := FromDomain(&dom); model != nil {
//   		models = append(models, *model)
//      }
// 	}
// 	return models
// }