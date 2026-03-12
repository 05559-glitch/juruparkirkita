package usecase

import (
	"arena-ban/internal/domain"
	"arena-ban/internal/repository"
	util "arena-ban/pkg"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo   *repository.UserRepository
	logger *logrus.Logger
}

func NewAuthUsecase(userRepo *repository.UserRepository, logger *logrus.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo:   userRepo,
		logger: logger,
	}
}

func (a *AuthUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := a.repo.GetByIdentifier(req.Email)
	if err != nil {
		return nil, errors.New("akun tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("password salah")
	}

	token, err := util.JwtUtil.CreateToken(user)
	if err != nil {
		return nil, errors.New("gagal membuat token akses")
	}

	return &domain.LoginResponse{Token: token}, nil
}

func (a *AuthUsecase) Register(req *domain.CreateAccountRequest) (string, error) {
	isExists, _ := a.repo.IsEmailExists(req.Email)
	if isExists {
		return "", errors.New("email sudah terdaftar")
	}

	otpToken, err := util.GenerateOTP()
	if err != nil {
		return "", errors.New("gagal generate OTP")
	}

	verificationPayload := &domain.RegisterVerification{
		Name:      req.Name,
		Email:     req.Email,
		Role:      string(req.Role),
		Token:     otpToken,
		IsUsed:    false,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	_, err = a.repo.InsertVerifyToken(verificationPayload)
	if err != nil {
		return "", errors.New("gagal menyimpan data verifikasi")
	}

	go func(email, token string) {
		body := fmt.Sprintf("Kode OTP Anda: %s", token)
		_ = util.NewSMTP().SendMail("Verifikasi Akun", body, email)
	}(req.Email, otpToken)

	return "Kode verifikasi telah dikirim", nil
}

func (a *AuthUsecase) RegisterPassword(req *domain.CreatePasswordRequest) error {
	tokenData, err := a.repo.GetValidTokenByString(req.Token)
	if err != nil {
		return errors.New("token tidak valid atau kedaluwarsa")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &domain.User{
		Name:     tokenData.Name,
		Email:    tokenData.Email,
		Role:     tokenData.Role,
		Password: string(hashedPassword),
	}

	if err := a.repo.CreateUserWithTransaction(user, tokenData.ID); err != nil {
		return err
	}

	return nil
}

func (a *AuthUsecase) VerifyToken(token string) error {
	_, err := a.repo.GetValidTokenByString(token)
	return err
}

func (a *AuthUsecase) ForgotPassword(req *domain.ForgotPasswordRequest) error {
	user, err := a.repo.GetByIdentifier(req.Email)
	if err != nil {
		return nil // Untuk security, jangan beri tahu jika email tidak ada
	}

	rawToken, _ := util.GenerateSecureToken()
	hashedToken := util.HashSHA256(rawToken)

	resetPayload := &domain.PasswordReset{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IsUsed:    false,
	}

	if err := a.repo.InsertPasswordReset(resetPayload); err != nil {
		return err
	}

	go func(email, token string) {
		body := fmt.Sprintf("Reset password di: https://arena-ban.com/reset?token=%s", token)
		_ = util.NewSMTP().SendMail("Reset Password", body, email)
	}(user.Email, rawToken)

	return nil
}

func (a *AuthUsecase) VerifyResetPassword(token string) error {
	hashedToken := util.HashSHA256(token)
	_, err := a.repo.GetValidResetToken(hashedToken)
	return err
}

func (a *AuthUsecase) ResetPassword(req *domain.VerifyResetPasswordRequest) error {
	hashedToken := util.HashSHA256(req.Token)
	resetData, err := a.repo.GetValidResetToken(hashedToken)
	if err != nil {
		return errors.New("permintaan tidak valid")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	return a.repo.UpdatePasswordWithTransaction(
		resetData.UserID,
		string(hashedPassword),
		resetData.ID,
	)
}
