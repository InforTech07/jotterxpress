package cli

import (
	"fmt"
	"jotterxpress/internal/domain/entities"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#25A065")).
			Padding(1, 2)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	previewTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	previewInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return previewTitleStyle.Copy().BorderStyle(b)
	}()
)

// ContextMenuMsg is a message to open/close context menu
type ContextMenuMsg struct {
	Action string
}

// OpenTextareaMsg is a message to open textarea for editing
type OpenTextareaMsg struct {
	Note *entities.Note
}

// OpenContactFormMsg is a message to open contact form for editing
type OpenContactFormMsg struct {
	Note *entities.Note
}

// OpenTaskFormMsg is a message to open task form for editing
type OpenTaskFormMsg struct {
	Note *entities.Note
}

// OpenReminderFormMsg is a message to open reminder form for editing
type OpenReminderFormMsg struct {
	Note *entities.Note
}

// CompleteTaskMsg is a message to complete a task
type CompleteTaskMsg struct {
	Note *entities.Note
}

// CompleteReminderMsg is a message to complete a reminder
type CompleteReminderMsg struct {
	Note *entities.Note
}

// PreviewNoteMsg is a message to show preview of a note
type PreviewNoteMsg struct {
	Note *entities.Note
}

// DeleteNoteMsg is a message to delete a note
type DeleteNoteMsg struct {
	Note *entities.Note
}

// NoteItem represents a note item in the list
type NoteItem struct {
	note *entities.Note
}

func (i NoteItem) Title() string {
	// Return only the main content (note text, task description, or contact name)
	content := i.note.Content
	if len(content) > 20 {
		return content[:20] + "..."
	}
	return content
}

func (i NoteItem) Description() string {
	var meta []string

	// Add date/time first
	timeStr := i.note.CreatedAt.Format("15:04:05")
	meta = append(meta, timeStr)

	switch i.note.Type {
	case entities.NoteTypeTask:
		if i.note.Metadata.Priority != "" {
			priority := string(i.note.Metadata.Priority)
			if len(priority) > 20 {
				priority = priority[:20] + "..."
			}
			meta = append(meta, priority)
		}
		if i.note.Metadata.Status != "" {
			status := string(i.note.Metadata.Status)
			if len(status) > 20 {
				status = status[:20] + "..."
			}
			meta = append(meta, status)
		}
		if i.note.Metadata.Assignee != "" {
			assignee := i.note.Metadata.Assignee
			if len(assignee) > 20 {
				assignee = assignee[:20] + "..."
			}
			meta = append(meta, assignee)
		}
	case entities.NoteTypeContact:
		if i.note.Metadata.Phone != "" {
			phone := i.note.Metadata.Phone
			if len(phone) > 20 {
				phone = phone[:20] + "..."
			}
			meta = append(meta, phone)
		}
		if i.note.Metadata.Email != "" {
			email := i.note.Metadata.Email
			if len(email) > 20 {
				email = email[:20] + "..."
			}
			meta = append(meta, email)
		}
	case entities.NoteTypeReminder:
		if i.note.Metadata.ReminderTime != "" {
			reminderTime := i.note.Metadata.ReminderTime
			if len(reminderTime) > 20 {
				reminderTime = reminderTime[:20] + "..."
			}
			meta = append(meta, reminderTime)
		}
		if i.note.Metadata.Status != "" {
			status := string(i.note.Metadata.Status)
			if len(status) > 20 {
				status = status[:20] + "..."
			}
			meta = append(meta, status)
		}
	default:
		// For text notes, just show the date
	}

	return strings.Join(meta, " • ")
}

func (i NoteItem) FilterValue() string {
	return i.note.Content
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	quit             key.Binding
	selectItem       key.Binding
	editItem         key.Binding
	completeItem     key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		selectItem: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "select"),
		),
		editItem: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit item"),
		),
		completeItem: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "complete"),
		),
	}
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

// ShortHelp returns short help entries
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// FullHelp returns full help entries
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

type ListModel struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	notes        []*entities.Note
	title        string
	selected     map[int]struct{} // Track selected items
	showMenu     bool
	showPreview  bool
	selectedNote *entities.Note
	cli          *CLI // Reference to CLI for calling update methods
}

