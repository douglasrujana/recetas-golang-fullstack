package utils

import (
    "fmt"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "github.com/gin-gonic/gin"
)

// buscarArchivoEnUploads recorre recursivamente el directorio "uploads/" para encontrar el archivo.
// Devuelve la ruta relativa desde uploads/ si lo encuentra, o error si no.
func buscarArchivoEnUploads(nombreArchivo string) (string, error) {
    baseDir := "uploads"
    var rutaRelativa string
    encontrado := false

    err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && info.Name() == nombreArchivo {
            // Devuelve la ruta relativa desde "uploads/"
            rutaRelativa = strings.TrimPrefix(path, baseDir+string(filepath.Separator))
            encontrado = true
            return filepath.SkipDir // Detener búsqueda
        }

        return nil
    })

    if err != nil {
        return "", fmt.Errorf("error al buscar archivo en %s: %w", baseDir, err)
    }

    if !encontrado {
        return "", fmt.Errorf("archivo '%s' no encontrado en '%s'", nombreArchivo, baseDir)
    }

    return rutaRelativa, nil
}

// ProcesarFotoURL convierte un nombre de imagen a una URL pública accesible desde el navegador.
// Si la imagen no existe, retorna una imagen por defecto.
func ProcesarFotoURL(c *gin.Context, foto string) string {
    // Si ya es una URL absoluta, no procesar
    parsedURL, err := url.Parse(foto)
    if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
        return foto
    }

    // Construye el esquema base (http o https)
    scheme := "http"
    if c.Request.TLS != nil {
        scheme = "https"
    }
    baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

    // Buscar el archivo real dentro de /uploads/
    rutaRelativa, err := buscarArchivoEnUploads(foto)
    if err != nil {
        fmt.Printf("⚠️ Imagen no encontrada: %s. Usando imagen por defecto.\n", foto)
        rutaRelativa = "default.png"
    }

    // Reemplaza las barras invertidas por barras normales para URLs válidas
    rutaRelativa = filepath.ToSlash(rutaRelativa)
    // Retorna la URL absoluta final
    return fmt.Sprintf("%s/uploads/%s", baseURL, rutaRelativa)
}
