// Package models define as estruturas de dados (modelos) da aplicação.
package models

import (
	"gorm.io/gorm"
	"time"
)

// Fueling representa um registro de abastecimento de um veículo.
// Cada registro está associado a um usuário.
type Fueling struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Cost      float64   `gorm:"not null"`
	FuelType  string    `gorm:"not null"`
	Location  string
	CarKM     float64
	Timestamp time.Time `gorm:"not null"`
}