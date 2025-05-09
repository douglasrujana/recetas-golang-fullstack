// backend/internal/domain/errors.go (o categoria.go)
package categorias

import "errors"

// Errores específicos del dominio de negocio
var (
	ErrCategoriaNotFound       = errors.New("categoría no encontrada")
	ErrCategoriaNombreYaExiste = errors.New("ya existe una categoría con ese nombre")
	// Puedes añadir otros errores de validación de negocio aquí si son necesarios
	// ErrCategoriaNombreInvalido = errors.New("el nombre de la categoría no es válido")
)