package ports

import (
	"jotterxpress/internal/domain/entities"
)

// NoteRepository defines the interface for note persistence
type NoteRepository interface {
	// Save saves a note to the repository
	Save(note *entities.Note) error

	// GetNotesByDate retrieves all notes for a specific date
	GetNotesByDate(date string) ([]*entities.Note, error)

	// GetNotesByDateRange retrieves notes within a date range
	GetNotesByDateRange(startDate, endDate string) ([]*entities.Note, error)

	// GetTodayNotes retrieves all notes for today
	GetTodayNotes() ([]*entities.Note, error)

	// GetNotesByMonth retrieves notes for a specific month (format: "2025-10")
	GetNotesByMonth(monthStr string) ([]*entities.Note, error)

	// DeleteNote deletes a note by ID
	DeleteNote(id string) error
}
