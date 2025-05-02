// backend/internal/handler/routes.go

package handler

import (
	"github.com/gin-gonic/gin"
	// [✨ IMPORTANTE] Necesitamos importar los otros archivos dentro de 'handler'
	// para poder referenciar los tipos de los handlers (AuthHandler, CategoriaHandler, etc.)
	// y sus métodos. Asegúrate de que los nombres de paquete y archivo sean correctos.

	// [🔵 PENDIENTE] Aún necesitamos el middleware de autenticación.
	// Lo ideal sería que también fuera un método de AuthHandler o un struct separado
	// que reciba dependencias (como AuthServicio para validar el token).
	// Por ahora, asumiremos que existe una función o método `authMiddlewareFunc`.
)

// [✅ BUENA PRÁCTICA] Agrupar las rutas por recurso/dominio usando Gin Groups.

// RegisterRoutes configura todas las rutas de la API en el router de Gin.
// Recibe el router y las instancias de los handlers necesarios.
// [✅ BUENA PRÁCTICA] Dependencias explícitas pasadas como argumentos.
func RegisterRoutes(
	router *gin.Engine,
	authHandler *AuthHandler, // Handler para autenticación (register/login)
	categoriaHandler *CategoriaHandler, // Handler para categorías
	recetaHandler *RecetaHandler, // Handler para recetas
	uploadHandler *UploadHandler, // Handler para subida de archivos
	// ...añadir otros handlers aquí si los hubiera ...
) {
	// --- Middleware Global (Si aplica) ---
	// router.Use(middleware.LoggingMiddleware()) // Ejemplo si tuvieras middleware global

	// --- Rutas Públicas ---
	api := router.Group("/api/v1") // [✅ BUENA PRÁCTICA] Usar un prefijo base para la API

	// Rutas de Autenticación
	authRoutes := api.Group("/auth") // [✨ SUGERENCIA] Agrupar auth bajo /auth
	{
		authRoutes.POST("/register", authHandler.Register) // Llama al método Register del AuthHandler
		authRoutes.POST("/login", authHandler.Login)       // Llama al método Login del AuthHandler
	}

	// Rutas de Subida (Upload) - Asumiendo pública por ahora, o necesita auth
	// [❓ PREGUNTA] ¿La subida requiere autenticación? Si sí, mover dentro de `apiAuthRequired`.
	api.POST("/upload", uploadHandler.HandleUpload) // Llama al método del UploadHandler

	// --- Rutas Protegidas por Autenticación ---
	apiAuthRequired := api.Group("/") // Grupo para rutas que requieren autenticación

	// [🏗️ MIDDLEWARE] Aquí aplicamos el middleware de autenticación.
	// Necesitas crear este middleware. Podría ser un método en AuthHandler
	// o un struct/función separada en un paquete `middleware`.
	// Ejemplo: apiAuthRequired.Use(authHandler.AuthMiddleware())
	// Por ahora, usaremos una función placeholder:
	apiAuthRequired.Use(placeholderAuthMiddleware()) // ¡¡REEMPLAZAR CON TU MIDDLEWARE REAL!!

	// Rutas de Categorías (Protegidas)
	categoriaRoutes := apiAuthRequired.Group("/categorias")
	{
		categoriaRoutes.GET("", categoriaHandler.GetAll)          // GET /api/v1/categorias
		categoriaRoutes.GET("/:id", categoriaHandler.GetByID)     // GET /api/v1/categorias/:id
		categoriaRoutes.POST("", categoriaHandler.Create)         // POST /api/v1/categorias
		categoriaRoutes.PUT("/:id", categoriaHandler.Update)      // PUT /api/v1/categorias/:id
		categoriaRoutes.DELETE("/:id", categoriaHandler.Delete)   // DELETE /api/v1/categorias/:id
	}

	// Rutas de Recetas (Protegidas)
	recetaRoutes := apiAuthRequired.Group("/recetas")
	{
		recetaRoutes.GET("", recetaHandler.GetAll)         // GET /api/v1/recetas
		recetaRoutes.GET("/:id", recetaHandler.GetByID)    // GET /api/v1/recetas/:id
		recetaRoutes.POST("", recetaHandler.Create)        // POST /api/v1/recetas
		// recetaRoutes.PUT("/:id", recetaHandler.Update)   // PUT /api/v1/recetas/:id (Si existe)
		// recetaRoutes.DELETE("/:id", recetaHandler.Delete) // DELETE /api/v1/recetas/:id (Si existe)
	}

	// ... Añadir otros grupos de rutas protegidas aquí ...

	log.Println("🛣️ Rutas de la API configuradas")
}

// --- Placeholder para el Middleware ---
// [🔴 TEMPORAL] ¡¡DEBES REEMPLAZAR ESTO!!
// Crea tu middleware real, probablemente en `auth_handler.go` o `middleware/auth.go`
// que valide el token JWT y añada la información del usuario al contexto de Gin.
func placeholderAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("⚠️ ADVERTENCIA: Usando Placeholder Auth Middleware. ¡Implementar validación JWT real!")
		// Aquí iría la lógica:
		// 1. Extraer token (e.g., "Bearer <token>") del header "Authorization".
		// 2. Validar el token (firma, expiración) usando la SECRET_KEY.
		// 3. Si es válido, extraer claims (e.g., user ID).
		// 4. Añadir claims/user ID al contexto Gin: c.Set("userID", userID)
		// 5. Llamar a c.Next() para pasar al siguiente handler.
		// 6. Si no es válido, llamar a c.AbortWithStatusJSON(http.StatusUnauthorized, ...)
		c.Next() // Por ahora, simplemente deja pasar todo.
	}
}