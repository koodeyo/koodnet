package models

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	OwnerID    uuid.UUID `json:"ownerId" gorm:"type:uuid;not null;index"`
	OwnerType  string    `json:"ownerType" gorm:"not null"`
	NotBefore  time.Time `json:"notBefore" gorm:"not null"`
	NotAfter   time.Time `json:"notAfter" gorm:"not null"`
	Crt        []byte    `json:"crt" swaggertype:"string"`
	Key        []byte    `json:"key" swaggertype:"string"`
	Pub        []byte    `json:"pub" swaggertype:"string"`
	Passphrase string    `json:"passphrase" gorm:"size:255"`
	IsCA       bool      `json:"isCa" gorm:"default:false"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (c Certificate) Expired() bool {
	t := time.Now()
	return c.NotBefore.After(t) || c.NotAfter.Before(t)
}
