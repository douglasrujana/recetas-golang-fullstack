package rutas_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"backend/rutas" // ⚠️ ajusta el import a tu ruta real
)

// TestEjemploUpload_SubidaExitosa 🔥
// Simula la subida exitosa de un archivo PNG y valida que el handler lo procese bien
func TestEjemploUpload_SubidaExitosa(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crea un buffer para construir el cuerpo del request tipo multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// ⚙️ Simula un archivo temporal de imagen (.png)
	testFilePath := "./testdata/test-image.png" // crea este archivo de prueba
	fileWriter, err := writer.CreateFormFile("file", filepath.Base(testFilePath))
	assert.NoError(t, err)

	file, err := os.Open(testFilePath)
	assert.NoError(t, err)
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	assert.NoError(t, err)

	// Finaliza la construcción del cuerpo
	writer.Close()

	// 🚀 Crea la solicitud POST simulada
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 🔧 Crea el contexto y response recorder
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// 🧪 Ejecuta el handler
	rutas.EjemploUpload(ctx)

	// ✅ Validaciones
	assert.Equal(t, http.StatusOK, w.Code)

	expectedMsg := "Archivo subido con éxito"
	assert.Contains(t, w.Body.String(), expectedMsg)

	t.Logf("🔥 Test exitoso: %s", w.Body.String())
}
