package main

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func parseFile(path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file %q: %w", path, err)
		}
	}()

	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		return err
	}
	fmt.Println(cell)

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	return nil
}
