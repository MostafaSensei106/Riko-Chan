package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/MostafaSensei106/Riko-Chan/internal/models"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) GetByID(id uuid.UUID) (*models.Message, error) {
	var message models.Message
	if err := r.db.First(&message, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) GetUserMessages(userID int64, status models.MessageStatus, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	query := r.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("scheduled_time ASC").Limit(limit).Offset(offset).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) GetPendingMessages(before time.Time) ([]*models.Message, error) {
	var messages []*models.Message
	if err := r.db.
		Where("status = ? AND scheduled_time <= ?", "pending", before).
		Order("scheduled_time ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *MessageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Message{}, "id = ?", id).Error
}

func (r *MessageRepository) UpdateStatus(id uuid.UUID, status models.MessageStatus) error {
	return r.db.Model(&models.Message{}).
		Where("id = ?", id).
		Update("status", status).Error
}
