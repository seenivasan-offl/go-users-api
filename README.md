Go Users API – User with DOB and Calculated Age
RESTful backend service built in Go to manage users with their name and date of birth (DOB). The API stores DOB in PostgreSQL and calculates a user’s age dynamically whenever user details are fetched.

Table of Contents
Overview

Architecture

Tech Stack

Project Structure

Getting Started

Prerequisites

Setup

Running the Server

Running Tests

API Endpoints

Design Decisions

Future Improvements

Overview
This project implements a Go backend service for managing users and their date of birth.
Key behaviour: the API does not store age in the database. Instead, it calculates age on the fly using Go’s time package whenever a user (or list of users) is fetched.

The implementation follows a layered architecture to keep HTTP, business logic, and database access clearly separated and easy to maintain.

Architecture
High-level flow:

HTTP Request → Middleware → Handler → Service → Repository → SQLC (generated) → PostgreSQL

Handler layer: HTTP concerns (routing, JSON input/output, HTTP status codes).

Service layer: Business logic (DOB parsing, age calculation, mapping to response models).

Repository layer: Database access via SQLC-generated queries.

SQLC: Type-safe, compiled SQL query layer for PostgreSQL.

This structure keeps concerns separate and makes the codebase easier to test and evolve.

Tech Stack
Language: Go (Golang)

Web Framework: Fiber

Database: PostgreSQL

DB Access Layer: SQLC (PostgreSQL + pgx/v5)

Logging: Uber Zap

Validation: go-playground/validator

DB Driver: pgx/v5 (via pgxpool)

Project Structure
text
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
Getting Started
Prerequisites
Go 1.21+ installed

Docker (for running SQLC via container) or native SQLC installed

PostgreSQL running locally

A database named users_db with the users table:

sql
CREATE TABLE IF NOT EXISTS users (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob  DATE NOT NULL
);
Setup
Clone the repository

bash
git clone https://github.com/<your-username>/go-users-api.git
cd go-users-api
Generate SQLC code

Using Docker (recommended):

bash
docker run --rm -v "$PWD:/src" -w /src sqlc/sqlc generate
Or if you have sqlc installed natively:

bash
sqlc generate
Set environment variables

On Windows (cmd):

text
set DATABASE_URL=postgres://postgres:postgres@localhost:5432/users_db?sslmode=disable
set SERVER_ADDR=:8080
On Unix shells:

bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/users_db?sslmode=disable
export SERVER_ADDR=:8080
Adjust username, password, and database name as needed.

Running the Server
From the project root:

bash
go mod tidy
go run ./cmd/server
The server will start and listen on SERVER_ADDR (default :8080).

Running Tests
To run only the age calculation tests:

bash
go test ./internal/service -run TestCalculateAge -v
To run all tests:

bash
go test ./... -v
API Endpoints
Base URL: http://localhost:8080

1. Create User
POST /users

Request body:

json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
Response 201 Created:

json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10T00:00:00Z"
}
Validated with go-playground/validator

dob must be in YYYY-MM-DD format

2. Get User by ID (with dynamic age)
GET /users/:id

Response 200 OK:

json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10T00:00:00Z",
  "age": 35
}
age is calculated dynamically using Go’s time package

No age column exists in the database

3. Update User
PUT /users/:id

Request body:

json
{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
Response 200 OK:

json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15T00:00:00Z"
}
4. Delete User
DELETE /users/:id

Response:

204 No Content on success

5. List All Users (with dynamic age)
GET /users

Response 200 OK:

json
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10T00:00:00Z",
    "age": 34
  }
]
Age is computed per user using their dob and current date.

Design Decisions
Dynamic Age Calculation
Age is computed using a pure function in the service layer:

go
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
This avoids storing redundant data and ensures age is always up to date.

The function is unit-tested with multiple edge cases (birthday today, before/after birthday, leap years, future DOB).

SQLC + PostgreSQL
All DB access is defined in .sql files and compiled into Go code by SQLC.

This gives:

Compile-time type checking of queries.

Centralised SQL definitions.

No manual rows.Scan or string-building.

Validation
Input models use validator tags to enforce:

Required fields

Name length

Correct date format (YYYY-MM-DD)

Logging and Middleware
Zap is used for structured, leveled logging.

Middleware:

Injects a X-Request-ID header for each request.

Logs method, path, status, and request latency.

This matches production logging best practices and simplifies debugging.

Future Improvements
Add pagination parameters to GET /users (e.g. ?page=1&size=10).

Add OpenAPI/Swagger documentation for the API.

Add Docker Compose for fully containerised setup (app + Postgres).

Extend validation with custom messages or localization.
