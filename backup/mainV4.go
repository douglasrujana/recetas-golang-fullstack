// backend/cmd/app/main.go

// Punto de entrada de la aplicación.
//
// Este archivo contiene la lógica principal de la aplicación, que incluye:
// - Carga de configuración
// - Inicialización de base de datos
// - Inyección de dependencias
// - Configuración del router
// - Registro de rutas
// - Arranque del servidor
//
// Nota: Se recomienda mantener este archivo lo más simple posible para facilitar
// la comprensión y mantenimiento de la aplicación.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	// --- Paquetes Compartidos ---
	"backend/shared/config"
	"backend/shared/database"
	"backend/shared/middleware"

	// --- Paquetes de Características ---
	"backend/categorias" // Asegúrate de que esta línea esté presente y correcta
	"backend/recetas"    // <-- ¡NUEVO! Importar el paquete de recetas

	// --- Paquetes de Swagger ---
	_ "backend/docs" // <-- ¡NUEVO! Importar el paquete de docs

	swaggerFiles "github.com/swaggo/files"     // <-- ¡NUEVO! Importar el paquete de archivos de Swagger
	ginSwagger "github.com/swaggo/gin-swagger" // <-- ¡NUEVO! Importar el paquete de Swagger para Gin
)

func main() {
	log.Println("🚀 Iniciando aplicación...")

	// --- 1. Carga de Configuración ---
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("❌ Error cargando config: %v", err)
	}
	fmt.Println("✅ Configuración cargada para entorno:", cfg.AppEnv)

	// --- 2. Inicialización de Base de Datos ---
	dbInstance, err := database.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Error inicializando BD (%s): %v", cfg.Database.Name, err)
	}
	log.Printf("✅ Conexión a base de datos '%s' establecida.\n", cfg.Database.Name)

	// --- 3. AutoMigrate ---
	log.Println("🔧 Ejecutando AutoMigrate...")
	err = dbInstance.AutoMigrate(
		&categorias.CategoriaModel{},
		&recetas.RecetaModel{}, // <-- ¡NUEVO! Añadir RecetaModel a AutoMigrate
		// &mysql.IngredienteModel{}, // Cuando exista
		// &mysql.UserModel{},        // Cuando exista
	)

	if err != nil {
		log.Fatalf("❌ ERROR CRÍTICO durante AutoMigrate: %v", err)
	}
	log.Println("✅ AutoMigrate completado.")

	// --- 4. Inyección de Dependencias ---
	log.Println("🏗️  Inicializando dependencias...")

	// Categorías (sin cambios)
	categoriaRepo := categorias.NewCategoriaRepository(dbInstance)
	categoriaService := categorias.NewCategoriaService(categoriaRepo)
	categoriaHandler := categorias.NewCategoriaHandler(categoriaService)
	log.Println("   - Dependencias de Categorías inicializadas.")

	// Recetas <-- ¡NUEVO!
	recetaRepo := recetas.NewRecetaRepository(dbInstance) // Usa constructor del paquete 'recetas'
	// RecetaService necesita CategoriaService para validar CategoriaID
	recetaService := recetas.NewRecetaService(recetaRepo, categoriaService) // Inyectar categoriaService
	recetaHandler := recetas.NewRecetaHandler(recetaService)
	log.Println("   - Dependencias de Recetas inicializadas.")
	log.Println("✅ Todas las dependencias inicializadas.")

	// --- 5. Inicialización del Router Gin ---
	if cfg.AppEnv != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	//router := gin.Default()
	router := gin.New()        // Usar gin.New() para control total
	router.Use(gin.Logger())   // Añadir logger de Gin
	router.Use(gin.Recovery()) // Añadir middleware de recuperación de panics

	// --- ✨ REGISTRAR MIDDLEWARE DE ERRORES GLOBALMENTE ---
	// Debe ir DESPUÉS de Logger y Recovery, pero ANTES de que se definan las rutas
	// para que pueda capturar errores de todos los handlers.
	router.Use(middleware.ErrorHandler())
	log.Printf("✅ Router Gin inicializado en modo: %s.\n", gin.Mode())

	// --- 6. Configuración de Rutas ---
	apiV1 := router.Group("/api/v1")

	if categoriaHandler != nil {
		categorias.RegisterCategoriaRoutes(apiV1, categoriaHandler) // Usa el nombre de función que definiste en el paquete
	}

	// Recetas <-- ¡NUEVO!
	if recetaHandler != nil { // <-- ¡NUEVO!
		recetas.RegisterRoutes(apiV1, recetaHandler) // Usa el nombre de función que definiste en el paquete
	}
	log.Println("✅ Todas las rutas de API registradas.")

	// --- ✨ ENDPOINT PARA SWAGGER UI ✨ ---
	// Asegúrate de que esta línea esté presente y sea accesible
	// Se registra en la raíz del router, NO bajo /api/v1
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("✅ Swagger UI disponible en /swagger/index.html")
	// -----------------------------------------------------------------------

	// --- Rutas Base / No API (sin cambios) ---
	router.Static("/public", "./public")
	router.Static("/uploads", "./uploads")
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "✅ API Recetas Go Funcionando!",
		})
	})
	log.Println("✅ Rutas adicionales configuradas.")

	// --- 7. Arranque del Servidor (sin cambios) ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("🚀 Servidor escuchando en http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor Gin: %v", err)
	}
}
