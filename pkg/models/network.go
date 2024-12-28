package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model
type Network struct {
	ID               uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;"`                                  // Unique identifier for the network (UUID).
	Name             string        `json:"name" gorm:"size:255;uniqueIndex:idx_name_cidr"`                    // Name of the network, must be unique in combination with the CIDR.
	IPs              []string      `json:"ips" gorm:"serializer:json;default:'[]'"`                           // List of IPv4 addresses and networks in CIDR notation. Limits the addresses for subordinate certificates.
	Subnets          []string      `json:"subnets" gorm:"serializer:json;default:'[]'"`                       // List of IPv4 subnets in CIDR notation. Defines subnets that subordinate certificates can use.
	Groups           []string      `json:"groups" gorm:"serializer:json;default:'[]'"`                        // List of groups for access control, restricting subordinate certificates' groups.
	Encrypt          bool          `json:"encrypt" gorm:"default:false"`                                      // Enables passphrase encryption for private keys. Default: true.
	Passphrase       string        `json:"passphrase" gorm:"size:255"`                                        // Passphrase used for encrypting the private key.
	ArgonMemory      uint          `json:"argonMemory" gorm:"default:2097152"`                                // Argon2 memory parameter in KiB for encrypted private key passphrase. Default: 2 MiB. (2*1024*1024)
	ArgonIterations  uint          `json:"argonIterations" gorm:"default:2"`                                  // Number of Argon2 iterations for encrypting private key passphrase. Default: 2.
	ArgonParallelism uint          `json:"argonParallelism" gorm:"default:4"`                                 // Argon2 parallelism parameter for encrypting private key passphrase. Default: 4.
	Curve            string        `json:"curve" gorm:"default:25519"`                                        // Cryptographic curve for key generation. Options include "25519" (default) and "P256".
	Duration         time.Duration `json:"duration" gorm:"default:17531" swaggertype:"number"`                // Certificate validity duration. Default: 2 years (17,531 hours). (time.Duration(time.Hour*8760))
	Ca               []Certificate `json:"ca,omitempty" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"` // Associated Certificate Authorities (CA) for the network.
	Hosts            []Host        `json:"hosts,omitempty" gorm:"constraint:OnDelete:CASCADE"`                // Associated hosts for the network.
	CreatedAt        time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
}

// DTO for create/update operations
type NetworkDto struct {
	Name             string        `json:"name,omitempty" example:"my-network"`
	IPs              []string      `json:"ips,omitempty" example:"100.100.0.0/22"`
	Subnets          []string      `json:"subnets,omitempty" example:"192.168.1.0/24"`
	Groups           []string      `json:"groups,omitempty" example:"laptop,ssh,servers"`
	Duration         time.Duration `json:"duration,omitempty" example:"17531" swaggertype:"number"`
	Encrypt          bool          `json:"encrypt,omitempty" example:"false"`
	Passphrase       string        `json:"passphrase" example:"orange-duck-walks-happy-sunset-92"`
	ArgonMemory      uint          `json:"argonMemory,omitempty" example:"2097152"`
	ArgonIterations  uint          `json:"argonIterations,omitempty" example:"2"`
	ArgonParallelism uint          `json:"argonParallelism,omitempty" example:"4"`
	Curve            string        `json:"curve,omitempty" example:"25519" enums:"25519,X25519,Curve25519,CURVE25519,P256"`
}

func (n *Network) ValidCAs() []Certificate {
	var validCAs []Certificate

	for _, ca := range n.Ca {
		if ca.Expired() {
			continue
		}

		validCAs = append(validCAs, ca)
	}

	return validCAs
}

// Hooks
func (n *Network) BeforeCreate(tx *gorm.DB) error {
	n.ID = uuid.New()

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
