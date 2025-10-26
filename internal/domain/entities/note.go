package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

// NoteType represents the type of note
type NoteType string

const (
	NoteTypeText     NoteType = "text"
	NoteTypeTask     NoteType = "task"
	NoteTypeContact  NoteType = "contact"
	NoteTypeIdea     NoteType = "idea"
	NoteTypeReminder NoteType = "reminder"
)

// Priority represents task priority
type Priority string

const (
	PriorityLow  Priority = "low"
	PriorityHigh Priority = "high"
)

// Status represents task status
type Status string

const (
	StatusToDo      Status = "por_hacer"
	StatusCompleted Status = "completed"
)

// Metadata represents additional fields for different note types
type Metadata struct {
	// Task fields
	Priority       Priority   `json:"priority,omitempty"`
	Status         Status     `json:"status,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	Assignee       string     `json:"assignee,omitempty"`
	EstimatedHours int        `json:"estimated_hours,omitempty"`

	// Contact fields
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Address string `json:"address,omitempty"`

	// Reminder fields
	ReminderTime string `json:"reminder_time,omitempty"`

	// General fields
	Tags     []string `json:"tags,omitempty"`
	Category string   `json:"category,omitempty"`
}

// Note represents a note entity in our domain
type Note struct {
	ID        string    `json:"id"`
	Type      NoteType  `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Date      string    `json:"date"` // Format: YYYY-MM-DD
	Metadata  Metadata  `json:"metadata,omitempty"`
}

// NewNote creates a new note with the given content
func NewNote(content string) *Note {
	now := time.Now()
	return &Note{
		ID:        generateID(),
		Type:      NoteTypeText,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      now.Format("2006-01-02"),
		Metadata:  Metadata{},
	}
}

// NewTask creates a new task note
func NewTask(content string, priority Priority) *Note {
	now := time.Now()
	return &Note{
		ID:        generateID(),
		Type:      NoteTypeTask,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      now.Format("2006-01-02"),
		Metadata: Metadata{
			Priority: priority,
			Status:   StatusToDo,
		},
	}
}

// NewContact creates a new contact note
func NewContact(content, phone, email string) *Note {
	now := time.Now()
	return &Note{
		ID:        generateID(),
		Type:      NoteTypeContact,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      now.Format("2006-01-02"),
		Metadata: Metadata{
			Phone: phone,
			Email: email,
		},
	}
}

// NewReminder creates a new reminder note
func NewReminder(content string, reminderTime string, status Status) *Note {
	now := time.Now()
	return &Note{
		ID:        generateID(),
		Type:      NoteTypeReminder,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      now.Format("2006-01-02"),
		Metadata: Metadata{
			ReminderTime: reminderTime,
			Status:       status,
		},
	}
}

// NewNoteWithDate creates a new note with specific date
func NewNoteWithDate(content string, date time.Time) *Note {
	now := time.Now()
	return &Note{
		ID:        generateID(),
		Type:      NoteTypeText,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      date.Format("2006-01-02"),
		Metadata:  Metadata{},
	}
}

// ToJSON converts the note to JSON
func (n *Note) ToJSON() ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}

// FromJSON creates a note from JSON
func FromJSON(data []byte) (*Note, error) {
	var note Note
	err := json.Unmarshal(data, &note)
	return &note, err
}

// String returns a string representation of the note
func (n *Note) String() string {
	base := fmt.Sprintf("[%s] %s", n.CreatedAt.Format("15:04:05"), n.Content)

	switch n.Type {
	case NoteTypeTask:
		return fmt.Sprintf("%s [%s, %s]", base, n.Metadata.Priority, n.Metadata.Status)
	case NoteTypeContact:
		if n.Metadata.Phone != "" {
			return fmt.Sprintf("%s [%s]", base, n.Metadata.Phone)
		}
		if n.Metadata.Email != "" {
			return fmt.Sprintf("%s [%s]", base, n.Metadata.Email)
		}
	case NoteTypeReminder:
		timeStr := n.Metadata.ReminderTime
		if timeStr == "" {
			timeStr = "09:00"
		}
		return fmt.Sprintf("%s [%s, %s]", base, timeStr, n.Metadata.Status)
	}

	return base
}

// generateID generates a simple ID based on timestamp and random
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
