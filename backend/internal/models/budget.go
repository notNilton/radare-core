// Package models defines the data structures (models) for the application.
package models

import (
	"gorm.io/gorm"
	"time"
)

// Transaction represents a single financial transaction recorded by a user.
// It is the core entity of the budget management system.
type Transaction struct {
	gorm.Model
	UserID      uint      `gorm:"not null;index"` // Foreign key to the User model.
	Description string    `gorm:"not null"`       // A brief description of the transaction.
	Amount      float64   `gorm:"not null"`       // The monetary value of the transaction.
	Date        time.Time `gorm:"not null"`       // The date and time when the transaction occurred.
	CategoryID  uint      `gorm:"not null;index"` // Foreign key to the Category model.
	Category    Category  `gorm:"foreignKey:CategoryID"` // Association with the Category model.
}

// Category represents a user-defined category for organizing transactions.
// This allows users to group expenses, such as "Groceries," "Utilities," or "Entertainment."
type Category struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"` // Foreign key to the User model, making categories user-specific.
	Name   string `gorm:"not null;uniqueIndex:idx_user_category_name"` // The name of the category (e.g., "Food", "Transport").
	User   User   `gorm:"foreignKey:UserID"` // Association with the User model.
}