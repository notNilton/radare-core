// Package handlers contains the HTTP request handlers for the API.
package handlers

import (
	"encoding/json"
	"net/http"
	"radare-datarecon/backend/internal/database"
	"radare-datarecon/backend/internal/middleware"
	"radare-datarecon/backend/internal/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CreateFuelingRequest defines the structure for creating a new fueling record.
type CreateFuelingRequest struct {
	Cost      float64   `json:"cost"`
	FuelType  string    `json:"fuel_type"`
	Location  string    `json:"location"`
	CarKM     float64   `json:"car_km"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateFuelingRecord creates a new fueling record for the authenticated user.
func CreateFuelingRecord(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid user ID in token"}
	}

	var req CreateFuelingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid request body: " + err.Error()}
	}

	fueling := models.Fueling{
		UserID:    uint(userID),
		Cost:      req.Cost,
		FuelType:  req.FuelType,
		Location:  req.Location,
		CarKM:     req.CarKM,
		Timestamp: req.Timestamp,
	}

	if result := database.DB.Create(&fueling); result.Error != nil {
		return middleware.HTTPError{Code: http.StatusInternalServerError, Message: "Error creating fueling record: " + result.Error.Error()}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fueling)
	return nil
}

// GetUserFuelingRecords retrieves all fueling records for the authenticated user.
func GetUserFuelingRecords(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid user ID in token"}
	}

	var fuelings []models.Fueling
	if result := database.DB.Where("user_id = ?", uint(userID)).Find(&fuelings); result.Error != nil {
		return middleware.HTTPError{Code: http.StatusInternalServerError, Message: "Error fetching fueling records: " + result.Error.Error()}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fuelings)
	return nil
}

// UpdateFuelingRecord updates an existing fueling record.
func UpdateFuelingRecord(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid user ID in token"}
	}

	vars := mux.Vars(r)
	recordID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid record ID"}
	}

	var fueling models.Fueling
	if result := database.DB.First(&fueling, recordID); result.Error != nil {
		return middleware.HTTPError{Code: http.StatusNotFound, Message: "Fueling record not found"}
	}

	if fueling.UserID != uint(userID) {
		return middleware.HTTPError{Code: http.StatusForbidden, Message: "You are not authorized to update this record"}
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid request body"}
	}

	if result := database.DB.Model(&fueling).Updates(updates); result.Error != nil {
		return middleware.HTTPError{Code: http.StatusInternalServerError, Message: "Error updating fueling record"}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fueling)
	return nil
}

// DeleteFuelingRecord deletes a fueling record.
func DeleteFuelingRecord(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid user ID in token"}
	}

	vars := mux.Vars(r)
	recordID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		return middleware.HTTPError{Code: http.StatusBadRequest, Message: "Invalid record ID"}
	}

	var fueling models.Fueling
	if result := database.DB.First(&fueling, recordID); result.Error != nil {
		return middleware.HTTPError{Code: http.StatusNotFound, Message: "Fueling record not found"}
	}

	if fueling.UserID != uint(userID) {
		return middleware.HTTPError{Code: http.StatusForbidden, Message: "You are not authorized to delete this record"}
	}

	if result := database.DB.Delete(&fueling); result.Error != nil {
		return middleware.HTTPError{Code: http.StatusInternalServerError, Message: "Error deleting fueling record"}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Fueling record deleted successfully"})
	return nil
}