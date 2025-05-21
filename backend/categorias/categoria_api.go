// backend/categorias/handler.go
// Implementación con Gin de CategoriaHandler.

//
// @Tags Categorias
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /categorias/:id [get]
// @Param id path uint true "ID de la categoria"
// @Success 200 {object} categorias.CategoriaResponseDTO "Categoria"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Categoria no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"

package categorias // Paquete de la característica

import (
	// Imports estándar y de terceros
	//"backend/shared/apitypes" // Para ErrorResponseDTO en anotaciones @Failure
	// Para errors.Is
	// "log" // Ya no es necesario aquí si el middleware loguea
	"fmt"
	"net/http"
	"strconv" // Para convertir ID de URL

	// "time" // No se usa directamente aquí si los mapeadores se simplifican

	"github.com/gin-gonic/gin"
	// "backend/shared/apitypes" // Ya no es necesario aquí, el middleware usa esto
)

// CategoriaHandler maneja las peticiones HTTP relacionadas con categorías.
type CategoriaHandler struct {
	service CategoriaService // Dependencia de la interfaz (definida en este paquete)
}

// NewCategoriaHandler es la Factory Function para crear el handler.
func NewCategoriaHandler(s CategoriaService) *CategoriaHandler {
	return &CategoriaHandler{
		service: s,
	}
}

// --- Mapeadores Helper (Internos al Handler o en un archivo separado dentro de este paquete) ---
// Usan los tipos Categoria (dominio) y CategoriaResponseDTO (API DTO) definidos en este paquete.

func mapDomainToResponseDTO(cat Categoria) CategoriaResponseDTO {
	return CategoriaResponseDTO{
		ID:     cat.ID,
		Nombre: cat.Nombre,
		Slug:   cat.Slug,
		// CreatedAt: cat.CreatedAt.Format(time.RFC3339), // Quitado para simplicidad o si el DTO no lo tiene
		// UpdatedAt: cat.UpdatedAt.Format(time.RFC3339),
	}
}

func mapDomainsToResponseDTOs(cats []Categoria) []CategoriaResponseDTO {
	responseDTOs := make([]CategoriaResponseDTO, 0, len(cats))
	for _, cat := range cats {
		responseDTOs = append(responseDTOs, mapDomainToResponseDTO(cat))
	}
	return responseDTOs
}

// --- Métodos del Handler (Refactorizados para usar c.Error) ---

// GetAll maneja GET /categorias
// Anotaciones Swagger se mantienen igual, pero el cuerpo del error será ErrorResponseDTO
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
func (h *CategoriaHandler) GetAll(c *gin.Context) {
	domainCategorias, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		_ = c.Error(err) // Pasar error al middleware
		return
	}
	c.JSON(http.StatusOK, mapDomainsToResponseDTOs(domainCategorias))
}

// GetByID maneja GET /categorias/:id
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Categoría no encontrada"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
func (h *CategoriaHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// Este es un error de parseo de la entrada del cliente.
		// El middleware podría no tener contexto para "ID de categoría inválido".
		// Podríamos devolver un 400 directamente o un error específico que el middleware entienda.
		// Por ahora, para que el middleware lo maneje como un 400 genérico o 500:
		_ = c.Error(fmt.Errorf("parámetro ID inválido: %s - %w", idStr, err)) // Error más descriptivo
		return
	}
	id := uint(idUint64)

	domainCategoria, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err) // Pasar error al middleware (ej: ErrCategoriaNotFound)
		return
	}
	c.JSON(http.StatusOK, mapDomainToResponseDTO(*domainCategoria))
}

// Create maneja POST /categorias
// @Failure 400 {object} apitypes.ErrorResponse "Datos de entrada inválidos"
// @Failure 409 {object} apitypes.ErrorResponse "Conflicto - El nombre de la categoría ya existe"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor"
func (h *CategoriaHandler) Create(c *gin.Context) {
	var requestBody CategoriaRequestDTO // DTO de este paquete
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		_ = c.Error(err) // Pasar error de validación de Gin al middleware
		return
	}

	serviceInput := CategoriaInputDTO{ // DTO de servicio de este paquete
		Nombre: requestBody.Nombre,
	}

	nuevaDomainCategoria, err := h.service.Create(c.Request.Context(), serviceInput)
	if err != nil {
		_ = c.Error(err) // Pasar error del servicio al middleware (ej: ErrCategoriaNombreYaExiste)
		return
	}
	c.JSON(http.StatusCreated, mapDomainToResponseDTO(*nuevaDomainCategoria))
}

// Update maneja PUT /categorias/:id
// (Anotaciones @Failure similares usando apitypes.ErrorResponse)
func (h *CategoriaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		_ = c.Error(fmt.Errorf("parámetro ID inválido para update: %s - %w", idStr, err))
		return
	}
	id := uint(idUint64)

	var requestBody CategoriaRequestDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		_ = c.Error(err) // Pasar error de validación de Gin
		return
	}

	serviceInput := CategoriaInputDTO{
		Nombre: requestBody.Nombre,
	}

	domainCategoriaActualizada, err := h.service.Update(c.Request.Context(), id, serviceInput)
	if err != nil {
		_ = c.Error(err) // Pasar error del servicio
		return
	}
	c.JSON(http.StatusOK, mapDomainToResponseDTO(*domainCategoriaActualizada))
}

// Delete maneja DELETE /categorias/:id
// (Anotaciones @Failure similares usando apitypes.ErrorResponse)
func (h *CategoriaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		_ = c.Error(fmt.Errorf("parámetro ID inválido para delete: %s - %w", idStr, err))
		return
	}
	id := uint(idUint64)

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err) // Pasar error del servicio
		return
	}
	c.Status(http.StatusNoContent)
}
