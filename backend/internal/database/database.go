package database

import (
	"fmt"
	"log"
	"radare-datarecon/backend/internal/config"
	"radare-datarecon/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSslMode,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Fueling{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
}