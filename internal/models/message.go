package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypePhoto    MessageType = "photo"
	MessageTypeDocument MessageType = "document"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeLocation MessageType = "location"
)
type RecurrenceType string

const (
	RecurrenceNone    RecurrenceType = "none"
	RecurrenceDaily   RecurrenceType = "daily"
	RecurrenceWeekly  RecurrenceType = "weekly"
	RecurrenceMonthly RecurrenceType = "monthly"
	RecurrenceYearly  RecurrenceType = "yearly"
)

type MessageStatus string

const (
	MessageStatusPending MessageStatus = "pending"
	MessageStatusSent    MessageStatus = "sent"
	MessageStatusCancelled MessageStatus = "cancelled"
	MessageStatusFailed  MessageStatus = "failed"
)


type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Title     string  `json:"title,omitempty"`
	Address   string  `json:"address,omitempty"`
}


type Message struct {
	ID uuid.UUID `json:"id" db:"id"`
	UserID int64 `json:"user_id" db:"user_id"`
	RecipientID *int64 `json:"recipient_id" db:"recipient_id"`
	GroupID *string `json:"group_id" db:"group_id"`
	ChannelID *string `json:"channel_id" db:"channel_id"`
	MessageType MessageType `json:"message_type" db:"message_type"`
	Content string `json:"content" db:"content"`
	MediaFileID *string `json:"media_file_id" db:"media_file_id"`
	Location *Location `json:"location" db:"location"`
	ScheduledTime time.Time `json:"scheduled_time" db:"scheduled_time"`
	Status MessageStatus `json:"status" db:"status"`
	RecurrenceType RecurrenceType `json:"recurrence_type" db:"recurrence_type"`
	RecurrenceCount int `json:"recurrence_count" db:"recurrence_count"`
	MaxRecurrences *int `json:"max_recurrences" db:"max_recurrences"`
	NotifyBefore *time.Duration `json:"notify_before" db:"notify_before"`
	PrivateViewMode bool `json:"private_view_mode" db:"private_view_mode"`
	GoogleCalendarID *string `json:"google_calendar_id" db:"google_calendar_id"`
	NotionPageID *string `json:"notion_page_id" db:"notion_page_id"`
	TrelloCardID *string `json:"trello_card_id" db:"trello_card_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}


func NewMessage (userID int64, messageType MessageType, content string) *Message {
	return &Message{
		ID: uuid.New(),
		UserID: userID,
		MessageType: messageType,
		Content: content,
		Status: MessageStatusPending,
		RecurrenceType: RecurrenceNone,
		RecurrenceCount: 0,
		PrivateViewMode: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
