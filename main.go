package main

import (
	"fmt"
	"os"
)

func getWorkingDirectory() (string, error) {
	// Проверяем, переопределил ли пользователь рабочую папку.
	// Первый аргумент всегда - название самой программы, а вот если есть второй - значит переопределил
	// И мы его берем используем
	if len(os.Args) == 2 {
		return os.Args[1], nil
	}

	// Если пользователь ничего не указывал, то используем текущую рабочую папку (из которой открыта программа)
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return wd, nil
}

func main() {
	// Получить путь к рабочей папке, из которой будет начат поиск исходного файла.
	workingDirectory, err := getWorkingDirectory()
	if err != nil {
		panic(err)
	}

	// Запускаем диалог поиска нужного исходного файла в заданной рабочей папке.
	filePath, choosed := chooseFile(workingDirectory)
	if choosed {
		fmt.Println("no file choose, exiting")
		return
	}

	// Считываем исходные данные из найденного файла.
	floors, dividers, err := parseFile(filePath)
	if err != nil {
		panic(err)
	}

	// Распределяем квартиры по разделителям.
	dividers = calculate(floors, dividers)

	displayer := NewDisplayer(floors, dividers)

	// Отобразить план этажей в терминале
	err = displayer.displayFloors()
	if err != nil {
		panic(err)
	}

	// Отобразить разделители с соответствующими квартирами
	err = displayer.displayDividers()
	if err != nil {
		panic(err)
	}
}
