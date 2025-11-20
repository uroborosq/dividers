package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filepicker.Model
	selectedFile string
	quitting     bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	if didSelect, path := m.DidSelectFile(msg); didSelect {
		m.selectedFile = path
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("Pick a file or press Q to quit")
	s.WriteString("\n\n" + m.Model.View() + "\n")
	return s.String()
}

func chooseFile(path string) (string, bool) {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".xlsx"}
	fp.ShowHidden = false
	fp.CurrentDirectory = path

	m := model{Model: fp}

	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)

	return mm.selectedFile, mm.selectedFile == ""
}