// NewListModel creates a new list model
func NewListModel(notes []*entities.Note, title string, cli *CLI) ListModel {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Convert notes to list items
	items := make([]list.Item, len(notes))
	for i, note := range notes {
		items[i] = NoteItem{note: note}
	}

	// Create the model first
	model := ListModel{
		keys:         listKeys,
		delegateKeys: delegateKeys,
		notes:        notes,
		title:        title,
		selected:     make(map[int]struct{}),
		showMenu:     false,
		showPreview:  false,
		selectedNote: nil,
		cli:          cli,
	}

	// Setup list with default delegate
	delegate := newItemDelegate(delegateKeys)
	noteList := list.New(items, delegate, 0, 0)
	noteList.Title = title
	noteList.Styles.Title = listTitleStyle
	noteList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
			listKeys.quit,
			listKeys.selectItem,
			listKeys.editItem,
			listKeys.completeItem,
		}
	}

	model.list = noteList
	return model
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Use full width and height
		m.list.SetSize(msg.Width, msg.Height)

	case ContextMenuMsg:
		if msg.Action == "open" {
			// Get the currently selected note
			currentIndex := m.list.Index()
			if currentIndex < len(m.notes) {
				m.selectedNote = m.notes[currentIndex]
				m.showMenu = true
			}
		} else if msg.Action == "update" {
			// Handle update action
			if m.selectedNote != nil {
				if m.selectedNote.Type == entities.NoteTypeText {
					// For text notes, open textarea directly
					return m, tea.Batch(
						tea.Cmd(func() tea.Msg {
							return OpenTextareaMsg{Note: m.selectedNote}
						}),
					)
				} else if m.selectedNote.Type == entities.NoteTypeTask {
					// For tasks, open task form
					return m, tea.Batch(
						tea.Cmd(func() tea.Msg {
							return OpenTaskFormMsg{Note: m.selectedNote}
						}),
					)
				} else if m.selectedNote.Type == entities.NoteTypeContact {
					// For contacts, open contact form
					return m, tea.Batch(
						tea.Cmd(func() tea.Msg {
							return OpenContactFormMsg{Note: m.selectedNote}
						}),
					)
				} else if m.selectedNote.Type == entities.NoteTypeReminder {
					// For reminders, open reminder form
					return m, tea.Batch(
						tea.Cmd(func() tea.Msg {
							return OpenReminderFormMsg{Note: m.selectedNote}
						}),
					)
				}
				// For other types, we'll implement later
			}
		} else if msg.Action == "complete" {
			// Handle complete task or reminder action
			if m.selectedNote != nil && m.selectedNote.Type == entities.NoteTypeTask {
				return m, tea.Batch(
					tea.Cmd(func() tea.Msg {
						return CompleteTaskMsg{Note: m.selectedNote}
					}),
				)
			}
			if m.selectedNote != nil && m.selectedNote.Type == entities.NoteTypeReminder {
				return m, tea.Batch(
					tea.Cmd(func() tea.Msg {
						return CompleteReminderMsg{Note: m.selectedNote}
					}),
				)
			}
		}

	case OpenTextareaMsg:
		// Handle opening textarea for editing
		if m.cli != nil {
			// Call the CLI method to update the note
			m.cli.updateNoteInteractive(msg.Note)
			return m, tea.Quit
		}

	case OpenContactFormMsg:
		// Handle opening contact form for editing
		if m.cli != nil {
			// Call the CLI method to update the contact
			m.cli.updateContactInteractive(msg.Note)
			return m, tea.Quit
		}

	case OpenTaskFormMsg:
		// Handle opening task form for editing
		if m.cli != nil {
			// Call the CLI method to update the task
			m.cli.updateTaskInteractive(msg.Note)
			return m, tea.Quit
		}

	case OpenReminderFormMsg:
		// Handle opening reminder form for editing
		if m.cli != nil {
			// Call the CLI method to update the reminder
			m.cli.updateReminderInteractive(msg.Note)
			return m, tea.Quit
		}

	case CompleteTaskMsg:
		// Handle completing a task
		if m.cli != nil {
			// Call the CLI method to complete the task
			m.cli.completeTaskInteractive(msg.Note)
			return m, tea.Quit
		}

	case CompleteReminderMsg:
		// Handle completing a reminder
		if m.cli != nil {
			// Call the CLI method to complete the reminder
			m.cli.completeReminderInteractive(msg.Note)
			return m, tea.Quit
		}

	case PreviewNoteMsg:
		// Show preview of the note
		m.selectedNote = msg.Note
		m.showPreview = true
		return m, nil

	case DeleteNoteMsg:
		// Handle deleting a note
		if m.cli != nil {
			m.cli.deleteNotePermanently(msg.Note)
			// Rebuild the list without the deleted note
			var filteredNotes []*entities.Note
			for _, note := range m.notes {
				if note.ID != msg.Note.ID {
					filteredNotes = append(filteredNotes, note)
				}
			}
			m.notes = filteredNotes

			// Rebuild list items
			items := make([]list.Item, len(filteredNotes))
			for i, note := range filteredNotes {
				items[i] = NoteItem{note: note}
			}

			newList := list.New(items, newItemDelegate(m.delegateKeys), m.list.Width(), m.list.Height())
			newList.Title = m.list.Title
			newList.Styles = m.list.Styles
			newList.AdditionalFullHelpKeys = m.list.AdditionalFullHelpKeys
			m.list = newList

			return m, nil
		}
		return m, nil

	case tea.KeyMsg:
		// Handle preview keys first
		if m.showPreview {
			if k := msg.String(); k == "ctrl+c" || k == "esc" || k == "q" {
				m.showPreview = false
				m.selectedNote = nil
				return m, nil
			}
			return m, nil
		}
		// Handle context menu keys first
		if m.showMenu {
			switch msg.String() {
			case "esc", "q":
				m.showMenu = false
				m.selectedNote = nil
				return m, nil
			case "1":
				// Update action
				m.showMenu = false
				return m, tea.Batch(
					tea.Cmd(func() tea.Msg {
						return ContextMenuMsg{Action: "update"}
					}),
				)
			case "2":
				// Complete action (for tasks and reminders)
				if m.selectedNote != nil && (m.selectedNote.Type == entities.NoteTypeTask || m.selectedNote.Type == entities.NoteTypeReminder) {
					m.showMenu = false
					return m, tea.Batch(
						tea.Cmd(func() tea.Msg {
							return ContextMenuMsg{Action: "complete"}
						}),
					)
				}
			}
			return m, nil
		}

		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.editItem):
			// Edit the currently selected item
			currentIndex := m.list.Index()
			if currentIndex < len(m.notes) {
				selectedNote := m.notes[currentIndex]
				// Handle different note types
				if selectedNote.Type == entities.NoteTypeText {
					m.cli.updateNoteInteractive(selectedNote)
					return m, tea.Quit
				} else if selectedNote.Type == entities.NoteTypeTask {
					m.cli.updateTaskInteractive(selectedNote)
					return m, tea.Quit
				} else if selectedNote.Type == entities.NoteTypeReminder {
					m.cli.updateReminderInteractive(selectedNote)
					return m, tea.Quit
				} else if selectedNote.Type == entities.NoteTypeContact {
					m.cli.updateContactInteractive(selectedNote)
					return m, tea.Quit
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.completeItem):
			// Complete the currently selected item (tasks and reminders)
			currentIndex := m.list.Index()
			if currentIndex < len(m.notes) {
				selectedNote := m.notes[currentIndex]
				if selectedNote.Type == entities.NoteTypeTask {
					m.cli.completeTaskInteractive(selectedNote)
					return m, tea.Quit
				} else if selectedNote.Type == entities.NoteTypeReminder {
					m.cli.completeReminderInteractive(selectedNote)
					return m, tea.Quit
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.selectItem):
			// Toggle selection of current item
			currentIndex := m.list.Index()
			if _, selected := m.selected[currentIndex]; selected {
				delete(m.selected, currentIndex)
			} else {
				m.selected[currentIndex] = struct{}{}
			}
			return m, nil

		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListModel) View() string {
	// If preview is open, show it
	if m.showPreview && m.selectedNote != nil {
		return m.renderPreview()
	}

	// If context menu is open, show it
	if m.showMenu && m.selectedNote != nil {
		return m.renderContextMenu()
	}

	// Create a prominent title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(1, 2).
		MarginBottom(1).
		Width(m.list.Width()).
		Align(lipgloss.Center).
		Render(m.title)

	return title + "\n" + m.list.View()
}

func (m ListModel) renderContextMenu() string {
	if m.selectedNote == nil {
		return ""
	}

	// Create menu title
	menuTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1).
		Width(m.list.Width()).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Options for: %s", m.selectedNote.Content))

	// Create menu options based on note type
	var options []string
	switch m.selectedNote.Type {
	case entities.NoteTypeTask:
		options = []string{
			"1. Update",
			"2. Complete Task",
		}
	case entities.NoteTypeReminder:
		options = []string{
			"1. Update",
			"2. Complete Reminder",
		}
	case entities.NoteTypeContact:
		options = []string{
			"1. Update",
		}
	default: // Text notes
		options = []string{
			"1. Update",
		}
	}

	// Style the options
	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		MarginTop(1)

	// Render options
	var optionsText strings.Builder
	for _, option := range options {
		optionsText.WriteString(optionStyle.Render(option))
		optionsText.WriteString("\n")
	}

	// Add help text
	helpText := helpStyle.Render("Press Esc to cancel")

	return menuTitle + "\n" + optionsText.String() + "\n" + helpText
}

