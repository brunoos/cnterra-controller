package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/brunoos/cnterra-controller/config"
)

var DB *gorm.DB

func Initialize() {
	dbParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DbAddress, config.DbPort, config.DbUser, config.DbPassword, config.DbName, config.DbSslMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dbParams), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ERRO] Error openning database connection: %s", err)
	}
}
