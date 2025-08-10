package services

import (
	"github.com/MostafaSensei106/Riko-Chan/internal/db"
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
