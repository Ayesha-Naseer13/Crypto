package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TxID            string             `bson:"tx_id" json:"txId"`
	SenderWalletID  string             `bson:"sender_wallet_id" json:"senderWalletId"`
	ReceiverWalletID string            `bson:"receiver_wallet_id" json:"receiverWalletId"`
	Amount          float64            `bson:"amount" json:"amount"`
	Note            string             `bson:"note" json:"note"`
	Timestamp       time.Time          `bson:"timestamp" json:"timestamp"`
	SenderPublicKey string             `bson:"sender_public_key" json:"senderPublicKey"`
	Signature       string             `bson:"signature" json:"signature"`
	InputUTXOs      []UTXOInput        `bson:"input_utxos" json:"inputUtxos"`
	OutputUTXOs     []UTXOOutput       `bson:"output_utxos" json:"outputUtxos"`
	Type            string             `bson:"type" json:"type"` // transfer, zakat_deduction, mining_reward
	Status          string             `bson:"status" json:"status"` // pending, confirmed, rejected
	BlockHash       string             `bson:"block_hash,omitempty" json:"blockHash,omitempty"`
	Fee             float64            `bson:"fee" json:"fee"`
}

type UTXOInput struct {
	TxID        string  `bson:"tx_id" json:"txId"`
	OutputIndex int     `bson:"output_index" json:"outputIndex"`
	Amount      float64 `bson:"amount" json:"amount"`
}

type UTXOOutput struct {
	WalletID string  `bson:"wallet_id" json:"walletId"`
	Amount   float64 `bson:"amount" json:"amount"`
	Index    int     `bson:"index" json:"index"`
}
