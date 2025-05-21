// backend/cmd/api/main.go
// Punto de entrada principal de la aplicación API.
// Responsable de:
// - Cargar configuración.
// - Establecer conexión con la base de datos.
// - Ejecutar migraciones automáticas (para desarrollo/test).
// - Inicializar e inyectar todas las dependencias (repositorios, servicios, handlers, notificadores).
// - Configurar el router Gin y registrar todas las rutas.
// - Arrancar el servidor HTTP.
package main

import (
	// --- Paquetes Estándar de Go ---
	"fmt"
	"log"      // Para logging inicial y errores fatales
	"net/http" // Para http.StatusNotFound y http.StatusOK

	// --- Paquetes de Terceros ---
	"github.com/gin-gonic/gin"         // El framework web Gin
	_ "github.com/go-sql-driver/mysql" // Driver MySQL (importado por sus efectos secundarios)

	// --- Paquetes Internos del Proyecto (Nueva Estructura "Paquete por Característica" y "Shared") ---
	_ "backend/docs" // Paquete generado por Swagger para la documentación de la API (importado por efectos secundarios)

	"backend/categorias" // Paquete para la característica/dominio de Categorías
	"backend/contactos"  // Paquete para la característica/dominio de Contactos
	"backend/recetas"    // Paquete para la característica/dominio de Recetas
	// "backend/auth"       // Paquete para Autenticación (cuando se implemente)

	"backend/shared/config"       // Paquete compartido para la configuración de la aplicación
	"backend/shared/database"     // Paquete compartido para la conexión a la base de datos
	"backend/shared/middleware"   // Paquete compartido para middlewares (ej: ErrorHandler)
	"backend/shared/notifications" // Paquete compartido para notificaciones (ej: EmailNotifier)

	// Paquetes de Swagger (si no los has importado en otro lado y los necesitas aquí)
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Recetas Fullstack
// @version 1.0
// @description Esta es una API para gestionar recetas, categorías, contactos e ingredientes.
// @termsOfService http://swagger.io/terms/

// @contact.name Soporte API Douglas Rujana
// @contact.url http://www.douglasrujana.com/support
// @contact.email douglasrujana@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Escribe "Bearer " seguido de un espacio y tu token JWT. (Ej: Bearer eyJhbGciOi...)
func main() {
	log.Println("🚀 Iniciando aplicación (Estructura Paquete por Característica)...")

	// --- 1. Carga de Configuración ---
	cfg, err := config.LoadConfig("config") // "config" es la carpeta de archivos YAML
	if err != nil {
		log.Fatalf("❌ Error fatal al cargar la configuración: %v", err)
	}
	log.Println("✅ Configuración cargada exitosamente para el entorno:", cfg.AppEnv)
	log.Printf("   - Puerto del Servidor: %d\n", cfg.Server.Port)
	log.Printf("   - Host de Base de Datos: %s\n", cfg.Database.Host)
	log.Printf("   - Nombre de Base de Datos: %s\n", cfg.Database.Name)
	log.Printf("   - SMTP Host: %s (Puerto: %d, From: %s, AdminTo: %s)\n", cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.From, cfg.SMTP.AdminTo)

	// --- 2. Inicialización de Base de Datos ---
	dbInstance, err := database.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Error fatal al inicializar la base de datos (%s): %v", cfg.Database.Name, err)
	}
	log.Printf("✅ Conexión a base de datos '%s' establecida.\n", cfg.Database.Name)

	// --- 3. Ejecutar AutoMigrate (para Desarrollo y Tests) ---
	log.Println("🔧 Ejecutando AutoMigrate para todas las entidades GORM...")
	err = dbInstance.AutoMigrate(
		&categorias.CategoriaModel{},
		&recetas.RecetaModel{},
		&contactos.ContactoModel{}, // Añadido modelo de Contactos
		// &auth.UserModel{},
		// ...otros *Model GORM aquí...
	)
	if err != nil {
		log.Fatalf("❌ ERROR CRÍTICO durante AutoMigrate: %v", err)
	}
	log.Println("✅ AutoMigrate completado.")

	// --- 4. Inyección de Dependencias ---
	log.Println("🏗️  Inicializando dependencias de la aplicación...")

	// Componentes Compartidos
	emailNotifier, err := notifications.NewSMTPNotifier(cfg.SMTP)
	if err != nil {
		log.Fatalf("❌ ERROR CRÍTICO al crear el notificador de email: %v", err)
	}
	log.Println("   - Notificador de Email (SMTP) inicializado.")

	// Dependencias de Categorías
	categoriaRepo := categorias.NewCategoriaRepository(dbInstance)
	categoriaService := categorias.NewCategoriaService(categoriaRepo)
	categoriaHandler := categorias.NewCategoriaHandler(categoriaService)
	log.Println("   - Dependencias de 'Categorías' inicializadas.")

	// Dependencias de Recetas
	recetaRepo := recetas.NewRecetaRepository(dbInstance)
	recetaService := recetas.NewRecetaService(recetaRepo, categoriaService) // RecetaService depende de CategoriaService
	recetaHandler := recetas.NewRecetaHandler(recetaService)
	log.Println("   - Dependencias de 'Recetas' inicializadas.")

	// Dependencias de Contactos
	contactoRepo := contactos.NewContactoRepository(dbInstance)
	// ContactoService necesita el notificador y los emails de admin/from de la config
	contactoService := contactos.NewContactoService(contactoRepo, emailNotifier, cfg.SMTP.AdminTo, cfg.SMTP.From)
	contactoHandler := contactos.NewContactoHandler(contactoService)
	log.Println("   - Dependencias de 'Contactos' inicializadas.")

	log.Println("✅ Todas las dependencias necesarias inicializadas.")

	// --- 5. Inicialización del Router Gin ---
	if cfg.AppEnv != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler()) // Nuestro middleware de errores global
	log.Println("✅ Middlewares globales (Logger, Recovery, ErrorHandler) registrados.")
	log.Printf("✅ Router Gin inicializado en modo: %s.\n", gin.Mode())

	// --- 6. Configuración de Rutas ---
	apiV1 := router.Group("/api/v1") // Grupo base para la API versionada

	if categoriaHandler != nil {
		categorias.RegisterCategoriaRoutes(apiV1, categoriaHandler)
	}
	if recetaHandler != nil {
		recetas.RegisterRecetaRoutes(apiV1, recetaHandler)
	}
	if contactoHandler != nil {
		contactos.RegisterContactoRoutes(apiV1, contactoHandler) // Registrar rutas de contactos
	}
	log.Println("✅ Rutas de API de características registradas.")

	// Endpoint para Swagger UI
	// La URL será http://localhost:PORT/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("✅ Swagger UI disponible en /swagger/index.html")

	// Rutas Base / Estáticas / No API
	router.Static("/public", "./public")
	router.Static("/uploads", "./uploads")
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "✅ API Recetas Go Funcionando! (Paquete por Característica)"})
	})
	log.Println("✅ Rutas base y de fallback configuradas.")

	// --- 7. Arranque del Servidor HTTP ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("🚀 Servidor escuchando en http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error fatal al iniciar el servidor Gin: %v", err)
	}
}