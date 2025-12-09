package services

import (
	"context"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogService struct {
	db *mongo.Database
}

func NewLogService(db *mongo.Database) *LogService {
	return &LogService{db: db}
}

// LogSystemEvent logs a system event
func (s *LogService) LogSystemEvent(ctx context.Context, action, userID, walletID, details, ipAddress, status string) error {
	collection := s.db.Collection("system_logs")

	log := models.SystemLog{
		ID:        primitive.NewObjectID(),
		Action:    action,
		UserID:    userID,
		WalletID:  walletID,
		Details:   details,
		IPAddress: ipAddress,
		Status:    status,
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(ctx, log)
	return err
}

// LogTransaction logs a transaction event
func (s *LogService) LogTransaction(ctx context.Context, txID, action, walletID string, amount float64, blockHash, status, note, ipAddress string) error {
	collection := s.db.Collection("transaction_logs")

	log := models.TransactionLog{
		ID:        primitive.NewObjectID(),
		TxID:      txID,
		Action:    action,
		WalletID:  walletID,
		Amount:    amount,
		BlockHash: blockHash,
		Status:    status,
		Note:      note,
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	}

	_, err := collection.InsertOne(ctx, log)
	return err
}

// GetSystemLogs returns system logs
func (s *LogService) GetSystemLogs(ctx context.Context, limit int64) ([]models.SystemLog, error) {
	collection := s.db.Collection("system_logs")

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.SystemLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetTransactionLogs returns transaction logs for a wallet
func (s *LogService) GetTransactionLogs(ctx context.Context, walletID string, limit int64) ([]models.TransactionLog, error) {
	collection := s.db.Collection("transaction_logs")

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{"wallet_id": walletID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.TransactionLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
