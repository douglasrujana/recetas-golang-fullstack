//go:build ignore
// +build ignore

package handler

import (
	"path/filepath"
	"time"
	"backend/internal/handler/dto"
	"backend/internal/utils" // Importa el paquete de tiempo para formatear la fecha
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

//GET:
func RecetasListar(c *gin.Context) {
	// Inicializa un slice de recetas
	// Se utiliza un slice para almacenar los resultados de la consulta a la base de datos
	datos := modelos.Recetas{}

	// Realiza una consulta a la base de datos para obtener todas las recetas
	// Preload("categoria") carga la relación "categoria" para cada receta
	// Order("id desc") ordena los resultados por ID en orden descendente
	database.Database.
		Preload("Categoria").
		Order("id desc").
		Find(&datos) // Realiza una consulta a la base de datos para obtener todas las recetas

	var receta []dto.RecetaResponse
	for _, recetaItem := range datos {
		receta = append(receta, dto.RecetaResponse{
			ID:          recetaItem.ID,
			Nombre:      recetaItem.Nombre,
			Slug:        recetaItem.Slug,
			CategoriaID: recetaItem.CategoriaId,
			Categoria:   recetaItem.Categoria.Nombre,
			Tiempo:      recetaItem.Tiempo,
			Descripcion: recetaItem.Descripcion,
			Foto:        utils.ProcesarFotoURL(c, recetaItem.Foto),
			Fecha:       utils.FormatearFecha(recetaItem.Fecha),
		})
	}

	if len(datos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No recipes found"})
	} else {
		c.JSON(http.StatusOK, receta)
	}
}

//GET/id
func RecetasListarById(c *gin.Context) {
	id := c.Param("id")
	datos := modelos.Receta{}
	// SQL
	if err := database.Database.
		Preload("Categoria").
		First(&datos, id).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":  "error",
			"mensaje": "No se encontró la receta con ID: " + id,
			"error":   err.Error(),
		})
		return
	}
	// Construimos la respuesta con URL de imagen y fecha formateada
	respuesta := dto.RecetaResponse{
		ID:          datos.ID,
		Nombre:      datos.Nombre,
		Slug:        datos.Slug,
		CategoriaID: datos.CategoriaId,
		Categoria:   datos.Categoria.Nombre,
		Tiempo:      datos.Tiempo,
		Descripcion: datos.Descripcion,
		Foto:        utils.ProcesarFotoURL(c, datos.Foto),
		Fecha:       utils.FormatearFecha(datos.Fecha),
	}
	c.JSON(http.StatusOK, respuesta)
}

//POST:
func RecetasCrear(c *gin.Context) {
	nombre := c.PostForm("nombre")
	slug := c.PostForm("slug")
	tiempo := c.PostForm("tiempo")
	descripcion := c.PostForm("descripcion")
	categoriaIDStr := c.PostForm("categoria_id")

	// Validaciones
	if nombre == "" || slug == "" || tiempo == "" || descripcion == "" || categoriaIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos los campos son obligatorios"})
		return
	}

	categoriaID, err := strconv.Atoi(categoriaIDStr)
	if err != nil || categoriaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}

	// Imagen
	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagen no recibida correctamente"})
		return
	}

	//Function: save image
	rutaRelativa, err := utils.GuardarImagen(c, file, "recetas")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Realiza una consulta a la base de datos para obtener todas las recetas
	datos := modelos.Recetas{}
	database.Database.Preload("Categoria").Order("id desc").Find(&datos) 
	receta := modelos.Receta{
		Nombre:      nombre,
		Slug:        slug,
		CategoriaId: uint(categoriaID),
		Tiempo:      tiempo,
		Descripcion: descripcion,
		Foto:        utils.ProcesarFotoURL(c, filepath.Base(rutaRelativa)),
		Fecha:       time.Now(),
	}

	if err := database.Database.Create(&receta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar receta"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"estado": "ok",
		"data":   receta,
	})
}
