package main

import "fmt"

// Этаж
type Floor struct {
	// Номер этажа
	Number int
	// Стояки на этаже
	Risers []Riser
	// Номера квартир
	Flats FlatRange
}

type FlatRange struct {
	// Начальный номер квартиры
	FlatStart int
	// Конечный номер квартиры
	FlatEnd int
}

func (r FlatRange) String() string {
	if r.FlatStart == r.FlatEnd {
		return fmt.Sprintf("%dкв", r.FlatStart)
	}

	return fmt.Sprintf("%d-%dкв", r.FlatStart, r.FlatEnd)
}

// Стояк отопления
type Riser struct {
	// Количество квартир, относящихся к стояку
	FlatNumber int
}

type Splitter struct {
	// Количество портов в разделителе.
	PortNumber int
	Flats      []FlatRange
}

func (d Splitter) GetPortNumber() int {
	return d.PortNumber
}
