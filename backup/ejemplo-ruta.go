package rutas // rutas de tipo handler

import (
	"backend/dto" // Importa el paquete dto para usar el DTO EjemploDTO
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// EjemploRuta es un manejador de ruta de ejemplo que responde con un mensaje JSON.

func EjemploRuta(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"estado":  "Ok",
		"mensaje": "Sapin!",
	})
}

func EjemploGet(c *gin.Context) {
	// Simulando una respuesta de ejemplo
	response := map[string]string{
		"message": "Método Get",
	}
	c.JSON(http.StatusOK, response)
}

func EjemploGetId(c *gin.Context) {
	// Muestra el id en la consola
	fmt.Print("Id=", c.Param("id"), "\n")
	// Muestra el tipo de dato del id en la consola
	fmt.Println(reflect.TypeOf(c.Param("id")))
	// Simulando una respuesta de ejemplo
	response := map[string]string{
		"message": "Método GeId  | " + c.Param("id"),
	}
	c.JSON(http.StatusOK, response)
}

func EjemploPost(c *gin.Context) {
	// Simulando una respuesta de ejemplo
	response := map[string]string{
		"message": "Método Post",
	}
	c.JSON(http.StatusCreated, response)
}

// EjemploPostDTO es un manejador de ruta de ejemplo que recibe un DTO y responde con un mensaje JSON.
func EjemploPostDTO(c *gin.Context) {
	// DTO (Data Transfer Object) para la ruta de ejemplo
	// Este objeto se utiliza para transferir datos entre el cliente y el servidor
	var body dto.EjemploDTO
	// Bind JSON a la estructura DTO
	// El método ShouldBindJSON se utiliza para enlazar el cuerpo de la solicitud JSON a la estructura DTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulando una respuesta de ejemplo
	response := map[string]string{
		"message":  "Método Post DTO",
		"correo":   body.Correo,
		"password": body.Password,
	}
	c.JSON(http.StatusCreated, response)
}

func EjemploPut(c *gin.Context) {
	// Simulando una respuesta de ejemplo
	response := map[string]string{
		"message": "Método Put: | " + c.Param("id"),
	}
	c.JSON(http.StatusOK, response)
}

func EjemploDelete(c *gin.Context) {
	// Simulando una respuesta de ejemplo
	response := map[string]string{
		"message": "Métod delete: | " + c.Param("id"),
	}
	c.JSON(http.StatusOK, response)
}

func EjemploQueryString(c *gin.Context) {
	response := map[string]string{
		"message": "Método: GET/QueryString: | " + c.Query("id") + "| slug: " + c.Query("slug"),
	}
	c.JSON(http.StatusOK, response)
}

// EjemploUpload maneja la carga de archivos desde una petición HTTP POST.
// Recibe un archivo en el campo "file", valida su tamaño, tipo y lo renombra
// con una marca de tiempo antes de guardarlo en el servidor.
func EjemploUpload(c *gin.Context) {
	// Obtener el archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al subir el archivo"})
		return
	}
	// Validar que el archivo no esté vacío
	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo está vacío"})
		return
	}
	// Validar tamaño máximo permitido (2MB)
	const maxSize = 2 << 20 // 2MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo excede el tamaño máximo permitido (2MB)"})
		return
	}
	// Validar extensión permitida (.jpg, .png)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de archivo no permitido. Solo .jpg y .png"})
		return
	}
	// Renombrar archivo con marca de tiempo única (ej: 20250411_143523.jpg)
	newFilename := fmt.Sprintf("%s%s", time.Now().Format("20060102_150405"), ext)
	// Ruta destino
	savePath := filepath.Join("./uploads", newFilename)
	// Guardar el archivo
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el archivo"})
		return
	}
	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"message":  "Archivo subido con éxito",
		"filename": newFilename,
	})
}
