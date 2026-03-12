package domain

import "github.com/golang-jwt/jwt/v5"

type UserRole string

const (
	ADMIN     UserRole = "ADMIN"
	CASHIER   UserRole = "CASHIER"
	CS        UserRole = "CS"
	WAREHOUSE UserRole = "WAREHOUSE"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"` // Disamakan menjadi min=8
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateAccountRequest struct {
	Name  string   `json:"name" validate:"required"`
	Email string   `json:"email" validate:"required,email"`
	Role  UserRole `json:"role" validate:"required,oneof=ADMIN CASHIER WAREHOUSE CS"`
}

type CreatePasswordRequest struct {
	Token    string `json:"token" validate:"required"` // Didapat dari mobile app saat menangkap Deep Link
	Password string `json:"password" validate:"required,min=8"` // Disamakan menjadi min=8
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// --- PASSWORD MANAGEMENT ---

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	// Menghapus min=8 karena validasi 'email' sudah cukup dan lebih aman
	Email string `json:"email" validate:"required,email"`
}

type VerifyResetPasswordRequest struct {
	// REVISI FATAL: Mengganti Email menjadi Token
	// Saat user mengklik link reset dari HP, mereka hanya mengirimkan token dan password baru
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type SetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=8"`
}

// ==========================================
// 4. JWT CLAIMS
// ==========================================

type TokenClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}