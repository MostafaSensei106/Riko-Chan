package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MostafaSensei106/Riko-Chan/config"
	"github.com/MostafaSensei106/Riko-Chan/internal/bot"
	"github.com/MostafaSensei106/Riko-Chan/internal/cache"
	"github.com/MostafaSensei106/Riko-Chan/internal/db"
	"github.com/MostafaSensei106/Riko-Chan/internal/services"
	"github.com/MostafaSensei106/Riko-Chan/internal/utils"
)

func Execute() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(cfg.LogLevel)

	// Initialize database
	database, err := db.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis
	redisClient := cache.NewRedisClient(cfg.Redis)
	defer redisClient.Close()

	// Initialize repositories
	userRepo := db.NewUserRepository(database)
	messageRepo := db.NewMessageRepository(database)

	// Initialize services
	messageService := services.NewMessageService(messageRepo, redisClient, logger)
	notificationService := services.NewNotificationService(cfg, logger)

	// Initialize scheduler
	scheduler := cache.NewScheduler(redisClient, messageService, logger)
	go scheduler.Start(context.Background())

	// Initialize bot
	telegramBot, err := bot.NewBot(cfg.Telegram, userRepo, messageService, notificationService, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize bot: %v", err)
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		cancel()
	}()

	// Start bot
	logger.Info("Starting Future Message Bot...")
	if err := telegramBot.Start(ctx); err != nil {
		logger.Fatalf("Bot stopped with error: %v", err)
	}
}
