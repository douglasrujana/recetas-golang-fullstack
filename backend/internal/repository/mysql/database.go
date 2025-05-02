package database

import (
	"fmt"
	"log"
	"backend/internal/config"
	"backend/internal/database" // Cambia "tu-modulo" por el nombre de tu módulo real
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database = func() *gorm.DB {
	// Cargamos las variables de entorno si aún no se han cargado
	config.CargarVariablesEntorno()

	// Construir el DSN de conexión para MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	// Abrimos la conexión con GORM y MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("🚨 Error conectando a la base de datos: %v", err)
	}

	log.Println("✅ Conexión a la base de datos establecida.")
	return db
}()
