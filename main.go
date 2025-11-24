package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
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
	// workingDirectory, err := getWorkingDirectory()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Запускаем диалог поиска нужного исходного файла в заданной рабочей папке.
	// filePath, choosed := chooseFile(workingDirectory)
	// if choosed {
	// 	fmt.Println("no file choose, exiting")
	// 	return
	// }

	// Считываем исходные данные из найденного файла.
	floors, dividers, err := parseFile("./Кв-емкость.xlsx")
	if err != nil {
		panic(err)
	}

	dividers = calculate(floors, dividers)

	table := tablewriter.NewWriter(os.Stdout)

	for j := range len(dividers[0]) {
		var row []any
		for _, riser := range dividers {
			row = append(row, strings.Join(lo.Map(riser[j].Flats, func(item FlatRange, _ int) string { return item.String() }), ","))
		}

		err = table.Append(row...)
		if err != nil {
			panic(err)
		}

	}

	err = table.Render()
	if err != nil {
		panic(err)
	}

	w := tablewriter.NewWriter(os.Stdout)

	flatToDivider := make(map[int]int)
	for _, riser := range dividers {
		for j, divider := range riser {
			for _, flatRange := range divider.Flats {
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

			flats = append(flats, color.New(color.Attribute(31+flatToDivider[i]%7)).Sprint(i))
		}
		row = append(row, strings.Join(flats, " "))

		err := w.Append(row...)
		if err != nil {
			panic(err)
		}
	}

	err = w.Render()
	if err != nil {
		panic(err)
	}
}
