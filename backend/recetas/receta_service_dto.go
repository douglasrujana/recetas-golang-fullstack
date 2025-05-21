// backend/recetas/service_dto.go

// Este archivo define los DTOs (Data Transfer Objects) específicos para la entidad Receta.
// RecetaInputDTO define la estructura para crear o actualizar una receta
// a nivel de la capa de servicio.

package recetas // Paquete de la característica 'recetas'

type RecetaInputDTO struct {
	Nombre            string
	CategoriaID       uint
	TiempoPreparacion string
	Descripcion       string
	Foto              string // Nombre del archivo o URL (el servicio podría procesar esto)
	// Ingredientes      []uint // Ejemplo: Podríamos recibir IDs de ingredientes aquí a futuro
}

// Podrías tener otros DTOs de servicio si fueran necesarios, por ejemplo, para filtros:
// type RecetaFiltroDTO struct {
//     NombreContiene *string
//     CategoriaID    *uint
//     // ...
// }