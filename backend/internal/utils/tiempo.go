package utils

import (
    "fmt"
    "time"
)

// Formatea una fecha en el formato DD/MM/YYYY
func FormatearFecha(fecha time.Time) string {
    return fmt.Sprintf("%02d/%02d/%d", fecha.Day(), fecha.Month(), fecha.Year())
}
