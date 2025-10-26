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

type (
	errMsg error
)

const (
	taskContent = iota
	taskPriority
	taskAssignee
)

const (
	hotPink   = lipgloss.Color("#04B575")
	darkGray  = lipgloss.Color("#767676")
	lightGray = lipgloss.Color("#626262")
)

var (
	inputStyle     = lipgloss.NewStyle().Foreground(hotPink)
	labelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true)
	continueStyle  = lipgloss.NewStyle().Foreground(darkGray)
	formTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
)

type TaskFormModel struct {
	inputs       []textinput.Model
	focused      int
	err          error
	done         bool
	task         *entities.Note
	existingTask *entities.Note // For updating existing tasks
}

// NewTaskFormModel creates a new task form model
func NewTaskFormModel() *TaskFormModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	// Task content input
	inputs[taskContent] = textinput.New()
	inputs[taskContent].Placeholder = "Enter task description..."
	inputs[taskContent].Focus()
	inputs[taskContent].CharLimit = 200
	inputs[taskContent].Width = 50

	// Priority input
	inputs[taskPriority] = textinput.New()
	inputs[taskPriority].Placeholder = "low, high"
	inputs[taskPriority].CharLimit = 10
	inputs[taskPriority].Width = 15
	inputs[taskPriority].Validate = priorityValidator

	// Assignee input
	inputs[taskAssignee] = textinput.New()
	inputs[taskAssignee].Placeholder = "Enter assignee name (optional)"
	inputs[taskAssignee].CharLimit = 50
	inputs[taskAssignee].Width = 30

	return &TaskFormModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
		done:    false,
		task:    nil,
	}
}

// NewTaskFormModelWithData creates a new task form model with existing data
func NewTaskFormModelWithData(task *entities.Note) *TaskFormModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	// Task content input
	inputs[taskContent] = textinput.New()
	inputs[taskContent].Placeholder = "Enter task description..."
	inputs[taskContent].Focus()
	inputs[taskContent].CharLimit = 200
	inputs[taskContent].Width = 50
	inputs[taskContent].SetValue(task.Content) // Set existing content

	// Priority input
	inputs[taskPriority] = textinput.New()
	inputs[taskPriority].Placeholder = "low, high"
	inputs[taskPriority].CharLimit = 10
	inputs[taskPriority].Width = 15
	inputs[taskPriority].Validate = priorityValidator
	inputs[taskPriority].SetValue(string(task.Metadata.Priority)) // Set existing priority

	// Assignee input
	inputs[taskAssignee] = textinput.New()
	inputs[taskAssignee].Placeholder = "Enter assignee name (optional)"
	inputs[taskAssignee].CharLimit = 50
	inputs[taskAssignee].Width = 30
	inputs[taskAssignee].SetValue(task.Metadata.Assignee) // Set existing assignee

	return &TaskFormModel{
		inputs:       inputs,
		focused:      0,
		err:          nil,
		done:         false,
		task:         nil,
		existingTask: task,
	}
}

// Init initializes the model
func (m TaskFormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m TaskFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				// Create the task and finish
				m.createTask()
				if m.err != nil {
					// If there's an error, don't quit, show the error
					return &m, nil
				}
				// Task created successfully, quit immediately
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
func (m TaskFormModel) View() string {
	title := formTitleStyle.Render("Create New Task")

	content := fmt.Sprintf(`
%s

 %s
 %s

 %s  %s
 %s  %s

 %s
`,
		title,
		labelStyle.Width(50).Render("Task Description"),
		m.inputs[taskContent].View(),
		labelStyle.Width(15).Render("Priority"),
		labelStyle.Width(30).Render("Assignee"),
		m.inputs[taskPriority].View(),
		m.inputs[taskAssignee].View(),
		continueStyle.Render("Press Enter to create task, Tab to navigate, Ctrl+C to cancel"),
	)

	if m.err != nil {
		content += fmt.Sprintf("\n %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return content + "\n"
}

// createTask creates the task from form data
func (m *TaskFormModel) createTask() {
	content := strings.TrimSpace(m.inputs[taskContent].Value())
	if content == "" {
		m.err = fmt.Errorf("task content is required")
		return
	}

	priorityStr := strings.TrimSpace(m.inputs[taskPriority].Value())
	if priorityStr == "" {
		priorityStr = "low"
	}

	var priority entities.Priority
	switch strings.ToLower(priorityStr) {
	case "low":
		priority = entities.PriorityLow
	case "high":
		priority = entities.PriorityHigh
	default:
		priority = entities.PriorityLow
	}

	task := entities.NewTask(content, priority)

	// Set optional fields
	assignee := strings.TrimSpace(m.inputs[taskAssignee].Value())
	if assignee != "" {
		task.Metadata.Assignee = assignee
	}

	if m.existingTask != nil {
		// Update existing task
		m.existingTask.Content = content
		m.existingTask.Metadata.Priority = priority
		m.existingTask.Metadata.Assignee = assignee
		m.existingTask.UpdatedAt = time.Now()
		m.task = m.existingTask
	} else {
		// Create new task
		m.task = task
	}

	m.done = true
	m.err = nil // Clear any previous errors

}

// GetTask returns the created task
func (m *TaskFormModel) GetTask() *entities.Note {
	return m.task
}

// IsDone returns true if the form is completed
func (m *TaskFormModel) IsDone() bool {
	return m.done
}

// nextInput focuses the next input field
func (m *TaskFormModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *TaskFormModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

// priorityValidator validates priority input
func priorityValidator(s string) error {
	if s == "" {
		return nil // Optional field
	}

	validPriorities := []string{"low", "high"}
	s = strings.ToLower(strings.TrimSpace(s))

	for _, valid := range validPriorities {
		if s == valid {
			return nil
		}
	}

	return fmt.Errorf("priority must be: low or high")
}
