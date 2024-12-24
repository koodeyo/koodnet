package database

import "github.com/koodeyo/koodnet/pkg/models"

func Migrate() {
	Conn.AutoMigrate(&models.Network{})
}
