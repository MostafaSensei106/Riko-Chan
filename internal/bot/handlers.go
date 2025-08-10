package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/MostafaSensei106/Riko-Chan/internal/models"
)

func (b *Bot) handleCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	command := message.Command()
	args := message.CommandArguments()

	switch strings.ToLower(command) {
	case "start":
		b.handleStartCommand(ctx, message, user)
	case "new":
		b.handleNewCommand(ctx, message, user, args)
	case "list":
		b.handleListCommand(ctx, message, user)
	case "cancel":
		b.handleCancelCommand(ctx, message, user, args)
	case "delete":
		b.handleDeleteCommand(ctx, message, user, args)
	case "settings":
		b.handleSettingsCommand(ctx, message, user)
	case "help":
		b.handleHelpCommand(ctx, message, user)
	default:
		b.sendMessage(message.Chat.ID, b.getText("unknown_command", user.Language), nil)
	}
}

func (b *Bot) handleStartCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	welcomeText := b.getText("welcome", user.Language)
	helpText := b.getText("help_text", user.Language)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(b.getText("new_message", user.Language)),
			tgbotapi.NewKeyboardButton(b.getText("my_messages", user.Language)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(b.getText("settings", user.Language)),
			tgbotapi.NewKeyboardButton(b.getText("help", user.Language)),
		),
	)

	b.sendMessage(message.Chat.ID, welcomeText+"\n\n"+helpText, &keyboard)
}

func (b *Bot) handleNewCommand(ctx context.Context, message *tgbotapi.Message, user *models.User, args string) {
	if args == "" {
		b.sendMessage(message.Chat.ID, b.getText("new_message_help", user.Language), nil)
		return
	}

	// Parse message and time
	parts := strings.SplitN(args, " at ", 2)
	if len(parts) != 2 {
		parts = strings.SplitN(args, " في ", 2) // Arabic support
	}

	if len(parts) != 2 {
		b.sendMessage(message.Chat.ID, b.getText("invalid_format", user.Language), nil)
		return
	}

	content := strings.TrimSpace(parts[0])
	timeStr := strings.TrimSpace(parts[1])

	// Parse time
	timeParser, err := utils.NewTimeParser(user.Timezone)
	if err != nil {
		b.logger.Error("Failed to create time parser", "error", err, "user_id", user.ID)
		b.sendMessage(message.Chat.ID, b.getText("error_occurred", user.Language), nil)
		return
	}

	scheduledTime, err := timeParser.ParseRelativeTime(timeStr)
	if err != nil {
		b.sendMessage(message.Chat.ID, b.getText("invalid_time_format", user.Language), nil)
		return
	}

	// Create message
	msg := models.NewMessage(user.ID, models.MessageTypeText, content)
	msg.ScheduledTime = scheduledTime
	msg.RecipientID = &user.ID // Send to self by default

	if err := b.messageService.CreateMessage(ctx, msg); err != nil {
		b.logger.Error("Failed to create message", "error", err, "user_id", user.ID)
		b.sendMessage(message.Chat.ID, b.getText("error_occurred", user.Language), nil)
		return
	}

	confirmText := fmt.Sprintf(b.getText("message_scheduled", user.Language),
		scheduledTime.Format("2006-01-02 15:04"), msg.ID)

	// Add inline keyboard for message options
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.getText("add_notification", user.Language), "notify_"+msg.ID.String()),
			tgbotapi.NewInlineKeyboardButtonData(b.getText("make_recurring", user.Language), "recur_"+msg.ID.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.getText("send_to_other", user.Language), "recipient_"+msg.ID.String()),
		),
	)

	b.sendMessage(message.Chat.ID, confirmText, &keyboard)
}

func (b *Bot) handleListCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	messages, err := b.messageService.GetUserMessages(ctx, user.ID, models.MessageStatusPending, 10, 0)
	if err != nil {
		b.logger.Error("Failed to get user messages", "error", err, "user_id", user.ID)
		b.sendMessage(message.Chat.ID, b.getText("error_occurred", user.Language), nil)
		return
	}

	if len(messages) == 0 {
		b.sendMessage(message.Chat.ID, b.getText("no_pending_messages", user.Language), nil)
		return
	}

	var responseText strings.Builder
	responseText.WriteString(b.getText("your_pending_messages", user.Language) + "\n\n")

	for i, msg := range messages {
		timeStr := msg.ScheduledTime.Format("2006-01-02 15:04")
		preview := msg.Content
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}

		responseText.WriteString(fmt.Sprintf("%d. %s\n📅 %s\n💬 %s\n🆔 %s\n\n",
			i+1, b.getText("message", user.Language), timeStr, preview, msg.ID.String()[:8]))
	}

	b.sendMessage(message.Chat.ID, responseText.String(), nil)
}

