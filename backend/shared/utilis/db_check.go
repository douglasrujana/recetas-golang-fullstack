// backend/internal/utils/db_check.go

// Este archivo contiene la lógica de chequeo de conexión a la base de datos.
// Utiliza GORM para la definición de la tabla y el mapeo de campos.
// No contiene tags de Gorm. Representa el negocio, no la base de datos.

package utils

import (
	"fmt"
	"log"
)

// CheckDatabaseConnection intenta cargar la configuración y establecer una conexión
// básica con la base de datos para verificar la conectividad.
// Cierra la conexión inmediatamente después de la prueba.
// Es útil para diagnósticos rápidos al inicio o en scripts.
//
// Parámetros:
//   - configPath: La ruta al directorio que contiene 'config.yaml'.
//
// Retorna:
//   - error: nil si la conexión es exitosa, de lo contrario, el error encontrado.
func CheckDatabaseConnection(configPath string) error {
	log.Println("🩺 Iniciando chequeo de conexión a base de datos...")

	// --- 1. Cargar Configuración ---
	log.Printf("   - Cargando configuración desde: %s\n", configPath)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// Envolver el error para dar contexto
		return fmt.Errorf("falló al cargar la configuración para el chequeo de BD: %w", err)
	}
	log.Println("   - Configuración cargada para chequeo.")
	// [Seguridad] No loguear datos sensibles aquí tampoco

	// --- 2. Intentar Conectar a la Base de Datos ---
	log.Printf("   - Intentando conectar a: %s:%d/%s como usuario %s...\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.User)
	dbInstance, err := database.ConnectDB(cfg.Database)
	if err != nil {
		// Envolver el error
		return fmt.Errorf("falló la conexión a la BD durante el chequeo: %w", err)
	}
	log.Println("   - ¡Conexión de prueba exitosa!")

	// --- 3. Cerrar la Conexión de Prueba ---
	// Es importante cerrar esta conexión de prueba para no dejarla colgando.
	log.Println("   - Cerrando conexión de prueba...")
	sqlDB, dbErr := dbInstance.DB()
	closeErr := fmt.Errorf("la instancia GORM es válida, pero no se pudo obtener sql.DB para cerrar: %w", dbErr) // Error por defecto si falla dbInstance.DB()

	if dbErr == nil { // Solo intentar cerrar si obtuvimos sqlDB
		closeErr = sqlDB.Close()
	}

	if closeErr != nil {
		log.Printf("   - ⚠️ Advertencia: Error al cerrar la conexión de prueba: %v", closeErr)
		// Nota: Decidimos NO devolver este error como fatal, ya que la conexión SÍ se estableció.
		// Si el cierre falla, es una advertencia, no un fallo del chequeo de conexión en sí.
	} else {
		log.Println("   - Conexión de prueba cerrada.")
	}

	log.Println("✅ Chequeo de conexión a base de datos finalizado con éxito.")
	return nil // Éxito
}