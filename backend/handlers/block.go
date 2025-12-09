package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type BlockHandler struct {
	blockchainService *services.BlockchainService
}

func NewBlockHandler(blockchainService *services.BlockchainService) *BlockHandler {
	return &BlockHandler{
		blockchainService: blockchainService,
	}
}

func (h *BlockHandler) GetAllBlocks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blocks, err := h.blockchainService.GetAllBlocks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
		"count":  len(blocks),
	})
}

func (h *BlockHandler) GetLatestBlock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	block, err := h.blockchainService.GetLatestBlock(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No blocks found"})
		return
	}

	c.JSON(http.StatusOK, block)
}

func (h *BlockHandler) GetBlockByHash(c *gin.Context) {
	hash := c.Param("hash")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	block, err := h.blockchainService.GetBlockByHash(ctx, hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(http.StatusOK, block)
}
