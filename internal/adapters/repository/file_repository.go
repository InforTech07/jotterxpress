package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"jotterxpress/internal/domain/entities"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// fileRepository implements the NoteRepository interface using file system
type fileRepository struct {
	notesDir string
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(notesDir string) *fileRepository {
	return &fileRepository{
		notesDir: notesDir,
	}
}

// GetNotesDir returns the notes directory path
func (r *fileRepository) GetNotesDir() string {
	return r.notesDir
}

// Save saves a note to a file
func (r *fileRepository) Save(note *entities.Note) error {
	// Ensure notes directory exists
	if err := os.MkdirAll(r.notesDir, 0755); err != nil {
		return fmt.Errorf("failed to create notes directory: %w", err)
	}

	// Create filename based on date
	filename := fmt.Sprintf("%s.json", note.Date)
	filepath := filepath.Join(r.notesDir, filename)

	// Read existing notes
	notes, err := r.readNotesFromFile(filepath, note.Date)
	if err != nil {
		return fmt.Errorf("failed to read existing notes: %w", err)
	}

	// Check if note already exists (by ID) and update it, otherwise add as new
	found := false
	for i, existingNote := range notes {
		if existingNote.ID == note.ID {
			// Update existing note
			notes[i] = note
			found = true
			break
		}
	}

	// If note doesn't exist, add it as new
	if !found {
		notes = append(notes, note)
	}

	// Write all notes back to file
	if err := r.writeNotesToFile(filepath, notes); err != nil {
		return fmt.Errorf("failed to write notes to file: %w", err)
	}

	return nil
}

// GetNotesByDate retrieves all notes for a specific date
func (r *fileRepository) GetNotesByDate(date string) ([]*entities.Note, error) {
	filename := fmt.Sprintf("%s.json", date)
	filepath := filepath.Join(r.notesDir, filename)

	return r.readNotesFromFile(filepath, date)
}

// GetNotesByDateRange retrieves notes within a date range
func (r *fileRepository) GetNotesByDateRange(startDate, endDate string) ([]*entities.Note, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	var allNotes []*entities.Note

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		notes, err := r.GetNotesByDate(dateStr)
		if err != nil {
			// Skip files that don't exist
			continue
		}
		allNotes = append(allNotes, notes...)
	}

	// Sort notes by update time (most recently updated first)
	sort.Slice(allNotes, func(i, j int) bool {
		return allNotes[i].UpdatedAt.After(allNotes[j].UpdatedAt)
	})

	return allNotes, nil
}

// GetTodayNotes retrieves all notes for today
func (r *fileRepository) GetTodayNotes() ([]*entities.Note, error) {
	today := time.Now().Format("2006-01-02")
	return r.GetNotesByDate(today)
}

// GetNotesByMonth retrieves notes for a specific month
// monthStr format: "2025-10" (YYYY-MM)
func (r *fileRepository) GetNotesByMonth(monthStr string) ([]*entities.Note, error) {
	// Parse the month string to get the first and last day of the month
	monthTime, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid month format, expected YYYY-MM: %w", err)
	}

	// Get the first day of the month
	startDate := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, monthTime.Location())

	// Get the last day of the month
	lastDay := time.Date(monthTime.Year(), monthTime.Month()+1, 0, 0, 0, 0, 0, monthTime.Location())
	endDate := lastDay

	// Use the existing date range function
	return r.GetNotesByDateRange(
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
}

// DeleteNote deletes a note by ID (not implemented for file-based storage)
func (r *fileRepository) DeleteNote(id string) error {
	return fmt.Errorf("delete operation not supported in file-based repository")
}

// readNotesFromFile reads notes from a specific file
func (r *fileRepository) readNotesFromFile(filepath, date string) ([]*entities.Note, error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*entities.Note{}, nil
		}
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	var notes []*entities.Note
	decoder := json.NewDecoder(file)

	// Try to decode as array first (new format)
	if err := decoder.Decode(&notes); err != nil {
		// If that fails, try to migrate from old text format
		return r.migrateFromTextFormat(filepath, date)
	}

	// Sort notes with custom logic:
	// 1. Pending reminders first (NoteTypeReminder with status pending)
	// 2. Rest sorted by UpdatedAt (most recently updated first)
	sort.Slice(notes, func(i, j int) bool {
		ni := notes[i]
		nj := notes[j]

		// Check if either is a pending reminder
		isReminderPendingI := ni.Type == entities.NoteTypeReminder &&
			(ni.Metadata.Status == entities.StatusToDo || ni.Metadata.Status == "")
		isReminderPendingJ := nj.Type == entities.NoteTypeReminder &&
			(nj.Metadata.Status == entities.StatusToDo || nj.Metadata.Status == "")

		// If i is pending reminder and j is not, i comes first
		if isReminderPendingI && !isReminderPendingJ {
			return true
		}
		// If j is pending reminder and i is not, j comes first
		if !isReminderPendingI && isReminderPendingJ {
			return false
		}
		// If both are pending reminders, sort by UpdatedAt
		if isReminderPendingI && isReminderPendingJ {
			return ni.UpdatedAt.After(nj.UpdatedAt)
		}
		// For non-reminders or completed reminders, sort by UpdatedAt
		return ni.UpdatedAt.After(nj.UpdatedAt)
	})

	return notes, nil
}

// writeNotesToFile writes notes to a specific file
func (r *fileRepository) writeNotesToFile(filepath string, notes []*entities.Note) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(notes); err != nil {
		return fmt.Errorf("failed to encode notes to JSON: %w", err)
	}

	return nil
}

// migrateFromTextFormat migrates old text format to JSON format
func (r *fileRepository) migrateFromTextFormat(filepath, date string) ([]*entities.Note, error) {
	// Try to read as text file
	textFilepath := strings.Replace(filepath, ".json", ".txt", 1)
	file, err := os.Open(textFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*entities.Note{}, nil
		}
		return nil, fmt.Errorf("failed to open text file %s: %w", textFilepath, err)
	}
	defer file.Close()

	var notes []*entities.Note
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse line format: [HH:MM:SS] content
		if strings.HasPrefix(line, "[") {
			parts := strings.SplitN(line, "] ", 2)
			if len(parts) == 2 {
				timeStr := strings.TrimPrefix(parts[0], "[")
				content := parts[1]

				// Parse time
				createdAt, err := time.Parse("15:04:05", timeStr)
				if err != nil {
					continue // Skip malformed lines
				}

				// Create note with new structure
				note := &entities.Note{
					ID:        fmt.Sprintf("%s-%s", date, timeStr),
					Type:      entities.NoteTypeText,
					Content:   content,
					CreatedAt: createdAt,
					UpdatedAt: createdAt,
					Date:      date,
					Metadata:  entities.Metadata{},
				}

				notes = append(notes, note)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading text file: %w", err)
	}

	// Write notes in new JSON format
	if len(notes) > 0 {
		if err := r.writeNotesToFile(filepath, notes); err != nil {
			return nil, fmt.Errorf("failed to migrate notes to JSON: %w", err)
		}
		// Remove old text file
		os.Remove(textFilepath)
	}

	return notes, nil
}
