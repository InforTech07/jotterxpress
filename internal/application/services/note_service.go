package services

import (
	"fmt"
	"jotterxpress/internal/domain/entities"
	"jotterxpress/internal/domain/ports"
	"strings"
	"time"
)

// noteService implements the NoteService interface
type noteService struct {
	repository ports.NoteRepository
}

// NewNoteService creates a new note service
func NewNoteService(repository ports.NoteRepository) ports.NoteService {
	return &noteService{
		repository: repository,
	}
}

// CreateNote creates a new note with the given content
func (s *noteService) CreateNote(content string) (*entities.Note, error) {
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("note content cannot be empty")
	}

	note := entities.NewNote(content)

	if err := s.repository.Save(note); err != nil {
		return nil, fmt.Errorf("failed to save note: %w", err)
	}

	return note, nil
}

// SaveNote saves an existing note
func (s *noteService) SaveNote(note *entities.Note) error {
	if strings.TrimSpace(note.Content) == "" {
		return fmt.Errorf("note content cannot be empty")
	}

	// Update the updated_at timestamp
	note.UpdatedAt = time.Now()

	if err := s.repository.Save(note); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}

	return nil
}

// GetTodayNotes retrieves all notes for today
func (s *noteService) GetTodayNotes() ([]*entities.Note, error) {
	notes, err := s.repository.GetTodayNotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get today's notes: %w", err)
	}

	return notes, nil
}

// GetNotesByDate retrieves notes for a specific date
func (s *noteService) GetNotesByDate(date string) ([]*entities.Note, error) {
	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	notes, err := s.repository.GetNotesByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes for date %s: %w", date, err)
	}

	return notes, nil
}

// GetNotesByMonth retrieves notes for a specific month
func (s *noteService) GetNotesByMonth(monthStr string) ([]*entities.Note, error) {
	// Parse month string to validate format
	if _, err := time.Parse("2006-01", monthStr); err != nil {
		return nil, fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
	}

	notes, err := s.repository.GetNotesByMonth(monthStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes for month %s: %w", monthStr, err)
	}

	return notes, nil
}

// ListNotes formats and returns notes for display
func (s *noteService) ListNotes(notes []*entities.Note) string {
	if len(notes) == 0 {
		return "No notes found."
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📝 Notes (%d found):\n\n", len(notes)))

	for i, note := range notes {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, note.String()))
	}

	return result.String()
}
