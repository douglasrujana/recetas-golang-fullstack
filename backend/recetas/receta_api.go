// backend/recetas/handler.go
// Implementación con Gin de RecetaHandler.
//
// @Tags Recetas
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /recetas/:id [get]
// @Param id path uint true "ID de la receta"
// @Success 200 {object} recetas.RecetaResponseDTO "Receta"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Receta no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"

package recetas // Paquete de la característica 'recetas'

import (
	// Importar paquetes necesarios
	"backend/categorias" // Para tipos como categorias.CategoriaResponseDTO y errores como categorias.ErrCategoriaNotFound
	//"errors"         // Para errors.Is
	"fmt"            // Para formatear errores si es necesario pasarlos al middleware
	// "log" // Ya no es tan necesario aquí si el middleware loguea centralmente
	"net/http"
	"strconv"
	//"strings" // Para strings.Contains en el manejo de errores específico
	"time"    // Para formatear CreatedAt/UpdatedAt
	"github.com/gin-gonic/gin" // El framework web
	// "backend/shared/apitypes" // No es necesario importar aquí si el middleware lo usa
)

// RecetaHandler maneja las peticiones HTTP relacionadas con Recetas.
type RecetaHandler struct {
	service RecetaService // Dependencia de la interfaz RecetaService (de este paquete)
}

// NewRecetaHandler es la Factory Function para crear el handler.
func NewRecetaHandler(s RecetaService) *RecetaHandler {
	return &RecetaHandler{service: s}
}

// --- Mapeadores Helper (Internos al Handler) ---

// mapDomainRecetaToResponseDTO convierte un domain.Receta a un RecetaResponseDTO.
func mapDomainRecetaToResponseDTO(receta Receta) RecetaResponseDTO {
	var catDTO categorias.CategoriaResponseDTO // Tipo del paquete 'categorias'
	if receta.Categoria != nil {
		catDTO = categorias.CategoriaResponseDTO{
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
		Foto:              receta.Foto,
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

// --- Métodos del Handler (Refactorizados para delegar errores) ---

// GetAll maneja GET /recetas
// Las anotaciones Swagger deberían actualizar sus @Failure para usar apitypes.ErrorResponse
func (h *RecetaHandler) GetAll(c *gin.Context) {
	domainRecetas, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		_ = c.Error(err) // Pasar error al middleware
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetasToResponseDTOs(domainRecetas))
}

// GetByID maneja GET /recetas/:id
// GetByID godoc
// @Summary Obtiene una receta por ID
// @Description Devuelve los detalles de una receta específica por su ID, con su categoría.
// @Tags Recetas
// @Accept  json
// @Produce json
// @Param   id path uint true "ID de la Receta" example:"1"
// @Success 200 {object} RecetaResponseDTO "Receta encontrada"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Receta no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
// @Router /recetas/{id} [get]
// @Security ApiKeyAuth
func (h *RecetaHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, errConv := strconv.ParseUint(idStr, 10, 32)
	if errConv != nil {
		// Crear un error específico para que el middleware lo interprete como Bad Request
		// o simplemente pasar el error de parseo.
		_ = c.Error(fmt.Errorf("ID de receta inválido en URL: %s - %w", idStr, errConv))
		return
	}
	id := uint(idUint64)

	domainReceta, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err) // Pasar error (ej: ErrRecetaNotFound) al middleware
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetaToResponseDTO(*domainReceta))
}

// Create maneja POST /recetas
// Create godoc
// @Summary Crea una nueva receta
// @Description Crea una nueva receta con los datos proporcionados. La foto se maneja como un nombre de archivo o URL.
// @Tags Recetas
// @Accept  json
// @Produce json
// @Param   receta body RecetaRequestDTO true "Datos de la Receta a Crear"
// @Success 201 {object} RecetaResponseDTO "Receta creada exitosamente"
// @Failure 400 {object} apitypes.ErrorResponse "Datos de entrada inválidos (ej: validación, categoría ID no existe)"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
// @Router /recetas [post]
// @Security ApiKeyAuth
func (h *RecetaHandler) Create(c *gin.Context) {
	var req RecetaRequestDTO // DTO de API de este paquete
	if err := c.ShouldBindJSON(&req); err != nil {
		// El middleware puede manejar validator.ValidationErrors si está configurado
		_ = c.Error(err) // Pasar error de validación de Gin
		return
	}

	serviceInput := RecetaInputDTO{ // DTO de Servicio de este paquete
		Nombre:            req.Nombre,
		CategoriaID:       req.CategoriaID,
		TiempoPreparacion: req.TiempoPreparacion,
		Descripcion:       req.Descripcion,
		Foto:              req.Foto,
	}

	nuevaDomainReceta, err := h.service.Create(c.Request.Context(), serviceInput)
	if err != nil {
		_ = c.Error(err) // Pasar error del servicio (ej: ErrRecetaSinCategoria, ErrRecetaNombreInvalido)
		return
	}
	c.JSON(http.StatusCreated, mapDomainRecetaToResponseDTO(*nuevaDomainReceta))
}

