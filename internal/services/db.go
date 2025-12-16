package services

import (
    "fmt"
    "time"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/kadyrbayev2005/studysync/internal/utils"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
    // Read database configuration from environment variables with defaults
    user := utils.GetEnv("DB_USER", "postgres")
    password := utils.GetEnv("DB_PASSWORD", "postgres")
    dbname := utils.GetEnv("DB_NAME", "studysync")
    host := utils.GetEnv("DB_HOST", "localhost")
    port := utils.GetEnv("DB_PORT", "5432")

    // Build DSN for PostgreSQL
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        user, password, host, port, dbname,
    )

    // Run migrations
    m, err := migrate.New(
        "file://migrations",
        dsn,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create migrate instance: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

    // Connect with GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        PrepareStmt: true,
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })

    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    // Set connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    Info("✅ Connected to database and migrations applied successfully")

    return db, nil
}