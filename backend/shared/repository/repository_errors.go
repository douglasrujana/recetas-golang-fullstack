// backend/internal/repository/errors.go

// Este archivo define errores comunes que pueden devolver las implementaciones de repositorios.
// Pertenece al paquete de la característica 'recetas'

package repository

import "errors"

// Errores comunes que pueden devolver las implementaciones de repositorios.
var (
	ErrRecordNotFound      = errors.New("registro no encontrado en la base de datos")
	ErrDuplicateRecord     = errors.New("registro duplicado viola restricción única")
	ErrForeignKeyViolation = errors.New("violación de llave foránea")
	// Puedes añadir otros errores específicos de DB si los necesitas mapear.
)