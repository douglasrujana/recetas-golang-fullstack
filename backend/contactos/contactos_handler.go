// backend/contactos/handler.go
package contactos

import (
	// "backend/shared/apitypes" // Para el ErrorResponseDTO
	// "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ContactoHandler maneja las peticiones HTTP para Contactos.
type ContactoHandler struct {
	service ContactoService
}

// NewContactoHandler crea una nueva instancia de ContactoHandler.
func NewContactoHandler(s ContactoService) *ContactoHandler {
	return &ContactoHandler{service: s}
}

// EnviarMensaje maneja la creación de un nuevo mensaje de contacto.
// @Summary Envía un mensaje de contacto
// @Description Recibe los datos del formulario de contacto, los guarda y notifica al admin.
// @Tags Contactos
// @Accept json
// @Produce json
// @Param contacto_form body ContactoRequestDTO true "Datos del Formulario de Contacto"
// @Success 201 {object} ContactoResponseDTO "Mensaje enviado y guardado exitosamente"
// @Failure 400 {object} apitypes.ErrorResponse "Datos de entrada inválidos"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno del servidor al procesar el mensaje"
// @Router /contactos [post]
func (h *ContactoHandler) EnviarMensaje(c *gin.Context) {
	var req ContactoRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err) // El middleware de errores manejará esto (validator.ValidationErrors)
		return
	}

	// Obtener UserID del contexto si el usuario está autenticado (esto vendría de un middleware de Auth)
	var userIDPtr *uint
	// userIDVal, exists := c.Get("userID") // Ejemplo si el middleware de Auth setea "userID"
	// if exists {
	//  uid, ok := userIDVal.(uint)
	//  if ok {
	//      userIDPtr = &uid
	//  }
	// }

	input := EnviarContactoInput{
		Nombre:    req.Nombre,
		Email:     req.Email,
		Telefono:  req.Telefono,
		Asunto:    req.Asunto,
		Mensaje:   req.Mensaje,
		IPOrigen:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		UserID:    userIDPtr, // Será nil si no hay usuario autenticado
	}

	contactoGuardado, err := h.service.ProcesarNuevoContacto(c.Request.Context(), input)
	if err != nil {
		// Pasar el error al middleware global de errores.
		// El middleware decidirá el código de estado basado en el tipo de error
		// (ej: ErrNombreRemitenteVacio -> 400, error de email -> 500).
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, ContactoResponseDTO{
		ID:        contactoGuardado.ID, // Devolver el ID del mensaje guardado
		Status:    "Mensaje enviado y guardado correctamente.",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// --- Handlers para administración (ejemplos, requerirían autenticación de admin) ---

// GetAllContactos godoc
// @Summary (Admin) Obtiene todos los mensajes de contacto
// @Description (Admin) Devuelve una lista de todos los mensajes recibidos.
// @Tags Contactos_Admin
// @Produce json
// @Success 200 {array} ContactoListItemDTO "Lista de mensajes"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno"
// @Router /admin/contactos [get]
// @Security ApiKeyAuth // Asumiendo que las rutas admin están protegidas
func (h *ContactoHandler) GetAllContactos(c *gin.Context) {
	// TODO: Implementar lógica de paginación/filtros si es necesario
	domainForms, err := h.service.ObtenerTodosLosContactos(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	responseDTOs := make([]ContactoListItemDTO, 0, len(domainForms))
	for _, df := range domainForms {
		responseDTOs = append(responseDTOs, ContactoListItemDTO{
			ID:            df.ID,
			NombreRemitente: df.NombreRemitente,
			EmailRemitente: df.EmailRemitente,
			Asunto:        df.Asunto,
			FechaContacto: df.FechaContacto.Format("2006-01-02 15:04"),
			Leido:         df.Leido,
		})
	}
	c.JSON(http.StatusOK, responseDTOs)
}

// MarcarComoLeido godoc
// @Summary (Admin) Marca un mensaje como leído
// @Description (Admin) Actualiza el estado de un mensaje a 'leído'.
// @Tags Contactos_Admin
// @Produce json
// @Param id path uint true "ID del Mensaje de Contacto"
// @Success 200 {object} gin.H "Mensaje marcado como leído"
// @Failure 400 {object} apitypes.ErrorResponse "ID inválido"
// @Failure 404 {object} apitypes.ErrorResponse "Mensaje no encontrado"
// @Failure 500 {object} apitypes.ErrorResponse "Error interno"
// @Router /admin/contactos/{id}/leido [patch]
// @Security ApiKeyAuth
func (h *ContactoHandler) MarcarComoLeido(c *gin.Context) {
	idStr := c.Param("id")
	idUint64, errConv := strconv.ParseUint(idStr, 10, 32)
	if errConv != nil {
		_ = c.Error(fmt.Errorf("parámetro ID inválido para marcar como leído: %s - %w", idStr, errConv))
		return
	}
	id := uint(idUint64)

	if err := h.service.MarcarContactoComoLeido(c.Request.Context(), id); err != nil {
		_ = c.Error(err) // El middleware manejará ErrContactoInvalido (como 404)
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensaje": fmt.Sprintf("Mensaje ID %d marcado como leído.", id)})
}