// backend/recetas/receta_api.go
package recetas // Paquete de la característica 'recetas'

import (
	// Importar paquetes necesarios
	"backend/categorias" // Para los errores de dominio de categoría (ErrCategoriaNotFound)
	"errors"             // Para errors.Is
	"fmt"                // Para logs/errores formateados
	"log"                // Logging temporal
	"net/http"           // Para status codes
	"strconv"            // Para convertir ID de URL
	"strings"            // Para formatear errores
	"time"               // Para formatear CreatedAt/UpdatedAt

	"github.com/gin-gonic/gin" // El framework web
)

// RecetaHandler maneja las peticiones HTTP relacionadas con Recetas.
type RecetaHandler struct {
	service RecetaService // Dependencia de la interfaz RecetaService (de este paquete)
	// logger *zap.Logger
}

// NewRecetaHandler es la Factory Function para crear el handler.
func NewRecetaHandler(s RecetaService) *RecetaHandler {
	return &RecetaHandler{service: s}
}

// --- Mapeadores Helper (Internos al Handler) ---

// mapDomainRecetaToResponseDTO convierte un domain.Receta a un RecetaResponseDTO.
func mapDomainRecetaToResponseDTO(receta Receta) RecetaResponseDTO {
	// El campo receta.Categoria es *categorias.Categoria (dominio)
	// Necesitamos convertirlo a categorias.CategoriaResponseDTO
	var catDTO categorias.CategoriaResponseDTO
	if receta.Categoria != nil { // Asegurarse que la categoría fue cargada
		catDTO = categorias.CategoriaResponseDTO{ // Usar el constructor/mapeador de categorias si existe, o hacerlo manual
			ID:     receta.Categoria.ID,
			Nombre: receta.Categoria.Nombre,
			Slug:   receta.Categoria.Slug,
		}
	}

	return RecetaResponseDTO{
		ID:                receta.ID,
		Nombre:            receta.Nombre,
		Slug:              receta.Slug,
		TiempoPreparacion: receta.TiempoPreparacion,
		Descripcion:       receta.Descripcion,
		Foto:              receta.Foto, // Aquí podrías añadir lógica para construir URL completa si Foto es solo nombre de archivo
		CreatedAt:         receta.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         receta.UpdatedAt.Format(time.RFC3339),
		Categoria:         catDTO,
	}
}

// mapDomainRecetasToResponseDTOs convierte un slice de domain.Receta a un slice de RecetaResponseDTO.
func mapDomainRecetasToResponseDTOs(recetas []Receta) []RecetaResponseDTO {
	responseDTOs := make([]RecetaResponseDTO, 0, len(recetas))
	for _, r := range recetas {
		responseDTOs = append(responseDTOs, mapDomainRecetaToResponseDTO(r))
	}
	return responseDTOs
}

// --- Métodos del Handler ---

// GetAll maneja GET /recetas
func (h *RecetaHandler) GetAll(c *gin.Context) {
	domainRecetas, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Handler Recetas: Error service.GetAll: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener recetas"})
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetasToResponseDTOs(domainRecetas))
}

// GetByID maneja GET /recetas/:id
func (h *RecetaHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de receta inválido"})
		return
	}
	id := uint(idUint64)

	domainReceta, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRecetaNotFound) { // Error de dominio de este paquete
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler Recetas: Error service.GetByID(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la receta"})
		}
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetaToResponseDTO(*domainReceta))
}

// Create maneja POST /recetas
func (h *RecetaHandler) Create(c *gin.Context) {
	var req RecetaRequestDTO // DTO de API de este paquete
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Handler Recetas: Error BindJSON en Create: %v\n", err)
		// Devolver errores de validación más específicos sería ideal con un middleware de errores
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	// Mapear DTO de API a DTO de Servicio (definido en service_dto.go de este paquete)
	serviceInput := RecetaInputDTO{
		Nombre:            req.Nombre,
		CategoriaID:       req.CategoriaID,
		TiempoPreparacion: req.TiempoPreparacion,
		Descripcion:       req.Descripcion,
		Foto:              req.Foto, // Asumimos que la foto es una URL o se maneja más tarde
	}

	nuevaDomainReceta, err := h.service.Create(c.Request.Context(), serviceInput)
	if err != nil {
		// Chequear errores específicos devueltos por el servicio
		if errors.Is(err, ErrRecetaNombreInvalido) || errors.Is(err, ErrRecetaSinCategoria) ||
			strings.Contains(err.Error(), "la categoría ID") { // Para el error formateado del servicio
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler Recetas: Error service.Create: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la receta"})
		}
		return
	}
	c.JSON(http.StatusCreated, mapDomainRecetaToResponseDTO(*nuevaDomainReceta))
}

// Update maneja PUT /recetas/:id
func (h *RecetaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de receta inválido"})
		return
	}
	id := uint(idUint64)

	var req RecetaRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Handler Recetas: Error BindJSON en Update: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	serviceInput := RecetaInputDTO{
		Nombre:            req.Nombre,
		CategoriaID:       req.CategoriaID,
		TiempoPreparacion: req.TiempoPreparacion,
		Descripcion:       req.Descripcion,
		Foto:              req.Foto,
	}

	domainRecetaActualizada, err := h.service.Update(c.Request.Context(), id, serviceInput)
	if err != nil {
		if errors.Is(err, ErrRecetaNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, ErrRecetaNombreInvalido) || errors.Is(err, ErrRecetaSinCategoria) ||
			strings.Contains(err.Error(), "la categoría ID") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler Recetas: Error service.Update(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la receta"})
		}
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetaToResponseDTO(*domainRecetaActualizada))
}

// Delete maneja DELETE /recetas/:id
func (h *RecetaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de receta inválido"})
		return
	}
	id := uint(idUint64)

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRecetaNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Handler Recetas: Error service.Delete(%d): %v\n", id, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la receta"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// FindByCategoria maneja GET /categorias/:categoria_id/recetas (EJEMPLO DE RUTA ADICIONAL)
// Necesitaremos registrar esta ruta si la queremos.
func (h *RecetaHandler) FindByCategoria(c *gin.Context) {
	catIdStr := c.Param("categoria_id") // Asumiendo que el param se llama así en la ruta
	catIdUint64, err := strconv.ParseUint(catIdStr, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID de categoría inválido"})
		return
	}
	catId := uint(catIdUint64)

	domainRecetas, err := h.service.FindByCategoriaID(c.Request.Context(), catId)
	if err != nil {
		// El servicio ya valida si la categoría existe, así que un error aquí es probablemente interno
		// o un error de formato en la respuesta del servicio.
		// Si ErrRecetaSinCategoria (porque la categoría no existe) se devuelve, sería un 404 o 400.
		if strings.Contains(err.Error(), "la categoría ID") && strings.Contains(err.Error(), "no existe") {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No se encontraron recetas para la categoría ID %d o la categoría no existe.", catId)})
		} else {
			log.Printf("Handler Recetas: Error service.FindByCategoriaID(%d): %v\n", catId, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener recetas por categoría"})
		}
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetasToResponseDTOs(domainRecetas))
}
