package models

import (
	"gorm.io/gorm"
)

// Used within the Nebula networking tool to retrieve
// detailed information about a specific host and its associated network. This includes:
// - The host's configuration and certificate.
// - The network's certificate authority (CA).
// - Other hosts in the network (excluding the current host), but only those configured
// as lighthouses or relays based on their configuration.
func PreloadHostWithFullDetails(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Configuration").
			Preload("Certificate").
			Preload("Network.Ca").
			Preload("Network.Hosts", func(db *gorm.DB) *gorm.DB {
				return db.Where("hosts.id != ?", id). // Ignore it, self
									Preload("Configuration").
									Joins("JOIN configurations ON configurations.id = hosts.configuration_id").
									Where("configurations.lighthouse_am_lighthouse = ? OR configurations.relay_am_relay = ?", true, true)
			}).
			Where("hosts.id = ?", id)
	}
}
