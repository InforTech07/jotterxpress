package cli

import (
	"fmt"
	"jotterxpress/internal/domain/entities"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	reminderContent = iota
	reminderTime
)

const (
	reminderPink = lipgloss.Color("#FFD700")
)

var (
	reminderInputStyle    = lipgloss.NewStyle().Foreground(reminderPink)
	reminderLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true)
	reminderContinueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	reminderTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
)

type ReminderFormModel struct {
	inputs           []textinput.Model
	focused          int
	err              error
	done             bool
	reminder         *entities.Note
	existingReminder *entities.Note // For updating existing reminders
}

// NewReminderFormModel creates a new reminder form model
func NewReminderFormModel() *ReminderFormModel {
	var inputs []textinput.Model = make([]textinput.Model, 2)

	// Reminder content input
	inputs[reminderContent] = textinput.New()
	inputs[reminderContent].Placeholder = "Enter reminder description..."
	inputs[reminderContent].Focus()
	inputs[reminderContent].CharLimit = 200
	inputs[reminderContent].Width = 50

	// Time input
	inputs[reminderTime] = textinput.New()
	inputs[reminderTime].Placeholder = "09:00"
	inputs[reminderTime].CharLimit = 5
	inputs[reminderTime].Width = 10
	inputs[reminderTime].Validate = timeValidator

	return &ReminderFormModel{
		inputs:   inputs,
		focused:  0,
		err:      nil,
		done:     false,
		reminder: nil,
	}
}

// NewReminderFormModelWithData creates a new reminder form model with existing data
func NewReminderFormModelWithData(reminder *entities.Note) *ReminderFormModel {
	var inputs []textinput.Model = make([]textinput.Model, 2)

	// Reminder content input
	inputs[reminderContent] = textinput.New()
	inputs[reminderContent].Placeholder = "Enter reminder description..."
	inputs[reminderContent].Focus()
	inputs[reminderContent].CharLimit = 200
	inputs[reminderContent].Width = 50
	inputs[reminderContent].SetValue(reminder.Content) // Set existing content

	// Time input
	inputs[reminderTime] = textinput.New()
	inputs[reminderTime].Placeholder = "09:00"
	inputs[reminderTime].CharLimit = 5
	inputs[reminderTime].Width = 10
	inputs[reminderTime].Validate = timeValidator
	if reminder.Metadata.ReminderTime != "" {
		inputs[reminderTime].SetValue(reminder.Metadata.ReminderTime) // Set existing time
	}

	return &ReminderFormModel{
		inputs:           inputs,
		focused:          0,
		err:              nil,
		done:             false,
		reminder:         nil,
		existingReminder: reminder,
	}
}

// Init initializes the model
func (m ReminderFormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m ReminderFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			// Create the reminder and finish
			m.createReminder()
			if m.err != nil {
				// If there's an error, don't quit, show the error
				return &m, nil
			}
			return &m, tea.Quit
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				// Create the reminder and finish
				m.createReminder()
				if m.err != nil {
					// If there's an error, don't quit, show the error
					return &m, nil
				}
				return &m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return &m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return &m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return &m, tea.Batch(cmds...)
}

// View renders the form
func (m ReminderFormModel) View() string {
	title := reminderTitleStyle.Render("Create New Reminder")

	content := fmt.Sprintf(`
%s

 %s
 %s

 %s
 %s

 %s
`,
		title,
		reminderLabelStyle.Width(50).Render("Reminder Description"),
		m.inputs[reminderContent].View(),
		reminderLabelStyle.Width(10).Render("Time"),
		m.inputs[reminderTime].View(),
		reminderContinueStyle.Render("Press Ctrl+S to create reminder, Tab to navigate, Ctrl+C to cancel"),
	)

	if m.err != nil {
		content += fmt.Sprintf("\n %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return content + "\n"
}

// createReminder creates the reminder from form data
func (m *ReminderFormModel) createReminder() {
	content := strings.TrimSpace(m.inputs[reminderContent].Value())
	if content == "" {
		m.err = fmt.Errorf("reminder content is required")
		return
	}

	timeStr := strings.TrimSpace(m.inputs[reminderTime].Value())
	if timeStr == "" {
		// Default to 09:00 if no time provided
		timeStr = "09:00"
	}

	// Status is always pending for new reminders
	status := entities.StatusToDo

	if m.existingReminder != nil {
		// Update existing reminder
		m.existingReminder.Content = content
		m.existingReminder.Metadata.ReminderTime = timeStr
		// Don't change the status when updating
		m.existingReminder.UpdatedAt = time.Now()
		m.reminder = m.existingReminder
	} else {
		// Create new reminder - always pending
		reminder := entities.NewReminder(content, timeStr, status)
		m.reminder = reminder
	}

	m.done = true
	m.err = nil // Clear any previous errors
}

// GetReminder returns the created reminder
func (m *ReminderFormModel) GetReminder() *entities.Note {
	return m.reminder
}

// IsDone returns true if the form is completed
func (m *ReminderFormModel) IsDone() bool {
	return m.done
}

// nextInput focuses the next input field
func (m *ReminderFormModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *ReminderFormModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

// timeValidator validates time input (HH:MM format)
func timeValidator(s string) error {
	if s == "" {
		return nil // Optional field
	}

	// Basic time validation
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return fmt.Errorf("time must be in HH:MM format")
	}

	// Validate hours
	if len(parts[0]) != 2 {
		return fmt.Errorf("hours must be 2 digits")
	}

	// Validate minutes
	if len(parts[1]) != 2 {
		return fmt.Errorf("minutes must be 2 digits")
	}

	// Check if it parses as a valid time
	_, err := time.Parse("15:04", s)
	if err != nil {
		return fmt.Errorf("invalid time format")
	}

	return nil
}

// statusValidator validates status input
func statusValidator(s string) error {
	if s == "" {
		return nil // Optional field
	}

	validStatuses := []string{"pending", "completed"}
	s = strings.ToLower(strings.TrimSpace(s))

	for _, valid := range validStatuses {
		if s == valid {
			return nil
		}
	}

	return fmt.Errorf("status must be: pending or completed")
}
