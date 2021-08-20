/*
Copyright (C) 2021 sleal (Stuart Leal Quesada)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package sudoku

import (
	"fmt"
	"strconv"
	"strings"
)

// Complete sudoku Board
type Board struct {
	boxContainer    [3][3]container
	columnContainer [9]container
	rowContainer    [9]container
	data            *[][]byte
	possibleValues  [9][9][]string
}

func (b *Board) add(i, j int, value string) {
	b.addToContainers(i, j, value)
	b.rmRestrictedFromContainers(i, j, value)
	(*b.data)[i][j] = byte(value[0])
}

func (b *Board) simpleAdd(i, j int, value string) {
	b.boxContainer[i/3][j/3].simpleAdd(i, j, value)
	b.columnContainer[j].simpleAdd(i, j, value)
	b.rowContainer[i].simpleAdd(i, j, value)

	(*b.data)[i][j] = byte(value[0])
}

func (b *Board) simpleRm(i, j int, value string) {
	b.boxContainer[i/3][j/3].simpleRm(i, j, value)
	b.columnContainer[j].simpleRm(i, j, value)
	b.rowContainer[i].simpleRm(i, j, value)

	(*b.data)[i][j] = byte("."[0])
}

func (b *Board) rmRestrictedFromContainers(i, j int, value string) {
	// All containers need to remove the possible value, in
	// the corresponding i,j row column combination, and all containers
	// related
	b.rowContainer[i].rmRestrictedPoint(i, j, value)
	b.columnContainer[j].rmRestrictedPoint(i, j, value)

	for iVar := 0; iVar < 3; iVar++ {
		b.boxContainer[iVar][j/3].rmRestrictedPoint(iVar, j, value)
	}
	for jVar := 0; jVar < 3; jVar++ {
		b.boxContainer[i/3][jVar].rmRestrictedPoint(i, jVar, value)
	}
}

func (b *Board) addToContainers(i, j int, value string) {
	iIndexBox := i / 3
	jIndexBox := j / 3

	b.boxContainer[iIndexBox][jIndexBox].add(i, j, value)
	b.columnContainer[j].add(i, j, value)
	b.rowContainer[i].add(i, j, value)

	// Update restricted values
	// b.boxContainer[iIndexBox][jIndexBox].
}

func (b *Board) addIdToContainers() {
	for i := 0; i < 9; i++ {
		b.rowContainer[i].addID(fmt.Sprintf("row: %d", i))
		for j := 0; j < 9; j++ {
			b.boxContainer[i/3][j/3].addID(fmt.Sprintf("box: %d,%d", i/3, j/3))
		}
	}
	for j := 0; j < 9; j++ {
		b.columnContainer[j].addID(fmt.Sprintf("col: %d", j))
	}
}

func (b *Board) createBoard(Board *[][]byte) {
	b.boxContainer = [3][3]container{}
	b.columnContainer = [9]container{}
	b.rowContainer = [9]container{}
	b.data = Board
	b.addIdToContainers()

	for i, row := range *Board {
		for j, ijthValue := range row {
			b.addToContainers(i, j, string(ijthValue))
		}
	}
}

// func (b *Board) getBoard() *[][]byte {
//     return &b.data
// }

func (b *Board) isValid() bool {
	for _, boxRow := range b.boxContainer {
		for _, ijthContainer := range boxRow {
			if !ijthContainer.isValid() {
				return false
			}
		}
	}

	for _, ithContainer := range b.columnContainer {
		if !ithContainer.isValid() {
			return false
		}
	}

	for _, jthContainer := range b.rowContainer {
		if !jthContainer.isValid() {
			return false
		}
	}

	return true
}

// func (b *Board) updatePossibleValues() {
// 	for i := 0; i < 9; i++ {
// 		for j := 0; j < 9; j++ {
// 			if string(b.data[i][j]) != "." {
// 				continue
// 			}

// 			iIndexBox := i / 3
// 			jIndexBox := j / 3

// 			boxPossibleValues := b.boxContainer[iIndexBox][jIndexBox].getPossibleValues()
// 			columnPossibleValues := b.columnContainer[j].getPossibleValues()
// 			rowPossibleValues := b.rowContainer[i].getPossibleValues()

// 			result := []string{}
// 		}
// 	}
// }

func (b *Board) addRestrictedToContainer(i, j int, value string) {
	iIndexBox := i / 3
	jIndexBox := j / 3

	b.boxContainer[iIndexBox][jIndexBox].addRestricted(i, j, value)
	b.columnContainer[j].addRestricted(i, j, value)
	b.rowContainer[i].addRestricted(i, j, value)
}

func (b *Board) calculatePossibleValues() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			psv := b.calculatePossibleValuesInCoordinate(i, j)
			b.possibleValues[i][j] = *psv

			// update the restrictedValues in each container
			for _, val := range *psv {
				b.addRestrictedToContainer(i, j, val)
			}
		}
	}
}

func (b *Board) getUniqueRestrictedFromBox(i, j int) map[string]Point {
	iIndexBox := i / 3
	jIndexBox := j / 3

	// s.restrictedValues = map[string]map[Point]bool{}
	//map[string]Point
	return b.boxContainer[iIndexBox][jIndexBox].getUniqueRestricted()

}

func (b *Board) getUniqueRestrictedFromRow(i int) map[string]Point {
	return b.rowContainer[i].getUniqueRestricted()
}

func (b *Board) getUniqueRestrictedFromCol(j int) map[string]Point {
	return b.columnContainer[j].getUniqueRestricted()
}

func (b *Board) calculatePossibleValuesInCoordinate(i, j int) *[]string {
	if string((*b.data)[i][j]) != "." {
		// fmt.Printf("This place is filled with: %s\n", string(b.data[i][j]))
		return &[]string{}
	}

	iIndexBox := i / 3
	jIndexBox := j / 3

	boxPossibleValues := b.boxContainer[iIndexBox][jIndexBox].getPossibleValues()
	columnPossibleValues := b.columnContainer[j].getPossibleValues()
	rowPossibleValues := b.rowContainer[i].getPossibleValues()

	result := []string{}
	for value, _ := range allValues {
		if (*boxPossibleValues)[value] && (*columnPossibleValues)[value] && (*rowPossibleValues)[value] {
			result = append(result, value)
		}
	}

	return &result
}

func (b *Board) getPossibleValues(i, j int) []string {
	return b.possibleValues[i][j]
}

func (b *Board) spacesLeft() int {
	var spacesLeft int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			// fmt.Printf("Place %d, %d, value %s\n", i, j, string(b.data[i][j]))
			if string((*b.data)[i][j]) == "." {
				spacesLeft++
			}
		}
	}
	return spacesLeft
}

func (b *Board) GetFirstEmptyPlace() Point {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if string((*b.data)[i][j]) == "." {
				return Point{i, j}
			}
		}
	}
	// This should not happen
	return Point{-1, -1}
}

// func (b *Board) ApplyTranslations(translations []Fill) {
//     for _, fill := range traslations {
//         (*b.data)[fill.point.X][fill.point.Y] = byte(string(fill.value)[0])
//     }
// }

// func (b *Board) ReverseTranslations(translations []Fill) {
//     for _, fill := range traslations {
//         (*b.data)[fill.point.X][fill.point.Y] = byte(".")
//     }
// }

func (b *Board) ApplyTranslation(translation Fill) {
	//fmt.Printf("value to be saved %s\n", strconv.Itoa(translation.value))
	//(*b.data)[translation.point.X][translation.point.Y] = byte(strconv.Itoa(translation.value)[0])
	b.simpleAdd(translation.point.X, translation.point.Y, strconv.Itoa(translation.value))
}

func (b *Board) ReverseTranslation(translation Fill) {
	//(*b.data)[translation.point.X][translation.point.Y] = byte("."[0])
	b.simpleRm(translation.point.X, translation.point.Y, strconv.Itoa(translation.value))
}

func (b *Board) Backtrack() {
	translationInOrder := []Fill{}

	//currentValue := 0
	currentPos := 0
	BackTracked := false

	for b.spacesLeft() != 0 || b.isValid() == false {

		if !BackTracked {
			tempPoint := b.GetFirstEmptyPlace()
			fill := Fill{value: 1, point: tempPoint}
			translationInOrder = append(translationInOrder, fill)
			b.ApplyTranslation(fill)
		}

		if b.isValid() && translationInOrder[currentPos].value < 9 {
			// continue back tracking
			currentPos++
			BackTracked = false
		} else {
			BackTracked = true
			if translationInOrder[currentPos].value == 9 {
				if len(translationInOrder) <= 1 {
					fmt.Printf("Backtracking went wrong\n")
					break
				}
				// remove this element
				b.ReverseTranslation(translationInOrder[currentPos])
				translationInOrder = translationInOrder[:len(translationInOrder)-1]
				currentPos--

				// increase the value of the previous
			}

			// this needs to be done always
			b.ReverseTranslation(translationInOrder[currentPos])
			translationInOrder[currentPos].value++
			b.ApplyTranslation(translationInOrder[currentPos])
		}

		//         fmt.Printf("Filled: %d\n", len(translationInOrder))

		//         fmt.Printf("Backtracing\n")
		//         fmt.Printf("%s\n", b.String())

		// fill first place that is emtpy
		// check if its valid
		// if valid, repeat
		// if not valid, go back one place (And remove current from map)
		//   if previous value is 9, go back one place (and remove current from map)
		//   if not, increase value, check if valid
	}
}

func (b *Board) String() string {
	var sb strings.Builder

	firstRow := true

	for i := 0; i < 9; i++ {
		if firstRow {
			fmt.Fprintf(&sb, "  | ")
			for k := 0; k < 9; k++ {
				fmt.Fprintf(&sb, "%d ", k)
			}
			fmt.Fprintf(&sb, "\n")
			fmt.Fprintf(&sb, "  | ")
			for k := 0; k < 9; k++ {
				fmt.Fprintf(&sb, "__")
			}
			fmt.Fprintf(&sb, "\n")
			firstRow = false
		}

		fmt.Fprintf(&sb, "%d | ", i)

		for j := 0; j < 9; j++ {
			fmt.Fprintf(&sb, "%s ", string((*b.data)[i][j]))
		}
		fmt.Fprintf(&sb, "\n")
	}

	return b.String()
}
