package services

import (
	"context"
	"sync"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MiningService struct {
	db          *mongo.Database
	blockchain  *BlockchainService
	transaction *TransactionService
	isMining    bool
	mutex       sync.Mutex
}

func NewMiningService(db *mongo.Database, blockchain *BlockchainService, transaction *TransactionService) *MiningService {
	return &MiningService{
		db:          db,
		blockchain:  blockchain,
		transaction: transaction,
		isMining:    false,
	}
}

// MineBlock mines pending transactions into a new block
func (s *MiningService) MineBlock(ctx context.Context, minerWalletID string) (*models.Block, error) {
	s.mutex.Lock()
	if s.isMining {
		s.mutex.Unlock()
		return nil, nil
	}
	s.isMining = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		s.isMining = false
		s.mutex.Unlock()
	}()

	// Get pending transactions
	pendingTxs, err := s.transaction.GetPendingTransactions(ctx)
	if err != nil {
		return nil, err
	}

	if len(pendingTxs) == 0 {
		return nil, nil
	}

	// Get latest block
	latestBlock, err := s.blockchain.GetLatestBlock(ctx)
	if err != nil {
		return nil, err
	}

	// Create new block
	newBlock := &models.Block{
		ID:           primitive.NewObjectID(),
		Index:        latestBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: pendingTxs,
		PreviousHash: latestBlock.Hash,
		Nonce:        0,
		Difficulty:   s.blockchain.GetDifficulty(),
		Miner:        minerWalletID,
		CreatedAt:    time.Now(),
	}

	// Calculate merkle root
	newBlock.MerkleRoot = s.blockchain.CalculateMerkleRoot(pendingTxs)

	// Perform Proof of Work
	s.blockchain.ProofOfWork(newBlock)

	// Add block to chain
	if err := s.blockchain.AddBlock(ctx, newBlock); err != nil {
		return nil, err
	}

	// Confirm transactions
	var txIDs []string
	for _, tx := range pendingTxs {
		txIDs = append(txIDs, tx.TxID)
	}
	if err := s.transaction.ConfirmTransactions(ctx, txIDs, newBlock.Hash); err != nil {
		return nil, err
	}

	// Mark UTXOs as spent and create new UTXOs
	if err := s.processTransactionUTXOs(ctx, pendingTxs, newBlock.Hash); err != nil {
		return nil, err
	}

	return newBlock, nil
}

// processTransactionUTXOs handles UTXO updates after mining
func (s *MiningService) processTransactionUTXOs(ctx context.Context, transactions []models.Transaction, blockHash string) error {
	utxoCollection := s.db.Collection("utxos")
	walletCollection := s.db.Collection("wallets")

	for _, tx := range transactions {
		// Mark input UTXOs as spent
		for _, input := range tx.InputUTXOs {
			_, err := utxoCollection.UpdateOne(ctx,
				primitive.M{"tx_id": input.TxID, "output_index": input.OutputIndex},
				primitive.M{"$set": primitive.M{
					"is_spent":    true,
					"spent_in_tx": tx.TxID,
				}},
			)
			if err != nil {
				return err
			}
		}

		// Create output UTXOs
		for _, output := range tx.OutputUTXOs {
			newUTXO := models.UTXO{
				ID:          primitive.NewObjectID(),
				TxID:        tx.TxID,
				OutputIndex: output.Index,
				WalletID:    output.WalletID,
				Amount:      output.Amount,
				IsSpent:     false,
				BlockHash:   blockHash,
				CreatedAt:   time.Now(),
			}
			_, err := utxoCollection.InsertOne(ctx, newUTXO)
			if err != nil {
				return err
			}

			// Update cached balance
			_, err = walletCollection.UpdateOne(ctx,
				primitive.M{"wallet_id": output.WalletID},
				primitive.M{"$inc": primitive.M{"cached_balance": output.Amount}},
			)
			if err != nil {
				return err
			}
		}

		// Update sender's cached balance (subtract spent amount)
		if tx.SenderWalletID != "system" {
			var totalSpent float64
			for _, input := range tx.InputUTXOs {
				totalSpent += input.Amount
			}
			_, err := walletCollection.UpdateOne(ctx,
				primitive.M{"wallet_id": tx.SenderWalletID},
				primitive.M{"$inc": primitive.M{"cached_balance": -totalSpent}},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IsMining returns current mining status
func (s *MiningService) IsMining() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.isMining
}

// GetMiningStatus returns mining status info
func (s *MiningService) GetMiningStatus(ctx context.Context) (map[string]interface{}, error) {
	pendingTxs, err := s.transaction.GetPendingTransactions(ctx)
	if err != nil {
		return nil, err
	}

	latestBlock, err := s.blockchain.GetLatestBlock(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"isMining":            s.IsMining(),
		"pendingTransactions": len(pendingTxs),
		"currentDifficulty":   s.blockchain.GetDifficulty(),
		"latestBlockIndex":    latestBlock.Index,
		"latestBlockHash":     latestBlock.Hash,
	}, nil
}