func (b *Bot) handleCancelCommand(ctx context.Context, message *tgbotapi.Message, user *models.User, args string) {
	if args == "" {
		b.sendMessage(message.Chat.ID, b.getText("cancel_help", user.Language), nil)
		return
	}

	// Try to parse UUID from args (could be short form)
	messageID, err := b.findMessageByShortID(ctx, user.ID, args)
	if err != nil {
		b.sendMessage(message.Chat.ID, b.getText("message_not_found", user.Language), nil)
		return
	}

	if err := b.messageService.CancelMessage(ctx, messageID, user.ID); err != nil {
		b.logger.Error("Failed to cancel message", "error", err, "user_id", user.ID, "message_id", messageID)
		b.sendMessage(message.Chat.ID, b.getText("error_occurred", user.Language), nil)
		return
	}

	b.sendMessage(message.Chat.ID, b.getText("message_cancelled", user.Language), nil)
}

func (b *Bot) handleDeleteCommand(ctx context.Context, message *tgbotapi.Message, user *models.User, args string) {
	if args == "" {
		b.sendMessage(message.Chat.ID, b.getText("delete_help", user.Language), nil)
		return
	}

	messageID, err := b.findMessageByShortID(ctx, user.ID, args)
	if err != nil {
		b.sendMessage(message.Chat.ID, b.getText("message_not_found", user.Language), nil)
		return
	}

	if err := b.messageService.DeleteMessage(ctx, messageID, user.ID); err != nil {
		b.logger.Error("Failed to delete message", "error", err, "user_id", user.ID, "message_id", messageID)
		b.sendMessage(message.Chat.ID, b.getText("error_occurred", user.Language), nil)
		return
	}

	b.sendMessage(message.Chat.ID, b.getText("message_deleted", user.Language), nil)
}

func (b *Bot) handleSettingsCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌍 "+b.getText("change_language", user.Language), "lang"),
			tgbotapi.NewInlineKeyboardButtonData("🕒 "+b.getText("change_timezone", user.Language), "timezone"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔗 "+b.getText("integrations", user.Language), "integrations"),
		),
	)

	settingsText := fmt.Sprintf(b.getText("current_settings", user.Language),
		user.Language, user.Timezone)

	b.sendMessage(message.Chat.ID, settingsText, &keyboard)
}

func (b *Bot) handleHelpCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	helpText := b.getText("detailed_help", user.Language)
	b.sendMessage(message.Chat.ID, helpText, nil)
}

func (b *Bot) handleTextMessage(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	text := strings.ToLower(message.Text)

	switch {
	case strings.Contains(text, b.getText("new_message", user.Language)):
		b.sendMessage(message.Chat.ID, b.getText("new_message_prompt", user.Language), nil)
	case strings.Contains(text, b.getText("my_messages", user.Language)):
		b.handleListCommand(ctx, message, user)
	case strings.Contains(text, b.getText("settings", user.Language)):
		b.handleSettingsCommand(ctx, message, user)
	case strings.Contains(text, b.getText("help", user.Language)):
		b.handleHelpCommand(ctx, message, user)
	default:
		// Try to parse as a quick message format: "content at time"
		if strings.Contains(message.Text, " at ") || strings.Contains(message.Text, " في ") {
			b.handleNewCommand(ctx, message, user, message.Text)
		} else {
			b.sendMessage(message.Chat.ID, b.getText("unclear_message", user.Language), nil)
		}
	}
}

func (b *Bot) findMessageByShortID(ctx context.Context, userID int64, shortID string) (uuid.UUID, error) {
	messages, err := b.messageService.GetUserMessages(ctx, userID, "", 50, 0)
	if err != nil {
		return uuid.Nil, err
	}

	for _, msg := range messages {
		if strings.HasPrefix(msg.ID.String(), shortID) {
			return msg.ID, nil
		}
	}

	return uuid.Nil, fmt.Errorf("message not found")
}
