package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

// func (m model) Init() tea.Cmd {
// 	return m.filepicker.Init()
// }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) { case tea.KeyMsg:
	// 	switch msg.String() {
	// 	case "ctrl+c", "q":
	// 		m.quitting = true
	// 		return m, tea.Quit
	// 	}
	// case clearErrorMsg:
	// 	m.err = nil
	// }
	//
	// var cmd tea.Cmd
	// m.Model, cmd = m.Model.Update(msg)
	//
	// // Did the user select a file?
	// if didSelect, path := m.Model.DidSelectFile(msg); didSelect {
	// 	// Get the path of the selected file.
	// 	m.selectedFile = path
	// }

	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	return m, cmd
}

// func (m model) View() string {
// 	if m.quitting {
// 		return ""
// 	}
// 	var s strings.Builder
// 	s.WriteString("\n  ")
// 	if m.err != nil {
// 		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
// 	} else if m.selectedFile == "" {
// 		s.WriteString("Pick a file:")
// 	} else {
// 		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
// 	}
// 	s.WriteString("\n\n" + m.filepicker.View() + "\n")
// 	return s.String()
// }

func main() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".xlsx"}
	fp.ShowHidden = false
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := model{Model: fp}

	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)
	if mm.selectedFile != "" {
		fmt.Println("\n  You selected: " + mm.selectedFile + "\n")
	} else {
		fmt.Println("Not selected anything")
	}
}
