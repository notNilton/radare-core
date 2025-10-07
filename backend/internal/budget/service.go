// Package budget contains the business logic for managing budgets, transactions, and categories.
package budget

import (
	"fmt"
	"net/http"
	"radare-datarecon/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// Service provides the core business logic for budget management.
// It encapsulates all operations related to transactions and categories.
type Service struct {
	DB *gorm.DB
}

// NewService creates a new budget service with a database connection.
func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// --- Category Management ---

// CreateCategory creates a new category for a specific user.
func (s *Service) CreateCategory(userID uint, name string) (*models.Category, error) {
	category := &models.Category{
		UserID: userID,
		Name:   name,
	}
	if err := s.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// GetCategoriesByUser retrieves all categories belonging to a specific user.
func (s *Service) GetCategoriesByUser(userID uint) ([]models.Category, error) {
	var categories []models.Category
	if err := s.DB.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// --- Transaction Management ---

// CreateTransaction creates a new transaction for a specific user.
func (s *Service) CreateTransaction(userID uint, description string, amount float64, date time.Time, categoryID uint) (*models.Transaction, error) {
	transaction := &models.Transaction{
		UserID:      userID,
		Description: description,
		Amount:      amount,
		Date:        date,
		CategoryID:  categoryID,
	}
	if err := s.DB.Create(transaction).Error; err != nil {
		return nil, err
	}
	// Preload the category to return it in the response.
	s.DB.Preload("Category").First(transaction, transaction.ID)
	return transaction, nil
}

// GetTransactionsByUser retrieves all transactions for a specific user, with optional filtering.
func (s *Service) GetTransactionsByUser(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := s.DB.Preload("Category").Where("user_id = ?", userID).Order("date desc").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetMonthlySummary calculates the total expenses for a user in a given month and year.
func (s *Service) GetMonthlySummary(userID uint, year int, month time.Month) (float64, error) {
	var total float64
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	err := s.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date < ?", userID, startOfMonth, endOfMonth).
		Select("sum(amount)").
		Row().
		Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetUserIDFromContext is a helper function to extract the user's ID from the request context.
// This is crucial for ensuring that users can only access their own data.
func GetUserIDFromContext(r *http.Request) (uint, error) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context or is of an invalid type")
	}
	return userID, nil
}