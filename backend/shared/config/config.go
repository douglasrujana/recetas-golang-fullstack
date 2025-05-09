// backend/internal/config/config.go (REFACTORIZADO para Entornos)
package config

import (
	"fmt"
	"os"      // Para leer variables de entorno
	"strings" // Para el reemplazo de claves de Env Var
	"github.com/spf13/viper"
)

// --- Structs Config, ServerConfig, DatabaseConfig (sin cambios) ---
type Config struct {
	AppEnv    string         `mapstructure:"app_env"`
	SecretKey string         `mapstructure:"secret_key"`
	Server    ServerConfig   `mapstructure:"server"`
	Database  DatabaseConfig `mapstructure:"database"`
}
type ServerConfig struct {
	Port int `mapstructure:"port"`
}
type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Params   string `mapstructure:"params"`
}
// --- Fin Structs ---

// LoadConfig carga la configuración basada en el entorno.
// Busca un archivo config.{APP_ENV}.yaml o config.yaml en configPath.
// Las variables de entorno (con prefijo APP_) tienen precedencia.
func LoadConfig(configPath string) (config Config, err error) {
	// 1. Determinar el Entorno y Nombre del Archivo
	configName := "config" // Nombre base por defecto (cargará config.yaml)
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development" // Default environment si APP_ENV no está seteada
		fmt.Println("ℹ️ APP_ENV no definida, usando entorno por defecto: 'development'. Cargando 'config.yaml'.")
	} else {
		// Si APP_ENV está definida (ej: "test", "production"), buscar config.{appEnv}.yaml
		configName = fmt.Sprintf("config.%s", appEnv)
		fmt.Printf("ℹ️ Detectado APP_ENV='%s'. Intentando cargar '%s.yaml'.\n", appEnv, configName)
	}

	// 2. Configurar Viper
	viper.AddConfigPath(configPath)     // Directorio donde buscar los archivos (ej: "config/")
	viper.SetConfigName(configName)     // Nombre del archivo a buscar (sin extensión)
	viper.SetConfigType("yaml")         // Tipo de archivo

	// 3. Configurar Lectura desde Variables de Entorno (¡PRIORIDAD ALTA!)
	// Las variables de entorno sobrescribirán los valores del archivo .yaml
	viper.SetEnvPrefix("APP") // ej: APP_SERVER_PORT, APP_DATABASE_PASSWORD
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // server.port -> SERVER_PORT
	viper.AutomaticEnv()              // Habilitar lectura automática

	// 4. Intentar Leer el Archivo de Configuración
	err = viper.ReadInConfig()
	if err != nil {
		// Manejar el caso en que el archivo no se encuentre
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Si falta config.yaml en desarrollo, es un problema.
			// Si falta config.test.yaml o config.production.yaml, podría ser normal
			// si se depende exclusivamente de variables de entorno.
			if configName == "config" { // Es decir, appEnv era "development" o default
                fmt.Printf("⚠️ Advertencia Crítica: Archivo de configuración base 'config.yaml' no encontrado en '%s'.\n", configPath)
                fmt.Println("    Continuando con valores por defecto y variables de entorno, pero esto puede fallar.")
                // Considera retornar un error aquí si config.yaml es mandatorio para desarrollo.
                // return Config{}, fmt.Errorf("archivo de configuración base 'config.yaml' no encontrado")
            } else {
                 fmt.Printf("⚠️ Archivo de configuración '%s.yaml' no encontrado en '%s'. Se usarán variables de entorno y/o defaults.\n", configName, configPath)
            }
		} else {
			// Otro error al leer el archivo
			return Config{}, fmt.Errorf("❌ error fatal al leer archivo de configuración '%s.yaml': %w", configName, err)
		}
	} else {
		fmt.Printf("✅ Archivo de configuración '%s' cargado exitosamente.\n", viper.ConfigFileUsed())
	}

	// 5. Establecer Valores por Defecto (se aplican si no están en archivo NI en ENV)
	viper.SetDefault("app_env", appEnv) // Usar el entorno detectado o 'development'
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.params", "parseTime=true")
    // Es mejor definir User/Name en los archivos de config o ENV que confiar en defaults
	// viper.SetDefault("database.user", "root")
	// viper.SetDefault("database.name", "recetas_dev")

	// 6. Mapear la Configuración Final (Archivo + Env + Defaults) al Struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("❌ incapaz de decodificar configuración en struct: %w", err)
	}

    // 7. Asegurar que AppEnv en el struct refleje el entorno real detectado
    // (por si el archivo .yaml tenía un valor diferente hardcodeado)
    config.AppEnv = appEnv

	// 8. Validación de Secretos (importante para producción)
    // En producción, forzar que las contraseñas/claves vengan de ENV VARS
    if config.AppEnv == "production" {
         if config.Database.Password == "" || config.SecretKey == "" {
              return Config{}, fmt.Errorf("❌ ERROR FATAL: En producción, APP_DATABASE_PASSWORD y APP_SECRET_KEY deben definirse como variables de entorno")
         }
    } else if config.Database.Password == "" { // Advertencia para dev/test
         fmt.Println("🚨 ¡Advertencia! La contraseña de la base de datos no está definida (ni en archivo ni como APP_DATABASE_PASSWORD).")
    }
    if config.SecretKey == "" && config.AppEnv != "production" {
         fmt.Println("🚨 ¡Advertencia! SECRET_KEY no está definida (ni en archivo ni como APP_SECRET_KEY).")
    }

	return config, nil
}

// --- Función DSN() (sin cambios) ---
// Devuelve el string de conexión a la base de datos
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