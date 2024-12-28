package database

import "github.com/koodeyo/koodnet/pkg/models"

func Migrate() {
	Conn.AutoMigrate(&models.Network{})
	Conn.AutoMigrate(&models.Certificate{})
	Conn.AutoMigrate(&models.Host{})
	Conn.AutoMigrate(&models.Configuration{})
}
