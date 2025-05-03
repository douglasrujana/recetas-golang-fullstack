// backend/internal/handler/categoria_handler.go
package handler

import (
	// --- Imports necesarios ---
	"backend/internal/domain"  // Para mapear errores de dominio
	"backend/internal/handler/dto" // DTOs para request y response HTTP
	"backend/internal/service" // La INTERFAZ del servicio que usaremos
	"errors"
	// "fmt" // Para formateo básico de errores si no usamos helper/middleware
	"log" // Temporal, idealmente usar logger estructurado inyectado o global
	"net/http"
	"strconv" // Para convertir ID de URL

	"github.com/gin-gonic/gin" // El framework web
)

// CategoriaHandler maneja las peticiones HTTP relacionadas con categorías.
type CategoriaHandler struct {
	service service.CategoriaService // Dependencia de la INTERFAZ del servicio
	// logger *zap.Logger            // Idealmente inyectar logger
}

// NewCategoriaHandler es la Factory Function para crear el handler.
// Recibe la dependencia (servicio) y devuelve un puntero al handler.
func NewCategoriaHandler(s service.CategoriaService /*, logger *zap.Logger*/) *CategoriaHandler {
	return &CategoriaHandler{
		service: s,
		// logger: logger,
	}
}

// --- Mapeadores Helper (Internos al Handler o en un paquete 'mapper') ---
// Estos ayudan a convertir entre domain.Categoria y dto.CategoriaResponseDTO

func mapCategoriaToResponseDTO(cat domain.Categoria) dto.CategoriaResponseDTO {
	return dto.CategoriaResponseDTO{
		ID:     cat.ID,
		Nombre: cat.Nombre,
		Slug:   cat.Slug,
	}
}

func mapCategoriasToResponseDTOs(cats []domain.Categoria) []dto.CategoriaResponseDTO {
	responseDTOs := make([]dto.CategoriaResponseDTO, 0, len(cats))
	for _, cat := range cats {
		responseDTOs = append(responseDTOs, mapCategoriaToResponseDTO(cat))
	}
	return responseDTOs
}


// --- Métodos del Handler (corresponden a las rutas) ---

// GetAll godoc
// @Summary Obtiene todas las categorías
// @Description Devuelve una lista de todas las categorías existentes
// @Tags Categorias
// @Accept json
// @Produce json
// @Success 200 {array} dto.CategoriaResponseDTO "Lista de categorías"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias [get]
// @Security ApiKeyAuth
func (h *CategoriaHandler) GetAll(c *gin.Context) {
	categorias, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Handler: Error al llamar a service.GetAll: %v\n", err) // Loguear error interno
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
		return
	}

	// Mapear el slice de domain.Categoria a un slice de dto.CategoriaResponseDTO
	response := mapCategoriasToResponseDTOs(categorias)
	c.JSON(http.StatusOK, response)
}

// GetByID godoc
// @Summary Obtiene una categoría por ID
// @Description Devuelve los detalles de una categoría específica por su ID
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path uint true "ID de la Categoría"
// @Success 200 {object} dto.CategoriaResponseDTO "Categoría encontrada"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Categoría no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias/{id} [get]
// @Security ApiKeyAuth
func (h *CategoriaHandler) GetByID(c *gin.Context) {
	// 1. Obtener y validar ID de la URL
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32) // Parsear como uint
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	id := uint(idUint64) // Convertir a uint

	// 2. Llamar al servicio
	categoria, err := h.service.GetByID(c.Request.Context(), id)

	// 3. Manejar error del servicio
	if err != nil {
		if errors.Is(err, domain.ErrCategoriaNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler: Error al llamar a service.GetByID(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la categoría"})
		}
		return
	}

	// 4. Mapear y responder
	response := mapCategoriaToResponseDTO(*categoria)
	c.JSON(http.StatusOK, response)
}

