package models

import "time"

type UserLanguage string

const (
	LanguageEnglish  UserLanguage = "en"
	LanguageArabic   UserLanguage = "ar"
	LanguageJapanese UserLanguage = "ja"
)

type User struct {
	ID           int64        `json:"id"`
	Username     *string      `json:"username"`
	FirstName    string       `json:"first_name"`
	LastName     *string      `json:"last_name"`
	Language     UserLanguage `json:"language"`
	Timezone     string       `json:"timezone"`
	GoogleTokens *string      `json:"google_tokens"`
	NotionToken  *string      `json:"notion_token"`
	TrelloToken  *string      `json:"trello_token"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func NewUser(id int64, firstName string) *User {
	return &User{
		ID:        id,
		FirstName: firstName,
		Language:  LanguageEnglish,
		Timezone:  "UTC",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
