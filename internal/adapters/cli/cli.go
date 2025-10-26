package cli

import (
	"encoding/json"
	"fmt"
	"jotterxpress/internal/adapters/repository"
	"jotterxpress/internal/application/services"
	"jotterxpress/internal/domain/entities"
	"jotterxpress/internal/domain/ports"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Styles for beautiful CLI output
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B5CF6")).
			Bold(true)
)

// CLI represents the command line interface
type CLI struct {
	noteService ports.NoteService
	repository  ports.NoteRepository
}

// NewCLI creates a new CLI instance
func NewCLI() *CLI {
	// Get home directory for notes storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	notesDir := filepath.Join(homeDir, ".jotterxpress", "notes")

	// Create repository and service
	noteRepo := repository.NewFileRepository(notesDir)
	noteService := services.NewNoteService(noteRepo)

	return &CLI{
		noteService: noteService,
		repository:  noteRepo,
	}
}

// SetupCommands sets up all CLI commands
func (cli *CLI) SetupCommands() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "jtx [note content]",
		Short: "JotterXpress - Quick note taking CLI tool",
		Long:  "A fast and simple CLI tool for taking notes.",
		Args:  cobra.ArbitraryArgs,
		Run:   cli.addNote,
	}

	// Add flags for all commands
	var listFlag, noteFlag, taskFlag, contactFlag, reminderFlag, interactiveFlag bool
	var listDateStr, listMonthStr string
	rootCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List today's notes")
	rootCmd.Flags().StringVar(&listDateStr, "list-date", "", "List notes for a specific date (format: YYYY-MM-DD)")
	rootCmd.Flags().StringVar(&listMonthStr, "list-month", "", "List notes for a specific month (format: MM)")
	rootCmd.Flags().BoolVarP(&noteFlag, "note", "n", false, "Open interactive textarea to create a note")
	rootCmd.Flags().BoolVarP(&taskFlag, "task", "t", false, "Create a new task (interactive mode)")
	rootCmd.Flags().BoolVarP(&contactFlag, "contact", "c", false, "Create a new contact (interactive mode)")
	rootCmd.Flags().BoolVarP(&reminderFlag, "reminder", "r", false, "Create a new reminder (interactive mode)")
	rootCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Open interactive list view")

	// Override the Run function to handle flags
	rootCmd.Run = cli.handleRootCommand

	return rootCmd
}

// handleRootCommand handles the root command with flags
func (cli *CLI) handleRootCommand(cmd *cobra.Command, args []string) {
	// Check flags first
	listFlag, _ := cmd.Flags().GetBool("list")
	listDateStr, _ := cmd.Flags().GetString("list-date")
	listMonthStr, _ := cmd.Flags().GetString("list-month")
	noteFlag, _ := cmd.Flags().GetBool("note")
	taskFlag, _ := cmd.Flags().GetBool("task")
	contactFlag, _ := cmd.Flags().GetBool("contact")
	reminderFlag, _ := cmd.Flags().GetBool("reminder")
	interactiveFlag, _ := cmd.Flags().GetBool("interactive")

	// Count how many flags are set
	flagCount := 0
	if listDateStr != "" {
		flagCount++
	}
	if listMonthStr != "" {
		flagCount++
	}
	if listFlag {
		flagCount++
	}
	if noteFlag {
		flagCount++
	}
	if taskFlag {
		flagCount++
	}
	if contactFlag {
		flagCount++
	}
	if reminderFlag {
		flagCount++
	}
	if interactiveFlag {
		flagCount++
	}

	// If multiple flags are set, show error
	if flagCount > 1 {
		fmt.Println(errorStyle.Render("Error: Only one command flag can be used at a time"))
		cmd.Help()
		return
	}

	// Handle each flag
	if listDateStr != "" {
		cli.listNotesByDate(cmd, []string{listDateStr})
		return
	}
	if listMonthStr != "" {
		cli.listNotesByMonth(cmd, listMonthStr)
		return
	}
	if listFlag {
		cli.listTodayNotes(cmd, args)
		return
	}
	if noteFlag {
		cli.createNoteInteractive()
		return
	}
	if taskFlag {
		cli.createTaskInteractive()
		return
	}
	if contactFlag {
		cli.createContactInteractive()
		return
	}
	if reminderFlag {
		cli.createReminderInteractive()
		return
	}
	if interactiveFlag {
		cli.openInteractiveList(cmd, args)
		return
	}

	// If no flags are set, show help in a nice modal
	if len(args) == 0 {
		cli.showHelp()
		return
	}

	// Require content to have at least 2 words or be explicitly quoted
	content := strings.Join(args, " ")
	if !strings.Contains(content, " ") {
		// Single word without spaces - reject as it's likely a command
		fmt.Println(errorStyle.Render(fmt.Sprintf("'%s' is not a valid command", content)))
		fmt.Println(infoStyle.Render("Use quotes for note content: jtx \"your note content\""))
		fmt.Println(infoStyle.Render("Or use flags: jtx -l (list), jtx -n (note), jtx -t (task), jtx -r (reminder), etc."))
		return
	}

	// Create note with the provided content
	cli.addNote(cmd, args)
}

