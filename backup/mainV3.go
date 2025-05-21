// backend/cmd/api/main.go
package main

import (
	// --- Paquetes Estándar de Go ---
	"fmt"
	"log"
	"net/http" // Para http.StatusNotFound y http.StatusOK

	// --- Paquetes de Terceros ---
	"github.com/gin-gonic/gin"         // El framework web Gin
	_ "github.com/go-sql-driver/mysql" // Importar el driver MySQL por sus efectos secundarios (registra el driver)

	// --- Paquetes Internos del Proyecto (Nueva Estructura "Paquete por Característica") ---
	// Asumiendo que tu go.mod define el módulo como "backend"
	// Si es github.com/nombreusuario/proyecto/backend, ajusta estos paths.
	// "backend/recetas"         // Ejemplo: Paquete para Recetas (cuando se migre)
	// "backend/auth"            // Ejemplo: Paquete para Autenticación (cuando se migre)
	// ...otros paquetes de características...
	"backend/categorias"      // Paquete para la característica/dominio de Categorías
	"backend/shared/config"   // Paquete compartido para la configuración
	"backend/shared/database" // Paquete compartido para la conexión a la base de datos
	// "backend/shared/middleware" // Ejemplo: Paquete para middleware global (cuando se cree)
)

func main() {
	log.Println("🚀 Iniciando aplicación (Estructura Paquete por Característica)...")

	// --- 1. Carga de Configuración ---
	// 'config' es el nombre de la carpeta donde están los archivos config.yaml, config.test.yaml, etc.
	// LoadConfig la buscará relativa al directorio de ejecución del binario o 'go run'.
	cfg, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalf("❌ Error cargando config: %v", err)
	}
	fmt.Println("✅ Configuración cargada para entorno:", cfg.AppEnv)
	fmt.Println("   - Puerto Servidor:", cfg.Server.Port)
	fmt.Println("   - Host BD:", cfg.Database.Host)
	// ... (otros logs seguros de la configuración) ...

	// --- 2. Inicialización de Base de Datos ---
	dbInstance, err := database.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Error inicializando BD (%s): %v", cfg.Database.Name, err)
	}
	log.Printf("✅ Conexión a base de datos '%s' establecida.\n", cfg.Database.Name)

	// --- 3. EJECUTAR AutoMigrate SIEMPRE al inicio ---
	// Esto crea/actualiza las tablas en la BD conectada
	// basándose en los structs GORM (*Model) que le pasemos.
	log.Println("🔧 Ejecutando AutoMigrate...")
	err = dbInstance.AutoMigrate(
		// Los modelos GORM ahora están DENTRO de sus paquetes de característica
		// y deben ser EXPORTADOS (empezar con Mayúscula) desde ese paquete.
		&categorias.CategoriaModel{}, // Modelo GORM del paquete 'categorias'

		// --- AÑADE AQUÍ LOS OTROS MODELOS GORM A MEDIDA QUE LOS CREES ---
		// &recetas.RecetaModel{},     // Cuando exista el paquete 'recetas' y su modelo
		// &auth.UserModel{},        // Cuando exista el paquete 'auth' y su modelo
		// -----------------------------------------------------------------
	)
	if err != nil {
		// Si AutoMigrate falla, la aplicación no puede continuar porque las tablas son incorrectas.
		log.Fatalf("❌ ERROR CRÍTICO durante AutoMigrate: %v", err)
	}
	log.Println("✅ AutoMigrate completado.")

	// --- 4. Inyección de Dependencias ---
	log.Println("🏗️  Inicializando dependencias...")

	// Repositorio de Categorías
	// El constructor NewCategoriaRepository está en el paquete 'categorias'
	// y devuelve categorias.CategoriaRepository (la interfaz definida en ese paquete).
	var categoriaRepo categorias.CategoriaRepository = categorias.NewCategoriaRepository(dbInstance)
	log.Println("   - Repositorio Categorías inicializado.")

	// Servicio de Categorías
	// El constructor NewCategoriaService está en el paquete 'categorias'
	// y devuelve categorias.CategoriaService (la interfaz definida en ese paquete).
	var categoriaService categorias.CategoriaService = categorias.NewCategoriaService(categoriaRepo)
	log.Println("   - Servicio Categorías inicializado.")

	// Handler de Categorías
	// El constructor NewCategoriaHandler está en el paquete 'categorias'.
	var categoriaHandler *categorias.CategoriaHandler = categorias.NewCategoriaHandler(categoriaService)
	log.Println("   - Handler Categorías inicializado.")
	log.Println("✅ Dependencias para Categorías inicializadas.")

	// ... (Aquí iría la inicialización de dependencias para Recetas, Auth, etc., cuando las migremos)

	// --- 5. Inicialización del Router Gin ---
	if cfg.AppEnv != "production" { // Usar modo debug para dev/test por defecto
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default() // Incluye logger y recovery middleware por defecto
	// Aquí podrías añadir tu middleware de errores global cuando lo tengas:
	// router.Use(middleware.ErrorHandler())
	log.Printf("✅ Router Gin inicializado en modo: %s.\n", gin.Mode())

	// --- 6. Configuración de Rutas ---
	// Creamos el grupo base para la API
	apiV1 := router.Group("/api/v1")

	// Registramos las rutas de Categorías usando la función del paquete 'categorias'
	// Asumimos que RegisterCategoriaRoutes está exportada desde el paquete 'categorias'
	// y que categoriaHandler es del tipo correcto esperado por esa función.
	if categoriaHandler != nil { // Chequeo por si algo falló arriba (aunque usamos Fatalf)
		categorias.RegisterCategoriaRoutes(apiV1, categoriaHandler)
		log.Println("✅ Rutas de Categorías registradas.")
	}
	// ... (Aquí irían las llamadas para registrar rutas de Recetas, Auth, etc.) ...
	// ejemplo: if recetaHandler != nil { recetas.RegisterRecetaRoutes(apiV1, recetaHandler) }

	// --- Rutas Base / No API ---
	router.Static("/public", "./public")   // Relativo al directorio de ejecución del binario
	router.Static("/uploads", "./uploads") // Relativo al directorio de ejecución
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "✅ API Recetas Go Funcionando! (Paquete por Característica)"})
	})
	log.Println("✅ Rutas adicionales configuradas.")
	// --- 7. Arranque del Servidor ---
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("🚀 Servidor escuchando en http://localhost%s\n", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor Gin: %v", err)
	}
}
