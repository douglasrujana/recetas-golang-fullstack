// backend/internal/handler/auth_handler.go (O en un archivo types.go/claims.go)

package handler

import (
	// ... otros imports ...
	"github.com/golang-jwt/jwt/v5" // Usaremos la versión v5, la más reciente
)

// Claims define la estructura de los datos que almacenamos en el token JWT.
// [✅ BUENA PRÁCTICA] Incrustar jwt.RegisteredClaims para incluir claims estándar (exp, iat, etc.)
type Claims struct {
	UserID   string `json:"user_id"` // Identificador único del usuario
	Username string `json:"username"` // Nombre de usuario (opcional, pero útil)
	// Puedes añadir otros campos si los necesitas, como roles:
	// Role string `json:"role"`
	jwt.RegisteredClaims // Incluye campos estándar como ExpiresAt, IssuedAt, etc.
}

// [✅ BUENA PRÁCTICA] Usar una constante para la clave del contexto Gin
// evita errores tipográficos al guardar/recuperar el valor.
const GinContextKeyUserClaims = "userClaims"