// addNote adds a new note
func (cli *CLI) addNote(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		// Show help when no arguments provided
		cmd.Help()
		return
	}

	content := strings.Join(args, " ")

	_, err := cli.noteService.CreateNote(content)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating note: %v", err)))
		os.Exit(1)
	}

	fmt.Println(successStyle.Render("Note saved successfully!"))
}

// createContactInteractive creates a contact using the interactive form
func (cli *CLI) createContactInteractive() {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive contact form
	model := NewContactFormModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive contact form: %v", err)))
		os.Exit(1)
	}

	// Get the created contact
	var contact *entities.Note
	if contactModel, ok := finalModel.(*ContactFormModel); ok {
		contact = contactModel.GetContact()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if contact == nil {
		// User cancelled, exit silently
		os.Exit(0)
	}

	// Save the contact
	if err := cli.noteService.SaveNote(contact); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error saving contact: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Contact created successfully!"))

	// Exit successfully
	os.Exit(0)
}

// listTodayNotes lists all notes for today
func (cli *CLI) listTodayNotes(cmd *cobra.Command, args []string) {
	notes, err := cli.noteService.GetTodayNotes()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error retrieving notes: %v", err)))
		os.Exit(1)
	}

	if len(notes) == 0 {
		fmt.Println(infoStyle.Render("No notes found for today."))
		fmt.Println(infoStyle.Render("Start taking notes with: jx \"your note content\""))
		return
	}

	// Check if we're in a TTY environment
	if cli.isTTY() {
		// Create and run the interactive list
		model := NewListModel(notes, "Today's Notes", cli)
		program := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := program.Run(); err != nil {
			// Fall back to text mode if interactive fails
			cli.showTextList(notes, "Today's Notes")
		}
	} else {
		// Use text mode for non-TTY environments
		cli.showTextList(notes, "Today's Notes")
	}
}

// listNotesByDate lists notes for a specific date
func (cli *CLI) listNotesByDate(cmd *cobra.Command, args []string) {
	date := args[0]

	notes, err := cli.noteService.GetNotesByDate(date)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error retrieving notes: %v", err)))
		os.Exit(1)
	}

	if len(notes) == 0 {
		fmt.Println(infoStyle.Render(fmt.Sprintf("No notes found for %s.", date)))
		return
	}

	// Check if we're in a TTY environment
	if cli.isTTY() {
		// Create and run the interactive list
		title := fmt.Sprintf("Notes for %s", date)
		model := NewListModel(notes, title, cli)
		program := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := program.Run(); err != nil {
			// Fall back to text mode if interactive fails
			cli.showTextList(notes, title)
		}
	} else {
		// Use text mode for non-TTY environments
		title := fmt.Sprintf("Notes for %s", date)
		cli.showTextList(notes, title)
	}
}

