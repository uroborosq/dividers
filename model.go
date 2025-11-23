package main

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

// Стояк отопления
type Riser struct {
	// Количество квартир, относящихся к стояку
	FlatNumber int
}

type Divider struct {
	// Количество портов в разделителе.
	PortNumber int
	Flats      []FlatRange
}

func (d Divider) GetPortNumber() int {
	return d.PortNumber
}
