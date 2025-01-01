package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"dario.cat/mergo"
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
	InPub           []byte         `json:"inPub,omitempty" swaggertype:"string"`
	NetworkID       uuid.UUID      `json:"networkId" gorm:"type:uuid"`
	Network         *Network       `json:"network,omitempty"`
	ConfigurationID uuid.UUID      `json:"configurationId" gorm:"type:uuid"`
	Configuration   *Configuration `json:"configuration,omitempty" gorm:"foreignKey:ConfigurationID;constraint:OnDelete:CASCADE"`
	Certificate     *Certificate   `json:"certificate,omitempty" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
}

type HostDto struct {
	Name            string         `json:"name,omitempty" example:"host-1"`
	IP              string         `json:"ip,omitempty" example:"100.100.0.1/24"`
	InPub           string         `json:"inPub,omitempty"`
	StaticAddresses []string       `json:"staticAddresses,omitempty" example:"109.243.69.39"`
	Subnets         []string       `json:"subnets,omitempty" example:"192.168.1.0/24"`
	Groups          []string       `json:"groups,omitempty" example:"laptop,servers,ssh"`
	NetworkID       uuid.UUID      `json:"networkId,omitempty" example:"c6d6c4c4-b65b-40e1-bcf2-1fd3122c653d"`
	Configuration   *Configuration `json:"configuration,omitempty"`
}

func (h *Host) GetIp() string {
	return strings.Split(h.IP, "/")[0]
}

func (h *Host) BeforeCreate(db *gorm.DB) error {
	h.ID = uuid.New()

	// Sign host
	if h.Certificate == nil {
		if err := h.Sign(db); err != nil {
			return err
		}
	}

	// Save default host config
	if h.Configuration != nil {
		mergo.Merge(h.Configuration, newConfig(), mergo.WithOverride, mergo.WithAppendSlice, mergo.WithOverrideEmptySlice)
	} else {
		h.Configuration = newConfig()
	}

	h.Configuration.ID = uuid.New()

	return nil
}

// Marshal serializes the Host configuration into either YAML or JSON format.
// Parameters:
//   - yml: if true, marshals to YAML; if false, marshals to JSON
//
// Returns:
//   - string: the marshaled configuration
//   - error: any error that occurred during marshaling
func (h *Host) Marshal(yml bool) (string, error) {
	if h == nil {
		return "", fmt.Errorf("cannot marshal nil host")
	}

	cfg := h.Configuration
	if err := mergo.Merge(cfg, newConfig(),
		mergo.WithOverride,
		mergo.WithOverrideEmptySlice); err != nil {
		return "", fmt.Errorf("failed to merge config: %w", err)
	}

	// Update PKI configuration
	cfg.PKI.CA = h.Network.MarshalCAs()
	cfg.PKI.Cert = string(h.Certificate.Crt)
	cfg.PKI.Key = string(h.Certificate.Key)

	var (
		cfgBytes []byte
		err      error
	)

	if yml {
		cfgBytes, err = yaml.Marshal(cfg)
	} else {
		cfgBytes, err = json.Marshal(cfg)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(cfgBytes), nil
}

func (h *Host) Sign(db *gorm.DB) error {
	// Attempt to find the network
	var n Network
	if err := db.Preload("Ca").First(&n, "id = ?", h.NetworkID).Error; err != nil {
		return errors.New("host network not found")
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

	h.Certificate = cert

	return nil
}
