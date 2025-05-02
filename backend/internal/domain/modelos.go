package modelos

import (
	"backend/internal/database" // Importa el paquete de base de datos
	"time"
)

// Definición de la estructura de la tabla Categoria
type Categoria struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"` //  Mysql  bigint(20) unsigned NOT NULL AUTO_INCREMENT
	Nombre string `gorm:"type:varchar(100)" json:"nombre"`
	Slug   string `gorm:"type:varchar(100)" json:"slug"`
}
type Categorias []Categoria // Definición de un slice de Categoría

// Definición de la estructura de la tabla Receta
type Receta struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"` // Mysql  bigint(20) unsigned NOT NULL AUTO_INCREMENT
	CategoriaId uint      `json:"categoria_id"`             // Mysql  bigint(20) unsigned NOT NULL AUTO_INCREMENT
	Categoria   Categoria `gorm:"foreignKey:CategoriaId"`   // Relación con la tabla Categoria
	Nombre      string    `gorm:"type:varchar(100)" json:"nombre"`
	Slug        string    `gorm:"type:varchar(100)" json:"slug"`
	Tiempo      string    `gorm:"type:varchar(100)" json:"tiempo"`
	Foto        string    `gorm:"type:varchar(100)" json:"foto"`
	Descripcion string    `gorm:"type:varchar(100)" json:"descripcion"`
	Fecha       time.Time `json:"fecha"` // Mysql datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
}
type Recetas []Receta // Definición de un slice de Categoría

// Definición de la estructura de la tabla Contacto
type Contacto struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"` //  Mysql  bigint(20) unsigned NOT NULL AUTO_INCREMENT
	Nombre   string    `gorm:"type:varchar(100)" json:"nombre"`
	Correo   string    `gorm:"type:varchar(100)" json:"correo"`
	Telefono string    `gorm:"type:varchar(100)" json:"telefono"`
	Mensaje  string    `gorm:"type:varchar(100)" json:"mensaje"`
	Fecha    time.Time `json:"fecha"` // Mysql datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
}
type Contactos []Contacto // Definición de un slice de Categoría

func Migraciones() {
	// Migraciones de las tablas
	database.Database.AutoMigrate(&Categoria{}, &Receta{}, &Contacto{})
}
