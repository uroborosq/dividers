package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

func parseFile(path string) ([]Floor, [][]Divider, error) {
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
		dividers [][]Divider
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

	dividers = make([][]Divider, (len(rows[offset+5])+1)/2)

	for i := offset + 5; i < len(rows); i++ {
		for j := 0; j < len(rows[i]); j += 2 {
			portNumber, err := strconv.Atoi(strings.TrimSpace(rows[i][j]))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse cell %d:%d: %w", i, j, err)
			}
			dividers[j/2] = append(dividers[j/2], Divider{PortNumber: portNumber})
		}
	}

	return floors, dividers, nil
}
