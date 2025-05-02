package dto
// dto: Para poder recibir los datos del json en el body de la peticion
// y poder enviarlos como respuesta en el json de la peticion
// Patron de diseño estrucutural DTO (Data Transfer Object)
// Contenedor simple de datos
// DTO (Data Transfer Object) para la ruta de ejemplo
// Este objeto se utiliza para transferir datos entre el cliente y el servidor
// en la ruta de ejemplo.
// Se utiliza para recibir datos en el cuerpo de la solicitud (request body) y
// para enviar datos en la respuesta (response).

type EjemploDTO struct {
	Correo string `json:"correo"`
	Password string `json:"password"`
}

type CategoriaDTO struct {
	Nombre string `json:"nombre" binding:"required"`
}

type RecetaResponse struct {
	ID uint 		 	`json:"id"`
	Nombre string 	 	`json:"nombre"`
	Slug string 	 	`json:"slug"`
	CategoriaID uint 	`json:"categoria_id"`
	Categoria string 	`json:"categoria"`
	Tiempo string    	`json:"tiempo"`
	Descripcion string  `json:"descripcion"`
	Foto string         `json:"foto"`
	Fecha string        `json:"fecha"`
}

type RecetaResponses []RecetaResponse

type RecetaDTO struct {
	Nombre string `json:"nombre"`
	Slug string `json:"slug"`
	CategoriaId uint `json:"categoria_id"`
	Tiempo string `json:"tiempo"`
	Descripcion string `json:"descripcion"`
	Foto string `json:"foto"`
} 