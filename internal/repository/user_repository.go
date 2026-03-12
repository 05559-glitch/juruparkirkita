package repository

import (
	"arena-ban/internal/domain"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		db:    db,
		redis: redis,
	}
}

func (r *UserRepository) GetByIdentifier(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) InsertVerifyToken(req *domain.RegisterVerification) (*domain.RegisterVerification, error) {
	err := r.db.Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Mengganti GetValidTokenByOTP menjadi parameter Token karena kita memakai string panjang
func (r *UserRepository) GetValidTokenByString(tokenString string) (*domain.RegisterVerification, error) {
	var token domain.RegisterVerification
	err := r.db.Where("token = ? AND is_used = ? AND expires_at > ?",
		tokenString, false, time.Now()).First(&token).Error
	return &token, err
}

// Transaksi untuk membuat user dan menandai token verifikasi sebagai telah digunakan
func (r *UserRepository) CreateUserWithTransaction(user *domain.User, tokenID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan user baru
		if err := tx.Create(user).Error; err != nil {
			return err // Otomatis rollback
		}

		// 2. Update status token
		if err := tx.Model(&domain.RegisterVerification{}).
			Where("id = ?", tokenID).
			Update("is_used", true).Error; err != nil {
			return err // Otomatis rollback
		}

		return nil // Commit transaksi
	})
}

func (r *UserRepository) InsertPasswordReset(req *domain.PasswordReset) error {
	return r.db.Create(req).Error
}

func (r *UserRepository) GetValidResetToken(hashedToken string) (*domain.PasswordReset, error) {
	var resetData domain.PasswordReset
	err := r.db.Where("token_hash = ? AND is_used = ? AND expires_at > ?", 
		hashedToken, false, time.Now()).First(&resetData).Error
	
	if err != nil {
		return nil, err
	}
	return &resetData, nil
}

// Transaksi untuk memperbarui password dan menandai token reset sebagai telah digunakan
func (r *UserRepository) UpdatePasswordWithTransaction(userID uint, newPassword string, tokenID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update password pada tabel user
		if err := tx.Model(&domain.User{}).Where("id = ?", userID).Update("password", newPassword).Error; err != nil {
			return err
		}

		// 2. Tandai token reset menjadi tidak aktif
		if err := tx.Model(&domain.PasswordReset{}).Where("id = ?", tokenID).Update("is_used", true).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) SetCache(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return r.redis.Set(ctx, key, value, exp).Err()
}

func (r *UserRepository) GetCache(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

func (r *UserRepository) DeleteCache(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}