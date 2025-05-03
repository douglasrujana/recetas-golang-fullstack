// backend/internal/database/database.go

package database

import (
	"backend/internal/config" // Necesitamos la configuración para conectar
	"fmt"
	"log"

	"time" // [✨ OPCIONAL PERO ÚTIL] Para configurar el logger de Gorm

	"gorm.io/driver/mysql" // Asegúrate de que esto coincide con tu DB real
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // [✨ OPCIONAL PERO ÚTIL] Para configurar el logger de Gorm
)

// ConnectDB encapsula la lógica para establecer la conexión con la base de datos.
// Recibe la configuración y devuelve la instancia de Gorm o un error.
// [✅ BUENA PRÁCTICA] Función dedicada y exportada para la conexión.
func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	// Construir el Data Source Name (DSN) para la conexión MySQL.
	// [⚠️ ATENCIÓN] Tu .env tiene DB_PORT=5432 (PostgreSQL default) pero usas el driver mysql.
	// Asegúrate de que el puerto y el driver coincidan con tu base de datos real.
	// Usaré el cfg.DBPort leído del .env.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword, // [⚠️ SEGURIDAD] DB_PASS="" en .env es inseguro excepto para dev local muy específico.
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// [✨ OPCIONAL PERO RECOMENDADO] Configuración de Logger de Gorm
	// Esto te permite ver las queries SQL que Gorm ejecuta, muy útil para debug.
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Umbral de query lenta
			LogLevel:      logger.Info, // Nivel de Log (Silent, Error, Warn, Info)
			Colorful:      true,        // Habilitar salida colorida
		},
	)

	// Intentar conectar usando Gorm con la configuración de logger.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger, // [✨] Aplicar logger configurado
	})

	if err != nil {
		// Devolvemos el error para que el llamador (main) lo maneje.
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	log.Println("✅ Conexión exitosa a la base de datos")

	// [🔴 CONFIRMACIÓN] Nos aseguramos de que NO hay AutoMigrate aquí.

	// [✨ OPCIONAL PERO RECOMENDADO] Configurar Pool de Conexiones
	// Ayuda a gestionar las conexiones de forma eficiente.
	sqlDB, err := db.DB()
	if err != nil {
		// Si no se puede obtener el pool subyacente, es un error grave.
		return nil, fmt.Errorf("error al obtener el pool de conexiones sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)           // Número máximo de conexiones inactivas en el pool.
	sqlDB.SetMaxOpenConns(100)          // Número máximo de conexiones abiertas a la base de datos.
	sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo máximo que una conexión puede ser reutilizada.
	log.Println("⚙️ Pool de conexiones configurado")

	return db, nil
}
