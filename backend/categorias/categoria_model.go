// Archivo: backend/categorias/categoria_model.go
// Funcionalidad: Modelo de dominio para Categoría
// Capa: Dominio / Lógica de negocio

// Descripción:
// Define la estructura base de una Categoría como entidad del dominio.
// No contiene tags de frameworks (GORM, JSON, etc). Representa el modelo de negocio, no la persistencia.
//
// Uso:
// - Usado por servicios del dominio para gestionar categorías.
// - Convertido a entidades GORM para almacenamiento, o a DTOs para respuestas HTTP.
//
// Ciclo de Vida:
// [✔] Definido según requerimientos del negocio turístico.
// [✔] Usado internamente en lógica de negocio y reglas de validación.
// [✘] No usado directamente por la capa de base de datos ni exposición HTTP.
//
// Responsabilidades:
// - Representar el concepto de Categoría.
// - Permitir validaciones como unicidad de slug o obligatoriedad del nombre.
//
// Reglas de Negocio:
// - Nombre debe ser obligatorio y legible.
// - Slug debe ser único por idioma o contexto.
//
// Posibles extensiones futuras:
// - Asociar categorías jerárquicamente (padre/hijo).
// - Agregar campo Descripción o multilenguaje.
//
// Alias para listas de categorías, útil en casos donde se definen métodos sobre conjuntos de categorías.
// type Categorias []Categoria

package categorias

import "time" // Solo si necesitas CreatedAt/UpdatedAt aquí

// Categoria representa la entidad de negocio pura para una categoría.
// [✅ BUENA PRÁCTICA] Sin tags de Gorm. Representa el negocio, no la base de datos.
type Categoria struct {
	ID     uint   // Identificador único de la categoria
	Nombre string // Nombre visible y legible de la categoria
	Slug   string // Slug URL-amigable único
	CreatedAt time.Time // Fecha de creación
	UpdatedAt time.Time // Última fecha de modificación
}

// Type alias para slices, si lo prefieres (opcional)
// type Categorias []Categoria
