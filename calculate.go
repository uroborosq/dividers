package main

import (
	"math"
	"slices"

	"github.com/samber/lo"
)

func calculate(floors []Floor, dividers [][]Divider) [][]Divider {
	riserToFlatNumber := make([]int, len(floors[0].Risers))
	riserToPortNumber := make([]int, len(floors[0].Risers))

	for floor := range slices.Values(floors) {
		for i, riser := range floor.Risers {
			riserToFlatNumber[i] += riser.FlatNumber
		}
	}
	for i := range dividers {
		riserToPortNumber[i] = lo.SumBy(dividers[i], Divider.GetPortNumber)
	}

	for i, riser := range dividers {
		var flatPointer, floorPointer int

		for j, divider := range riser {
			flatLeft := int(math.Round(float64(divider.PortNumber) * float64(riserToFlatNumber[i]) / float64(riserToPortNumber[i])))

			riserToFlatNumber[i] -= flatLeft
			riserToPortNumber[i] -= divider.PortNumber

			for _, floor := range floors[floorPointer:] {
				previousRiserFlats := lo.SumBy(floor.Risers[:i], func(item Riser) int { return item.FlatNumber })
				nextRiserFlats := lo.SumBy(floor.Risers[i+1:], func(item Riser) int { return item.FlatNumber })

				flatLimit := max(0, floor.Risers[i].FlatNumber-flatLeft)
				dividers[i][j].Flats = append(dividers[i][j].Flats, FlatRange{
					FlatStart: floor.Flats.FlatStart + previousRiserFlats + flatPointer,
					FlatEnd:   floor.Flats.FlatEnd - nextRiserFlats - flatLimit,
				})

				flatLeft = flatLeft - floor.Risers[i].FlatNumber + flatPointer
				flatPointer = 0

				if flatLimit == 0 {
					floorPointer++
				}

				if flatLeft < 0 {
					flatPointer = flatLeft + floor.Risers[i].FlatNumber
					break
				} else if flatLeft == 0 {
					break
				}
			}
		}
	}

	return dividers
}
