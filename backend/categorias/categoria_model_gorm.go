// backend/categorias/categoria_model_gorm.go

// Este archivo define el modelo de persistencia para una categoría.
// Utiliza GORM para la definición de la tabla y el mapeo de campos.
// No contiene tags de Gorm. Representa el negocio, no la base de datos.

package categorias // Asegúrate que el paquete sea 'categorias'

import (
	// YA NO necesitas: "backend/internal/domain" si Categoria (dominio) está en este paquete
	"time"
	"gorm.io/gorm"
)

// CategoriaModel representa la tabla 'categorias' en la BD y usa GORM.
type CategoriaModel struct {
	ID        uint           `gorm:"primaryKey"`
	Nombre    string         `gorm:"type:varchar(100);uniqueIndex:uk_categorias_nombre,priority:1"`
	Slug      string         `gorm:"type:varchar(120);uniqueIndex:uk_categorias_slug,priority:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CategoriaModel) TableName() string {
	return "categorias"
}

// --- Funciones de Mapeo ---

// ToDomain convierte el modelo de persistencia (GORM) al modelo de dominio (puro).
// Ahora Categoria (el tipo de dominio) está en el mismo paquete.
func (m *CategoriaModel) ToDomain() *Categoria { // <-- SIN 'domain.'
	if m == nil {
		return nil
	}
	return &Categoria{ // <-- SIN 'domain.'
		ID:        m.ID,
		Nombre:    m.Nombre,
		Slug:      m.Slug,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromDomain convierte un modelo de dominio (puro) al modelo de persistencia (GORM).
// Ahora Categoria (el tipo de dominio) está en el mismo paquete.
func FromDomain(d *Categoria) *CategoriaModel { // <-- SIN 'domain.'
	if d == nil {
		return nil
	}
	return &CategoriaModel{
		ID:     d.ID,
		Nombre: d.Nombre,
		Slug:   d.Slug,
	}
}

// ModelsToDomains convierte un slice de modelos GORM a un slice de modelos de dominio.
func ModelsToDomains(models []CategoriaModel) []Categoria { // <-- SIN 'domain.'
	if models == nil {
		return []Categoria{} // Devuelve slice vacío
	}
	domainCategorias := make([]Categoria, 0, len(models))
	for _, model := range models {
		if domainModel := model.ToDomain(); domainModel != nil {
			domainCategorias = append(domainCategorias, *domainModel)
		}
	}
	return domainCategorias
}