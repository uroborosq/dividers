package main

import (
	"io"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
)

func displayDividers(w io.Writer, splitters [][]Splitter) error {
	table := tablewriter.NewWriter(w)

	for j := range len(splitters[0]) {
		var row = []any{j}
		for _, riser := range splitters {
			flats := strings.Join(lo.Map(riser[j].Flats, func(item FlatRange, _ int) string { return item.String() }), ",")

			row = append(row, riser[j].PortNumber, randomColor(j, flats))
		}

		err := table.Append(row...)
		if err != nil {
			return err
		}

	}

	return table.Render()
}

func displayFloors(w io.Writer, floors []Floor, splitters [][]Splitter) error {
	table := tablewriter.NewWriter(w)

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
