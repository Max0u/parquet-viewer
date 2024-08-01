package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(50).
			Height(10).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(100).
				Height(10).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	selectedFile     string
	windows          []*Window
	focusWindow      *Window
	filePickerWindow Window
	tableWindow      Window
}

func newModel() mainModel {
	m := mainModel{}

	m.selectedFile = "/Users/maxime/perso/vibe-parquet-viewer/airflow_2024072902_chunk_13_0.parquet"

	m.filePickerWindow = Window{
		x:     0,
		y:     0,
		model: filepicker.New(),
	}
	if fp, ok := m.filePickerWindow.model.(filepicker.Model); ok {
		fp.AllowedTypes = []string{""}
		fp.ShowPermissions = false
		fp.ShowHidden = false
		fp.CurrentDirectory, _ = os.UserHomeDir()
	}

	// Initialize with empty table
	m.tableWindow = Window{
		x: 1,
		y: 0,
		model: table.New(
			table.WithColumns([]table.Column{}),
			table.WithRows([]table.Row{}),
			table.WithFocused(true),
			table.WithHeight(10),
		),
	}

	m.windows = []*Window{&m.filePickerWindow, &m.tableWindow}
	m.focusWindow = &m.filePickerWindow

	return m
}

func (m mainModel) Init() tea.Cmd {

	var fp filepicker.Model
	var ok bool

	if fp, ok = m.filePickerWindow.model.(filepicker.Model); !ok {
		log.Println("model in filePickerWindow1 is not of type filepicker.Model")
	}

	if ok {
		return tea.Batch(fp.Init())
	}

	return nil // Or an appropriate default `tea.Cmd`
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "h":
			m.focusWindow = FindNextWindow(m.windows, m.focusWindow, Left)
		case "j":
			m.focusWindow = FindNextWindow(m.windows, m.focusWindow, Down)
		case "k":
			m.focusWindow = FindNextWindow(m.windows, m.focusWindow, Up)
		case "l":
			m.focusWindow = FindNextWindow(m.windows, m.focusWindow, Right)
		}
	}

	switch model := m.focusWindow.model.(type) {
	case filepicker.Model:
		// Update the file picker model with the incoming message
		model, cmd := model.Update(msg)

		m.filePickerWindow.model = model

		// Did the user select a file?
		if didSelect, path := model.DidSelectFile(msg); didSelect {
			// Get the path of the selected file.
			m.selectedFile = path
		}
		cmds = append(cmds, cmd)
	case table.Model:
		model, cmd = model.Update(msg)
		m.tableWindow.model = model
		cmds = append(cmds, cmd)
	default:
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {

	yGroups := make(map[int][]string)

	// Group windows by y coordinate for horizontal joining
	for _, win := range m.windows {
		yGroups[win.y] = append(yGroups[win.y], modelStyle.Render(win.model.View()))
	}

	// Sort yGroups into slices for ordered joining
	var yKeys []int
	for y := range yGroups {
		yKeys = append(yKeys, y)
	}

	// Join models horizontally for each y group
	var horizontalGroups []string
	for _, y := range yKeys {
		horizontalGroup := lipgloss.JoinHorizontal(lipgloss.Top, yGroups[y]...)
		horizontalGroups = append(horizontalGroups, horizontalGroup)
	}

	// Join horizontally joined groups vertically
	result := lipgloss.JoinVertical(lipgloss.Top, horizontalGroups...)

	// Add help text
	helpText := helpStyle.Render("\ntab: focus next • n: new model • q: exit\n")
	result += helpText

	return result

}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
