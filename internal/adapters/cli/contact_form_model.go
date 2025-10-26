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
	contactName = iota
	contactPhone
	contactEmail
)

const (
	contactPink = lipgloss.Color("#8B5CF6")
)

var (
	contactInputStyle    = lipgloss.NewStyle().Foreground(contactPink)
	contactLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true)
	contactContinueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676"))
	contactTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1)
)

type ContactFormModel struct {
	inputs          []textinput.Model
	focused         int
	err             error
	done            bool
	contact         *entities.Note
	existingContact *entities.Note // For updating existing contacts
}

// NewContactFormModel creates a new contact form model
func NewContactFormModel() *ContactFormModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	// Contact name input
	inputs[contactName] = textinput.New()
	inputs[contactName].Placeholder = "Enter contact name..."
	inputs[contactName].Focus()
	inputs[contactName].CharLimit = 100
	inputs[contactName].Width = 50
	inputs[contactName].Prompt = ""

	// Phone input
	inputs[contactPhone] = textinput.New()
	inputs[contactPhone].Placeholder = "Enter phone number (optional)"
	inputs[contactPhone].CharLimit = 20
	inputs[contactPhone].Width = 30
	inputs[contactPhone].Prompt = ""
	inputs[contactPhone].Validate = phoneValidator

	// Email input
	inputs[contactEmail] = textinput.New()
	inputs[contactEmail].Placeholder = "Enter email address (optional)"
	inputs[contactEmail].CharLimit = 100
	inputs[contactEmail].Width = 40
	inputs[contactEmail].Prompt = ""
	inputs[contactEmail].Validate = emailValidator

	return &ContactFormModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
		done:    false,
		contact: nil,
	}
}

// NewContactFormModelWithData creates a new contact form model with existing data
func NewContactFormModelWithData(contact *entities.Note) *ContactFormModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	// Contact name input
	inputs[contactName] = textinput.New()
	inputs[contactName].Placeholder = "Enter contact name..."
	inputs[contactName].Focus()
	inputs[contactName].CharLimit = 100
	inputs[contactName].Width = 50
	inputs[contactName].Prompt = ""
	inputs[contactName].SetValue(contact.Content) // Set existing name

	// Phone input
	inputs[contactPhone] = textinput.New()
	inputs[contactPhone].Placeholder = "Enter phone number (optional)"
	inputs[contactPhone].CharLimit = 20
	inputs[contactPhone].Width = 30
	inputs[contactPhone].Prompt = ""
	inputs[contactPhone].Validate = phoneValidator
	inputs[contactPhone].SetValue(contact.Metadata.Phone) // Set existing phone

	// Email input
	inputs[contactEmail] = textinput.New()
	inputs[contactEmail].Placeholder = "Enter email address (optional)"
	inputs[contactEmail].CharLimit = 100
	inputs[contactEmail].Width = 40
	inputs[contactEmail].Prompt = ""
	inputs[contactEmail].Validate = emailValidator
	inputs[contactEmail].SetValue(contact.Metadata.Email) // Set existing email

	return &ContactFormModel{
		inputs:          inputs,
		focused:         0,
		err:             nil,
		done:            false,
		contact:         nil,
		existingContact: contact,
	}
}

// Init initializes the model
func (m ContactFormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (m ContactFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			// Create the contact and finish
			m.createContact()
			if m.err != nil {
				// If there's an error, don't quit, show the error
				return &m, nil
			}
			return &m, tea.Quit
		case tea.KeyEnter:
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
func (m ContactFormModel) View() string {
	title := contactTitleStyle.Render("Create New Contact")

	content := fmt.Sprintf(`
%s

 %s
 %s

 %s  %s
 %s  %s

 %s
`,
		title,
		contactLabelStyle.Width(50).Render("Contact Name"),
		m.inputs[contactName].View(),
		contactLabelStyle.Width(30).Render("Phone"),
		contactLabelStyle.Width(40).Render("Email"),
		m.inputs[contactPhone].View(),
		m.inputs[contactEmail].View(),
		contactContinueStyle.Render("Press Ctrl+S to create contact, Tab to navigate, Ctrl+C to cancel"),
	)

	if m.err != nil {
		content += fmt.Sprintf("\n %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return content + "\n"
}

// createContact creates the contact from form data
func (m *ContactFormModel) createContact() {
	name := strings.TrimSpace(m.inputs[contactName].Value())
	if name == "" {
		m.err = fmt.Errorf("contact name is required")
		return
	}

	phone := strings.TrimSpace(m.inputs[contactPhone].Value())
	email := strings.TrimSpace(m.inputs[contactEmail].Value())

	if phone == "" && email == "" {
		m.err = fmt.Errorf("at least phone or email is required")
		return
	}

	if m.existingContact != nil {
		// Update existing contact
		m.existingContact.Content = name
		m.existingContact.Metadata.Phone = phone
		m.existingContact.Metadata.Email = email
		m.existingContact.UpdatedAt = time.Now()
		m.contact = m.existingContact
	} else {
		// Create new contact
		contact := entities.NewContact(name, phone, email)
		m.contact = contact
	}

	m.done = true
	m.err = nil // Clear any previous errors
}

// GetContact returns the created contact
func (m *ContactFormModel) GetContact() *entities.Note {
	return m.contact
}

// IsDone returns true if the form is completed
func (m *ContactFormModel) IsDone() bool {
	return m.done
}

// nextInput focuses the next input field
func (m *ContactFormModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *ContactFormModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

// phoneValidator validates phone input
func phoneValidator(s string) error {
	if s == "" {
		return nil // Optional field
	}

	// Basic phone validation - should contain only digits, spaces, +, -, (, )
	allowedChars := "0123456789+-() "
	for _, char := range s {
		if !strings.ContainsRune(allowedChars, char) {
			return fmt.Errorf("phone number contains invalid characters")
		}
	}

	return nil
}

// emailValidator validates email input
func emailValidator(s string) error {
	if s == "" {
		return nil // Optional field
	}

	// Basic email validation
	if !strings.Contains(s, "@") || !strings.Contains(s, ".") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
