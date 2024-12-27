package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Certificate struct {
	ID         uuid.UUID                   `json:"id" gorm:"type:uuid;primary_key;"`
	OwnerID    uuid.UUID                   `json:"owner_id" gorm:"type:uuid;not null;index"`
	OwnerType  string                      `json:"owner_type" gorm:"not null"`
	NotBefore  time.Time                   `json:"not_before" gorm:"not null"`
	NotAfter   time.Time                   `json:"not_after" gorm:"not null"`
	Crt        []byte                      `json:"crt" swaggertype:"string" format:"base64"`
	Key        []byte                      `json:"key" swaggertype:"string" format:"base64"`
	Pub        []byte                      `json:"pub" swaggertype:"string" format:"base64"`
	Passphrase string                      `json:"passphrase" gorm:"size:255"`
	IsCA       bool                        `json:"is_ca" gorm:"default:false"`
	IPs        datatypes.JSONSlice[string] `json:"ips" swaggertype:"array,string"`     // List of IPv4 addresses and networks in CIDR notation. Limits the addresses for subordinate certificates.
	Subnets    datatypes.JSONSlice[string] `json:"subnets" swaggertype:"array,string"` // List of IPv4 subnets in CIDR notation. Defines subnets that subordinate certificates can use.
	Groups     datatypes.JSONSlice[string] `json:"groups" swaggertype:"array,string"`  // List of groups for access control, restricting subordinate certificates' groups.
	CreatedAt  time.Time                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time                   `json:"updated_at" gorm:"autoUpdateTime"`
}
