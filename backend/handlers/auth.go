package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authService *services.AuthService
	logService  *services.LogService
}

func NewAuthHandler(authService *services.AuthService, logService *services.LogService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logService:  logService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required"`
	CNIC     string `json:"cnic" binding:"required"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.authService.Register(ctx, req.Email, req.FullName, req.CNIC)
	if err != nil {
		h.logService.LogSystemEvent(ctx, "register_failed", "", "", err.Error(), c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logService.LogSystemEvent(ctx, "user_registered", user.ID.Hex(), user.WalletID, "User registered successfully", c.ClientIP(), "success")

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Registration successful. OTP sent to email.",
		"walletId": user.WalletID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.authService.Login(ctx, req.Email)
	if err != nil {
		h.logService.LogSystemEvent(ctx, "login_failed", "", "", err.Error(), c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logService.LogSystemEvent(ctx, "login_initiated", user.ID.Hex(), user.WalletID, "OTP sent", c.ClientIP(), "success")

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email",
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, user, err := h.authService.VerifyOTP(ctx, req.Email, req.OTP)
	if err != nil {
		h.logService.LogSystemEvent(ctx, "otp_verification_failed", "", "", err.Error(), c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logService.LogSystemEvent(ctx, "login_success", user.ID.Hex(), user.WalletID, "OTP verified", c.ClientIP(), "success")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID.Hex(),
			"email":     user.Email,
			"fullName":  user.FullName,
			"walletId":  user.WalletID,
			"publicKey": user.PublicKey,
		},
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID.Hex(),
		"email":         user.Email,
		"fullName":      user.FullName,
		"cnic":          user.CNIC,
		"walletId":      user.WalletID,
		"publicKey":     user.PublicKey,
		"beneficiaries": user.Beneficiaries,
		"zakatTracking": user.ZakatTracking,
		"isVerified":    user.IsVerified,
		"createdAt":     user.CreatedAt,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updates := make(map[string]interface{})
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}

	if err := h.authService.UpdateUser(ctx, userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	h.logService.LogSystemEvent(ctx, "profile_updated", userID.Hex(), "", "Profile updated", c.ClientIP(), "success")

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
