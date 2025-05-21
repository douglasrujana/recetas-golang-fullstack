// backend/recetas/receta_model_gorm.go

// Este archivo define el modelo de persistencia para una receta.
// Utiliza GORM para la definición de la tabla y el mapeo de campos.
// No contiene tags de Gorm. Representa el negocio, no la base de datos.

package recetas // El paquete es 'recetas'

import (
	// Necesitamos importar el paquete 'categorias' para referenciar 'categorias.CategoriaModel'
	// y el 'categorias.Categoria' (struct de dominio) en los mapeadores.
	"backend/categorias"
	"time"

	"gorm.io/gorm"
)

// RecetaModel representa la tabla 'recetas' en la BD y usa GORM.
// Define la estructura de la tabla y las reglas de mapeo del ORM.
type RecetaModel struct {
	ID                uint           `gorm:"primaryKey"`
	Nombre            string         `gorm:"type:varchar(150);not null"`
	Slug              string         `gorm:"type:varchar(180);uniqueIndex:uk_recetas_slug"` // Asumo que slug debe ser único
	TiempoPreparacion string         `gorm:"type:varchar(50)"`
	Descripcion       string         `gorm:"type:text"`
	Foto              string         `gorm:"type:varchar(100);default:null"` // Permitir NULL si la foto es opcional
	CreatedAt         time.Time      // GORM maneja esto
	UpdatedAt         time.Time      // GORM maneja esto
	DeletedAt         gorm.DeletedAt `gorm:"index"` // Para soft delete (opcional)

	// --- Relación con Categoria (Belongs To) ---
	CategoriaID       uint                    // Columna de clave foránea explícita
	// Usamos el tipo CategoriaModel del paquete 'categorias' para la relación.
	// GORM usará el TableName() de categorias.CategoriaModel para saber a qué tabla unirse.
	Categoria         categorias.CategoriaModel `gorm:"foreignKey:CategoriaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	                                          // constraint: opcional, define comportamiento de FK
	// --- Relación con Ingredientes (Many To Many) - SE AÑADIRÁ LUEGO ---
	// Ingredientes      []categorias.IngredienteModel `gorm:"many2many:receta_ingredientes;"`
	// (Asumiendo que IngredienteModel vivirá en el paquete 'categorias' o 'ingredientes')
}

// TableName especifica el nombre de la tabla física en la base de datos.
func (RecetaModel) TableName() string {
	return "recetas" // Nombre de la tabla en plural, siguiendo convenciones o tu preferencia
}

// --- Funciones de Mapeo ---

// ToDomain convierte el modelo de persistencia (RecetaModel) al modelo de dominio (Receta de este paquete).
func (m *RecetaModel) ToDomain() *Receta {
	if m == nil {
		return nil
	}

	var domainCategoria *categorias.Categoria // Tipo de dominio del paquete 'categorias'
	// Solo mapear Categoria si m.Categoria (el CategoriaModel anidado) fue cargada
	// por GORM (usualmente con Preload) y tiene un ID válido.
	if m.Categoria.ID != 0 { // Asumiendo que CategoriaModel tiene ID como primaryKey
		// Llamamos al mapeador ToDomain del CategoriaModel del paquete 'categorias'
		domainCategoria = m.Categoria.ToDomain()
	}

	return &Receta{ // Crear instancia de Receta (dominio de este paquete)
		ID:                m.ID,
		Nombre:            m.Nombre,
		Slug:              m.Slug,
		TiempoPreparacion: m.TiempoPreparacion,
		Descripcion:       m.Descripcion,
		Foto:              m.Foto,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		CategoriaID:       m.CategoriaID,
		Categoria:         domainCategoria, // Asignar el *categorias.Categoria (dominio) mapeado
		// Ingredientes:    // Se mapearán cuando implementemos la relación con ingredientes
	}
}

// FromRecetaDomain convierte un modelo de dominio (Receta de este paquete) al modelo de persistencia (RecetaModel).
func FromRecetaDomain(d *Receta) *RecetaModel { // Recibe *Receta (dominio de este paquete)
	if d == nil {
		return nil
	}
	// La Categoria anidada (d.Categoria) no se mapea directamente a RecetaModel.Categoria.
	// GORM maneja la asociación principalmente a través de RecetaModel.CategoriaID.
	// Si quisiéramos crear/actualizar la categoría asociada en la misma operación,
	// esa lógica compleja usualmente residiría en la capa de servicio.
	return &RecetaModel{
		ID:                d.ID,
		Nombre:            d.Nombre,
		Slug:              d.Slug,
		TiempoPreparacion: d.TiempoPreparacion,
		Descripcion:       d.Descripcion,
		Foto:              d.Foto,
		CategoriaID:       d.CategoriaID, // Muy importante para la FK
		// Categoria (el struct categorias.CategoriaModel) no se asigna desde d.Categoria (el struct domain) aquí.
		// GORM la asociará si CategoriaID está presente y CategoriaModel ya existe con ese ID.
	}
}

// RecetaModelsToDomains convierte un slice de RecetaModel a un slice de Receta (dominio de este paquete).
func RecetaModelsToDomains(models []RecetaModel) []Receta {
	if models == nil {
		return []Receta{} // Devolver slice vacío es más seguro
	}
	domainRecetas := make([]Receta, 0, len(models))
	for _, model := range models {
		if domainModel := model.ToDomain(); domainModel != nil {
			domainRecetas = append(domainRecetas, *domainModel)
		}
	}
	return domainRecetas
}