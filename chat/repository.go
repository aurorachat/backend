package chat

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type EventPayload map[string]string

func (p EventPayload) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p EventPayload) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &p)
}

type User struct {
	*gorm.Model
	ID            int
	Conversations []Conversation `gorm:"many2many:user_conversations;"`
}

type Event struct {
	*gorm.Model
	ID string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	// Sender can be -1 if event is sent by the system
	Sender         int
	SenderName     string
	Type           string
	Payload        EventPayload
	ConversationID int
}

type Conversation struct {
	*gorm.Model
	ID     int
	Users  []User `gorm:"many2many:user_conversations;"`
	Events []Event
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	instance := &Repository{db: db}
	err := instance.initialize()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) initialize() error {
	err := r.db.AutoMigrate(&Event{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = r.db.AutoMigrate(&Conversation{})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddEvent(event *Event) error {
	err := r.db.Create(event).Error
	if err != nil {
		return err
	}
	return nil
}
