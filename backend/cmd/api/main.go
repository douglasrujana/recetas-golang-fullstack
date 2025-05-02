// backend/cmd/api/main.go

package main

import (
	"fmt"
	"log"
	"net/http" // Necesario para router.GET("/") simple

	"backend/internal/config" // Configuración (¡Bien!)
	// --- [✨ NUEVA ARQUITECTURA] ---
	"backend/internal/database"  // Nuestro nuevo paquete para la conexión DB
	"backend/internal/repository/mysql" // Implementaciones del repositorio MySQL/Gorm
	"backend/internal/service"   // Lógica de negocio
	"backend/internal/handler"   // Handlers HTTP (el antiguo 'rutas')
	// ---------------------------------

	// "backend/internal/domain" // Domain usualmente no se importa en main, sino en repo/service/handler

	"github.com/gin-gonic/gin"
)

func main() {
	// --- 1. Carga de Configuración ---
	config.CargarVariablesEntorno()
	cfg := config.AppConfig // Accedemos a la configuración cargada
	fmt.Println("Corriendo en entorno:", cfg.AppEnv)
	// [⚠️ VALIDACIÓN] Asegúrate que el puerto DB_PORT en .env (5432) es correcto para MySQL o cámbialo.

	// --- 2. Inicialización de Dependencias Core (Base de Datos) ---
	// [✨ ELEGANCIA] Llamamos a la función del paquete 'database'.
	dbInstance, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error fatal al inicializar la base de datos: %v", err)
	}
	// sqlDB, _ := dbInstance.DB() // Podrías querer cerrar el pool al final
	// defer sqlDB.Close() // Cierre ordenado al terminar main (opcional pero buena práctica)

	// --- 3. Inyección de Dependencias (¡La Magia!) ---
	// Aquí creamos las instancias de cada capa, inyectando las dependencias necesarias.
	// El orden es importante: Repos -> Services -> Handlers

	// [🏗️ REPOSITORIOS] Crear instancias de los repositorios, pasando la conexión DB.
	// (Necesitarás crear estas funciones constructoras `New...Repository` en `internal/repository/mysql/`)
	userRepo := mysql.NewUserRepository(dbInstance)
	categoriaRepo := mysql.NewCategoriaRepository(dbInstance)
	recetaRepo := mysql.NewRecetaRepository(dbInstance)
	// ... otros repositorios si los hubiera ...

	// [🏗️ SERVICIOS] Crear instancias de los servicios, pasando los repositorios.
	// (Necesitarás crear estas funciones constructoras `New...Service` en `internal/service/`)
	authService := service.NewAuthService(userRepo, cfg.SecretKey) // Auth necesita el repo de user y la secret key
	categoriaService := service.NewCategoriaService(categoriaRepo)
	recetaService := service.NewRecetaService(recetaRepo)
	// ... otros servicios ...

	// [🏗️ HANDLERS] Crear instancias de los handlers, pasando los servicios.
	// (Necesitarás crear estas funciones constructoras `New...Handler` en `internal/handler/`)
	authHandler := handler.NewAuthHandler(authService)
	categoriaHandler := handler.NewCategoriaHandler(categoriaService)
	recetaHandler := handler.NewRecetaHandler(recetaService)
	uploadHandler := handler.NewUploadHandler() // Asumiendo que no necesita servicios por ahora
	// ... otros handlers ...

	// --- 4. Inicialización del Router Gin ---
	gin.SetMode(gin.ReleaseMode) // O gin.DebugMode
	router := gin.Default()      // Incluye logger y recovery middleware

	// --- 5. Configuración de Rutas (Delegada) ---
	// [✨ ELEGANCIA] Toda la definición de rutas se mueve a una función dedicada en el paquete handler.
	// Le pasamos el router y las instancias de los handlers que necesita.
	// (Necesitarás crear esta función `RegisterRoutes` en `internal/handler/routes.go`)
	handler.RegisterRoutes(router, authHandler, categoriaHandler, recetaHandler, uploadHandler)

	// [✅ BUENA PRÁCTICA] Configurar rutas estáticas y 404 sigue estando bien aquí o en RegisterRoutes.
	router.Static("/public", "./public")
	router.Static("/uploads", "./uploads")
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"}) // Usar http status codes
	})
	// Simple health check/root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API Recetas Go Funcionando!"})
	})

	// --- 6. Arranque del Servidor ---
	fmt.Printf("🚀 Servidor escuchando en http://localhost:%s\n", cfg.AppPort)
	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		panic(fmt.Sprintf("❌ Error al iniciar el servidor Gin: %v", err))
	}
}