package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

func parseFile(path string) ([]Floor, [][]Splitter, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file %q: %s", path, err.Error())
		}
	}()

	var (
		floors   []Floor
		dividers [][]Splitter
	)
	sheetName := f.GetSheetList()[0]

	rows, err := f.GetRows(sheetName, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sheet %q: %w", sheetName, err)
	}

	var offset = 2

	for i := 3; ; i++ {
		if len(rows[i]) == 0 {
			break
		}

		var numbers []int
		for j := 0; j < len(rows[i]); j++ {
			number, err := strconv.Atoi(strings.TrimSpace(rows[i][j]))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse digits in cell %d:%d : %w", i, j, err)
			}

			numbers = append(numbers, number)
		}

		offset++
		floors = append(floors, Floor{
			Number: numbers[0],
			Risers: lo.Map(numbers[3:], func(number int, _ int) Riser {
				return Riser{FlatNumber: number}
			}),
			Flats: FlatRange{
				FlatStart: numbers[1],
				FlatEnd:   numbers[2],
			},
		})
	}

	dividers = make([][]Splitter, (len(rows[offset+5])+1)/2)

	for i := offset + 5; i < len(rows); i++ {
		if len(rows[i]) == 0 || rows[i][0] == "" {
			continue
		}

		for j := 0; j < len(rows[i]); j += 2 {
			portNumber, err := strconv.Atoi(strings.TrimSpace(rows[i][j]))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse cell %d:%d: %w", i, j, err)
			}

			dividers[j/2] = append(dividers[j/2], Splitter{PortNumber: portNumber})
		}
	}

	return floors, dividers, nil
}

func writeResults(path string, splitters [][]Splitter) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file %q: %s", path, err.Error())
		}
	}()

	sheetName := f.GetSheetList()[0]

	baseColumn := 'B'
	baseRow := 1

	for {
		value, err := f.GetCellValue(sheetName, string(baseColumn)+strconv.Itoa(baseRow))
		if err != nil {
			return err
		}

		if value == "Кв-ры" {
			break
		}

		baseRow++
	}
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: false,
		},
	})
	if err != nil {
		return err
	}
	for i, riser := range splitters {
		column := string(baseColumn + 2*rune(i))
		var width int

		for j, splitter := range riser {
			formatted := strings.Join(lo.Map(splitter.Flats, func(item FlatRange, index int) string { return item.String() }), ",")
			cell := column + strconv.Itoa(j+baseRow+1)

			if len([]rune(formatted)) > width {
				width = len([]rune(formatted))
			}

			err = f.SetCellValue(sheetName, cell, formatted)
			if err != nil {
				return err
			}
		}

		err = f.SetColStyle(sheetName, column, style)
		if err != nil {
			return err
		}

		err = f.SetColWidth(sheetName, column, column, float64(width))
		if err != nil {
			return err
		}

	}

	return f.Save()
}
