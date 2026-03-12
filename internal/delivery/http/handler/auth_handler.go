package handler

import (
	"arena-ban/internal/domain"
	"arena-ban/internal/usecase"
	util "arena-ban/pkg"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
	logger  *logrus.Logger
}

func NewAuthHandler(u *usecase.AuthUsecase, l *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		usecase: u,
		logger:  l,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	req := new(domain.LoginRequest)

	if err := c.Bind().JSON(req); err != nil {
		h.logger.Warnf("Invalid JSON format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "format data tidak valid",
		})
	}

	if errs := util.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi gagal",
			"errors":  errs,
		})
	}

	response, err := h.usecase.Login(req)
	if err != nil {
		h.logger.Errorf("Login error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	h.logger.Infof("User %s successfully logged in", req.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "login berhasil",
		"data":    response,
	})
}

// Register adalah proses admin/sistem mendaftarkan email, lalu mengirimkan Deep Link
func (h *AuthHandler) Register(c fiber.Ctx) error {
	req := new(domain.CreateAccountRequest)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "format data tidak valid",
		})
	}

	if errs := util.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi gagal",
			"errors":  errs,
		})
	}

	msg, err := h.usecase.Register(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": msg,
	})
}

// RegisterPassword dipanggil oleh aplikasi mobile setelah user mengklik link di email
func (h *AuthHandler) RegisterPassword(c fiber.Ctx) error {
	req := new(domain.CreatePasswordRequest)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "format data tidak valid",
		})
	}

	if errs := util.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi gagal",
			"errors":  errs,
		})
	}

	err := h.usecase.RegisterPassword(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Akun berhasil dibuat. Silakan login menggunakan email dan password Anda.",
	})
}

func (h *AuthHandler) ForgotPassword(c fiber.Ctx) error {
	req := new(domain.ForgotPasswordRequest)
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "format data tidak valid",
		})
	}

	if errs := util.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi gagal",
			"errors":  errs,
		})
	}

	err := h.usecase.ForgotPassword(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Jika email terdaftar, instruksi reset password telah dikirimkan.",
	})
}

func (h *AuthHandler) ResetPassword(c fiber.Ctx) error {
	req := new(domain.VerifyResetPasswordRequest) // Kita gabungkan penerimaan token dari URL atau body
	
	// Untuk keamanan, sebaiknya token dikirim via body request dari mobile app, bukan URL params
	req.Token = c.Query("token") 
	
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "format data tidak valid",
		})
	}

	if errs := util.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi gagal",
			"errors":  errs,
		})
	}

	err := h.usecase.ResetPassword(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password berhasil diperbarui. Silakan login.",
	})
}

func (h *AuthHandler) VerifyToken(c fiber.Ctx) error {
	req := new(domain.VerifyTokenRequest)
	
	if err := c.Bind().JSON(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "format data tidak valid",
		})
	}

	if errs := util.ValidateStruct(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "validasi gagal",
			"errors":  errs,
		})
	}

	// Memanggil usecase untuk verifikasi
	err := h.usecase.VerifyToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "token valid",
	})
}