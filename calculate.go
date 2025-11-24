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

				// fmt.Println(flatLeft)
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

// 12| 536 537 538 539 | 540 541 542 543 544
// 11| 527 528 529 530 | 531 532 533 534 535
// 10| 518 519 520 521 | 522 523 524 525 526
// 9| 509 510 511 512 | 513 514 515 516 517
// 8| 500 501 502 503 | 504 505 506 507 508
// 7| 491 492 493 494 | 495 496 497 498 499
// 6| 482 483 484 485 | 486 487 488 489 490
// 5| 469 470 471 472 473 474 | 475 476 477 478 479 480 481
// 4| 456 457 458 459 460 461 | 462 463 464 465 466 467 468
// 3| 443 444 445 446 447 448 | 449 450 451 452 453 454 455
// 2| 430 431 432 433 434 435 | 436 437 438 439 440 441 442
// 1| 417 418 419 420 421 422 | 423 424 425 426 427 428 429

// 509 - 536
// 482-483,491-510
