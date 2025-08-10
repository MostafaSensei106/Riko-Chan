package config


import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct{
	Telegram TelegramConfig
	Database DatabaseConfig
	Redis RedisConfig
	Integrations IntegrationsConfig
	Security SecurityConfig
	LogLevel string
}

type TelegramConfig struct {
	Token string
	Webhook string
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
	SSLMode string
}

type RedisConfig struct {
	Addr string
	Password string
	DB int
}

type IntegrationsConfig struct {
	GoogleCalendar GoogleCalendarConfig
	Notion NotionConfig
	Trello TrelloConfig
}

type GoogleCalendarConfig struct {
	ClientID string
	ClientSecret string
	RedirectURL string
}

type NotionConfig struct {
	Token string
}

type TrelloConfig struct {
	APIkey string
	Token string
}

type SecurityConfig struct {
	EncryptionKey string
}


func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


func Load() (*Config, error) {
	godotenv.Load();
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))


	config := &Config{
		Telegram: TelegramConfig{
			Token: getEnv("TELEGRAM_BOT_TOKEN", ""),
			Webhook: getEnv("TELEGRAM_WEBHOOK", ""),
		},
		Database: DatabaseConfig{
			Host : getEnv("DB_HOST", "localhost"),
			Port: strconv.Itoa(port),
			User: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName: getEnv("DB_NAME", "future_messages"),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB: redisDB,
		},
		Integrations: IntegrationsConfig{
			GoogleCalendar: GoogleCalendarConfig{
				ClientID: getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL: getEnv("GOOGLE_REDIRECT_URL", ""),
			},
			Notion: NotionConfig{
				Token: getEnv("NOTION_TOKEN", ""),
			},
			Trello: TrelloConfig{
				APIkey: getEnv("TRELLO_API_KEY", ""),
				Token: getEnv("TRELLO_TOKEN", ""),
			},
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
	return config, nil

}
