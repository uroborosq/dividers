package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "dividers",
		Short: "calculate",
		Args:  cobra.RangeArgs(0, 1),
		RunE:  execute,
	}

	cmd.Flags().Bool("output", false, "show colorized schema with floors and splitters")

	_ = cmd.Execute()
}

func execute(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if len(args) == 1 {
		dir = args[0]
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		} else if !strings.HasSuffix(file.Name(), ".xlsx") {
			continue
		}

		// Складываем имя файла и имя папки в один путь.
		path := filepath.Join(dir, file.Name())

		fmt.Println("Processing", path)

		// Считываем исходные данные из найденного файла.
		floors, splitters, err := parseFile(path)
		if err != nil {
			return err
		}

		// Распределяем квартиры по разделителям.
		splitters = calculate(floors, splitters)

		// Отобразить план этажей в терминале
		err = displayFloors(floors, splitters)
		if err != nil {
			return err
		}

		// Отобразить разделители с соответствующими квартирами
		err = displayDividers(splitters)
		if err != nil {
			return err
		}

		err = writeResults(path, splitters)
		if err != nil {
			return err
		}
	}

	return nil
}
