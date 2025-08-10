package services

import (
	"github.com/MostafaSensei106/Riko-Chan/config"
)

type NotificationService struct {
	config *config.Config
	//logger *utils.Logger
}

func NewNotification(ctf *config.Config,

// logger *utils.Logger
) *NotificationService {
	return &NotificationService{
		config: ctf,
		//	logger: logger,
	}
}

func (s *NotificationService) SendNotification(userID int64, message string) error {
	// TODO: Implement notification sending via Telegram API
	//s.logger.Info("Notification sent", "user_id", userID, "message", message)
	return nil
}
