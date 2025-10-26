package ports

import (
	"jotterxpress/internal/domain/entities"
)

// NoteService defines the interface for note business logic
type NoteService interface {
	// CreateNote creates a new note with the given content
	CreateNote(content string) (*entities.Note, error)

	// SaveNote saves an existing note
	SaveNote(note *entities.Note) error

	// GetTodayNotes retrieves all notes for today
	GetTodayNotes() ([]*entities.Note, error)

	// GetNotesByDate retrieves notes for a specific date
	GetNotesByDate(date string) ([]*entities.Note, error)

	// GetNotesByMonth retrieves notes for a specific month (format: "2025-10")
	GetNotesByMonth(monthStr string) ([]*entities.Note, error)

	// ListNotes formats and returns notes for display
	ListNotes(notes []*entities.Note) string
}
