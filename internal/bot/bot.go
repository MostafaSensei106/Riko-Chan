package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MostafaSensei106/Riko-Chan/config"
	"github.com/MostafaSensei106/Riko-Chan/internal/db"
	"github.com/MostafaSensei106/Riko-Chan/internal/models"
	"github.com/MostafaSensei106/Riko-Chan/internal/services"
	"github.com/MostafaSensei106/Riko-Chan/internal/utils"
)

type Bot struct {
	api                 *tgbotapi.BotAPI
	userRepo            *db.UserRepository
	messageService      *services.NotificationService
	notificationService *services.NotificationService
	logger              *utils.Logger
}

func NewBot(cfg config.TelegramConfig, userRepo *db.UserRepository, messageService *services.MessageService, notificationService *services.NotificationService, logger *utils.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	return &Bot{
		api:                 api,
		userRepo:            userRepo,
		messageService:      messageService,
		notificationService: notificationService,
		logger:              logger,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("Bot started", "username", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("Bot stopping...")
			return nil
		case update := <-updates:
			if update.Message != nil {
				go b.handleMessage(ctx, update.Message)
			} else if update.CallbackQuery != nil {
				go b.handleCallbackQuery(ctx, update.CallbackQuery)
			}
		}
	}
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	// Ensure user exists in database
	user, err := b.userRepo.GetByID(message.From.ID)
	if err != nil {
		// Create new user
		user = &models.User{
			ID:        message.From.ID,
			Username:  &message.From.UserName,
			FirstName: message.From.FirstName,
			LastName:  &message.From.LastName,
			Language:  models.LanguageEnglish,
			Timezone:  "UTC",
		}
		if err := b.userRepo.Create(user); err != nil {
			b.logger.Error("Failed to create user", "error", err, "user_id", message.From.ID)
			return
		}
	}

	// Handle different message types
	if message.IsCommand() {
		b.handleCommand(ctx, message, user)
	} else {
		b.handleTextMessage(ctx, message, user)
	}
}

func (b *Bot) handleCallbackQuery(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	// Handle inline keyboard callbacks
	b.logger.Info("Callback query received", "data", callbackQuery.Data, "user_id", callbackQuery.From.ID)

	// Send callback answer to prevent loading state
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		b.logger.Error("Failed to send callback answer", "error", err)
	}

	// Handle different callback types
	// TODO: Implement callback handling logic
}

func (b *Bot) sendMessage(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)

	switch k := keyboard.(type) {
	case *tgbotapi.ReplyKeyboardMarkup:
		msg.ReplyMarkup = k
	case *tgbotapi.InlineKeyboardMarkup:
		msg.ReplyMarkup = k
	}

	if _, err := b.api.Send(msg); err != nil {
		b.logger.Error("Failed to send message", "error", err, "chat_id", chatID)
	}
}
