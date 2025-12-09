package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MiningHandler struct {
	miningService *services.MiningService
	logService    *services.LogService
}

func NewMiningHandler(miningService *services.MiningService, logService *services.LogService) *MiningHandler {
	return &MiningHandler{
		miningService: miningService,
		logService:    logService,
	}
}

func (h *MiningHandler) Mine(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)
	walletID := c.MustGet("walletID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	h.logService.LogSystemEvent(ctx, "mining_started", userID.Hex(), walletID, "Mining initiated", c.ClientIP(), "success")

	block, err := h.miningService.MineBlock(ctx, walletID)
	if err != nil {
		h.logService.LogSystemEvent(ctx, "mining_failed", userID.Hex(), walletID, err.Error(), c.ClientIP(), "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mining failed: " + err.Error()})
		return
	}

	if block == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "No pending transactions to mine or mining already in progress",
		})
		return
	}

	h.logService.LogSystemEvent(ctx, "mining_success", userID.Hex(), walletID, "Block mined: "+block.Hash, c.ClientIP(), "success")

	c.JSON(http.StatusOK, gin.H{
		"message": "Block mined successfully",
		"block": gin.H{
			"index":        block.Index,
			"hash":         block.Hash,
			"previousHash": block.PreviousHash,
			"nonce":        block.Nonce,
			"merkleRoot":   block.MerkleRoot,
			"transactions": len(block.Transactions),
			"timestamp":    block.Timestamp,
		},
	})
}

func (h *MiningHandler) GetStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := h.miningService.GetMiningStatus(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mining status"})
		return
	}

	c.JSON(http.StatusOK, status)
}