func (m ListModel) renderPreview() string {
	if m.selectedNote == nil {
		return ""
	}

	note := m.selectedNote

	// Modal box style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(80).
		Align(lipgloss.Left)

	// Content
	var content strings.Builder
	content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Type:"), note.Type))
	content.WriteString(fmt.Sprintf("%s %s\n\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Content:"), note.Content))
	content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Created:"), note.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Updated:"), note.UpdatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Date:"), note.Date))

	// Add metadata based on note type
	switch note.Type {
	case entities.NoteTypeTask:
		if note.Metadata.Priority != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Priority:"), note.Metadata.Priority))
		}
		if note.Metadata.Status != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Status:"), note.Metadata.Status))
		}
		if note.Metadata.Assignee != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Assignee:"), note.Metadata.Assignee))
		}
	case entities.NoteTypeContact:
		if note.Metadata.Phone != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Phone:"), note.Metadata.Phone))
		}
		if note.Metadata.Email != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Email:"), note.Metadata.Email))
		}
	case entities.NoteTypeReminder:
		if note.Metadata.ReminderTime != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Time:"), note.Metadata.ReminderTime))
		}
		if note.Metadata.Status != "" {
			content.WriteString(fmt.Sprintf("%s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Status:"), note.Metadata.Status))
		}
	}

	// Modal content
	modalContent := modalStyle.Render(content.String())

	// Center the modal on screen
	return lipgloss.Place(80, 0, lipgloss.Center, lipgloss.Center, modalContent) + "\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("  Press Esc or Q to close")
}

