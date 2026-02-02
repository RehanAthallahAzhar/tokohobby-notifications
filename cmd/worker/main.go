package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	rmqLib "github.com/RehanAthallahAzhar/tokohobby-messaging/rabbitmq"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/configs"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/messaging"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/repositories"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/senders"
	"github.com/RehanAthallahAzhar/tokohobby-notifications/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger.Info("Starting Notification Worker...")

	// Load config
	cfg, err := configs.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load config")
	}

	// Set log level
	if level, err := logrus.ParseLevel(cfg.LogLevel); err == nil {
		logger.SetLevel(level)
	}

	// Connect to RabbitMQ
	rmqConfig := &rmqLib.RabbitMQConfig{
		URL:            cfg.RabbitMQ.URL,
		MaxRetries:     cfg.RabbitMQ.MaxRetries,
		RetryDelay:     cfg.RabbitMQ.RetryDelay,
		PrefetchCount:  cfg.RabbitMQ.PrefetchCount,
		ReconnectDelay: cfg.RabbitMQ.ReconnectDelay,
	}

	logger.WithField("url", rmqConfig.URL).Info("Connecting to RabbitMQ...")

	rmq, err := rmqLib.NewRabbitMQ(rmqConfig)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to RabbitMQ")
	}
	defer rmq.Close()

	logger.Info("Connected to RabbitMQ")

	// Setup exchanges
	if err := rmqLib.SetupOrderExchange(rmq); err != nil {
		logger.WithError(err).Fatal("Failed to setup order exchange")
	}

	logger.Info("Order exchange setup complete")

	// Connect to database
	db, err := configs.NewDatabaseConnection(&cfg.Database, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Initialize repository
	notifRepo := repositories.NewNotificationRepository(db, logger)

	// Initialize senders (mock mode)
	var emailSender, pushSender senders.Sender

	if cfg.MockMode {
		logger.Info("Using MOCK senders (demo mode)")
		emailSender = senders.NewMockEmailSender(logger)
		pushSender = senders.NewMockPushSender(logger)
	} else {
		logger.Info("Using REAL senders (production mode)")
		// TODO: Initialize real senders
		// emailSender = senders.NewSMTPSender(smtpConfig, logger)
		// pushSender = senders.NewFCMSender(fcmConfig, logger)
		logger.Fatal("Real senders not implemented yet")
	}

	// Initialize notification service
	notifService := services.NewNotificationService(notifRepo, emailSender, pushSender, logger)

	// Initialize consumers
	orderConsumer := messaging.NewOrderEventConsumer(rmq, notifService, logger)
	blogConsumer := messaging.NewBlogEventConsumer(rmq, notifService, logger)

	// Start consumers with context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start order consumer
	go func() {
		if err := orderConsumer.Start(ctx); err != nil {
			logger.WithError(err).Error("Order consumer error")
		}
	}()

	// Start blog consumer
	go func() {
		if err := blogConsumer.Start(ctx); err != nil {
			logger.WithError(err).Error("Blog consumer error")
		}
	}()

	logger.Info("Notification worker is running. Waiting for events... (Press Ctrl+C to exit)")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down notification worker...")
	cancel()

	// Give workers time to finish
	time.Sleep(2 * time.Second)

	logger.Info("Notification worker stopped gracefully")
}
