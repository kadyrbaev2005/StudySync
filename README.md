# StudySync - Student Task Management System

A robust Go-based backend web service for managing student tasks, deadlines, and study schedules.

## 🚀 Features

- **JWT Authentication & Authorization** with role-based access (Admin/User)
- **CRUD Operations** for Users, Subjects, Tasks, Deadlines, and Sprints
- **Email Notifications** via SMTP for upcoming deadlines
- **Redis Caching** for improved performance
- **Database Migrations** using golang-migrate
- **Structured Logging** with file rotation
- **Swagger API Documentation**
- **Integration & Unit Tests**
- **Docker Support** for easy deployment
- **Graceful Shutdown** and background workers

## 🏗️ Architecture

- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **Redis** - Caching layer
- **JWT** - Authentication tokens
- **Docker** - Containerization
- **golang-migrate** - Database migrations

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

## 🛠️ Installation

### Option 1: Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/kadyrbayev2005/studysync.git
cd studysync

# Start services using Docker Compose
docker-compose up -d

# Updating DB migrations
migrate -path migrations -database "postgres://postgres:password@localhost:5433/studysync?sslmode=disable" up

# Execute the program
go run cmd/main.go

# The API will be available at http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
```
### Option 2: Using MakerFile
```bash
# Clone the repository
git clone https://github.com/kadyrbayev2005/studysync.git
cd studysync

# Start services using Docker Compose
make docker up

# Updating DB migrations
make migrate-up

# Execute the program
make run

# For any other commands
make help
```