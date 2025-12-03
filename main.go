package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rtfFile, err := os.Create("splitter.rtf")
	if err != nil {
		panic(err)
	}

	rtf := newRTFWriter(rtfFile)
	defer rtfFile.Close()
	defer rtf.Close()

	writer := dualWriter{console: os.Stdout, rtf: rtf}
	_, err = fmt.Println("splitters")
	err = execute(writer)
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}
}

func execute(w io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if len(os.Args) == 2 {
		dir = os.Args[1]
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	w.Write([]byte(dir))

	for _, file := range files {
		if err := processFile(w, dir, file); err != nil {
			fmt.Fprintf(w, "processing file %q failed: %s", file.Name())
		}
	}

	return nil
}

func processFile(w io.Writer, dir string, file os.DirEntry) error {
	if file.IsDir() {
		// папки пропускаем
		return nil
	} else if !strings.HasSuffix(file.Name(), ".xlsx") {
		// файлы с другими расширениями тоже пропускаем
		return nil
	}

	// Складываем имя файла и имя папки в один путь.
	path := filepath.Join(dir, file.Name())

	// Пишем загадочную умную надпись
	_, err := fmt.Fprintln(w, "Processing", path)
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

	return nil
}
