package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockchainService struct {
	db         *mongo.Database
	difficulty int
}

func NewBlockchainService(db *mongo.Database) *BlockchainService {
	return &BlockchainService{
		db:         db,
		difficulty: 5, // Hash must start with 5 zeros
	}
}

// InitializeGenesisBlock creates the genesis block if not exists
func (s *BlockchainService) InitializeGenesisBlock(ctx context.Context) error {
	collection := s.db.Collection("blocks")
	
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil // Genesis block already exists
	}

	genesis := models.Block{
		ID:           primitive.NewObjectID(),
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: []models.Transaction{},
		PreviousHash: "0",
		Nonce:        0,
		Hash:         "",
		MerkleRoot:   "0",
		Difficulty:   s.difficulty,
		Miner:        "system",
		CreatedAt:    time.Now(),
	}

	// Calculate genesis hash
	genesis.Hash = s.CalculateHash(&genesis)

	_, err = collection.InsertOne(ctx, genesis)
	return err
}

// CalculateHash computes SHA-256 hash of block
func (s *BlockchainService) CalculateHash(block *models.Block) string {
	data := fmt.Sprintf("%d%s%s%s%d%s",
		block.Index,
		block.Timestamp.Format(time.RFC3339),
		block.PreviousHash,
		block.MerkleRoot,
		block.Nonce,
		s.transactionsToString(block.Transactions),
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CalculateMerkleRoot computes merkle root of transactions
func (s *BlockchainService) CalculateMerkleRoot(transactions []models.Transaction) string {
	if len(transactions) == 0 {
		return "0"
	}
	
	var hashes []string
	for _, tx := range transactions {
		hashes = append(hashes, tx.TxID)
	}
	
	for len(hashes) > 1 {
		var newHashes []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				combined := hashes[i] + hashes[i+1]
				hash := sha256.Sum256([]byte(combined))
				newHashes = append(newHashes, hex.EncodeToString(hash[:]))
			} else {
				newHashes = append(newHashes, hashes[i])
			}
		}
		hashes = newHashes
	}
	
	return hashes[0]
}

// ProofOfWork performs mining with difficulty
func (s *BlockchainService) ProofOfWork(block *models.Block) {
	prefix := strings.Repeat("0", s.difficulty)
	for {
		block.Hash = s.CalculateHash(block)
		if strings.HasPrefix(block.Hash, prefix) {
			break
		}
		block.Nonce++
	}
}

// GetLatestBlock returns the most recent block
func (s *BlockchainService) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	collection := s.db.Collection("blocks")
	
	opts := options.FindOne().SetSort(bson.D{{Key: "index", Value: -1}})
	var block models.Block
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&block)
	if err != nil {
		return nil, err
	}
	
	return &block, nil
}

// AddBlock adds a new block to the chain
func (s *BlockchainService) AddBlock(ctx context.Context, block *models.Block) error {
	collection := s.db.Collection("blocks")
	_, err := collection.InsertOne(ctx, block)
	return err
}

// GetAllBlocks returns all blocks in the chain
func (s *BlockchainService) GetAllBlocks(ctx context.Context) ([]models.Block, error) {
	collection := s.db.Collection("blocks")
	
	opts := options.Find().SetSort(bson.D{{Key: "index", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var blocks []models.Block
	if err := cursor.All(ctx, &blocks); err != nil {
		return nil, err
	}
	
	return blocks, nil
}

// GetBlockByHash returns a block by its hash
func (s *BlockchainService) GetBlockByHash(ctx context.Context, hash string) (*models.Block, error) {
	collection := s.db.Collection("blocks")
	
	var block models.Block
	err := collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&block)
	if err != nil {
		return nil, err
	}
	
	return &block, nil
}

// ValidateChain validates the entire blockchain
func (s *BlockchainService) ValidateChain(ctx context.Context) (bool, error) {
	blocks, err := s.GetAllBlocks(ctx)
	if err != nil {
		return false, err
	}
	
	for i := len(blocks) - 1; i > 0; i-- {
		currentBlock := blocks[i-1]
		previousBlock := blocks[i]
		
		// Verify hash
		if currentBlock.Hash != s.CalculateHash(&currentBlock) {
			return false, errors.New("invalid block hash")
		}
		
		// Verify chain link
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false, errors.New("invalid chain link")
		}
	}
	
	return true, nil
}

// GetDifficulty returns current mining difficulty
func (s *BlockchainService) GetDifficulty() int {
	return s.difficulty
}

func (s *BlockchainService) transactionsToString(txs []models.Transaction) string {
	var result string
	for _, tx := range txs {
		result += tx.TxID
	}
	return result
}
