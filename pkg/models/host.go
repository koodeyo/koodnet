package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Host struct {
	ID              uuid.UUID                   `json:"id" gorm:"type:uuid;primary_key;"`
	Name            string                      `json:"name" gorm:"size:255;not null;uniqueIndex:idx_name_network"`
	IP              string                      `json:"ip" gorm:"size:255;not null;uniqueIndex:idx_ip_network"`
	StaticAddresses datatypes.JSONSlice[string] `json:"static_addresses" gorm:"type:json;default:'[]'" swaggertype:"array,string"`
	Subnets         datatypes.JSONSlice[string] `json:"subnets" gorm:"type:json;default:'[]'" swaggertype:"array,string"`
	Groups          datatypes.JSONSlice[string] `json:"groups" gorm:"type:json;default:'[]'" swaggertype:"array,string"`
	ListenPort      uint                        `json:"listen_port" gorm:"default:4242"`
	IsLighthouse    bool                        `json:"is_lighthouse" gorm:"default:false"`
	NetworkID       uuid.UUID                   `json:"network_id" gorm:"type:uuid"`
	Network         Network                     `json:"network"`
	Certificates    []Certificate               `json:"Certificates,omitempty" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time                   `json:"updated_at" gorm:"autoUpdateTime"`
}

type HostDto struct {
	Name            string                      `json:"name,omitempty" example:"host-1"`
	IP              string                      `json:"ip,omitempty" example:"192.168.1.100"`
	StaticAddresses datatypes.JSONSlice[string] `json:"static_addresses,omitempty" swaggertype:"array,string" example:"109.243.69.39"`
	Subnets         datatypes.JSONSlice[string] `json:"subnets,omitempty" swaggertype:"array,string" example:"192.168.1.0/24,10.0.0.0/16"`
	Groups          datatypes.JSONSlice[string] `json:"groups,omitempty" swaggertype:"array,string" example:"servers,laptops"`
	NetworkID       uuid.UUID                   `json:"network_id,omitempty" example:"c6d6c4c4-b65b-40e1-bcf2-1fd3122c653d"`
	ListenPort      uint                        `json:"listen_port,omitempty" example:"4242"`
	IsLighthouse    bool                        `json:"is_lighthouse,omitempty" example:"false"`
}

func (n *Host) BeforeCreate(tx *gorm.DB) error {
	n.ID = uuid.New()

	return nil
}
