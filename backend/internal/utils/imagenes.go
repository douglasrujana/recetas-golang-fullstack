package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GuardarImagen guarda un archivo de imagen en una subcarpeta dentro de /uploads/
// Valida tamaño y formato. Devuelve ruta relativa y error.
func GuardarImagen(c *gin.Context, file *multipart.FileHeader, carpeta string) (string, error) {
	//Validar  mimetype
	contentType := file.Header.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(contentType, "image/") {
    	return "", fmt.Errorf("el archivo no es una imagen")
	}

	// Validaciones básicas
	if file.Size == 0 || file.Size > 2<<20 {
		return "", fmt.Errorf("la imagen debe pesar entre 1 byte y 2MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" {
		return "", fmt.Errorf("formato de imagen no permitido (solo .jpg, .png)")
	}

	// Generar nombre único y ruta de guardado
	filename := fmt.Sprintf("%s%s", time.Now().Format("20060102_150405"), ext)
	rutaRelativa := filepath.Join("uploads", carpeta, filename)

	// Guardar archivo
	if err := c.SaveUploadedFile(file, rutaRelativa); err != nil {
		return "", fmt.Errorf("no se pudo guardar la imagen: %w", err)
	}

	return rutaRelativa, nil
}

// Call
// rutaRelativa, err := utils.GuardarImagen(c, file, "usuarios")
