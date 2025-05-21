// backend/cmd/seeder/main.go
// Este archivo contiene la lógica del seeder para la base de datos.
// Se encarga de cargar la configuración, conectar a la base de datos,
// ejecutar AutoMigrate y seedear datos iniciales.
package main

import (
	"log"
	//"time" // Para CreatedAt/UpdatedAt si los seteamos manualmente
	// Importar paquetes necesarios
	"backend/categorias"         // Para CategoriaModel y sus constructores/tipos si es necesario
	"backend/recetas"            // Para RecetaModel y sus constructores/tipos
	"backend/shared/config"    // Para cargar configuración
	"backend/shared/database"  // Para conectar a la BD
	// "github.com/gosimple/slug" // Si necesitas generar slugs aquí también
	// "gorm.io/gorm" // No es estrictamente necesario importar gorm aquí si los modelos lo encapsulan
)

func main() {
	log.Println("🌱 Iniciando Seeding de Base de Datos...")

	// --- 1. Cargar Configuración (usualmente la de desarrollo) ---
	// Asumimos que este seeder se corre en un contexto donde config.yaml es accesible
	// o APP_ENV está seteado para desarrollo.
	// La ruta "config" es relativa al directorio donde se ejecuta 'go run' o el binario.
	// Si ejecutas desde 'backend/', y 'config' está en 'backend/config/', la ruta es "config".
	cfg, err := config.LoadConfig("config") // Carga config.yaml por defecto
	if err != nil {
		log.Fatalf("❌ Error cargando config para seeder: %v", err)
	}
	// Puedes forzar un entorno si es necesario, ej: os.Setenv("APP_ENV", "development") antes de LoadConfig

	log.Printf("   - Usando configuración para entorno: %s, BD: %s\n", cfg.AppEnv, cfg.Database.Name)

	// --- 2. Conectar a la Base de Datos ---
	db, err := database.ConnectDB(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Error conectando a la BD para seeder: %v", err)
	}
	log.Println("   - Conexión a BD exitosa.")

	// --- 3. Ejecutar AutoMigrate (Opcional pero Recomendado) ---
	// Asegura que las tablas existan ANTES de intentar insertar datos.
	// Es bueno tenerlo aquí para que el seeder sea más autocontenido.
	log.Println("   - Ejecutando AutoMigrate para asegurar tablas...")
	err = db.AutoMigrate(
		&categorias.CategoriaModel{},
		&recetas.RecetaModel{},
		// ... añadir TODOS tus otros *Model GORM aquí ...
	)
	if err != nil {
		log.Fatalf("❌ Error durante AutoMigrate en seeder: %v", err)
	}
	log.Println("   - AutoMigrate completado.")

	// --- 4. Seeding de Datos ---

	// --- Seed Categorías ---
	log.Println("   - Seedeando Categorías...")
	categoriasSeed := []categorias.CategoriaModel{
		{Nombre: "Postres", Slug: "postres"}, // GORM manejará CreatedAt/UpdatedAt
		{Nombre: "Ensaladas Frescas", Slug: "ensaladas-frescas"},
		{Nombre: "Sopas Caseras", Slug: "sopas-caseras"},
		{Nombre: "Platos Fuertes", Slug: "platos-fuertes"},
	}

	for i := range categoriasSeed {
		// Usar FirstOrCreate para evitar duplicados si el seeder se corre múltiples veces
		// Busca por un campo único (Nombre o Slug)
		result := db.FirstOrCreate(&categoriasSeed[i], categorias.CategoriaModel{Nombre: categoriasSeed[i].Nombre})
		if result.Error != nil {
			log.Printf("     ❌ Error seedeando categoría '%s': %v\n", categoriasSeed[i].Nombre, result.Error)
		} else {
			if result.RowsAffected > 0 {
				log.Printf("     ✅ Categoría '%s' creada/encontrada con ID: %d\n", categoriasSeed[i].Nombre, categoriasSeed[i].ID)
			} else {
				log.Printf("     ℹ️ Categoría '%s' ya existía con ID: %d\n", categoriasSeed[i].Nombre, categoriasSeed[i].ID)
			}
		}
	}
	log.Println("   - Seeding de Categorías finalizado.")


	// --- Seed Recetas ---
	// Para seedeear recetas, necesitamos los IDs de las categorías que acabamos de crear o que ya existían.
	// Primero, obtenemos las categorías de la BD para tener sus IDs.
	var catPostres categorias.CategoriaModel
	db.Where("slug = ?", "postres").First(&catPostres)

	var catPlatosFuertes categorias.CategoriaModel
	db.Where("slug = ?", "platos-fuertes").First(&catPlatosFuertes)

	log.Println("   - Seedeando Recetas...")
	recetasSeed := []recetas.RecetaModel{
		{
			Nombre:            "Tiramisú Clásico",
			Slug:              "tiramisu-clasico",
			CategoriaID:       catPostres.ID, // Usar ID de categoría obtenida
			TiempoPreparacion: "30 minutos + refrigeración",
			Descripcion:       "El auténtico tiramisú italiano, cremoso y delicioso.",
			Foto:              "tiramisu.jpg", // Asumimos nombres de archivo
		},
		{
			Nombre:            "Lomo Saltado Peruano",
			Slug:              "lomo-saltado-peruano",
			CategoriaID:       catPlatosFuertes.ID,
			TiempoPreparacion: "45 minutos",
			Descripcion:       "Un plato emblemático de la cocina peruana, lleno de sabor.",
			Foto:              "lomo-saltado.jpg",
		},
		{
			Nombre:            "Cheesecake de Fresa sin Horno",
			Slug:              "cheesecake-fresa-sin-horno",
			CategoriaID:       catPostres.ID,
			TiempoPreparacion: "20 minutos + refrigeración",
			Descripcion:       "Fácil, rápido y perfecto para cualquier ocasión.",
			Foto:              "cheesecake-fresa.jpg",
		},
	}

	for i := range recetasSeed {
		// Solo crear si CategoriaID es válida (mayor que 0)
		if recetasSeed[i].CategoriaID > 0 {
			result := db.FirstOrCreate(&recetasSeed[i], recetas.RecetaModel{Nombre: recetasSeed[i].Nombre}) // Busca por Nombre para evitar duplicados
			if result.Error != nil {
				log.Printf("     ❌ Error seedeando receta '%s': %v\n", recetasSeed[i].Nombre, result.Error)
			} else {
				if result.RowsAffected > 0 {
					log.Printf("     ✅ Receta '%s' creada/encontrada con ID: %d\n", recetasSeed[i].Nombre, recetasSeed[i].ID)
				} else {
					log.Printf("     ℹ️ Receta '%s' ya existía con ID: %d\n", recetasSeed[i].Nombre, recetasSeed[i].ID)
				}
			}
		} else {
			log.Printf("     ⚠️ Receta '%s' omitida porque su CategoriaID no fue encontrada/asignada.", recetasSeed[i].Nombre)
		}
	}
	log.Println("   - Seeding de Recetas finalizado.")

	// --- (A futuro) Seed Ingredientes ---
	// --- (A futuro) Seed Relaciones Receta-Ingredientes ---

	log.Println("✅ Seeding de Base de Datos Completado.")
}