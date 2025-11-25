package main

import (
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
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
		cmd = tea.Quit
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

func displayDividers(splitters [][]Splitter) error {
	table := tablewriter.NewWriter(os.Stdout)

	for j := range len(splitters[0]) {
		var row []any
		for _, riser := range splitters {
			flats := strings.Join(lo.Map(riser[j].Flats, func(item FlatRange, _ int) string { return item.String() }), ",")

			row = append(row, randomColor(j, flats))
		}

		err := table.Append(row...)
		if err != nil {
			return err
		}

	}

	return table.Render()
}

func displayFloors(floors []Floor, splitters [][]Splitter) error {
	table := tablewriter.NewWriter(os.Stdout)

	flatToDivider := make(map[int]int)
	for _, riser := range splitters {
		for j, splitter := range riser {
			for _, flatRange := range splitter.Flats {
				for i := flatRange.FlatStart; i <= flatRange.FlatEnd; i++ {
					flatToDivider[i] = j
				}
			}
		}
	}

	for floor := range slices.Values(floors) {
		var row []any

		row = append(row, floor.Number)
		var j int
		var number int

		var flats []string

		for i := floor.Flats.FlatStart; i <= floor.Flats.FlatEnd; i++ {
			if i-floor.Flats.FlatStart-number == floor.Risers[j].FlatNumber {
				number += floor.Risers[j].FlatNumber
				j++

				row = append(row, strings.Join(flats, " "))
				flats = nil
			}

			flats = append(flats, randomColor(flatToDivider[i], i))
		}
		row = append(row, strings.Join(flats, " "))

		err := table.Append(row...)
		if err != nil {
			return err
		}
	}

	return table.Render()
}

func randomColor(seed int, args any) string {
	return color.New(color.Bold, color.Attribute(31+seed%7)).Sprint(args)
}
