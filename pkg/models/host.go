package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type Host struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;"`
	Name            string         `json:"name" gorm:"size:255;not null;uniqueIndex:idx_name_network"`
	IP              string         `json:"ip" gorm:"size:255;not null;uniqueIndex:idx_ip_network"`
	StaticAddresses []string       `json:"staticAddresses" gorm:"serializer:json;default:'[]'"`
	Subnets         []string       `json:"subnets" gorm:"serializer:json;default:'[]'"`
	Groups          []string       `json:"groups" gorm:"serializer:json;default:'[]'"`
	ListenPort      uint           `json:"listenPort" gorm:"default:4242"`
	IsLighthouse    bool           `json:"isLighthouse" gorm:"default:false"`
	InPub           []byte         `json:"inPub,omitempty" swaggertype:"string"`
	NetworkID       uuid.UUID      `json:"networkId" gorm:"type:uuid"`
	Network         *Network       `json:"network,omitempty"`
	Configuration   *Configuration `json:"configuration,omitempty" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"`
	Certificates    []Certificate  `json:"certificates,omitempty" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
}

type HostDto struct {
	Name            string         `json:"name,omitempty" example:"host-1"`
	IP              string         `json:"ip,omitempty" example:"100.100.0.1/24"`
	InPub           string         `json:"inPub"`
	StaticAddresses []string       `json:"staticAddresses,omitempty" example:"109.243.69.39"`
	Subnets         []string       `json:"subnets,omitempty" example:"192.168.1.0/24"`
	Groups          []string       `json:"groups,omitempty" example:"laptop,servers,ssh"`
	NetworkID       uuid.UUID      `json:"networkId,omitempty" example:"c6d6c4c4-b65b-40e1-bcf2-1fd3122c653d"`
	ListenPort      uint           `json:"listenPort,omitempty" example:"4242"`
	IsLighthouse    bool           `json:"isLighthouse,omitempty" example:"false"`
	Configuration   *Configuration `json:"configuration,omitempty"`
}

func (h *Host) BeforeCreate(db *gorm.DB) error {
	h.ID = uuid.New()

	if len(h.Certificates) == 0 {
		if err := h.Sign(db); err != nil {
			return err
		}
	}

	return nil
}

func (h *Host) NebulaConfig() (string, error) {
	cfg := NewConfig()

	cfg.Listen.Port = uint(h.ListenPort)
	cfg.Lighthouse.AmLighthouse = h.IsLighthouse

	finalConfig, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	return string(finalConfig), nil
}

func (h *Host) Sign(db *gorm.DB) error {
	// Attempt to find the network
	var n Network
	if err := db.Preload("Ca").First(&n, "id = ?", h.NetworkID).Error; err != nil {
		return errors.New("host is not associated with any network")
	}

	CAs := n.ValidCAs()

	// Validate CA is present in the network
	if len(CAs) == 0 {
		return errors.New("no CA certificates found for the network")
	}

	ca := CAs[len(CAs)-1]

	cert, err := h.NewCert(ca)
	if err != nil {
		return err
	}

	h.Certificates = append(h.Certificates, *cert)

	return nil
}
