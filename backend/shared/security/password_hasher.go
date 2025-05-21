// backend/shared/security/password_hasher.go
// Funcionalidad: Define la abstracción para el hasheo y comparación de contraseñas.
// Capa: Compartida (Utilidad de Seguridad).
// Interfaz: PasswordHasher

// Descripción:
// Esta interfaz proporciona un contrato para cualquier implementación de hasheo de contraseñas.
// Permite desacoplar la lógica de negocio (ej: AuthService) de la implementación específica
// de hasheo (ej: bcrypt, scrypt, etc.), facilitando cambios futuros o pruebas.
//
// Uso:
// - Inyectada en servicios que manejan autenticación o registro de usuarios.
//
// Responsabilidades:
// - Definir cómo se hashea una contraseña en texto plano.
// - Definir cómo se compara una contraseña en texto plano con un hash existente.
package security

// PasswordHasher define los métodos para el hasheo y comparación de contraseñas.
type PasswordHasher interface {
	// Hash toma una contraseña en texto plano y devuelve su representación hasheada.
	// También puede devolver un error si el proceso de hasheo falla.
	Hash(password string) (string, error)

	// Compare compara una contraseña en texto plano con un hash previamente generado.
	// Devuelve nil si coinciden, o un error específico (ej: bcrypt.ErrMismatchedHashAndPassword)
	// si no coinciden, o cualquier otro error del proceso de comparación.
	Compare(hashedPassword string, password string) error
}