package services

import (
	"context"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ZakatRate         = 0.025 // 2.5%
	ZakatPoolWalletID = "zakat_pool_wallet_00000000000000000000"
)

type ZakatService struct {
	db          *mongo.Database
	transaction *TransactionService
	blockchain  *BlockchainService
}

func NewZakatService(db *mongo.Database, transaction *TransactionService, blockchain *BlockchainService) *ZakatService {
	return &ZakatService{
		db:          db,
		transaction: transaction,
		blockchain:  blockchain,
	}
}

// ProcessMonthlyZakat processes zakat deduction for all wallets
func (s *ZakatService) ProcessMonthlyZakat(ctx context.Context) error {
	walletCollection := s.db.Collection("wallets")
	userCollection := s.db.Collection("users")

	// Get all wallets with balance > 0
	cursor, err := walletCollection.Find(ctx, bson.M{"cached_balance": bson.M{"$gt": 0}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var wallets []models.Wallet
	if err := cursor.All(ctx, &wallets); err != nil {
		return err
	}

	for _, wallet := range wallets {
		// Calculate zakat amount (2.5%)
		zakatAmount := wallet.CachedBalance * ZakatRate

		if zakatAmount <= 0 {
			continue
		}

		// Get wallet UTXOs
		utxoCollection := s.db.Collection("utxos")
		utxoCursor, err := utxoCollection.Find(ctx, bson.M{
			"wallet_id": wallet.WalletID,
			"is_spent":  false,
		})
		if err != nil {
			continue
		}

		var utxos []models.UTXO
		if err := utxoCursor.All(ctx, &utxos); err != nil {
			utxoCursor.Close(ctx)
			continue
		}
		utxoCursor.Close(ctx)

		// Create input UTXOs
		var inputUTXOs []models.UTXOInput
		var totalInput float64
		for _, utxo := range utxos {
			inputUTXOs = append(inputUTXOs, models.UTXOInput{
				TxID:        utxo.TxID,
				OutputIndex: utxo.OutputIndex,
				Amount:      utxo.Amount,
			})
			totalInput += utxo.Amount
			if totalInput >= zakatAmount {
				break
			}
		}

		// Calculate change
		change := totalInput - zakatAmount

		// Create output UTXOs
		outputUTXOs := []models.UTXOOutput{
			{WalletID: ZakatPoolWalletID, Amount: zakatAmount, Index: 0},
		}
		if change > 0 {
			outputUTXOs = append(outputUTXOs, models.UTXOOutput{
				WalletID: wallet.WalletID,
				Amount:   change,
				Index:    1,
			})
		}

		// Create zakat transaction
		timestamp := time.Now()
		txID := s.transaction.GenerateTxID(wallet.WalletID, ZakatPoolWalletID, zakatAmount, timestamp)

		tx := &models.Transaction{
			ID:               primitive.NewObjectID(),
			TxID:             txID,
			SenderWalletID:   wallet.WalletID,
			ReceiverWalletID: ZakatPoolWalletID,
			Amount:           zakatAmount,
			Note:             "Monthly Zakat Deduction (2.5%)",
			Timestamp:        timestamp,
			SenderPublicKey:  wallet.PublicKey,
			Signature:        "system_zakat",
			InputUTXOs:       inputUTXOs,
			OutputUTXOs:      outputUTXOs,
			Type:             "zakat_deduction",
			Status:           "pending",
			Fee:              0,
		}

		txCollection := s.db.Collection("transactions")
		if _, err := txCollection.InsertOne(ctx, tx); err != nil {
			continue
		}

		// Update user's zakat tracking
		zakatRecord := models.ZakatRecord{
			Amount: zakatAmount,
			Date:   timestamp,
			TxID:   txID,
		}

		_, err = userCollection.UpdateOne(ctx,
			bson.M{"wallet_id": wallet.WalletID},
			bson.M{
				"$push": bson.M{"zakat_tracking": zakatRecord},
				"$set":  bson.M{"updated_at": time.Now()},
			},
		)
		if err != nil {
			continue
		}

		// Log the zakat deduction
		s.logZakatDeduction(ctx, wallet.WalletID, zakatAmount, txID)
	}

	return nil
}

// GetZakatHistory returns zakat history for a user
func (s *ZakatService) GetZakatHistory(ctx context.Context, walletID string) ([]models.ZakatRecord, error) {
	userCollection := s.db.Collection("users")

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"wallet_id": walletID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user.ZakatTracking, nil
}

// logZakatDeduction logs the zakat deduction event
func (s *ZakatService) logZakatDeduction(ctx context.Context, walletID string, amount float64, txID string) {
	logCollection := s.db.Collection("system_logs")

	log := models.SystemLog{
		ID:        primitive.NewObjectID(),
		Action:    "zakat_deduction",
		WalletID:  walletID,
		Details:   "Monthly zakat deduction processed",
		Status:    "success",
		Timestamp: time.Now(),
	}

	logCollection.InsertOne(ctx, log)

	txLogCollection := s.db.Collection("transaction_logs")
	txLog := models.TransactionLog{
		ID:        primitive.NewObjectID(),
		TxID:      txID,
		Action:    "zakat_deducted",
		WalletID:  walletID,
		Amount:    amount,
		Status:    "pending",
		Note:      "Monthly Zakat Deduction (2.5%)",
		Timestamp: time.Now(),
	}

	txLogCollection.InsertOne(ctx, txLog)
}
