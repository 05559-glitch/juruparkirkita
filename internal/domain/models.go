package domain

import (
	"time"
)


type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null;index"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time
}

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey"` 
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` 
	TokenHash string    `gorm:"type:varchar(64);unique;not null;index"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	IsUsed    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:current_timestamp"`
}

type RegisterVerification struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"type:varchar(100);not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Role      string    `gorm:"type:varchar(20);not null"`
	Token     string    `gorm:"type:varchar(255);not null;index"`
	IsUsed    bool      `gorm:"default:false"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
}


type Brand struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time
	Products  []Product `gorm:"foreignKey:BrandID"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID             uint      `gorm:"primaryKey"`
	SKU            string    `gorm:"type:varchar(50);unique;not null;index"`
	CategoryID     uint      `gorm:"not null"`
	BrandID        uint      `gorm:"not null"`
	Name           string    `gorm:"type:varchar(200);not null"`
	Price          float64   `gorm:"type:decimal(16,2);not null"`
	RackLocation   string    `gorm:"type:varchar(50)"`
	Specifications string    `gorm:"type:jsonb"` 
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time
	Stocks         []Stock   `gorm:"foreignKey:ProductID"`
}

type Stock struct {
	ID        uint      `gorm:"primaryKey"`
	ProductID uint      `gorm:"not null"`
	BatchCode string    `gorm:"type:varchar(50)"`
	Quantity  int       `gorm:"default:0"`
	BuyPrice  float64   `gorm:"type:decimal(16,2)"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time
}


type Customer struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null"`
	PhoneNumber string    `gorm:"type:varchar(20);unique"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time
	Vehicles    []VehiclePlate `gorm:"foreignKey:CustomerID"`
}

type VehiclePlate struct {
	ID           uint   `gorm:"primaryKey"`
	CustomerID   uint   `gorm:"not null"`
	PlateNumber  string `gorm:"type:varchar(20);unique;not null"`
	VehicleModel string `gorm:"type:varchar(100)"`
}


type Service struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(200);not null"`
	Price     float64   `gorm:"type:decimal(16,2);not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time
}

type Order struct {
	ID            uint      `gorm:"primaryKey"`
	InvoiceNumber string    `gorm:"type:varchar(50);unique;not null;index"`
	CashierID     uint      `gorm:"not null"`
	CustomerID    uint      `gorm:"not null"`
	Status        string    `gorm:"type:varchar(20);default:'PENDING'"`
	PaymentStatus string    `gorm:"type:varchar(20);default:'UNPAID'"`
	PaymentMethod string    `gorm:"type:varchar(50)"`
	CreatedAt     time.Time `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time
	Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	ItemType  string  `gorm:"type:varchar(20);not null"` 
	StockID   *uint  
	ServiceID *uint   
	Quantity  int     `gorm:"not null"`
	UnitPrice float64 `gorm:"type:decimal(16,2);not null"`
	CostPrice float64 `gorm:"type:decimal(16,2);not null"`
	Subtotal  float64 `gorm:"type:decimal(16,2);not null"`
}