// createTask creates a new task
func (cli *CLI) createTask(cmd *cobra.Command, args []string) {
	// Check if online mode is requested
	online, _ := cmd.Flags().GetBool("online")

	if !online {
		// Interactive mode is the default
		cli.createTaskInteractive()
		return
	}

	// Online mode (single line command)
	if len(args) == 0 {
		fmt.Println(errorStyle.Render("Error: Please provide task content"))
		fmt.Println(infoStyle.Render("Usage: jtx task \"your task content\" --priority high --online"))
		fmt.Println(infoStyle.Render("Or use: jtx task (interactive mode by default)"))
		os.Exit(1)
	}

	content := strings.Join(args, " ")

	// Parse priority
	var priority entities.Priority
	switch cmd.Flag("priority").Value.String() {
	case "low":
		priority = entities.PriorityLow
	case "high":
		priority = entities.PriorityHigh
	default:
		fmt.Println(errorStyle.Render("Error: Invalid priority. Use: low or high"))
		os.Exit(1)
	}

	task := entities.NewTask(content, priority)

	if err := cli.noteService.SaveNote(task); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating task: %v", err)))
		os.Exit(1)
	}

	fmt.Println(successStyle.Render("Task created successfully!"))
}

// createContact creates a new contact
func (cli *CLI) createContact(cmd *cobra.Command, args []string) {
	// Check if online mode is requested
	online, _ := cmd.Flags().GetBool("online")

	if !online {
		// Interactive mode is the default
		cli.createContactInteractive()
		return
	}

	// Online mode (single line command)
	if len(args) == 0 {
		fmt.Println(errorStyle.Render("Error: Please provide contact name"))
		fmt.Println(infoStyle.Render("Usage: jtx contact \"John Doe\" --phone +1234567890 --email john@example.com --online"))
		fmt.Println(infoStyle.Render("Or use: jtx contact (interactive mode by default)"))
		os.Exit(1)
	}

	name := strings.Join(args, " ")
	phone := cmd.Flag("phone").Value.String()
	email := cmd.Flag("email").Value.String()

	if phone == "" && email == "" {
		fmt.Println(errorStyle.Render("Error: Please provide at least phone or email"))
		fmt.Println(infoStyle.Render("Usage: jtx contact \"John Doe\" --phone +1234567890 --email john@example.com"))
		os.Exit(1)
	}

	contact := entities.NewContact(name, phone, email)

	if err := cli.noteService.SaveNote(contact); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error creating contact: %v", err)))
		os.Exit(1)
	}

	fmt.Println(successStyle.Render("Contact created successfully!"))
}

// openInteractiveList opens the interactive list view
func (cli *CLI) openInteractiveList(cmd *cobra.Command, args []string) {
	notes, err := cli.noteService.GetTodayNotes()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error retrieving notes: %v", err)))
		os.Exit(1)
	}

	if len(notes) == 0 {
		fmt.Println(infoStyle.Render("No notes found for today."))
		fmt.Println(infoStyle.Render("Start taking notes with: jx \"your note content\""))
		return
	}

	// Check if we're in a TTY environment
	if cli.isTTY() {
		// Create and run the interactive list
		model := NewListModel(notes, "JotterXpress - Interactive Notes", cli)
		program := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := program.Run(); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive list: %v", err)))
			os.Exit(1)
		}
	} else {
		// Use text mode for non-TTY environments
		cli.showTextList(notes, "JotterXpress - Interactive Notes")
	}
}

// isTTY checks if we're in a TTY environment
func (cli *CLI) isTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// showTextList shows notes in text mode
func (cli *CLI) showTextList(notes []*entities.Note, title string) {
	fmt.Println(titleStyle.Render(title))
	fmt.Println("")

	result := cli.noteService.ListNotes(notes)
	fmt.Println(result)
}

