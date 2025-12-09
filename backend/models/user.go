package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email             string             `bson:"email" json:"email"`
	FullName          string             `bson:"full_name" json:"fullName"`
	CNIC              string             `bson:"cnic" json:"cnic"`
	WalletID          string             `bson:"wallet_id" json:"walletId"`
	PublicKey         string             `bson:"public_key" json:"publicKey"`
	EncryptedPrivKey  string             `bson:"encrypted_priv_key" json:"-"`
	Beneficiaries     []Beneficiary      `bson:"beneficiaries" json:"beneficiaries"`
	ZakatTracking     []ZakatRecord      `bson:"zakat_tracking" json:"zakatTracking"`
	OTP               string             `bson:"otp" json:"-"`
	OTPExpiry         time.Time          `bson:"otp_expiry" json:"-"`
	IsVerified        bool               `bson:"is_verified" json:"isVerified"`
	CreatedAt         time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updatedAt"`
}

type Beneficiary struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WalletID string             `bson:"wallet_id" json:"walletId"`
	Name     string             `bson:"name" json:"name"`
	AddedAt  time.Time          `bson:"added_at" json:"addedAt"`
}

type ZakatRecord struct {
	Amount    float64   `bson:"amount" json:"amount"`
	Date      time.Time `bson:"date" json:"date"`
	BlockHash string    `bson:"block_hash" json:"blockHash"`
	TxID      string    `bson:"tx_id" json:"txId"`
}
