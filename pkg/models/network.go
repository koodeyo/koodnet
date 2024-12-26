package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Model
type Network struct {
	ID               string                      `json:"id" gorm:"type:text;primary_key;"`                                   // Unique identifier for the network (UUID).
	Name             string                      `json:"name" gorm:"size:255;uniqueIndex:idx_name_cidr"`                     // Name of the network, must be unique in combination with the CIDR.
	IPs              datatypes.JSONSlice[string] `json:"ips" gorm:"type:json;default:'[]'" swaggertype:"array,string"`       // List of IPv4 addresses and networks in CIDR notation. Limits the addresses for subordinate certificates.
	Subnets          datatypes.JSONSlice[string] `json:"subnets" gorm:"type:json;default:'[]'" swaggertype:"array,string"`   // List of IPv4 subnets in CIDR notation. Defines subnets that subordinate certificates can use.
	Groups           datatypes.JSONSlice[string] `json:"groups" gorm:"type:json;default:'[]'" swaggertype:"array,string"`    // List of groups for access control, restricting subordinate certificates' groups.
	Encrypt          bool                        `json:"encrypt" gorm:"default:false"`                                       // Enables passphrase encryption for private keys. Default: true.
	Passphrase       string                      `json:"passphrase" gorm:"size:255"`                                         // Passphrase used for encrypting the private key.
	ArgonMemory      uint                        `json:"argon_memory" gorm:"default:2097152"`                                // Argon2 memory parameter in KiB for encrypted private key passphrase. Default: 2 MiB. (2*1024*1024)
	ArgonIterations  uint                        `json:"argon_iterations" gorm:"default:2"`                                  // Number of Argon2 iterations for encrypting private key passphrase. Default: 2.
	ArgonParallelism uint                        `json:"argon_parallelism" gorm:"default:4"`                                 // Argon2 parallelism parameter for encrypting private key passphrase. Default: 4.
	Curve            string                      `json:"curve" gorm:"default:25519"`                                         // Cryptographic curve for key generation. Options include "25519" (default) and "P256".
	Duration         time.Duration               `json:"duration" gorm:"default:17531" swaggertype:"number"`                 // Certificate validity duration. Default: 2 years (17,531 hours). (time.Duration(time.Hour*8760))
	Ca               []Certificate               `json:"ca,omitempty" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE;"` // Associated Certificate Authorities (CA) for the network.
	CreatedAt        time.Time                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time                   `json:"updated_at" gorm:"autoUpdateTime"`
}

// DTO for create/update operations
type NetworkDto struct {
	Name             string                      `json:"name" example:"my-network"`
	IPs              datatypes.JSONSlice[string] `json:"ips" example:"100.100.0.0/22" swaggertype:"array,string"`
	Subnets          datatypes.JSONSlice[string] `json:"subnets" example:"192.168.1.0/24" swaggertype:"array,string"`
	Groups           datatypes.JSONSlice[string] `json:"groups" example:"laptop,ssh,servers" swaggertype:"array,string"`
	Duration         time.Duration               `json:"duration" example:"17531" swaggertype:"number"`
	Encrypt          bool                        `json:"encrypt" example:"false"`
	Passphrase       string                      `json:"passphrase" example:"orange-duck-walks-happy-sunset-92"`
	ArgonMemory      uint                        `json:"argon_memory" example:"2097152"`
	ArgonIterations  uint                        `json:"argon_iterations" example:"2"`
	ArgonParallelism uint                        `json:"argon_parallelism" example:"4"`
	Curve            string                      `json:"curve" example:"25519" enums:"25519,X25519,Curve25519,CURVE25519,P256"`
}

// Hooks
func (n *Network) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}

	if len(n.Ca) == 0 {
		ca, err := n.NewCA()
		if err != nil {
			return err
		}
		n.Ca = append(n.Ca, *ca)
	}

	return nil
}

func (n *Network) BeforeSave(tx *gorm.DB) error {
	if err := n.validate(); err != nil {
		return err
	}

	return nil
}
