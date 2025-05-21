// backend/shared/security/jwt_handler.go
// Funcionalidad: Provee una abstracción para generar y verificar JSON Web Tokens (JWT).
// Capa: Compartida (Utilidad de Seguridad).
package security

import (
	"backend/shared/config" // Para acceder a la SECRET_KEY y duración del token
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5" // La librería JWT
)

// Claims representa los datos personalizados que queremos incluir en el payload del JWT,
// además de los claims estándar (RegisteredClaims).
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	// Nombre string `json:"nombre,omitempty"` // Podrías añadir más
	jwt.RegisteredClaims // Embeber claims estándar
}

// TokenGenerator define el contrato para generar tokens JWT.
type TokenGenerator interface {
	GenerateToken(userID uint, email string) (string, error)
}

// TokenVerifier define el contrato para verificar y parsear tokens JWT.
type TokenVerifier interface {
	VerifyToken(tokenString string) (*Claims, error)
}

// jwtManager implementa ambas interfaces, TokenGenerator y TokenVerifier.
type jwtManager struct {
	secretKey    []byte        // Clave secreta para firmar/verificar (de config)
	tokenExpires time.Duration // Duración de validez del token (de config)
	issuer       string        // Emisor del token (de config, opcional)
}

// NewJWTManager es la factory function para crear una instancia de jwtManager.
// Recibe la configuración JWT (que contendrá la clave secreta, duración, issuer).
func NewJWTManager(cfg config.JWTConfig) (TokenGenerator, TokenVerifier, error) { // Devuelve ambas interfaces
	if len(cfg.SecretKey) == 0 {
		return nil, nil, fmt.Errorf("jwtManager: la clave secreta JWT no puede estar vacía")
	}
	if cfg.TokenExpiresInMinutes <= 0 {
		return nil, nil, fmt.Errorf("jwtManager: la duración de expiración del token debe ser positiva")
	}

	return &jwtManager{
		secretKey:    []byte(cfg.SecretKey),
		tokenExpires: time.Minute * time.Duration(cfg.TokenExpiresInMinutes),
		issuer:       cfg.Issuer,
	}, nil
}

// GenerateToken crea un nuevo token JWT firmado.
func (jm *jwtManager) GenerateToken(userID uint, email string) (string, error) {
	// Definir el tiempo de expiración
	expirationTime := time.Now().Add(jm.tokenExpires)

	// Crear los claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jm.issuer, // Opcional
			// Subject: fmt.Sprintf("%d", userID), // Opcional
		},
	}

	// Crear el token con el algoritmo de firma HS256 y los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(jm.secretKey)
	if err != nil {
		return "", fmt.Errorf("jwtManager: error al firmar el token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken verifica un token string, lo parsea y devuelve los claims si es válido.
func (jm *jwtManager) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parsear el token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Asegurarse de que el algoritmo de firma sea el esperado (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwtManager: algoritmo de firma inesperado: %v", token.Header["alg"])
		}
		return jm.secretKey, nil
	})

	if err != nil {
		// Aquí jwt.ParseWithClaims puede devolver varios tipos de errores:
		// - jwt.ErrTokenMalformed
		// - jwt.ErrTokenExpired
		// - jwt.ErrTokenNotValidYet
		// - errores de firma, etc.
		// El servicio que llame a VerifyToken deberá manejar estos errores adecuadamente.
		return nil, fmt.Errorf("jwtManager: error al parsear/validar token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("jwtManager: token inválido")
	}

	// El struct 'claims' ahora está populado con los datos del token.
	return claims, nil
}

// Necesitaremos añadir JWTConfig al struct Config en shared/config/config.go
// y a los archivos YAML.