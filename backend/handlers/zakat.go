package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type ZakatHandler struct {
	zakatService *services.ZakatService
}

func NewZakatHandler(zakatService *services.ZakatService) *ZakatHandler {
	return &ZakatHandler{
		zakatService: zakatService,
	}
}

func (h *ZakatHandler) GetHistory(c *gin.Context) {
	walletID := c.MustGet("walletID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	history, err := h.zakatService.GetZakatHistory(ctx, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch zakat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"zakatHistory": history,
	})
}

func (h *ZakatHandler) ProcessZakat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := h.zakatService.ProcessMonthlyZakat(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process zakat: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Zakat processing completed",
	})
}
