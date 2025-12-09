package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Block struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Index        int64              `bson:"index" json:"index"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
	Transactions []Transaction      `bson:"transactions" json:"transactions"`
	PreviousHash string             `bson:"previous_hash" json:"previousHash"`
	Nonce        int64              `bson:"nonce" json:"nonce"`
	Hash         string             `bson:"hash" json:"hash"`
	MerkleRoot   string             `bson:"merkle_root" json:"merkleRoot"`
	Difficulty   int                `bson:"difficulty" json:"difficulty"`
	Miner        string             `bson:"miner" json:"miner"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
}
