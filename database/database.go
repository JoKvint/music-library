package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// ConnectDB initializes and returns the PostgreSQL database connection.
// It uses the environment variables to configure the connection.
func ConnectDB() {
	// Используем sync.Once, чтобы подключиться к базе данных только один раз
	once.Do(func() {
		// Формируем строку подключения из переменных окружения
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"))

		var err error
		// Открываем подключение к базе данных
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		log.Println("Database connected successfully!") // Успешное подключение
	})
}

// GetDB returns the instance of the database connection.
func GetDB() *gorm.DB {
	if db == nil {
		ConnectDB() // Инициализируем подключение, если оно еще не установлено
	}
	return db
}
