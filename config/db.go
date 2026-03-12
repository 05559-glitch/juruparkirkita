package config

import (
	"arena-ban/internal/domain"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	
	godotenv.Load()
	host := os.Getenv("PSQL_HOST") 
	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASSWORD")
	dbName := os.Getenv("PSQL_DB_NAME")
	port := os.Getenv("PSQL_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	fmt.Println("Connected to Database!")
	return db
}


func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.Brand{},
		&domain.Category{},
		&domain.User{},
		&domain.Customer{},
		&domain.VehiclePlate{},
		&domain.Product{},
		&domain.Stock{},
		&domain.Service{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.RegisterVerification{},
		&domain.PasswordReset{},
	)

	if err != nil {
		log.Fatalf("Gagal melakukan migrasi: %v", err)
	}
	fmt.Println("✅ Database Migration Completed!")
}