// createTaskInteractive creates a task using the interactive form
func (cli *CLI) createTaskInteractive() {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive form
	model := NewTaskFormModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive form: %v", err)))
		os.Exit(1)
	}

	// Get the created task directly from the model
	var task *entities.Note
	if taskModel, ok := finalModel.(*TaskFormModel); ok {
		task = taskModel.GetTask()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if task == nil {
		fmt.Println(errorStyle.Render("Task creation was cancelled"))
		os.Exit(1)
	}

	// Save the task
	if err := cli.noteService.SaveNote(task); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error saving task: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Task created successfully!"))

	// Exit successfully
	os.Exit(0)
}

// createNoteInteractive creates a note using the interactive textarea
func (cli *CLI) createNoteInteractive() {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive textarea
	model := NewNoteTextareaModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive textarea: %v", err)))
		os.Exit(1)
	}

	// Get the created note
	var note *entities.Note
	if noteModel, ok := finalModel.(*NoteTextareaModel); ok {
		if noteModel.IsCancelled() {
			// User cancelled, exit silently
			os.Exit(0)
		}
		note = noteModel.GetNote()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if note == nil {
		// User cancelled, exit silently
		os.Exit(0)
	}

	// Save the note
	if err := cli.noteService.SaveNote(note); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error saving note: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Note created successfully!"))

	// Exit successfully
	os.Exit(0)
}

// updateNoteInteractive updates a note using the interactive textarea
func (cli *CLI) updateNoteInteractive(note *entities.Note) {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive textarea with existing content
	model := NewNoteTextareaModelWithContent(note)
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive textarea: %v", err)))
		os.Exit(1)
	}

	// Get the updated note
	var updatedNote *entities.Note
	if noteModel, ok := finalModel.(*NoteTextareaModel); ok {
		if noteModel.IsCancelled() {
			// User cancelled, exit silently
			os.Exit(0)
		}
		updatedNote = noteModel.GetNote()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if updatedNote == nil {
		// User cancelled, exit silently
		os.Exit(0)
	}

	// Update the existing note with new content
	note.Content = updatedNote.Content
	note.UpdatedAt = updatedNote.UpdatedAt

	// Save the updated note
	if err := cli.noteService.SaveNote(note); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error updating note: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Note updated successfully!"))

	// Exit successfully
	os.Exit(0)
}

// updateContactInteractive updates a contact using the interactive form
func (cli *CLI) updateContactInteractive(contact *entities.Note) {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive contact form with existing data
	model := NewContactFormModelWithData(contact)
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive contact form: %v", err)))
		os.Exit(1)
	}

	// Get the updated contact
	var updatedContact *entities.Note
	if contactModel, ok := finalModel.(*ContactFormModel); ok {
		updatedContact = contactModel.GetContact()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if updatedContact == nil {
		fmt.Println(errorStyle.Render("Contact update was cancelled"))
		os.Exit(1)
	}

	// Update the existing contact with new data
	contact.Content = updatedContact.Content
	contact.Metadata.Phone = updatedContact.Metadata.Phone
	contact.Metadata.Email = updatedContact.Metadata.Email
	contact.UpdatedAt = updatedContact.UpdatedAt

	// Save the updated contact
	if err := cli.noteService.SaveNote(contact); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error updating contact: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Contact updated successfully!"))

	// Exit successfully
	os.Exit(0)
}

// completeTaskInteractive completes a task by changing its status
func (cli *CLI) completeTaskInteractive(task *entities.Note) {
	// Update the task status to completed
	task.Metadata.Status = entities.StatusCompleted
	task.UpdatedAt = time.Now()

	// Save the updated task
	if err := cli.noteService.SaveNote(task); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error completing task: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Task completed successfully!"))

	// Exit successfully
	os.Exit(0)
}

