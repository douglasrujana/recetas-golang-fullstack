package handler

import (
	"backend/internal/database"
	"backend/internal/handler/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

// GET: Select all categories
func CategoriasLista(c *gin.Context) {
	datos := modelos.Categorias{} // Inicializa un slice de categorías
	// Se utiliza un slice para almacenar los resultados de la consulta a la base de datos
	database.Database.Order("id desc").Find(&datos) // Realiza una consulta a la base de datos para obtener todas las categorías
	if len(datos) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No categories found"})
	} else {
		c.JSON(http.StatusOK, datos)
	}
}

// GET: ById
func CategoriasListaId(c *gin.Context) {
	id := c.Param("id")          // Obtiene el ID de la categoría desde la ruta
	datos := modelos.Categoria{} // Inicializa una categoría
	// Realiza una consulta a la base de datos para obtener una categoría por ID
	if err := database.Database.First(&datos, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":  "error",
			"mensaje": "No se encontró la categoría con ID: " + id,
			"error":   err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, datos)
	}
}

// POST
func CategoriasPost(c *gin.Context) {
	var body dto.CategoriaDTO // Inicializa un DTO de categoría

	// BindJSON vincula el cuerpo de la solicitud JSON a la estructura de datos
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "Error",
			"mensaje": "Error al procesar la solicitud",
			"error":   err.Error(),
		})
		return
	}

	//Validar que no exita el nombre de la categoria
	var existente modelos.Categoria
	if err := database.Database.Where("nombre = ?", body.Nombre).First(&existente).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "Error",
			"mensaje": "Ya existe una categoría con ese nombre",
		})
		return
	}

	//Crear un nuevo objeto de categoría a partir del DTO
	datos := modelos.Categoria{
		Nombre: body.Nombre,
		Slug:   slug.Make(body.Nombre), // Genera un slug a partir del nombre
	}
	// Inserta la nueva categoría en la base de datos
	if err := database.Database.Create(&datos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"estado":  "Error",
			"mensaje": "Error al crear la categoría",
			"error":   err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"estado":  "ok",
		"mensaje": "Categoría creada correctamente",
	})
}

// PUT Categorias
func CategoriasPut(c *gin.Context) {
	id := c.Param("id")       // Obtiene el ID de la categoría desde la ruta
	var body dto.CategoriaDTO // Inicializa un DTO de categoría

	// BindJSON vincula el cuerpo de la solicitud JSON a la estructura de datos
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"estado":  "Error",
			"mensaje": "Error al procesar la solicitud",
			"error":   err.Error(),
		})
		return
	}

	datos := modelos.Categoria{} // Inicializa una categoría
	// Realiza una consulta a la base de datos para obtener una categoría por ID
	if err := database.Database.First(&datos, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":  "error",
			"mensaje": "No se encontró la categoría con ID: " + id,
			"error":   err.Error(),
		})
		return
	}

	datos.Nombre = body.Nombre          // Actualiza el nombre de la categoría
	datos.Slug = slug.Make(body.Nombre) // Genera un nuevo slug a partir del nuevo nombre

	if err := database.Database.Save(&datos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"estado":  "Error",
			"mensaje": "Error al actualizar la categoría",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Categoría actualizada correctamente",
	})
}

// DELETE Categorias
func CategoriasDelete(c *gin.Context) {
	id := c.Param("id")          // Obtiene el ID de la categoría desde la ruta
	datos := modelos.Categoria{} // Inicializa una categoría
	// Realiza una consulta a la base de datos para obtener una categoría por ID
	if err := database.Database.First(&datos, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"estado":  "error",
			"mensaje": "No se encontró la categoría con ID: " + id,
			"error":   err.Error(),
		})
		return
	}

	if err := database.Database.Delete(&datos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"estado":  "Error",
			"mensaje": "Error al eliminar la categoría",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"estado":  "ok",
		"mensaje": "Categoría eliminada correctamente",
	})
}
