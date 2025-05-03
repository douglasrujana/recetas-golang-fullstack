// backend/cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http" // Necesario para router.GET("/") simple

	"backend/internal/config" // Nuestro paquete de configuración con Viper
	"backend/internal/database"
	"backend/internal/handler"
	"github.com/gin-gonic/gin"
	// Asegúrate que el driver de MySQL está importado, usualmente en database.go o aquí con _
	// _ "github.com/go-sql-driver/mysql" // Si no está en database.go
)

func main() {
	// --- 1. Carga de Configuración con Viper ---
	// [✨ REFACTOR] Usamos LoadConfig que lee ./config.yaml y/o variables de entorno
	cfg, err := config.LoadConfig(".") // Busca config.yaml en el directorio actual (.)
	if err != nil {
		log.Fatalf("❌ Error fatal al cargar la configuración: %v", err)
	}
	// [✅ BUENA PRÁCTICA] Loguear el entorno, pero NUNCA secretos.
	fmt.Println("✅ Configuración cargada.")
	fmt.Println("   - Entorno de aplicación (APP_ENV):", cfg.AppEnv)
	fmt.Println("   - Puerto del servidor (SERVER_PORT):", cfg.Server.Port)
	fmt.Println("   - Host de BD (DATABASE_HOST):", cfg.Database.Host)
	fmt.Println("   - Puerto de BD (DATABASE_PORT):", cfg.Database.Port)
	fmt.Println("   - Nombre de BD (DATABASE_NAME):", cfg.Database.Name)
	// [⚠️ ADVERTENCIA] ¡No loguees cfg.SecretKey ni cfg.Database.Password en producción!

	// [⚠️ VALIDACIÓN] El puerto DB_PORT (ahora `cfg.Database.Port`) debe ser el correcto para tu MySQL.

	// --- 2. Inicialización de Dependencias Core (Base de Datos) ---
	// [✨ REFACTOR] Pasamos la sub-estructura `cfg.Database` a ConnectDB.
	//    Asegúrate que `database.ConnectDB` ahora acepte `config.DatabaseConfig`.
	//    Si `ConnectDB` espera un DSN, puedes generarlo aquí: dsn := cfg.Database.DSN()
	dbInstance, err := database.ConnectDB(cfg.Database) // Ajusta ConnectDB si es necesario
	if err != nil {
		log.Fatalf("❌ Error fatal al inicializar la base de datos: %v", err)
	}
	log.Println("✅ Conexión a base de datos establecida.")
	// sqlDB, err := dbInstance.DB() // Si usas GORM o sqlx, el método puede variar
	// if err == nil {
	// 	 defer sqlDB.Close() // Cierre ordenado
	// }

	// --- 3. Inyección de Dependencias (¡Se mantiene igual, excelente!) ---
	// Solo cambiamos de dónde viene el valor de `cfg.SecretKey`

	// [🏗️ REPOSITORIOS]
	userRepo := mysql.NewUserRepository(dbInstance)
	categoriaRepo := mysql.NewCategoriaRepository(dbInstance)
	recetaRepo := mysql.NewRecetaRepository(dbInstance)
	// ... otros repositorios ...
	log.Println("✅ Repositorios inicializados.")

	// [🏗️ SERVICIOS]
	// [✨ REFACTOR] Usamos cfg.SecretKey obtenido de Viper.
	authService := service.NewAuthService(userRepo, cfg.SecretKey)
	categoriaService := service.NewCategoriaService(categoriaRepo)
	recetaService := service.NewRecetaService(recetaRepo)
	// ... otros servicios ...
	log.Println("✅ Servicios inicializados.")

	// [🏗️ HANDLERS]
	authHandler := handler.NewAuthHandler(authService)
	categoriaHandler := handler.NewCategoriaHandler(categoriaService)
	recetaHandler := handler.NewRecetaHandler(recetaService)
	uploadHandler := handler.NewUploadHandler() // Asumiendo sin dependencias de servicio
	// ... otros handlers ...
	log.Println("✅ Handlers inicializados.")

	// --- 4. Inicialización del Router Gin ---
	// [✨ REFACTOR] Usamos cfg.AppEnv para decidir el modo de Gin.
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode) // O déjalo por defecto que es debug
	}
	router := gin.Default() // Incluye logger y recovery middleware
	log.Printf("✅ Router Gin inicializado en modo: %s.\n", gin.Mode())


	// --- 5. Configuración de Rutas (Delegada - ¡Se mantiene igual!) ---
	handler.RegisterRoutes(router, authHandler, categoriaHandler, recetaHandler, uploadHandler)
	log.Println("✅ Rutas de la API registradas.")

	// [✅ BUENA PRÁCTICA] Rutas estáticas y 404
	router.Static("/public", "./public") // Asegúrate que esta ruta es correcta
	router.Static("/uploads", "./uploads") // Asegúrate que esta ruta es correcta
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
	// Health check/root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "✅ API Recetas Go Funcionando!"})
	})
	log.Println("✅ Rutas estáticas y de health check configuradas.")

	// --- 6. Arranque del Servidor ---
	// [✨ REFACTOR] Usamos cfg.Server.Port obtenido de Viper (que es un int).
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("🚀 Servidor escuchando en http://localhost%s\n", serverAddr)

	// Usamos router.Run que maneja errores internamente y loguea.
	// Si falla al iniciar, logueará fatalmente.
	if err := router.Run(serverAddr); err != nil {
		// Este log es un poco redundante ya que Gin suele loguear fatal, pero no hace daño.
		log.Fatalf("❌ Error al iniciar el servidor Gin en %s: %v", serverAddr, err)
	}
}