// updateTaskInteractive updates a task using the interactive form
func (cli *CLI) updateTaskInteractive(task *entities.Note) {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive task form with existing data
	model := NewTaskFormModelWithData(task)
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive task form: %v", err)))
		os.Exit(1)
	}

	// Get the updated task
	var updatedTask *entities.Note
	if taskModel, ok := finalModel.(*TaskFormModel); ok {
		updatedTask = taskModel.GetTask()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if updatedTask == nil {
		fmt.Println(errorStyle.Render("Task update was cancelled"))
		os.Exit(1)
	}

	// Update the existing task with new data
	task.Content = updatedTask.Content
	task.Metadata.Priority = updatedTask.Metadata.Priority
	task.Metadata.Assignee = updatedTask.Metadata.Assignee
	task.UpdatedAt = updatedTask.UpdatedAt

	// Save the updated task
	if err := cli.noteService.SaveNote(task); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error updating task: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Task updated successfully!"))

	// Exit successfully
	os.Exit(0)
}

// createReminderInteractive creates a reminder using the interactive form
func (cli *CLI) createReminderInteractive() {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive form
	model := NewReminderFormModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive form: %v", err)))
		os.Exit(1)
	}

	// Get the created reminder directly from the model
	var reminder *entities.Note
	if reminderModel, ok := finalModel.(*ReminderFormModel); ok {
		reminder = reminderModel.GetReminder()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if reminder == nil {
		// User cancelled, exit silently
		os.Exit(0)
	}

	// Save the reminder
	if err := cli.noteService.SaveNote(reminder); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error saving reminder: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Reminder created successfully!"))

	// Exit successfully
	os.Exit(0)
}

// updateReminderInteractive updates a reminder using the interactive form
func (cli *CLI) updateReminderInteractive(reminder *entities.Note) {
	// Check if we're in a TTY environment
	if !cli.isTTY() {
		fmt.Println(errorStyle.Render("Interactive mode requires a TTY environment"))
		os.Exit(1)
	}

	// Create and run the interactive reminder form with existing data
	model := NewReminderFormModelWithData(reminder)
	program := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error running interactive reminder form: %v", err)))
		os.Exit(1)
	}

	// Get the updated reminder
	var updatedReminder *entities.Note
	if reminderModel, ok := finalModel.(*ReminderFormModel); ok {
		updatedReminder = reminderModel.GetReminder()
	} else {
		fmt.Println(errorStyle.Render("Invalid model type returned"))
		os.Exit(1)
	}

	if updatedReminder == nil {
		// User cancelled, exit silently
		os.Exit(0)
	}

	// Update the existing reminder with new data
	reminder.Content = updatedReminder.Content
	reminder.Metadata.ReminderTime = updatedReminder.Metadata.ReminderTime
	reminder.Metadata.Status = updatedReminder.Metadata.Status
	reminder.UpdatedAt = updatedReminder.UpdatedAt

	// Save the updated reminder
	if err := cli.noteService.SaveNote(reminder); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error updating reminder: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Reminder updated successfully!"))

	// Exit successfully
	os.Exit(0)
}

// completeReminderInteractive completes a reminder by changing its status
func (cli *CLI) completeReminderInteractive(reminder *entities.Note) {
	// Update the reminder status to completed
	reminder.Metadata.Status = entities.StatusCompleted
	reminder.UpdatedAt = time.Now()

	// Save the updated reminder
	if err := cli.noteService.SaveNote(reminder); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error completing reminder: %v", err)))
		os.Exit(1)
	}

	// Show success message and exit
	fmt.Println(successStyle.Render("Reminder completed successfully!"))

	// Exit successfully
	os.Exit(0)
}

// deleteNotePermanently deletes a note from the file permanently
func (cli *CLI) deleteNotePermanently(note *entities.Note) {
	// Get notes directory from the repository
	notesDir := cli.repository.(interface{ GetNotesDir() string }).GetNotesDir()
	filepath := filepath.Join(notesDir, fmt.Sprintf("%s.json", note.Date))

	file, err := os.Open(filepath)
	if err != nil {
		return // File doesn't exist or can't be opened
	}
	defer file.Close()

	var notes []*entities.Note
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&notes); err != nil {
		return // Can't decode
	}

	// Filter out the note to delete
	var filteredNotes []*entities.Note
	for _, n := range notes {
		if n.ID != note.ID {
			filteredNotes = append(filteredNotes, n)
		}
	}

	// Write back the filtered notes
	writeFile, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer writeFile.Close()

	encoder := json.NewEncoder(writeFile)
	encoder.SetIndent("", "  ")
	encoder.Encode(filteredNotes)
}

