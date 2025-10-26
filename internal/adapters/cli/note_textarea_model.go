package cli

import (
	"fmt"
	"jotterxpress/internal/domain/entities"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type noteErrMsg error

type NoteTextareaModel struct {
	textarea     textarea.Model
	err          error
	note         *entities.Note
	existingNote *entities.Note // For updating existing notes
	cancelled    bool           // Track if user cancelled
}

var (
	noteTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#04B575")).Padding(0, 1)
	noteInfoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
)

// NewNoteTextareaModel creates a new note textarea model
func NewNoteTextareaModel() *NoteTextareaModel {
	ti := textarea.New()
	ti.Placeholder = "Escribe tu nota aquí..."
	ti.Focus()
	ti.CharLimit = 2000
	ti.SetWidth(80)
	ti.SetHeight(15)
	ti.ShowLineNumbers = false

	return &NoteTextareaModel{
		textarea: ti,
		err:      nil,
	}
}

// NewNoteTextareaModelWithContent creates a new note textarea model with existing content
func NewNoteTextareaModelWithContent(note *entities.Note) *NoteTextareaModel {
	ti := textarea.New()
	ti.Placeholder = "Escribe tu nota aquí..."
	ti.Focus()
	ti.CharLimit = 2000
	ti.SetWidth(80)
	ti.SetHeight(15)
	ti.ShowLineNumbers = false
	ti.SetValue(note.Content) // Set existing content

	return &NoteTextareaModel{
		textarea:     ti,
		err:          nil,
		existingNote: note, // Store the existing note
	}
}

func (m NoteTextareaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m NoteTextareaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			// Ctrl+S to save
			content := strings.TrimSpace(m.textarea.Value())
			if content != "" {
				if m.existingNote != nil {
					// Update existing note
					m.existingNote.Content = content
					m.existingNote.UpdatedAt = time.Now()
					m.note = m.existingNote
				} else {
					// Create new note
					m.note = entities.NewNote(content)
				}
				return &m, tea.Quit
			}
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			m.cancelled = true
			return &m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case noteErrMsg:
		m.err = msg
		return &m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return &m, tea.Batch(cmds...)
}

func (m NoteTextareaModel) View() string {
	title := noteTitleStyle.Render("Nota de Texto")

	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		title,
		m.textarea.View(),
		noteInfoStyle.Render("Ctrl+S para guardar • Ctrl+C para cancelar"),
	)

	return content + "\n"
}

// GetNote returns the created note
func (m *NoteTextareaModel) GetNote() *entities.Note {
	return m.note
}

// IsCancelled returns true if the user cancelled
func (m *NoteTextareaModel) IsCancelled() bool {
	return m.cancelled
}
