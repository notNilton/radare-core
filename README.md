# BudgetBuddy: A Personal Budget Management System

BudgetBuddy is a streamlined and efficient personal finance management tool designed to help you take control of your expenses. It provides a secure, API-driven backend for tracking transactions, organizing them into categories, and analyzing your spending habits over time.

## Core Features

-   **Secure User Management**: Full support for user registration and JWT-based authentication to ensure your financial data remains private.
-   **Transaction Tracking**: A robust API for creating, retrieving, and managing all your financial transactions.
-   **Expense Categorization**: Organize your spending by creating custom categories, allowing for a clearer understanding of where your money goes.
-   **Monthly Summaries**: Automatically calculate and retrieve total expenses for any given month, helping you stay on top of your budget.
-   **API-First Architecture**: A powerful Go-based backend API enables seamless integration with potential future applications, such as a web or mobile client.

This `README.md` provides a comprehensive overview of the BudgetBuddy architecture, its components, and detailed instructions for getting started with the backend and API client.

## System Architecture

BudgetBuddy is designed with a modern, decoupled architecture, consisting of two primary components relevant to this refactoring:

-   **Backend**: A high-performance Go-based API that serves as the core of the system, handling all business logic for users, categories, and transactions.
-   **API Client**: A [Bruno](https://www.usebruno.com/) collection that facilitates direct interaction with the backend API for testing, development, and integration purposes.

## Getting Started

### Backend

The backend is a Go application that exposes a RESTful API for personal budget management.

**To run the backend server:**

1.  Navigate to the `backend` directory.
2.  Ensure you have a `.env` file configured with your database credentials (see `backend/internal/config/config.go` for required variables).
3.  Run the application:
    ```bash
    cd backend
    go run main.go
    ```

The server will start on port `8080` by default. This can be configured via the `PORT` environment variable.

### API Client

The `apiclient` directory contains a Bruno collection for interacting with the backend API. This is an invaluable tool for developers who need to test the API endpoints. The collection includes requests for:

-   User Registration & Login
-   Creating & Viewing Categories
-   Creating & Viewing Transactions
-   Retrieving Monthly Summaries

## Contributing

We welcome contributions to the BudgetBuddy project. If you have ideas for new features, bug fixes, or enhancements, please open an issue or submit a pull request on our GitHub repository.