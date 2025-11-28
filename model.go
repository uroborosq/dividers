package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// Этаж
type Floor struct {
	// Номер этажа
	Number int
	// Стояки на этаже
	Risers []Riser
	// Номера квартир
	Flats FlatRange
}

type FlatRanges []FlatRange

func (r FlatRanges) String() string {
	formatted := strings.Join(lo.Map(r, func(item FlatRange, index int) string { return item.String() }), ",")

	return formatted + "кв"
}

type FlatRange struct {
	// Начальный номер квартиры
	FlatStart int
	// Конечный номер квартиры
	FlatEnd int
}

func (r FlatRange) String() string {
	if r.FlatStart == r.FlatEnd {
		return strconv.Itoa(r.FlatEnd)
	}

	return fmt.Sprintf("%d-%d", r.FlatStart, r.FlatEnd)
}

// Стояк отопления
type Riser struct {
	// Количество квартир, относящихся к стояку
	FlatNumber int
}

type Splitter struct {
	// Количество портов в разделителе.
	PortNumber int
	Flats      FlatRanges
}

func (d Splitter) GetPortNumber() int {
	return d.PortNumber
}
