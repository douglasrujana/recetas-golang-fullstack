package handler

const prefijo = "/api/v1/" // Prefijo para las rutas
const id = "/:id"          // Ruta constante para "ejemplo/:id"

// Configurar el middleware de recuperación para manejar errores
router.Static("/public", "./public")
router.Static("/uploads", "./uploads") // Sirve la carpeta "uploads"

// GET
router.GET("/", func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello Go...!"})
})

// Ejemplo
router.GET(prefijo+"ejemplo", rutas.EjemploRuta)
router.GET(prefijo+"ejemplo"+id, rutas.EjemploGetId)
router.GET(prefijo+"ejemplo-querystring", rutas.EjemploQueryString)
router.POST(prefijo+"ejemplo", rutas.EjemploPost)
router.POST(prefijo+"ejemplo-dto", rutas.EjemploPostDTO)
router.PUT(prefijo+"ejemplo"+id, rutas.EjemploPut)
router.DELETE(prefijo+"ejemplo"+id, rutas.EjemploDelete)

// Categorias
router.GET(prefijo+"categorias", rutas.CategoriasLista)      // Ruta para listar categorías
router.GET(prefijo+"categorias"+id, rutas.CategoriasListaId) // Ruta para listar categorías por ID
router.POST(prefijo+"categorias", rutas.CategoriasPost) // Ruta para crear una nueva
router.PUT(prefijo+"categorias"+id, rutas.CategoriasPut)
router.DELETE(prefijo+"categorias"+id, rutas.CategoriasDelete)

// Recetas
router.GET(prefijo+"recetas", rutas.RecetasListar)          // Ruta para listar recetas
router.GET(prefijo+"recetas"+id, rutas.RecetasListarById)    // Ruta para listar recetas por ID
router.POST(prefijo+"recetas", rutas.RecetasCrear)    // Ruta para listar recetas por ID

// Upload
router.POST(prefijo+"upload", rutas.EjemploUpload)





