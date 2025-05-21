// backend/contactos/model_gorm.go
package contactos

import (
	// "backend/users" // Para users.UserModel, cuando exista
	"time"
	"gorm.io/gorm"
)

type ContactoModel struct {
	ID                uint           `gorm:"primaryKey"`
	UserID            *uint          `gorm:"index;default:null"` // FK Opcional
	NombreRemitente   string         `gorm:"type:varchar(150);not null"`
	EmailRemitente    string         `gorm:"type:varchar(255);not null;index"` // Añadido index
	TelefonoRemitente string         `gorm:"type:varchar(30);default:null"`
	Asunto            string         `gorm:"type:varchar(255);default:null"`
	Mensaje           string         `gorm:"type:text;not null"`
	Leido             bool           `gorm:"not null;default:false;index:idx_contactos_leido_fecha_model"`
	FechaContacto     time.Time      `gorm:"not null;index:idx_contactos_leido_fecha_model"`
	IPOrigen          string         `gorm:"type:varchar(45);default:null"`
	UserAgent         string         `gorm:"type:text;default:null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"` // Opcional para soft delete

	// User users.UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// Se definirá cuando 'users.UserModel' exista. AutoMigrate creará la columna user_id.
	// GORM añadirá la constraint FK si se corre AutoMigrate DE NUEVO después de crear UserModel.
}

func (ContactoModel) TableName() string {
	return "contactos"
}

func (m *ContactoModel) ToDomain() *ContactoForm {
	if m == nil { return nil }
	return &ContactoForm{
		ID: m.ID, UserID: m.UserID, NombreRemitente: m.NombreRemitente, EmailRemitente: m.EmailRemitente,
		TelefonoRemitente: m.TelefonoRemitente, Asunto: m.Asunto, Mensaje: m.Mensaje, Leido: m.Leido,
		FechaContacto: m.FechaContacto, IPOrigen: m.IPOrigen, UserAgent: m.UserAgent,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func FromContactoFormDomain(d *ContactoForm) *ContactoModel {
	if d == nil { return nil }
	return &ContactoModel{
		ID: d.ID, UserID: d.UserID, NombreRemitente: d.NombreRemitente, EmailRemitente: d.EmailRemitente,
		TelefonoRemitente: d.TelefonoRemitente, Asunto: d.Asunto, Mensaje: d.Mensaje, Leido: d.Leido,
		FechaContacto: d.FechaContacto, IPOrigen: d.IPOrigen, UserAgent: d.UserAgent,
	}
}

func ContactoModelsToDomains(models []ContactoModel) []ContactoForm {
	if models == nil { return []ContactoForm{} }
	domainForms := make([]ContactoForm, 0, len(models))
	for _, model := range models {
		if dm := model.ToDomain(); dm != nil { domainForms = append(domainForms, *dm) }
	}
	return domainForms
}