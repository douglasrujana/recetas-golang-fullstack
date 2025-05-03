package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

// Config es el struct central con todas las variables ya parseadas
type Config struct {
	AppPort    string
	AppEnv     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	UploadsDir string
	PublicHost string
}

var AppConfig Config

// CargarVariablesEntorno lee el archivo .env y asigna a AppConfig
func CargarVariablesEntorno() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ No se encontró .env. Usando variables del entorno del sistema.")
		}

		AppConfig = Config{
			AppPort:    getEnv("APP_PORT", "8085"),
			AppEnv:     getEnv("APP_ENV", "development"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", "recetas_db"),
			UploadsDir: getEnv("UPLOADS_DIR", "uploads"),
			PublicHost: getEnv("PUBLIC_HOST", "http://localhost:8085"),
		}
	})
}

func getEnv(clave string, valorPorDefecto string) string {
	if valor, ok := os.LookupEnv(clave); ok {
		return valor
	}
	return valorPorDefecto
}

//url := config.AppConfig.PublicHost + "/uploads/" + nombreArchivo
