package categorias

// --- DTOs específicos para la entrada del SERVICIO ---
// (Pueden ser iguales a los del handler, pero definirlos aquí desacopla)

// CategoriaInputDTO define la estructura para crear o actualizar una categoría a nivel de servicio.
type CategoriaInputDTO struct {
	Nombre string
}