// listNotesByMonth lists notes for a specific month
func (cli *CLI) listNotesByMonth(cmd *cobra.Command, monthStr string) {
	// Validate and parse month
	monthInt := 0
	_, err := fmt.Sscanf(monthStr, "%d", &monthInt)
	if err != nil || monthInt < 1 || monthInt > 12 {
		fmt.Println(errorStyle.Render("Error: Invalid month. Month must be between 1 and 12"))
		fmt.Println(infoStyle.Render("Usage: jtx --list-month 10"))
		os.Exit(1)
	}

	// Get current year
	currentYear := time.Now().Year()
	monthStrFormatted := fmt.Sprintf("%d-%02d", currentYear, monthInt)

	notes, err := cli.noteService.GetNotesByMonth(monthStrFormatted)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error retrieving notes: %v", err)))
		os.Exit(1)
	}

	if len(notes) == 0 {
		fmt.Println(infoStyle.Render(fmt.Sprintf("No notes found for %s/%d.", monthStr, currentYear)))
		return
	}

	// Check if we're in a TTY environment
	if cli.isTTY() {
		// Create and run the interactive list
		title := fmt.Sprintf("Notes for %s/%d", monthStr, currentYear)
		model := NewListModel(notes, title, cli)
		program := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := program.Run(); err != nil {
			// Fall back to text mode if interactive fails
			cli.showTextList(notes, title)
		}
	} else {
		// Use text mode for non-TTY environments
		title := fmt.Sprintf("Notes for %s/%d", monthStr, currentYear)
		cli.showTextList(notes, title)
	}
}

// showHelp shows a nice formatted help message
func (cli *CLI) showHelp() {
	// Modal box style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(70).
		Align(lipgloss.Left)

	// Content
	var content strings.Builder
	content.WriteString(fmt.Sprintf("%s\n\n", titleStyle.Render("JotterXpress - Quick Note Taking CLI")))

	content.WriteString(fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Usage:")))
	content.WriteString("  jtx \"your note\"              Quick note\n")
	content.WriteString("  jtx --list                   List today's notes\n")
	content.WriteString("  jtx --list-date YYYY-MM-DD  List notes for specific date\n")
	content.WriteString("  jtx --list-month MM          List notes for specific month\n\n")

	content.WriteString(fmt.Sprintf("%s\n\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Interactive Commands:")))
	content.WriteString("  jtx --note (-n)              Create note\n")
	content.WriteString("  jtx --task (-t)              Create task\n")
	content.WriteString("  jtx --reminder (-r)         Create reminder\n")
	content.WriteString("  jtx --contact (-c)          Create contact\n")
	content.WriteString("  jtx --interactive (-i)       Open interactive view\n\n")

	content.WriteString(fmt.Sprintf("%s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Flags:")))
	content.WriteString("  -l, --list                   List today's notes\n")
	content.WriteString("      --list-date              List by date\n")
	content.WriteString("      --list-month             List by month\n")
	content.WriteString("  -n, --note                   Create note\n")
	content.WriteString("  -t, --task                   Create task\n")
	content.WriteString("  -r, --reminder              Create reminder\n")
	content.WriteString("  -c, --contact                Create contact\n")
	content.WriteString("  -i, --interactive            Interactive list view\n\n")

	content.WriteString(fmt.Sprintf("%s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8B5CF6")).Render("Interactive Shortcuts:")))
	content.WriteString("  Enter  Preview note\n")
	content.WriteString("  e      Edit note\n")
	content.WriteString("  c      Complete (tasks/reminders)\n")
	content.WriteString("  x      Delete note\n")
	content.WriteString("  q      Quit\n")

	// Modal content
	modalContent := modalStyle.Render(content.String())

	// Center the modal on screen
	fmt.Println(lipgloss.Place(70, 0, lipgloss.Center, lipgloss.Center, modalContent))
}
