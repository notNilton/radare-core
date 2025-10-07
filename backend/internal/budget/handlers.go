// Package budget contains the handlers for the budget management API.
package budget

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Handlers encapsulates the budget service for use in HTTP handlers.
type Handlers struct {
	service *Service
}

// NewHandlers creates a new set of budget handlers.
func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// --- DTOs (Data Transfer Objects) ---

// CreateCategoryRequest defines the expected JSON for creating a category.
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// CreateTransactionRequest defines the expected JSON for creating a transaction.
type CreateTransactionRequest struct {
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	CategoryID  uint      `json:"category_id"`
}

// --- Combined Handlers ---

// HandleCategories routes requests for categories based on the HTTP method.
// GET /api/categories -> getCategories
// POST /api/categories -> createCategory
func (h *Handlers) HandleCategories(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return h.getCategories(w, r)
	case http.MethodPost:
		return h.createCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}
}

// HandleTransactions routes requests for transactions based on the HTTP method.
// GET /api/transactions -> getTransactions
// POST /api/transactions -> createTransaction
func (h *Handlers) HandleTransactions(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return h.getTransactions(w, r)
	case http.MethodPost:
		return h.createTransaction(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}
}

// --- Category Handler Logic ---

// createCategory handles the creation of a new category for the logged-in user.
func (h *Handlers) createCategory(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: User ID not found.", http.StatusUnauthorized)
		return nil
	}

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return nil
	}

	if req.Name == "" {
		http.Error(w, "Category name cannot be empty.", http.StatusBadRequest)
		return nil
	}

	category, err := h.service.CreateCategory(userID, req.Name)
	if err != nil {
		http.Error(w, "Failed to create category.", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(category)
}

// getCategories handles fetching all categories for the logged-in user.
func (h *Handlers) getCategories(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: User ID not found.", http.StatusUnauthorized)
		return nil
	}

	categories, err := h.service.GetCategoriesByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve categories.", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(categories)
}

// --- Transaction Handler Logic ---

// createTransaction handles the creation of a new transaction for the logged-in user.
func (h *Handlers) createTransaction(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: User ID not found.", http.StatusUnauthorized)
		return nil
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return nil
	}

	transaction, err := h.service.CreateTransaction(userID, req.Description, req.Amount, req.Date, req.CategoryID)
	if err != nil {
		http.Error(w, "Failed to create transaction.", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(transaction)
}

// getTransactions handles fetching all transactions for the logged-in user.
func (h *Handlers) getTransactions(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: User ID not found.", http.StatusUnauthorized)
		return nil
	}

	transactions, err := h.service.GetTransactionsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve transactions.", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(transactions)
}

// GetMonthlySummary handles calculating and returning the budget summary for a given month.
func (h *Handlers) GetMonthlySummary(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized: User ID not found.", http.StatusUnauthorized)
		return nil
	}

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid 'year' query parameter.", http.StatusBadRequest)
		return nil
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Invalid 'month' query parameter.", http.StatusBadRequest)
		return nil
	}

	total, err := h.service.GetMonthlySummary(userID, year, time.Month(month))
	if err != nil {
		http.Error(w, "Failed to calculate monthly summary.", http.StatusInternalServerError)
		return err
	}

	response := map[string]float64{"total_expenses": total}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}