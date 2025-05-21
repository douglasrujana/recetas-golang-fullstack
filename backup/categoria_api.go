// backend/categorias/categoria_handler.go
package categorias // Paquete de la característica

import (
	// Ya no necesitamos los imports "backend/internal/..." para domain, service, dto
	// si todos esos tipos ahora viven en este paquete 'categorias'.
	// Imports estándar y de terceros
	"errors" // Para errors.Is
	"log"    // Temporal, idealmente usar logger estructurado
	"net/http"
	"strconv" // Para convertir ID de URL

	"github.com/gin-gonic/gin"
)

// CategoriaHandler maneja las peticiones HTTP relacionadas con categorías.
type CategoriaHandler struct {
	// La dependencia es de la interfaz CategoriaService,
	// que ahora está definida en este mismo paquete 'categorias'
	// (probablemente en categoria_service.go).
	service CategoriaService
	// logger *zap.Logger // Idealmente inyectar logger
}

// NewCategoriaHandler es la Factory Function para crear el handler.
// Recibe la interfaz CategoriaService.
func NewCategoriaHandler(s CategoriaService /*, logger *zap.Logger*/) *CategoriaHandler {
	return &CategoriaHandler{
		service: s,
		// logger: logger,
	}
}

// --- Mapeadores Helper ---
// Usan los tipos Categoria (dominio) y CategoriaResponseDTO (API DTO)
// definidos en este paquete 'categorias'.

func mapCategoriaToResponseDTO(cat Categoria) CategoriaResponseDTO {
	return CategoriaResponseDTO{
		ID:     cat.ID,
		Nombre: cat.Nombre,
		Slug:   cat.Slug,
		// Asumiendo que CategoriaResponseDTO tiene estos campos
		// CreatedAt: cat.CreatedAt.Format(time.RFC3339),
		// UpdatedAt: cat.UpdatedAt.Format(time.RFC3339),
	}
}

func mapCategoriasToResponseDTOs(cats []Categoria) []CategoriaResponseDTO {
	responseDTOs := make([]CategoriaResponseDTO, 0, len(cats))
	for _, cat := range cats {
		responseDTOs = append(responseDTOs, mapCategoriaToResponseDTO(cat))
	}
	return responseDTOs
}

// --- Métodos del Handler ---

// GetAll maneja GET /categorias
// GetAll godoc
// @Summary Obtiene todas las categorías
// @Description Devuelve una lista de todas las categorías existentes.
// @Tags Categorias
// @Accept  json
// @Produce json
// @Success 200 {array} CategoriaResponseDTO "Lista de categorías"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias [get]
// @Security ApiKeyAuth // Indica que esta ruta requiere la seguridad ApiKeyAuth
func (h *CategoriaHandler) GetAll(c *gin.Context) {
	// service.GetAll() devuelve []Categoria (dominio)
	domainCategorias, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Handler: Error al llamar a service.GetAll: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
		return
	}

	// Mapear slice de dominio a slice de DTO de respuesta API
	response := mapCategoriasToResponseDTOs(domainCategorias)
	c.JSON(http.StatusOK, response)
}

// GetByID maneja GET /categorias/:id
// GetByID godoc
// @Summary Obtiene una categoría por ID
// @Description Devuelve los detalles de una categoría específica por su ID.
// @Tags Categorias
// @Accept  json
// @Produce json
// @Param   id path uint true "ID de la Categoría" example:"1"
// @Success 200 {object} CategoriaResponseDTO "Categoría encontrada"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Categoría no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias/{id} [get]
// @Security ApiKeyAuth
func (h *CategoriaHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}
	id := uint(idUint64)

	// service.GetByID() devuelve *Categoria (dominio)
	domainCategoria, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		// Usar errores de dominio definidos en este paquete 'categorias'
		if errors.Is(err, ErrCategoriaNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler: Error al llamar a service.GetByID(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la categoría"})
		}
		return
	}

	// Mapear dominio a DTO de respuesta API
	response := mapCategoriaToResponseDTO(*domainCategoria)
	c.JSON(http.StatusOK, response)
}

