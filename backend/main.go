package main

import (
	"context"
	"log"
	"os"
	"time"

	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := config.ConnectDB(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer config.DisconnectDB(ctx)

	// Initialize services
	blockchainService := services.NewBlockchainService(db)
	walletService := services.NewWalletService(db)
	transactionService := services.NewTransactionService(db, blockchainService)
	miningService := services.NewMiningService(db, blockchainService, transactionService)
	authService := services.NewAuthService(db, walletService)
	zakatService := services.NewZakatService(db, transactionService, blockchainService)
	logService := services.NewLogService(db)

	// Initialize genesis block if blockchain is empty
	if err := blockchainService.InitializeGenesisBlock(ctx); err != nil {
		log.Println("Genesis block initialization:", err)
	}

	// Setup Zakat scheduler (runs monthly on the 1st at midnight)
	c := cron.New()
	c.AddFunc("0 0 1 * *", func() {
		log.Println("Running monthly Zakat deduction...")
		if err := zakatService.ProcessMonthlyZakat(context.Background()); err != nil {
			log.Println("Zakat processing error:", err)
		}
	})
	c.Start()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logService)
	walletHandler := handlers.NewWalletHandler(walletService, logService)
	transactionHandler := handlers.NewTransactionHandler(transactionService, walletService, logService)
	miningHandler := handlers.NewMiningHandler(miningService, logService)
	blockHandler := handlers.NewBlockHandler(blockchainService)
	zakatHandler := handlers.NewZakatHandler(zakatService)
	logHandler := handlers.NewLogHandler(logService)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
			auth.PUT("/profile", middleware.AuthMiddleware(), authHandler.UpdateProfile)
		}

		// Wallet routes (protected)
		wallet := api.Group("/wallet")
		wallet.Use(middleware.AuthMiddleware())
		{
			wallet.GET("", walletHandler.GetWallet)
			wallet.GET("/balance", walletHandler.GetBalance)
			wallet.GET("/utxos", walletHandler.GetUTXOs)
			wallet.POST("/beneficiaries", walletHandler.AddBeneficiary)
			wallet.DELETE("/beneficiaries/:id", walletHandler.RemoveBeneficiary)
			wallet.GET("/beneficiaries", walletHandler.GetBeneficiaries)
		}

		// Transaction routes (protected)
		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware())
		{
			transactions.POST("/send", transactionHandler.SendMoney)
			transactions.GET("/history", transactionHandler.GetHistory)
			transactions.GET("/pending", transactionHandler.GetPending)
		}

		// Mining routes (protected)
		mining := api.Group("/mining")
		mining.Use(middleware.AuthMiddleware())
		{
			mining.POST("/mine", miningHandler.Mine)
			mining.GET("/status", miningHandler.GetStatus)
		}

		// Block explorer routes (public)
		blocks := api.Group("/blocks")
		{
			blocks.GET("", blockHandler.GetAllBlocks)
			blocks.GET("/latest", blockHandler.GetLatestBlock)
			blocks.GET("/:hash", blockHandler.GetBlockByHash)
		}

		// Zakat routes (protected)
		zakat := api.Group("/zakat")
		zakat.Use(middleware.AuthMiddleware())
		{
			zakat.GET("/history", zakatHandler.GetHistory)
			zakat.POST("/process", zakatHandler.ProcessZakat)
		}

		// Log routes (protected)
		logs := api.Group("/logs")
		logs.Use(middleware.AuthMiddleware())
		{
			logs.GET("/system", logHandler.GetSystemLogs)
			logs.GET("/transactions", logHandler.GetTransactionLogs)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
