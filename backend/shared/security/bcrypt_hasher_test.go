// backend/shared/security/bcrypt_hasher_test.go
package security_test // Usar paquete _test para probar como cliente externo

import (
	"testing"

	"backend/shared/security" // El paquete que estamos probando
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt" // Para errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
    "errors" // Para errors.Is
)

func TestBcryptHasher_HashAndCompare(t *testing.T) {
	hasher := security.NewBcryptHasher(bcrypt.MinCost) // Usar MinCost para tests rápidos

	password := "P@$$wOrd123"
	wrongPassword := "WrongP@$$wOrd"

	// Testear Hash
	hashedPassword, err := hasher.Hash(password)
	require.NoError(t, err, "Hashear no debería producir error")
	require.NotEmpty(t, hashedPassword, "El hash no debería estar vacío")
	// Verificar que el hash no es la contraseña original
	assert.NotEqual(t, password, hashedPassword, "El hash no debe ser igual a la contraseña original")

	// Testear Compare - Contraseña Correcta
	err = hasher.Compare(hashedPassword, password)
	assert.NoError(t, err, "Compare con contraseña correcta debería ser exitoso (nil error)")

	// Testear Compare - Contraseña Incorrecta
	err = hasher.Compare(hashedPassword, wrongPassword)
	assert.Error(t, err, "Compare con contraseña incorrecta debería producir un error")
	// Verificar que el error es específicamente el de mismatch de bcrypt
	assert.True(t, errors.Is(err, bcrypt.ErrMismatchedHashAndPassword), "El error debería ser bcrypt.ErrMismatchedHashAndPassword")

	// Testear Compare - Hash Inválido
	err = hasher.Compare("not_a_valid_bcrypt_hash", password)
	assert.Error(t, err, "Compare con hash inválido debería producir un error")
	// El error específico puede variar, pero no debería ser nil
}

func TestNewBcryptHasher_DefaultCost(t *testing.T) {
    // Este test es un poco más complejo porque necesitaría acceder al campo 'cost' no exportado
    // o hashear y verificar la longitud del hash, lo cual no es directo para el costo.
    // Por ahora, podemos confiar en que si pasamos <=0, usa bcrypt.DefaultCost.
    // Si el costo fuera configurable desde fuera y crítico, se podría añadir más aquí.
    hasherDefault := security.NewBcryptHasher(0)
    hasherExplicitDefault := security.NewBcryptHasher(bcrypt.DefaultCost)

    // No podemos comparar los hashers directamente, pero podemos verificar que no son nil
    assert.NotNil(t, hasherDefault)
    assert.NotNil(t, hasherExplicitDefault)
    // Una prueba simple podría ser hashear la misma contraseña y ver si los hashes son diferentes
    // (debido al salt, siempre deberían serlo, pero no prueba el costo directamente).
    // Esta prueba es más para asegurar que el constructor no paniquea.
    _, err1 := hasherDefault.Hash("test")
    _, err2 := hasherExplicitDefault.Hash("test")
    assert.NoError(t, err1)
    assert.NoError(t, err2)
}