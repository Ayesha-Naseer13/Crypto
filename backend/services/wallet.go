package services

import (
	"context"
	"errors"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DefaultWalletAmount = 100.0

type WalletService struct {
	db     *mongo.Database
	crypto *CryptoService
}

func NewWalletService(db *mongo.Database) *WalletService {
	return &WalletService{
		db:     db,
		crypto: NewCryptoService(),
	}
}

// CreateWallet creates a new wallet for a user
func (s *WalletService) CreateWallet(ctx context.Context, userID primitive.ObjectID) (*models.Wallet, string, error) {
	// Generate key pair
	privateKey, publicKey, err := s.crypto.GenerateKeyPair()
	if err != nil {
		return nil, "", err
	}

	pubKeyStr := s.crypto.PublicKeyToString(publicKey)
	privKeyStr := s.crypto.PrivateKeyToString(privateKey)
	walletID := s.crypto.GenerateWalletID(pubKeyStr)

	// Encrypt private key
	encryptedPrivKey, err := s.crypto.EncryptPrivateKey(privKeyStr)
	if err != nil {
		return nil, "", err
	}

	wallet := &models.Wallet{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		WalletID:      walletID,
		PublicKey:     pubKeyStr,
		CachedBalance: DefaultWalletAmount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	collection := s.db.Collection("wallets")
	_, err = collection.InsertOne(ctx, wallet)
	if err != nil {
		return nil, "", err
	}

	if DefaultWalletAmount > 0 {
		initialUTXO := &models.UTXO{
			ID:          primitive.NewObjectID(),
			TxID:        "genesis_" + walletID,
			OutputIndex: 0,
			WalletID:    walletID,
			Amount:      DefaultWalletAmount,
			IsSpent:     false,
			BlockHash:   "genesis",
			CreatedAt:   time.Now(),
		}
		utxoCollection := s.db.Collection("utxos")
		_, err = utxoCollection.InsertOne(ctx, initialUTXO)
		if err != nil {
			return nil, "", err
		}
	}

	return wallet, encryptedPrivKey, nil
}

// GetWalletByUserID retrieves wallet by user ID
func (s *WalletService) GetWalletByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Wallet, error) {
	collection := s.db.Collection("wallets")

	var wallet models.Wallet
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetWalletByWalletID retrieves wallet by wallet ID
func (s *WalletService) GetWalletByWalletID(ctx context.Context, walletID string) (*models.Wallet, error) {
	collection := s.db.Collection("wallets")

	var wallet models.Wallet
	err := collection.FindOne(ctx, bson.M{"wallet_id": walletID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid wallet ID")
		}
		return nil, err
	}

	return &wallet, nil
}

// ValidateWalletExists checks if wallet ID exists
func (s *WalletService) ValidateWalletExists(ctx context.Context, walletID string) bool {
	_, err := s.GetWalletByWalletID(ctx, walletID)
	return err == nil
}

// GetUTXOsForWallet retrieves all unspent UTXOs for a wallet
func (s *WalletService) GetUTXOsForWallet(ctx context.Context, walletID string) ([]models.UTXO, error) {
	collection := s.db.Collection("utxos")

	cursor, err := collection.Find(ctx, bson.M{
		"wallet_id": walletID,
		"is_spent":  false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var utxos []models.UTXO
	if err := cursor.All(ctx, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}

// CalculateBalance calculates balance from UTXOs
func (s *WalletService) CalculateBalance(ctx context.Context, walletID string) (float64, error) {
	utxos, err := s.GetUTXOsForWallet(ctx, walletID)
	if err != nil {
		return 0, err
	}

	var balance float64
	for _, utxo := range utxos {
		balance += utxo.Amount
	}

	return balance, nil
}

// UpdateCachedBalance updates the cached balance
func (s *WalletService) UpdateCachedBalance(ctx context.Context, walletID string) error {
	balance, err := s.CalculateBalance(ctx, walletID)
	if err != nil {
		return err
	}

	collection := s.db.Collection("wallets")
	_, err = collection.UpdateOne(ctx,
		bson.M{"wallet_id": walletID},
		bson.M{"$set": bson.M{
			"cached_balance": balance,
			"updated_at":     time.Now(),
		}},
	)

	return err
}

// CreateUTXO creates a new UTXO
func (s *WalletService) CreateUTXO(ctx context.Context, utxo *models.UTXO) error {
	collection := s.db.Collection("utxos")
	_, err := collection.InsertOne(ctx, utxo)
	return err
}

// MarkUTXOAsSpent marks a UTXO as spent
func (s *WalletService) MarkUTXOAsSpent(ctx context.Context, txID string, outputIndex int, spentInTx string) error {
	collection := s.db.Collection("utxos")
	_, err := collection.UpdateOne(ctx,
		bson.M{"tx_id": txID, "output_index": outputIndex},
		bson.M{"$set": bson.M{
			"is_spent":    true,
			"spent_in_tx": spentInTx,
		}},
	)
	return err
}

// GetCryptoService returns the crypto service
func (s *WalletService) GetCryptoService() *CryptoService {
	return s.crypto
}

// GetAllWallets returns all wallets (for zakat processing)
func (s *WalletService) GetAllWallets(ctx context.Context) ([]models.Wallet, error) {
	collection := s.db.Collection("wallets")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var wallets []models.Wallet
	if err := cursor.All(ctx, &wallets); err != nil {
		return nil, err
	}

	return wallets, nil
}
