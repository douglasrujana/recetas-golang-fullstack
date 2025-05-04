// backend/internal/config/config.go
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config struct holds all configuration for the application.
// Se mapea desde config.yaml o variables de entorno.
type Config struct {
	AppEnv    string         `mapstructure:"app_env"`    // Entorno (development, production, etc.)
	SecretKey string         `mapstructure:"secret_key"` // Clave secreta para JWT, etc.
	Server    ServerConfig   `mapstructure:"server"`
	Database  DatabaseConfig `mapstructure:"database"`
}

// ServerConfig holds server specific configuration.
type ServerConfig struct {
	Port int `mapstructure:"port"` // Puerto del servidor (ej: 8080)
}

// DatabaseConfig holds database specific configuration.
type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"` // Puerto de la BD (ej: 3306 para MySQL)
	Name     string `mapstructure:"name"`
	Params   string `mapstructure:"params"` // Ej: "parseTime=true"
}

// LoadConfig reads configuration from file or environment variables using Viper.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)     // Directorio donde buscar config.yaml
	viper.SetConfigName("config") // Nombre del archivo (sin extensión)
	viper.SetConfigType("yaml")   // Tipo del archivo

	viper.SetEnvPrefix("APP")                              // Prefijo para variables de entorno (ej: APP_SERVER_PORT)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Reemplaza . con _ (server.port -> SERVER_PORT)
	viper.AutomaticEnv()                                   // Lee variables de entorno que coincidan

	// Valores por defecto (opcional pero útil)
	viper.SetDefault("app_env", "development")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306) // Puerto por defecto MySQL
	viper.SetDefault("database.params", "parseTime=true")

	err = viper.ReadInConfig() // Intenta leer el archivo config.yaml
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// El archivo no existe, no es un error fatal si las variables de entorno están configuradas
			fmt.Println("⚠️ Archivo 'config.yaml' no encontrado, usando valores por defecto y variables de entorno.")
		} else {
			// Otro error al leer el archivo
			return Config{}, fmt.Errorf("❌ error fatal al leer archivo de configuración: %w", err)
		}
	}

	// Mapear la configuración leída (de archivo y/o env) al struct Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("❌ incapaz de decodificar configuración en struct: %w", err)
	}

	// [⚠️ VALIDACIÓN IMPORTANTE] Verifica si los secretos esenciales están presentes
	if config.SecretKey == "" {
		// Podrías devolver un error si es absolutamente necesario
		fmt.Println("🚨 ¡Advertencia! SECRET_KEY no está definida en config.yaml ni como variable de entorno (APP_SECRET_KEY).")
		// return Config{}, errors.New("secret_key es requerida")
	}
	if config.Database.Password == "" {
		fmt.Println("🚨 ¡Advertencia! La contraseña de la base de datos no está definida (database.password o APP_DATABASE_PASSWORD).")
	}

	return config, nil
}

// DSN constructs the Data Source Name string for MySQL connection.
// [✨ MEJORA] Esta función es útil tenerla aquí.
func (dbConfig DatabaseConfig) DSN() string {
	// user:password@tcp(host:port)/dbname?params
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Params,
	)
}