func (m ListModel) generateMarkdownForNote(note *entities.Note) string {
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# %s\n\n", note.Content))
	md.WriteString(fmt.Sprintf("**Type:** %s\n\n", note.Type))
	md.WriteString(fmt.Sprintf("**Created:** %s\n\n", note.CreatedAt.Format("2006-01-02 15:04:05")))
	md.WriteString(fmt.Sprintf("**Last Updated:** %s\n\n", note.UpdatedAt.Format("2006-01-02 15:04:05")))
	md.WriteString(fmt.Sprintf("**Date:** %s\n\n", note.Date))

	// Show metadata based on note type
	switch note.Type {
	case entities.NoteTypeTask:
		if note.Metadata.Priority != "" {
			md.WriteString(fmt.Sprintf("**Priority:** %s\n\n", note.Metadata.Priority))
		}
		if note.Metadata.Status != "" {
			md.WriteString(fmt.Sprintf("**Status:** %s\n\n", note.Metadata.Status))
		}
		if note.Metadata.Assignee != "" {
			md.WriteString(fmt.Sprintf("**Assignee:** %s\n\n", note.Metadata.Assignee))
		}
		if note.Metadata.DueDate != nil {
			md.WriteString(fmt.Sprintf("**Due Date:** %s\n\n", note.Metadata.DueDate.Format("2006-01-02")))
		}
	case entities.NoteTypeContact:
		if note.Metadata.Phone != "" {
			md.WriteString(fmt.Sprintf("**Phone:** %s\n\n", note.Metadata.Phone))
		}
		if note.Metadata.Email != "" {
			md.WriteString(fmt.Sprintf("**Email:** %s\n\n", note.Metadata.Email))
		}
		if note.Metadata.Address != "" {
			md.WriteString(fmt.Sprintf("**Address:** %s\n\n", note.Metadata.Address))
		}
	case entities.NoteTypeReminder:
		if note.Metadata.ReminderTime != "" {
			md.WriteString(fmt.Sprintf("**Time:** %s\n\n", note.Metadata.ReminderTime))
		}
		if note.Metadata.Status != "" {
			md.WriteString(fmt.Sprintf("**Status:** %s\n\n", note.Metadata.Status))
		}
	}

	return md.String()
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				// Show preview for all note types
				if i, ok := m.SelectedItem().(NoteItem); ok {
					return tea.Batch(
						m.NewStatusMessage(statusMessageStyle("Showing preview...")),
						tea.Cmd(func() tea.Msg {
							return PreviewNoteMsg{Note: i.note}
						}),
					)
				}

			case key.Matches(msg, keys.remove):
				if i, ok := m.SelectedItem().(NoteItem); ok {
					return tea.Batch(
						m.NewStatusMessage(statusMessageStyle("Deleting...")),
						tea.Cmd(func() tea.Msg {
							return DeleteNoteMsg{Note: i.note}
						}),
					)
				}
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
