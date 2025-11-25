package main

import (
	"fmt"
	"io"
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

	showOutput, err := cmd.Flags().GetBool("output")
	if err != nil {
		return err
	}

	w := io.Discard
	if showOutput {
		w = os.Stdout
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			// папки пропускаем
			continue
		} else if !strings.HasSuffix(file.Name(), ".xlsx") {
			// файлы с другими расширениями тоже пропускаем
			continue
		}

		// Складываем имя файла и имя папки в один путь.
		path := filepath.Join(dir, file.Name())

		_, err = fmt.Fprintln(w, "Processing", path)
		if err != nil {
			return err
		}

		// Считываем исходные данные из найденного файла.
		floors, splitters, err := parseFile(path)
		if err != nil {
			return err
		}

		// Распределяем квартиры по разделителям.
		splitters = calculate(floors, splitters)

		// Отобразить план этажей в терминале
		err = displayFloors(w, floors, splitters)
		if err != nil {
			return err
		}

		// Отобразить разделители с соответствующими квартирами
		err = displayDividers(w, splitters)
		if err != nil {
			return err
		}

		// Записываем результаты в файлек
		err = writeResults(path, splitters)
		if err != nil {
			return err
		}
	}

	return nil
}
