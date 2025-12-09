package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"backend/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	db            *mongo.Database
	walletService *WalletService
}

func NewAuthService(db *mongo.Database, walletService *WalletService) *AuthService {
	return &AuthService{
		db:            db,
		walletService: walletService,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, email, fullName, cnic string) (*models.User, error) {
	collection := s.db.Collection("users")

	// Check if email already exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Generate OTP
	otp := s.generateOTP()

	// Create user
	user := &models.User{
		ID:            primitive.NewObjectID(),
		Email:         email,
		FullName:      fullName,
		CNIC:          cnic,
		Beneficiaries: []models.Beneficiary{},
		ZakatTracking: []models.ZakatRecord{},
		OTP:           otp,
		OTPExpiry:     time.Now().Add(10 * time.Minute),
		IsVerified:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create wallet
	wallet, encryptedPrivKey, err := s.walletService.CreateWallet(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	user.WalletID = wallet.WalletID
	user.PublicKey = wallet.PublicKey
	user.EncryptedPrivKey = encryptedPrivKey

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// In production, send OTP via email
	// For now, we'll log it (stub implementation)
	fmt.Printf("OTP for %s: %s\n", email, otp)

	return user, nil
}

// Login initiates login and sends OTP
func (s *AuthService) Login(ctx context.Context, email string) (*models.User, error) {
	collection := s.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Generate new OTP
	otp := s.generateOTP()
	user.OTP = otp
	user.OTPExpiry = time.Now().Add(10 * time.Minute)

	_, err = collection.UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"otp":        otp,
			"otp_expiry": user.OTPExpiry,
		}},
	)
	if err != nil {
		return nil, err
	}

	// In production, send OTP via email
	fmt.Printf("OTP for %s: %s\n", email, otp)

	return &user, nil
}

// VerifyOTP verifies OTP and returns JWT token
func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) (string, *models.User, error) {
	collection := s.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", nil, errors.New("user not found")
	}

	if user.OTP != otp {
		return "", nil, errors.New("invalid OTP")
	}

	if time.Now().After(user.OTPExpiry) {
		return "", nil, errors.New("OTP expired")
	}

	// Mark as verified and clear OTP
	_, err = collection.UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"is_verified": true,
			"otp":         "",
			"updated_at":  time.Now(),
		}},
	)
	if err != nil {
		return "", nil, err
	}

	// Generate JWT
	token, err := s.generateJWT(&user)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	collection := s.db.Collection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates user profile
func (s *AuthService) UpdateUser(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) error {
	collection := s.db.Collection("users")

	updates["updated_at"] = time.Now()

	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": updates},
	)

	return err
}

// AddBeneficiary adds a beneficiary to user
func (s *AuthService) AddBeneficiary(ctx context.Context, userID primitive.ObjectID, walletID, name string) error {
	collection := s.db.Collection("users")

	beneficiary := models.Beneficiary{
		ID:       primitive.NewObjectID(),
		WalletID: walletID,
		Name:     name,
		AddedAt:  time.Now(),
	}

	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{
			"$push": bson.M{"beneficiaries": beneficiary},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)

	return err
}

// RemoveBeneficiary removes a beneficiary from user
func (s *AuthService) RemoveBeneficiary(ctx context.Context, userID, beneficiaryID primitive.ObjectID) error {
	collection := s.db.Collection("users")

	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{
			"$pull": bson.M{"beneficiaries": bson.M{"_id": beneficiaryID}},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)

	return err
}

// generateOTP generates a 6-digit OTP
func (s *AuthService) generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// generateJWT generates a JWT token
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID.Hex(),
		"email":     user.Email,
		"wallet_id": user.WalletID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
