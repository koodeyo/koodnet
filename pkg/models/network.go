package models

import (
	"time"
)

type Network struct {
	ID                  string    `json:"id" gorm:"type:uuid;primary_key;"`
	Name                string    `gorm:"type:varchar(255);uniqueIndex:idx_name_cidr" json:"name"` // Max 255 characters, part of unique constraint
	Description         string    `gorm:"type:varchar(255)" json:"description,omitempty"`          // Max 500 characters, optional
	Cidr                string    `gorm:"type:varchar(255)" json:"cidr"`                           // IP Range
	LighthousesAsRelays bool      `gorm:"default:false" json:"lighthousesAsRelays"`                // Default to false
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	// Tags                []NetworkTag         `gorm:"foreignKey:NetworkID" json:"tags"`                        // One-to-Many with NetworkTag
	// Roles               []NetworkRole        `gorm:"foreignKey:NetworkID" json:"roles"`                       // One-to-Many with NetworkRole
	// Routes              []NetworkRoute       `gorm:"foreignKey:NetworkID" json:"routes"`                      // One-to-Many with NetworkRoute
	// Hosts               []NetworkHost        `gorm:"foreignKey:NetworkID" json:"hosts"`                       // One-to-Many with NetworkHost
	// Ca                  []NetworkCertificate `gorm:"foreignKey:NetworkID" json:"ca"`                          // One-to-Many with NetworkCertificate
}