// Create godoc
// @Summary Crea una nueva categoría
// @Description Crea una nueva categoría con el nombre proporcionado
// @Tags Categorias
// @Accept json
// @Produce json
// @Param categoria body dto.CategoriaRequestDTO true "Datos de la categoría a crear"
// @Success 201 {object} dto.CategoriaResponseDTO "Categoría creada exitosamente"
// @Failure 400 {object} gin.H "Datos de entrada inválidos o nombre ya existe"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias [post]
// @Security ApiKeyAuth
func (h *CategoriaHandler) Create(c *gin.Context) {
	var requestBody dto.CategoriaRequestDTO

	// 1. Bind y validar JSON de entrada
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Handler: Error en BindJSON para crear categoría: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// 2. Mapear DTO de Handler a DTO de Servicio
	serviceInput := service.CategoriaInputDTO{
		Nombre: requestBody.Nombre,
	}

	// 3. Llamar al servicio
	nuevaCategoria, err := h.service.Create(c.Request.Context(), serviceInput)

	// 4. Manejar error del servicio
	if err != nil {
		if errors.Is(err, domain.ErrCategoriaNombreYaExiste) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict es apropiado
		} else {
			// Podríamos tener otros errores de validación del servicio aquí
			log.Printf("Handler: Error al llamar a service.Create: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la categoría"})
		}
		return
	}

	// 5. Mapear y responder con el objeto creado
	response := mapCategoriaToResponseDTO(*nuevaCategoria)
	c.JSON(http.StatusCreated, response) // 201 Created
}

// Update godoc
// @Summary Actualiza una categoría existente
// @Description Actualiza el nombre y slug de una categoría por su ID
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path uint true "ID de la Categoría a actualizar"
// @Param categoria body dto.CategoriaRequestDTO true "Nuevos datos para la categoría"
// @Success 200 {object} dto.CategoriaResponseDTO "Categoría actualizada exitosamente"
// @Failure 400 {object} gin.H "ID inválido o datos de entrada inválidos"
// @Failure 404 {object} gin.H "Categoría no encontrada"
// @Failure 409 {object} gin.H "El nuevo nombre ya existe en otra categoría"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias/{id} [put]
// @Security ApiKeyAuth
func (h *CategoriaHandler) Update(c *gin.Context) {
	// 1. Obtener y validar ID de la URL
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	id := uint(idUint64)

	// 2. Bind y validar JSON de entrada
	var requestBody dto.CategoriaRequestDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Handler: Error en BindJSON para actualizar categoría: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// 3. Mapear DTO de Handler a DTO de Servicio
	serviceInput := service.CategoriaInputDTO{
		Nombre: requestBody.Nombre,
	}

	// 4. Llamar al servicio
	categoriaActualizada, err := h.service.Update(c.Request.Context(), id, serviceInput)

	// 5. Manejar error del servicio
    if err != nil {
        if errors.Is(err, domain.ErrCategoriaNotFound) {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else if errors.Is(err, domain.ErrCategoriaNombreYaExiste) {
            c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
        } else {
            // Podríamos tener otros errores de validación del servicio aquí
            log.Printf("Handler: Error al llamar a service.Update(%d): %v\n", id, err)
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la categoría"})
        }
        return
    }

	// 6. Mapear y responder
	response := mapCategoriaToResponseDTO(*categoriaActualizada)
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Elimina una categoría por ID
// @Description Elimina permanentemente una categoría específica por su ID
// @Tags Categorias
// @Accept json
// @Produce json
// @Param id path uint true "ID de la Categoría a eliminar"
// @Success 204 "Categoría eliminada exitosamente (sin contenido)"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Categoría no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias/{id} [delete]
// @Security ApiKeyAuth
func (h *CategoriaHandler) Delete(c *gin.Context) {
	// 1. Obtener y validar ID de la URL
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	id := uint(idUint64)

	// 2. Llamar al servicio
	err = h.service.Delete(c.Request.Context(), id)

	// 3. Manejar error del servicio
	if err != nil {
		if errors.Is(err, domain.ErrCategoriaNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			// Podríamos tener otros errores (ej: foreign key si el servicio lo manejara)
			log.Printf("Handler: Error al llamar a service.Delete(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la categoría"})
		}
		return
	}

	// 4. Responder (Éxito para DELETE suele ser 204 No Content)
	c.Status(http.StatusNoContent)
}