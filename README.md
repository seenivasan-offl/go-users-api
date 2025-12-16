# Go Users API – User with DOB and Calculated Age

A RESTful backend service built in **Go** to manage users with their **name** and **date of birth (DOB)**.
The API stores DOB in **PostgreSQL** and calculates a user’s **age dynamically** whenever user details are fetched.

---

## 📑 Table of Contents

* [Overview](#overview)
* [Architecture](#architecture)
* [Tech Stack](#tech-stack)
* [Project Structure](#project-structure)
* [Getting Started](#getting-started)

  * [Prerequisites](#prerequisites)
  * [Setup](#setup)
  * [Running the Server](#running-the-server)
  * [Running Tests](#running-tests)
* [API Endpoints](#api-endpoints)
* [Design Decisions](#design-decisions)
* [Future Improvements](#future-improvements)

---

## Overview

This project implements a Go backend service for managing users and their date of birth.

**Key behaviour:** the API does **not** store age in the database. Instead, it calculates age on the fly using Go’s `time` package whenever a user (or list of users) is fetched.

The implementation follows a **layered architecture** to keep HTTP, business logic, and database access clearly separated and easy to maintain.

---

## Architecture

**High-level flow:**

```
HTTP Request → Middleware → Handler → Service → Repository → SQLC (generated) → PostgreSQL
```

**Layers:**

* **Handler layer**: HTTP concerns (routing, JSON input/output, HTTP status codes)
* **Service layer**: Business logic (DOB parsing, age calculation, mapping to response models)
* **Repository layer**: Database access via SQLC-generated queries
* **SQLC**: Type-safe, compiled SQL query layer for PostgreSQL

This structure keeps concerns separate and makes the codebase easier to test and evolve.

---

## Tech Stack

* **Language**: Go (Golang)
* **Web Framework**: Fiber
* **Database**: PostgreSQL
* **DB Access Layer**: SQLC (PostgreSQL + pgx/v5)
* **Logging**: Uber Zap
* **Validation**: go-playground/validator
* **DB Driver**: pgx/v5 (via pgxpool)

---

## Project Structure

```
go-users-api/
├── cmd/
│   └── server/
│       └── main.go            # Application entrypoint
│
├── config/
│   └── config.go              # Environment-based configuration
│
├── db/
│   ├── migrations/
│   │   ├── 000001_create_users_table.up.sql
│   │   └── 000001_create_users_table.down.sql
│   ├── query/
│   │   └── users.sql          # SQLC queries (CRUD for users)
│   └── sqlc/                  # SQLC-generated Go code (do not edit)
│
└── internal/
    ├── handler/
    │   └── user_handler.go    # HTTP handlers for /users endpoints
    ├── repository/
    │   └── user_repository.go # DB repository using SQLC
    ├── service/
    │   ├── user_service.go    # Business logic + age calculation
    │   └── age_test.go        # Unit tests for age calculation
    ├── routes/
    │   └── routes.go          # Route registration
    ├── middleware/
    │   ├── request_id.go      # Injects X-Request-ID
    │   └── request_logger.go  # Logs request duration & details
    ├── models/
    │   ├── user.go            # DTOs / API models
    │   └── validator.go       # Validator initialization
    └── logger/
        └── logger.go          # Zap logger configuration
```

---

## Getting Started

### Prerequisites

* Go **1.21+** installed
* Docker (for running SQLC via container) **or** native SQLC installed
* PostgreSQL running locally

Database schema:

```sql
CREATE TABLE IF NOT EXISTS users (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob  DATE NOT NULL
);
```

---

### Setup

**1. Clone the repository**

```bash
git clone https://github.com/<your-username>/go-users-api.git
cd go-users-api
```

**2. Generate SQLC code**

Using Docker (recommended):

```bash
docker run --rm -v "$PWD:/src" -w /src sqlc/sqlc generate
```

Or using native SQLC:

```bash
sqlc generate
```

**3. Set environment variables**

On Windows (cmd):

```text
set DATABASE_URL=postgres://postgres:postgres@localhost:5432/users_db?sslmode=disable
set SERVER_ADDR=:8080
```

On Unix shells:

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/users_db?sslmode=disable
export SERVER_ADDR=:8080
```

Adjust credentials as needed.

---

### Running the Server

```bash
go mod tidy
go run ./cmd/server
```

The server listens on `SERVER_ADDR` (default `:8080`).

---

### Running Tests

Run only age calculation tests:

```bash
go test ./internal/service -run TestCalculateAge -v
```

Run all tests:

```bash
go test ./... -v
```

---

## API Endpoints

**Base URL:** `http://localhost:8080`

### 1. Create User

**POST /users**

Request:

```json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

Response **201 Created**:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10T00:00:00Z"
}
```

---

### 2. Get User by ID (with dynamic age)

**GET /users/:id**

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10T00:00:00Z",
  "age": 35
}
```

---

### 3. Update User

**PUT /users/:id**

```json
{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

---

### 4. Delete User

**DELETE /users/:id**

Response: **204 No Content**

---

### 5. List All Users (with dynamic age)

**GET /users**

```json
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10T00:00:00Z",
    "age": 34
  }
]
```

---

## Design Decisions

### Dynamic Age Calculation

```go
func calculateAge(dob, now time.Time) int {
    years := now.Year() - dob.Year()
    if now.Month() < dob.Month() ||
       (now.Month() == dob.Month() && now.Day() < dob.Day()) {
        years--
    }
    if years < 0 {
        return 0
    }
    return years
}
```

* Avoids storing redundant data
* Ensures age is always up to date
* Covered by unit tests (edge cases included)

---

### SQLC + PostgreSQL

* SQL defined in `.sql` files
* Compile-time type safety
* No manual `rows.Scan`

---

### Logging & Middleware

* **Zap** for structured logging
* Middleware:

  * Injects `X-Request-ID`
  * Logs request method, path, status, and duration

---

## Future Improvements

* Pagination for `GET /users`
* OpenAPI / Swagger documentation
* Docker Compose (App + Postgres)
* Enhanced validation messages
