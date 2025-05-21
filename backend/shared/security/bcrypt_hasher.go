// backend/shared/security/bcrypt_hasher.go
// Funcionalidad: Implementación de PasswordHasher usando el algoritmo bcrypt.
// Capa: Compartida (Utilidad de Seguridad - Implementación).

// Descripción:
// Proporciona una forma segura de hashear y comparar contraseñas usando bcrypt,
// que es un algoritmo de hasheo adaptativo y resistente a ataques de fuerza bruta.
//
// Uso:
// - Se crea una instancia de este hasher (usando NewBcryptHasher) y se inyecta
//   donde se necesite la interfaz PasswordHasher.
//
// Referencias:
// - Paquete bcrypt: golang.org/x/crypto/bcrypt
package security

import (
	"errors" // Para errors.Is
	"fmt"    // Para envolver errores
	"golang.org/x/crypto/bcrypt"
)

// bcryptHasher implementa la interfaz PasswordHasher usando bcrypt.
type bcryptHasher struct {
	// cost es el costo computacional del hasheo. Un valor más alto es más seguro
	// pero más lento. bcrypt.DefaultCost (actualmente 10) es un buen punto de partida.
	cost int
}

// NewBcryptHasher es la factory function para crear un nuevo bcryptHasher.
// Se puede pasar un costo; si es <= 0, se usa bcrypt.DefaultCost.
func NewBcryptHasher(cost int) PasswordHasher { // Devuelve la interfaz
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

// Hash genera un hash bcrypt para la contraseña dada.
func (h *bcryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("bcryptHasher: error generando hash: %w", err)
	}
	return string(hashedBytes), nil
}

// Compare compara un hash bcrypt con una contraseña en texto plano.
func (h *bcryptHasher) Compare(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Devolver el error específico de bcrypt si es un mismatch,
		// o envolverlo para otros tipos de errores (ej: hash inválido).
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return err // El servicio que llama puede usar errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
		}
		return fmt.Errorf("bcryptHasher: error comparando hash: %w", err)
	}
	return nil // Coinciden
}