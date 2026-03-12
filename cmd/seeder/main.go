package main

import (
	"arena-ban/config"
	"arena-ban/internal/domain"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	db := config.ConnectDB()

	password := "admin123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	admin := domain.User{
		Name: "admin",
		Email:    "admin@arenaban.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	err := db.Create(&admin).Error
	if err != nil {
		log.Fatal("Gagal membuat data seeder: ", err)
	}

	fmt.Println("✅ Seeder Berhasil!")
	fmt.Printf("Username: %s\nPassword: %s\n", admin.Name, password)
}