package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/MostafaSensei106/Riko-Chan/internal/cache"
	"github.com/MostafaSensei106/Riko-Chan/internal/db"
	"github.com/MostafaSensei106/Riko-Chan/internal/models"
	"github.com/MostafaSensei106/Riko-Chan/internal/utils"
)

type MessageService struct {
	repo      *db.MessageRepository
	userRepo  *db.UserRepository
	redis     *cache.RedisClient
	scheduler *cache.Scheduler
	encryptor *utils.Encryptor
	logger    *utils.Logger
}

func NewMessageService(repo *db.MessageRepository, redis *cache.RedisClient, logger *utils.Logger) *MessageService {
	return &MessageService{
		repo:   repo,
		redis:  redis,
		logger: logger,
	}
}

func (s *MessageService) SetUserRepo(userRepo *db.UserRepository) {
	s.userRepo = userRepo
}

func (s *MessageService) SetScheduler(scheduler *cache.Scheduler) {
	s.scheduler = scheduler
}

func (s *MessageService) SetEncryptor(encryptor *utils.Encryptor) {
	s.encryptor = encryptor
}

func (s *MessageService) CreateMessage(ctx context.Context, message *models.Message) error {
	// Encrypt content if encryptor is available
	if s.encryptor != nil {
		encryptedContent, err := s.encryptor.Encrypt(message.Content)
		if err != nil {
			return fmt.Errorf("failed to encrypt message content: %w", err)
		}
		message.Content = encryptedContent
	}

	// Save to database
	if err := s.repo.Create(message); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Schedule message if scheduler is available
	if s.scheduler != nil {
		if err := s.scheduler.ScheduleMessage(ctx, message); err != nil {
			s.logger.Error("Failed to schedule message", "error", err, "message_id", message.ID)
		}
	}

	s.logger.Info("Message created", "message_id", message.ID, "user_id", message.UserID)
	return nil
}

func (s *MessageService) GetMessage(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	message, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Decrypt content if encryptor is available
	if s.encryptor != nil && message.Content != "" {
		decryptedContent, err := s.encryptor.Decrypt(message.Content)
		if err != nil {
			s.logger.Error("Failed to decrypt message content", "error", err, "message_id", id)
			return message, nil // Return message with encrypted content
		}
		message.Content = decryptedContent
	}

	return message, nil
}

func (s *MessageService) GetUserMessages(ctx context.Context, userID int64, status models.MessageStatus, limit, offset int) ([]*models.Message, error) {
	messages, err := s.repo.GetUserMessages(userID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user messages: %w", err)
	}

	// Decrypt content for each message
	if s.encryptor != nil {
		for _, message := range messages {
			if message.Content != "" {
				decryptedContent, err := s.encryptor.Decrypt(message.Content)
				if err != nil {
					s.logger.Error("Failed to decrypt message content", "error", err, "message_id", message.ID)
					continue
				}
				message.Content = decryptedContent
			}
		}
	}

	return messages, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	// Encrypt content if encryptor is available
	if s.encryptor != nil {
		encryptedContent, err := s.encryptor.Encrypt(message.Content)
		if err != nil {
			return fmt.Errorf("failed to encrypt message content: %w", err)
		}
		message.Content = encryptedContent
	}

	if err := s.repo.Update(message); err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	// Reschedule if needed
	if s.scheduler != nil && message.Status == models.MessageStatusPending {
		if err := s.scheduler.ScheduleMessage(ctx, message); err != nil {
			s.logger.Error("Failed to reschedule message", "error", err, "message_id", message.ID)
		}
	}

	s.logger.Info("Message updated", "message_id", message.ID)
	return nil
}

func (s *MessageService) CancelMessage(ctx context.Context, id uuid.UUID, userID int64) error {
	// Update status to cancelled
	if err := s.repo.UpdateStatus(id, models.MessageStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel message: %w", err)
	}

	// Remove from scheduler
	if s.scheduler != nil {
		if err := s.scheduler.CancelMessage(ctx, id, userID); err != nil {
			s.logger.Error("Failed to cancel scheduled message", "error", err, "message_id", id)
		}
	}

	s.logger.Info("Message cancelled", "message_id", id)
	return nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, id uuid.UUID, userID int64) error {
	// Cancel from scheduler first
	if s.scheduler != nil {
		if err := s.scheduler.CancelMessage(ctx, id, userID); err != nil {
			s.logger.Error("Failed to cancel scheduled message", "error", err, "message_id", id)
		}
	}

	// Delete from database
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	s.logger.Info("Message deleted", "message_id", id)
	return nil
}

func (s *MessageService) SendScheduledMessage(ctx context.Context, messageID uuid.UUID) error {
	message, err := s.GetMessage(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	if message.Status != models.MessageStatusPending {
		return fmt.Errorf("message is not in pending status")
	}

	// TODO: Implement actual message sending via Telegram API
	// This would involve sending the message to the recipient via telegram bot API

	// Update status to sent
	if err := s.repo.UpdateStatus(messageID, models.MessageStatusSent); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	// Handle recurrence
	if message.RecurrenceType != models.RecurrenceNone {
		if err := s.handleRecurrence(ctx, message); err != nil {
			s.logger.Error("Failed to handle message recurrence", "error", err, "message_id", messageID)
		}
	}

	s.logger.Info("Scheduled message sent", "message_id", messageID)
	return nil
}

func (s *MessageService) SendNotification(ctx context.Context, messageID uuid.UUID) error {
	_, err := s.GetMessage(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// TODO: Implement notification sending via Telegram API
	// Send a notification to the user about upcoming scheduled message

	s.logger.Info("Notification sent", "message_id", messageID)
	return nil
}

func (s *MessageService) handleRecurrence(ctx context.Context, message *models.Message) error {
	// Check if we've reached the max recurrences
	if message.MaxRecurrences != nil && message.RecurrenceCount >= *message.MaxRecurrences {
		return nil
	}

	// Create next occurrence
	nextMessage := *message
	nextMessage.ID = uuid.New()
	nextMessage.RecurrenceCount++
	nextMessage.Status = models.MessageStatusPending
	nextMessage.CreatedAt = time.Now()
	nextMessage.UpdatedAt = time.Now()

	// Calculate next scheduled time
	nextMessage.ScheduledTime = utils.GetNextRecurrenceTime(
		message.ScheduledTime,
		string(message.RecurrenceType),
		1,
	)

	// Create the recurring message
	if err := s.CreateMessage(ctx, &nextMessage); err != nil {
		return fmt.Errorf("failed to create recurring message: %w", err)
	}

	s.logger.Info("Recurring message created",
		"original_id", message.ID,
		"new_id", nextMessage.ID,
		"next_time", nextMessage.ScheduledTime)

	return nil
}
