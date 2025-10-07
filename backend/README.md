# Backend Documentation for BudgetBuddy

This document provides a comprehensive guide to the backend system for the BudgetBuddy personal finance application. The backend is a high-performance server built in Go, providing a robust REST API for user management, transaction tracking, and expense categorization.

## System Overview

The primary objective of this backend is to provide a secure and scalable API for a personal budgeting application. It allows users to register, log in, and manage their financial data, including tracking expenses and organizing them into personalized categories.

## Architectural Design

The backend is structured to separate concerns, promoting maintainability and scalability.

-   `main.go`: The application's entry point. It initializes the database connection, configures the HTTP server, registers API routes, and manages graceful shutdowns.
-   `internal/`: This directory houses the core logic of the application.
    -   `config/`: Manages application configuration from environment variables.
    -   `database/`: Handles the database connection (PostgreSQL) and GORM setup.
    -   `handlers/`: Contains HTTP request handlers for user-related endpoints (registration, login).
    -   `middleware/`: Provides middleware for logging, error handling, and JWT-based authentication.
    -   `models/`: Defines the data structures for `User`, `Category`, and `Transaction`.
    -   `budget/`: Implements the core business logic and handlers for all budget-related features (transactions, categories, summaries).
-   `go.mod`, `go.sum`: These files manage the project's dependencies.
-   `Dockerfile`: Defines the container for deploying the backend service.

## API Endpoints

All endpoints are prefixed with `/api`. Authenticated endpoints require a `Bearer` token in the `Authorization` header.

### User Management

-   `POST /register`: Registers a new user.
-   `POST /login`: Authenticates a user and returns a JWT.

### Categories

-   `POST /categories`: Creates a new expense category for the authenticated user.
-   `GET /categories`: Retrieves all categories for the authenticated user.

### Transactions

-   `POST /transactions`: Creates a new financial transaction for the authenticated user.
-   `GET /transactions`: Retrieves all transactions for the authenticated user.

### Budget Summary

-   `GET /budget/summary`: Calculates the total expenses for a given month and year (e.g., `/api/budget/summary?year=2024&month=7`).

### Health Check

-   `GET /healthz`: Checks the health of the server.

## Setup and Execution

### Prerequisites

-   Go 1.18 or later
-   A running PostgreSQL instance

### Steps

1.  **Navigate to the backend directory:**
    ```bash
    cd backend
    ```

2.  **Create a `.env` file:**
    Create a `.env` file in the `backend` directory with the necessary environment variables for your database connection. See `internal/config/config.go` for the required variables (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `JWT_SECRET`).

3.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Run the server:**
    ```bash
    go run main.go
    ```
    The server will start on port `8080` by default. You can set the `PORT` environment variable to use a different port.

## Running Tests

To run the backend's unit tests, navigate to the `backend` directory and execute the following command:
```bash
go test ./...
```