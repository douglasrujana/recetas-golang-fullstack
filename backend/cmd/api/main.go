package main

import (
	"fmt"
	"backend/internal/config" // Importa el paquete de configuración"
	"backend/internal/handler" // Importa el paquete de rutas
	"backend/internal/domain" // Importar los modelos de BD
	"backend/internal/repository/mysql" // Importa el paquete de repositorio
	"github.com/gin-gonic/gin"
)

func main() {
	//Cargar archivo de configutación
	config.CargarVariablesEntorno()
	fmt.Println("Corriendo en entorno:", config.AppConfig.AppEnv)

	//SEt Mode
	gin.SetMode(gin.ReleaseMode) // Establece el modo de Gin a ReleaseMode
	// gin.SetMode(gin.DebugMode) // Establece el modo de Gin a DebugMode

	//Inicializar Gin
	router := gin.Default()

	//Error 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Ruta no encontrada"})
	})

	//Ejecuatr las migraiones
	modelos.Migraciones() // Llama a la función de migraciones

    //listen and serve on 0.0.0.0:8085 (for windows "localhost:8085")
	fmt.Printf("🚀 Servidor escuchando en el puerto %s\n", config.AppConfig.AppPort)
	err := router.Run(":" + config.AppConfig.AppPort)
	if err != nil {
		// [✅ BUENA PRÁCTICA] Manejar el error si el servidor no puede iniciar. 👍
		panic(fmt.Sprintf("Error al iniciar el servidor Gin: %v", err))
	}
}
