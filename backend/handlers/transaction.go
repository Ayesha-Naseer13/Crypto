package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MinimumTransferAmount = 0.01 // Minimum amount that can be transferred

type TransactionHandler struct {
	transactionService *services.TransactionService
	walletService      *services.WalletService
	logService         *services.LogService
}

func NewTransactionHandler(transactionService *services.TransactionService, walletService *services.WalletService, logService *services.LogService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		walletService:      walletService,
		logService:         logService,
	}
}

type SendMoneyRequest struct {
	ReceiverWalletID string  `json:"receiverWalletId" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
	Note             string  `json:"note"`
	Signature        string  `json:"signature" binding:"required"`
}

func (h *TransactionHandler) SendMoney(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)
	walletID := c.MustGet("walletID").(string)

	var req SendMoneyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount < MinimumTransferAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Minimum transfer amount is %.2f", MinimumTransferAmount),
		})
		return
	}

	if walletID == req.ReceiverWalletID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send to your own wallet"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Validate receiver wallet exists
	if !h.walletService.ValidateWalletExists(ctx, req.ReceiverWalletID) {
		h.logService.LogSystemEvent(ctx, "invalid_wallet_id_attempt", userID.Hex(), walletID, "Attempted to send to invalid wallet", c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver Wallet ID"})
		return
	}

	// Get sender's wallet
	senderWallet, err := h.walletService.GetWalletByWalletID(ctx, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallet"})
		return
	}

	// Get available UTXOs
	utxos, err := h.walletService.GetUTXOsForWallet(ctx, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get UTXOs"})
		return
	}

	var totalAvailable float64
	for _, utxo := range utxos {
		totalAvailable += utxo.Amount
	}

	if totalAvailable < req.Amount {
		h.logService.LogSystemEvent(ctx, "insufficient_balance", userID.Hex(), walletID,
			fmt.Sprintf("Insufficient balance: has %.4f, needs %.4f", totalAvailable, req.Amount),
			c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Insufficient balance. Available: %.4f, Required: %.4f", totalAvailable, req.Amount),
		})
		return
	}

	// Select UTXOs for transaction
	var inputUTXOs []models.UTXOInput
	var totalInput float64
	for _, utxo := range utxos {
		inputUTXOs = append(inputUTXOs, models.UTXOInput{
			TxID:        utxo.TxID,
			OutputIndex: utxo.OutputIndex,
			Amount:      utxo.Amount,
		})
		totalInput += utxo.Amount
		if totalInput >= req.Amount {
			break
		}
	}

	if totalInput < req.Amount {
		h.logService.LogSystemEvent(ctx, "insufficient_balance", userID.Hex(), walletID, "Insufficient balance for transaction", c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Calculate change
	change := totalInput - req.Amount

	// Create output UTXOs
	outputUTXOs := []models.UTXOOutput{
		{WalletID: req.ReceiverWalletID, Amount: req.Amount, Index: 0},
	}
	if change > 0 {
		outputUTXOs = append(outputUTXOs, models.UTXOOutput{
			WalletID: walletID,
			Amount:   change,
			Index:    1,
		})
	}

	// Create transaction
	timestamp := time.Now()
	txID := h.transactionService.GenerateTxID(walletID, req.ReceiverWalletID, req.Amount, timestamp)

	tx := &models.Transaction{
		ID:               primitive.NewObjectID(),
		TxID:             txID,
		SenderWalletID:   walletID,
		ReceiverWalletID: req.ReceiverWalletID,
		Amount:           req.Amount,
		Note:             req.Note,
		Timestamp:        timestamp,
		SenderPublicKey:  senderWallet.PublicKey,
		Signature:        req.Signature,
		InputUTXOs:       inputUTXOs,
		OutputUTXOs:      outputUTXOs,
		Type:             "transfer",
		Status:           "pending",
		Fee:              0,
	}

	if err := h.transactionService.CreateTransaction(ctx, tx); err != nil {
		h.logService.LogSystemEvent(ctx, "transaction_rejected", userID.Hex(), walletID, err.Error(), c.ClientIP(), "failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logService.LogTransaction(ctx, txID, "sent", walletID, req.Amount, "", "pending", req.Note, c.ClientIP())
	h.logService.LogTransaction(ctx, txID, "received", req.ReceiverWalletID, req.Amount, "", "pending", req.Note, c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Transaction created successfully",
		"transactionId": txID,
		"status":        "pending",
	})
}

func (h *TransactionHandler) GetHistory(c *gin.Context) {
	walletID := c.MustGet("walletID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	transactions, err := h.transactionService.GetTransactionHistory(ctx, walletID, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction history"})
		return
	}

	// Add direction info
	var history []gin.H
	for _, tx := range transactions {
		direction := "received"
		if tx.SenderWalletID == walletID {
			direction = "sent"
		}
		history = append(history, gin.H{
			"id":               tx.ID,
			"txId":             tx.TxID,
			"senderWalletId":   tx.SenderWalletID,
			"receiverWalletId": tx.ReceiverWalletID,
			"amount":           tx.Amount,
			"note":             tx.Note,
			"timestamp":        tx.Timestamp,
			"type":             tx.Type,
			"status":           tx.Status,
			"blockHash":        tx.BlockHash,
			"direction":        direction,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": history,
	})
}

func (h *TransactionHandler) GetPending(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	transactions, err := h.transactionService.GetPendingTransactions(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pendingTransactions": transactions,
		"count":               len(transactions),
	})
}

// Helper to create signed payload
func CreateSignedPayload(senderID, receiverID string, amount float64, timestamp time.Time, note string) string {
	return fmt.Sprintf("%s%s%.8f%s%s", senderID, receiverID, amount, timestamp.Format(time.RFC3339), note)
}
