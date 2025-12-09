package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Action    string             `bson:"action" json:"action"`
	UserID    string             `bson:"user_id,omitempty" json:"userId,omitempty"`
	WalletID  string             `bson:"wallet_id,omitempty" json:"walletId,omitempty"`
	Details   string             `bson:"details" json:"details"`
	IPAddress string             `bson:"ip_address" json:"ipAddress"`
	Status    string             `bson:"status" json:"status"` // success, failed
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type TransactionLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TxID      string             `bson:"tx_id" json:"txId"`
	Action    string             `bson:"action" json:"action"` // sent, received, mined, zakat_deducted
	WalletID  string             `bson:"wallet_id" json:"walletId"`
	Amount    float64            `bson:"amount" json:"amount"`
	BlockHash string             `bson:"block_hash" json:"blockHash"`
	Status    string             `bson:"status" json:"status"`
	Note      string             `bson:"note" json:"note"`
	IPAddress string             `bson:"ip_address" json:"ipAddress"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