// Create maneja POST /categorias
// Create godoc
// @Summary Crea una nueva categoría
// @Description Crea una nueva categoría con el nombre proporcionado.
// @Tags Categorias
// @Accept  json
// @Produce json
// @Param   categoria body CategoriaRequestDTO true "Datos de la Categoría a Crear"
// @Success 201 {object} CategoriaResponseDTO "Categoría creada exitosamente"
// @Failure 400 {object} gin.H "Datos de entrada inválidos"
// @Failure 409 {object} gin.H "Conflicto - El nombre de la categoría ya existe"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias [post]
// @Security ApiKeyAuth
func (h *CategoriaHandler) Create(c *gin.Context) {
	// CategoriaRequestDTO está definido en este paquete 'categorias'
	var requestBody CategoriaRequestDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Handler: Error en BindJSON para crear categoría: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// CategoriaInputDTO (para el servicio) está definido en este paquete 'categorias'
	serviceInput := CategoriaInputDTO{
		Nombre: requestBody.Nombre,
	}

	// service.Create() devuelve *Categoria (dominio)
	nuevaDomainCategoria, err := h.service.Create(c.Request.Context(), serviceInput)
	if err != nil {
		// Usar errores de dominio definidos en este paquete 'categorias'
		if errors.Is(err, ErrCategoriaNombreYaExiste) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler: Error al llamar a service.Create: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la categoría"})
		}
		return
	}

	// Mapear dominio a DTO de respuesta API
	response := mapCategoriaToResponseDTO(*nuevaDomainCategoria)
	c.JSON(http.StatusCreated, response)
}

// Update maneja PUT /categorias/:id
// Update godoc
// @Summary Actualiza una categoría existente
// @Description Actualiza los datos de una categoría existente por su ID.
// @Tags Categorias
// @Accept  json
// @Produce json
// @Param   id path uint true "ID de la Categoría a Actualizar" example:"1"
// @Param   categoria body CategoriaRequestDTO true "Datos de la Categoría a Actualizar"
// @Success 200 {object} CategoriaResponseDTO "Categoría actualizada exitosamente"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Categoría no encontrada"
// @Failure 409 {object} gin.H "Conflicto - El nombre de la categoría ya existe"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias/{id} [put]
// @Security ApiKeyAuth
func (h *CategoriaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}
	id := uint(idUint64)

	var requestBody CategoriaRequestDTO // Tipo de este paquete
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Handler: Error en BindJSON para actualizar categoría: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	serviceInput := CategoriaInputDTO{ // Tipo de este paquete
		Nombre: requestBody.Nombre,
	}

	// service.Update() devuelve *Categoria (dominio)
	domainCategoriaActualizada, err := h.service.Update(c.Request.Context(), id, serviceInput)
	if err != nil {
		if errors.Is(err, ErrCategoriaNotFound) { // Error de dominio de este paquete
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, ErrCategoriaNombreYaExiste) { // Error de dominio de este paquete
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler: Error al llamar a service.Update(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la categoría"})
		}
		return
	}

	// Mapear dominio a DTO de respuesta API
	response := mapCategoriaToResponseDTO(*domainCategoriaActualizada)
	c.JSON(http.StatusOK, response)
}

// Delete maneja DELETE /categorias/:id
// Delete godoc
// @Summary Elimina una categoría
// @Description Elimina una categoría existente por su ID.
// @Tags Categorias
// @Accept  json
// @Produce json
// @Param   id path uint true "ID de la Categoría a Eliminar" example:"1"
// @Success 204 {object} gin.H "Categoría eliminada exitosamente"
// @Failure 400 {object} gin.H "ID inválido"
// @Failure 404 {object} gin.H "Categoría no encontrada"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /categorias/{id} [delete]
// @Security ApiKeyAuth
func (h *CategoriaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}
	id := uint(idUint64)

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrCategoriaNotFound) { // Error de dominio de este paquete
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler: Error al llamar a service.Delete(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la categoría"})
		}
		return
	}
	c.Status(http.StatusNoContent) // 204 No Content
}
