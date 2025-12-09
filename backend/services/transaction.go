package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionService struct {
	db         *mongo.Database
	blockchain *BlockchainService
	crypto     *CryptoService
}

func NewTransactionService(db *mongo.Database, blockchain *BlockchainService) *TransactionService {
	return &TransactionService{
		db:         db,
		blockchain: blockchain,
		crypto:     NewCryptoService(),
	}
}

// CreateTransaction creates and validates a new transaction
func (s *TransactionService) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	// Validate transaction
	if err := s.ValidateTransaction(ctx, tx); err != nil {
		return err
	}

	collection := s.db.Collection("transactions")
	_, err := collection.InsertOne(ctx, tx)
	return err
}

// ValidateTransaction validates transaction signature and UTXOs
func (s *TransactionService) ValidateTransaction(ctx context.Context, tx *models.Transaction) error {
	// Verify sender wallet exists
	walletCollection := s.db.Collection("wallets")
	var wallet models.Wallet
	err := walletCollection.FindOne(ctx, bson.M{"wallet_id": tx.SenderWalletID}).Decode(&wallet)
	if err != nil {
		return errors.New("invalid sender wallet ID")
	}

	// Verify receiver wallet exists (skip for system transactions)
	if tx.Type != "mining_reward" {
		err = walletCollection.FindOne(ctx, bson.M{"wallet_id": tx.ReceiverWalletID}).Decode(&wallet)
		if err != nil {
			return errors.New("invalid receiver wallet ID")
		}
	}

	// Verify digital signature
	if tx.Type != "mining_reward" && tx.Type != "zakat_deduction" {
		pubKey, err := s.crypto.StringToPublicKey(tx.SenderPublicKey)
		if err != nil {
			return errors.New("invalid public key")
		}

		signedPayload := fmt.Sprintf("%s%s%.8f%s%s",
			tx.SenderWalletID,
			tx.ReceiverWalletID,
			tx.Amount,
			tx.Timestamp.Format(time.RFC3339),
			tx.Note,
		)

		if !s.crypto.VerifySignature(pubKey, signedPayload, tx.Signature) {
			return errors.New("invalid digital signature")
		}
	}

	// Verify UTXOs are not double-spent
	for _, input := range tx.InputUTXOs {
		utxoCollection := s.db.Collection("utxos")
		var utxo models.UTXO
		err := utxoCollection.FindOne(ctx, bson.M{
			"tx_id":        input.TxID,
			"output_index": input.OutputIndex,
		}).Decode(&utxo)

		if err != nil {
			return errors.New("UTXO not found")
		}

		if utxo.IsSpent {
			return errors.New("double-spend detected: UTXO already spent")
		}
	}

	// Verify input amount >= output amount
	var inputSum float64
	for _, input := range tx.InputUTXOs {
		inputSum += input.Amount
	}

	var outputSum float64
	for _, output := range tx.OutputUTXOs {
		outputSum += output.Amount
	}

	if inputSum < outputSum {
		return errors.New("insufficient balance")
	}

	return nil
}

// GetPendingTransactions returns all pending transactions
func (s *TransactionService) GetPendingTransactions(ctx context.Context) ([]models.Transaction, error) {
	collection := s.db.Collection("transactions")

	cursor, err := collection.Find(ctx, bson.M{"status": "pending"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ConfirmTransactions marks transactions as confirmed
func (s *TransactionService) ConfirmTransactions(ctx context.Context, txIDs []string, blockHash string) error {
	collection := s.db.Collection("transactions")

	_, err := collection.UpdateMany(ctx,
		bson.M{"tx_id": bson.M{"$in": txIDs}},
		bson.M{"$set": bson.M{
			"status":     "confirmed",
			"block_hash": blockHash,
		}},
	)

	return err
}

// GetTransactionHistory returns transaction history for a wallet
func (s *TransactionService) GetTransactionHistory(ctx context.Context, walletID string, limit int64) ([]models.Transaction, error) {
	collection := s.db.Collection("transactions")

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"sender_wallet_id": walletID},
			{"receiver_wallet_id": walletID},
		},
		"status": "confirmed",
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// GenerateTxID generates a unique transaction ID
func (s *TransactionService) GenerateTxID(senderID, receiverID string, amount float64, timestamp time.Time) string {
	data := fmt.Sprintf("%s%s%.8f%s%d", senderID, receiverID, amount, timestamp.Format(time.RFC3339), time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CreateSystemTransaction creates a system transaction (mining reward, zakat)
func (s *TransactionService) CreateSystemTransaction(ctx context.Context, txType string, receiverWalletID string, amount float64, note string) (*models.Transaction, error) {
	timestamp := time.Now()
	txID := s.GenerateTxID("system", receiverWalletID, amount, timestamp)

	tx := &models.Transaction{
		ID:               primitive.NewObjectID(),
		TxID:             txID,
		SenderWalletID:   "system",
		ReceiverWalletID: receiverWalletID,
		Amount:           amount,
		Note:             note,
		Timestamp:        timestamp,
		SenderPublicKey:  "system",
		Signature:        "system",
		InputUTXOs:       []models.UTXOInput{},
		OutputUTXOs: []models.UTXOOutput{
			{WalletID: receiverWalletID, Amount: amount, Index: 0},
		},
		Type:   txType,
		Status: "pending",
		Fee:    0,
	}

	collection := s.db.Collection("transactions")
	_, err := collection.InsertOne(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
