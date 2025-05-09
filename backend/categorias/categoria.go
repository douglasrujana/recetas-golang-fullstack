// backend/categorias/categoria.go
package categorias

import "time" // Solo si necesitas CreatedAt/UpdatedAt aquí

// Categoria representa la entidad de negocio pura para una categoría.
// [✅ BUENA PRÁCTICA] Sin tags de Gorm. Representa el negocio, no la base de datos.
type Categoria struct {
	ID     uint   // Identificador único
	Nombre string // Nombre de la categoría
	Slug   string // Slug para URLs amigables
	// Puedes añadir CreatedAt/UpdatedAt si son relevantes para la lógica de negocio
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Type alias para slices, si lo prefieres (opcional)
// type Categorias []Categoria
