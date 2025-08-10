package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/MostafaSensei106/Riko-Chan/internal/models"
	"github.com/MostafaSensei106/Riko-Chan/internal/services"
	"github.com/MostafaSensei106/Riko-Chan/internal/utils"
)

const (
	ScheduledMessagesKey = "scheduled_messages"
	NotificationsKey     = "notifications"
)

type ScheduledMessage struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    int64     `json:"user_id"`
}

type Scheduler struct {
	redis          *RedisClient
	messageService *services.MessageService
	logger         *utils.Logger
}

func NewScheduler(redis *RedisClient, messageService *services.MessageService, logger *utils.Logger) *Scheduler {
	return &Scheduler{
		redis:          redis,
		messageService: messageService,
		logger:         logger,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	s.logger.Info("Scheduler started")

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Scheduler stopping...")
			return
		case <-ticker.C:
			if err := s.processScheduledMessages(ctx); err != nil {
				s.logger.Error("Failed to process scheduled messages", "error", err)
			}
			if err := s.processNotifications(ctx); err != nil {
				s.logger.Error("Failed to process notifications", "error", err)
			}
		}
	}
}

func (s *Scheduler) ScheduleMessage(ctx context.Context, message *models.Message) error {
	scheduledMsg := ScheduledMessage{
		MessageID: message.ID,
		UserID:    message.UserID,
	}

	score := float64(message.ScheduledTime.Unix())
	if err := s.redis.ZAdd(ctx, ScheduledMessagesKey, score, scheduledMsg); err != nil {
		return fmt.Errorf("failed to schedule message: %w", err)
	}

	// Schedule notification if needed
	if message.NotifyBefore != nil {
		notifyTime := message.ScheduledTime.Add(-*message.NotifyBefore)
		if notifyTime.After(time.Now()) {
			notificationScore := float64(notifyTime.Unix())
			if err := s.redis.ZAdd(ctx, NotificationsKey, notificationScore, scheduledMsg); err != nil {
				s.logger.Error("Failed to schedule notification", "error", err)
			}
		}
	}

	s.logger.Info("Message scheduled", "message_id", message.ID, "scheduled_time", message.ScheduledTime)
	return nil
}

func (s *Scheduler) CancelMessage(ctx context.Context, messageID uuid.UUID, userID int64) error {
	scheduledMsg := ScheduledMessage{
		MessageID: messageID,
		UserID:    userID,
	}

	if err := s.redis.ZRem(ctx, ScheduledMessagesKey, scheduledMsg); err != nil {
		return fmt.Errorf("failed to cancel scheduled message: %w", err)
	}

	if err := s.redis.ZRem(ctx, NotificationsKey, scheduledMsg); err != nil {
		s.logger.Error("Failed to cancel notification", "error", err)
	}

	s.logger.Info("Message cancelled", "message_id", messageID)
	return nil
}

func (s *Scheduler) processScheduledMessages(ctx context.Context) error {
	now := time.Now().Unix()
	members, err := s.redis.ZRangeByScore(ctx, ScheduledMessagesKey, "-inf", strconv.FormatInt(now, 10))
	if err != nil {
		return fmt.Errorf("failed to get scheduled messages: %w", err)
	}

	for _, member := range members {
		var scheduledMsg ScheduledMessage
		if err := json.Unmarshal([]byte(member), &scheduledMsg); err != nil {
			s.logger.Error("Failed to unmarshal scheduled message", "error", err, "member", member)
			continue
		}

		if err := s.messageService.SendScheduledMessage(ctx, scheduledMsg.MessageID); err != nil {
			s.logger.Error("Failed to send scheduled message", "error", err, "message_id", scheduledMsg.MessageID)
			continue
		}

		// Remove from scheduled messages
		if err := s.redis.ZRem(ctx, ScheduledMessagesKey, scheduledMsg); err != nil {
			s.logger.Error("Failed to remove sent message from schedule", "error", err, "message_id", scheduledMsg.MessageID)
		}
	}

	return nil
}

func (s *Scheduler) processNotifications(ctx context.Context) error {
	now := time.Now().Unix()
	members, err := s.redis.ZRangeByScore(ctx, NotificationsKey, "-inf", strconv.FormatInt(now, 10))
	if err != nil {
		return fmt.Errorf("failed to get notifications: %w", err)
	}

	for _, member := range members {
		var scheduledMsg ScheduledMessage
		if err := json.Unmarshal([]byte(member), &scheduledMsg); err != nil {
			s.logger.Error("Failed to unmarshal notification", "error", err, "member", member)
			continue
		}

		if err := s.messageService.SendNotification(ctx, scheduledMsg.MessageID); err != nil {
			s.logger.Error("Failed to send notification", "error", err, "message_id", scheduledMsg.MessageID)
			continue
		}

		// Remove from notifications
		if err := s.redis.ZRem(ctx, NotificationsKey, scheduledMsg); err != nil {
			s.logger.Error("Failed to remove sent notification", "error", err, "message_id", scheduledMsg.MessageID)
		}
	}

	return nil
}
