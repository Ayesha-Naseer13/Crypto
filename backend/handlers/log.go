package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler(logService *services.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

func (h *LogHandler) GetSystemLogs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logs, err := h.logService.GetSystemLogs(ctx, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch system logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

func (h *LogHandler) GetTransactionLogs(c *gin.Context) {
	walletID := c.MustGet("walletID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logs, err := h.logService.GetTransactionLogs(ctx, walletID, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
