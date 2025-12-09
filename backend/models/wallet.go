package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"userId"`
	WalletID      string             `bson:"wallet_id" json:"walletId"`
	PublicKey     string             `bson:"public_key" json:"publicKey"`
	CachedBalance float64            `bson:"cached_balance" json:"cachedBalance"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

type UTXO struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TxID        string             `bson:"tx_id" json:"txId"`
	OutputIndex int                `bson:"output_index" json:"outputIndex"`
	WalletID    string             `bson:"wallet_id" json:"walletId"`
	Amount      float64            `bson:"amount" json:"amount"`
	IsSpent     bool               `bson:"is_spent" json:"isSpent"`
	SpentInTx   string             `bson:"spent_in_tx,omitempty" json:"spentInTx,omitempty"`
	BlockHash   string             `bson:"block_hash" json:"blockHash"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}
