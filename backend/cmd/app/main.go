// backend/cmd/api/main.go (Enfoque: Solo Categorías)
package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/config"   // Configuración
	"backend/internal/database"  // Conexión BD
	// "backend/internal/domain" // Puede no ser necesario importar aquí directamente
	"backend/internal/handler"   // Handlers (el paquete)
	// "backend/internal/repository" // Puede no ser necesario importar aquí directamente
	"backend/internal/repository/mysql" // Implementación MySQL del Repositorio
	"backend/internal/service"   // Servicios (el paquete)

	"github.com/gin-gonic/gin"
	// Asegúrate de que el driver MySQL esté importado (usualmente en database.go)
	// _ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("🚀 Iniciando aplicación (Enfoque: Solo Categorías)...")

	// --- 1. Carga de Configuración ---
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ Error cargando config: %v", err)
	}
	fmt.Println("✅ Configuración cargada.")
	fmt.Println("   - Entorno:", cfg.AppEnv)
	fmt.Println("   - Puerto Servidor:", cfg.Server.Port)
	// ... otros logs de config seguros ...

	// --- 2. Inicialización de Base de Datos ---
	dbInstance, err := database.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Error inicializando BD: %v", err)
	}
	log.Println("✅ Conexión a base de datos establecida.")

	// --- 3. Inyección de Dependencias (Solo Categorías) ---
	log.Println("🏗️  Inicializando dependencias (Solo Categorías)...")

	// [REPOSITORIOS]
	categoriaRepo := mysql.NewCategoriaRepository(dbInstance)
	log.Println("   - Repositorio Categorías inicializado.")

	// [SERVICIOS]
	categoriaService := service.NewCategoriaService(categoriaRepo)
	log.Println("   - Servicio Categorías inicializado.")

	// [HANDLERS]
	categoriaHandler := handler.NewCategoriaHandler(categoriaService)
	log.Println("   - Handler Categorías inicializado.")
	log.Println("✅ Dependencias de Categorías inicializadas.")


	// --- 4. Inicialización del Router Gin ---
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	log.Printf("✅ Router Gin inicializado en modo: %s.\n", gin.Mode())


	// --- 5. Configuración de Rutas (Solo Categorías) ---
	// Pasamos SOLO el handler de categorías.
	handler.RegisterRoutes(router, categoriaHandler)
	log.Println("✅ Rutas de Categorías registradas.")

	// [Rutas Base / No API]
	router.Static("/public", "./public")
	router.Static("/uploads", "./uploads") // Dejar por ahora, puede ser necesario
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "✅ API Recetas Go Funcionando! (Enfoque Categorías)"})
	})
	log.Println("✅ Rutas adicionales configuradas.")


	// --- 6. Arranque del Servidor ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("🚀 Servidor (Categorías) escuchando en http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor Gin: %v", err)
	}
}