// Update maneja PUT /recetas/:id
// Update godoc
// @Summary Actualiza una receta existente
// @Description Actualiza una receta existente con los datos proporcionados. La foto se maneja como un nombre de archivo o URL.
// @Tags Recetas
// @Accept  json
// @Produce json
// @Param   id path uint true "ID de la Receta a Actualizar" example:"1"
// @Param   receta body RecetaRequestDTO true "Datos de la Receta a Actualizar"
// @Success 200 {object} RecetaResponseDTO "Receta actualizada exitosamente"
// @Failure 400 {object} apitypes.ErrorResponse "Datos de entrada inválidos (ej: validación, categoría ID no existe)"
// @Failure 404 {object} apitypes.ErrorResponse "Receta no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
// @Router /recetas/{id} [put]
// @Security ApiKeyAuth
func (h *RecetaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, errConv := strconv.ParseUint(idStr, 10, 32)
	if errConv != nil {
		_ = c.Error(fmt.Errorf("ID de receta inválido en URL para update: %s - %w", idStr, errConv))
		return
	}
	id := uint(idUint64)

	var req RecetaRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err) // Pasar error de validación de Gin
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
		_ = c.Error(err) // Pasar error del servicio
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetaToResponseDTO(*domainRecetaActualizada))
}

// Delete maneja DELETE /recetas/:id
// Delete godoc
// @Summary Elimina una receta por ID
// @Description Elimina permanentemente una receta específica por su ID.
// @Tags Recetas
// @Accept  json
// @Produce json
// @Param   id path uint true "ID de la Receta a Eliminar" example:"1"
// @Success 204 "Sin contenido (eliminación exitosa)"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Receta no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
// @Router /recetas/{id} [delete]
// @Security ApiKeyAuth
func (h *RecetaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, errConv := strconv.ParseUint(idStr, 10, 32)
	if errConv != nil {
		_ = c.Error(fmt.Errorf("ID de receta inválido en URL para delete: %s - %w", idStr, errConv))
		return
	}
	id := uint(idUint64)

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err) // Pasar error del servicio
		return
	}
	c.Status(http.StatusNoContent)
}

// FindByCategoria maneja GET /recetas/categoria/:categoria_id
// FindByCategoria godoc
// @Summary Obtiene recetas por ID de categoría
// @Description Devuelve una lista de recetas que pertenecen a una categoría específica.
// @Tags Recetas, Categorias
// @Accept  json
// @Produce json
// @Param   categoria_id path uint true "ID de la Categoría para filtrar recetas" example:"1"
// @Success 200 {array} RecetaResponseDTO "Lista de recetas para la categoría"
// @Failure 400 {object} apitypes.ErrorResponse "ID de categoría inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Categoría no encontrada o sin recetas"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
// @Router /categorias/{categoria_id}/recetas [get] // Ruta como la definimos
// @Security ApiKeyAuth
func (h *RecetaHandler) FindByCategoria(c *gin.Context) {
	catIdStr := c.Param("categoria_id")
	catIdUint64, errConv := strconv.ParseUint(catIdStr, 10, 32)
	if errConv != nil {
		_ = c.Error(fmt.Errorf("ID de categoría inválido en URL para filtro: %s - %w", catIdStr, errConv))
		return
	}
	catId := uint(catIdUint64)

	domainRecetas, err := h.service.FindByCategoriaID(c.Request.Context(), catId)
	if err != nil {
		_ = c.Error(err) // Pasar error del servicio
		return
	}
	c.JSON(http.StatusOK, mapDomainRecetasToResponseDTOs(domainRecetas))
}