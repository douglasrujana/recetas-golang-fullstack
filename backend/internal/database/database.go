// backend/internal/database/database.go
package database

import (
	// [✨ INYECCIÓN DE DEPENDENCIA] Importamos específicamente la configuración de BD
	// que necesitamos, no toda la configuración de la app. Menor acoplamiento.
	"backend/internal/config"
	"fmt"
	"log"
	"os" // Necesario para log.New(os.Stdout, ...)
	"time"

	"gorm.io/driver/mysql" // Driver GORM para MySQL
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Logger de GORM para depuración
)

// --- PATRONES DE DISEÑO UTILIZADOS ---
// 1. FACTORY FUNCTION: ConnectDB actúa como una función fábrica. Su responsabilidad
//    es crear y configurar un objeto complejo (*gorm.DB) basado en la configuración.
// 2. DEPENDENCY INJECTION (Indirecta): Esta función CREA una dependencia (*gorm.DB)
//    que LUEGO será inyectada por main.go en los repositorios.
// 3. CONFIGURATION MANAGEMENT: Dependemos de una estructura de configuración (`config.DatabaseConfig`)
//    que es gestionada externamente (por Viper en main.go y el paquete config),
//    separando la configuración de la lógica de conexión.

// ConnectDB encapsula la lógica para establecer la conexión con la base de datos usando GORM.
// [✨ DESACOPLAMIENTO] Recibe SOLAMENTE la configuración específica de la base de datos (`config.DatabaseConfig`),
// en lugar de toda la `config.Config`.
// Devuelve la instancia de Gorm o un error.
func ConnectDB(cfgDb config.DatabaseConfig) (*gorm.DB, error) {

	// [✨ REFACTOR] Usamos el método DSN() definido en config.DatabaseConfig.
	// Esto centraliza la lógica de construcción del DSN y sigue el principio DRY (Don't Repeat Yourself).
	dsn := cfgDb.DSN()

	// [⚠️ VALIDACIÓN] La configuración de Viper (puerto, usuario, pass, host, name) debe ser correcta
	// para tu instancia de MySQL. Viper ya debería haber cargado estos valores desde config.yaml o env vars.
	log.Printf("ℹ️ Intentando conectar a la base de datos: mysql://%s@%s:%d/%s\n", cfgDb.User, cfgDb.Host, cfgDb.Port, cfgDb.Name)
	// [🚨 SEGURIDAD] La advertencia sobre contraseña vacía se movió a `config.LoadConfig` donde es más apropiado validarla.

	// --- Configuración del Logger de Gorm ---
	// [✅ BUENA PRÁCTICA] Logger configurable para ver SQL, útil en desarrollo/debug.
	logLevel := logger.Silent // Por defecto, no mostrar logs de GORM en producción
	isDevelopment := os.Getenv("APP_ENV") != "production" // O usa cfg.AppEnv si lo pasas o lo lees aquí de nuevo
	if isDevelopment {
		logLevel = logger.Info // Mostrar queries SQL en desarrollo
	}

	gormLogger := logger.New(
		// Usar log.New para asegurar la salida estándar con prefijo si es necesario
		log.New(os.Stdout, "\r\n[GORM] ", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond, // Umbral para queries lentas (ajustable)
			LogLevel:      logLevel,               // Nivel de log (Silent, Error, Warn, Info)
			Colorful:      isDevelopment,          // Salida colorida solo en desarrollo
		},
	)
	log.Printf("⚙️ Logger de GORM configurado a nivel: %v\n", logLevel)

	// --- Conexión GORM ---
	// [✅ BUENA PRÁCTICA] Usar gorm.Open con el driver específico y la configuración.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger, // Aplicar el logger configurado
		// Podrías añadir aquí otras configuraciones de GORM si las necesitas
		// PrepareStmt: true, // Podría mejorar rendimiento en algunos casos
	})

	if err != nil {
		// [✅ BUENA PRÁCTICA] Envolver el error para dar contexto claro.
		return nil, fmt.Errorf("❌ error al conectar a la base de datos con GORM: %w", err)
	}

	log.Println("✅ Conexión GORM establecida exitosamente.")

	// [🔴 CONFIRMACIÓN IMPORTANTE]
	// ¡NO DEBES LLAMAR A db.AutoMigrate() AQUÍ!
	// Las migraciones de esquema deben manejarse de forma explícita y controlada,
	// preferiblemente a través de un sistema/comando de migración separado
	// (ej: migrate, goose, gorm-migrate, o scripts SQL manuales).
	// Hacerlo aquí acopla el inicio de la aplicación con cambios de esquema,
	// lo cual es peligroso y poco flexible.

	// --- Configuración del Pool de Conexiones ---
	// [✅ BUENA PRÁCTICA] Configurar el pool subyacente de sql.DB para eficiencia.
	sqlDB, err := db.DB()
	if err != nil {
		// Si GORM no puede devolver el *sql.DB subyacente, algo va muy mal.
		// Podrías cerrar la conexión gorm aquí si fuera necesario: db.Close() ?
		return nil, fmt.Errorf("❌ error crítico al obtener el pool sql.DB de GORM: %w", err)
	}

	// Valores de ejemplo, ajústalos según tus necesidades y pruebas de carga
	sqlDB.SetMaxIdleConns(10)           // Conexiones inactivas máximas
	sqlDB.SetMaxOpenConns(50)          // Conexiones abiertas máximas totales
	sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo máximo de vida de una conexión
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Tiempo máximo inactiva antes de cerrarse
	log.Printf("⚙️ Pool de conexiones SQL configurado: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v\n",
		50, 10, time.Hour, 10*time.Minute)

	// [✅ BUENA PRÁCTICA] Hacer un Ping para verificar la conexión real antes de devolverla.
	err = sqlDB.Ping()
	if err != nil {
		// Si el ping falla, cierra el pool antes de devolver el error.
		sqlDB.Close()
		return nil, fmt.Errorf("❌ la base de datos no respondió al ping inicial: %w", err)
	}
	log.Println("✅ Ping a la base de datos exitoso.")


	// Devuelve la instancia de GORM lista para ser usada por los repositorios.
	return db, nil
}