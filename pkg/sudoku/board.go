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

// Complete sudoku board
type SudokuBoard struct {
	boxContainer    [3][3]SudokuContainer
	columnContainer [9]SudokuContainer
	rowContainer    [9]SudokuContainer
	board           *[][]byte
	possibleValues  [9][9][]string
}

func (sc *SudokuBoard) add(i, j int, value string) {
	sc.addToContainers(i, j, value)
	sc.rmRestrictedFromContainers(i, j, value)
	(*sc.board)[i][j] = byte(value[0])
}

func (sc *SudokuBoard) simpleAdd(i, j int, value string) {
	sc.boxContainer[i/3][j/3].simpleAdd(i, j, value)
	sc.columnContainer[j].simpleAdd(i, j, value)
	sc.rowContainer[i].simpleAdd(i, j, value)

	(*sc.board)[i][j] = byte(value[0])
}

func (sc *SudokuBoard) simpleRm(i, j int, value string) {
	sc.boxContainer[i/3][j/3].simpleRm(i, j, value)
	sc.columnContainer[j].simpleRm(i, j, value)
	sc.rowContainer[i].simpleRm(i, j, value)

	(*sc.board)[i][j] = byte("."[0])
}

func (sc *SudokuBoard) rmRestrictedFromContainers(i, j int, value string) {
	// All containers need to remove the possible value, in
	// the corresponding i,j row column combination, and all containers
	// related
	sc.rowContainer[i].rmRestrictedPoint(i, j, value)
	sc.columnContainer[j].rmRestrictedPoint(i, j, value)

	for iVar := 0; iVar < 3; iVar++ {
		sc.boxContainer[iVar][j/3].rmRestrictedPoint(iVar, j, value)
	}
	for jVar := 0; jVar < 3; jVar++ {
		sc.boxContainer[i/3][jVar].rmRestrictedPoint(i, jVar, value)
	}
}

func (sc *SudokuBoard) addToContainers(i, j int, value string) {
	iIndexBox := i / 3
	jIndexBox := j / 3

	sc.boxContainer[iIndexBox][jIndexBox].add(i, j, value)
	sc.columnContainer[j].add(i, j, value)
	sc.rowContainer[i].add(i, j, value)

	// Update restricted values
	// sc.boxContainer[iIndexBox][jIndexBox].
}

func (sc *SudokuBoard) addIdToContainers() {
	for i := 0; i < 9; i++ {
		sc.rowContainer[i].addID(fmt.Sprintf("row: %d", i))
		for j := 0; j < 9; j++ {
			sc.boxContainer[i/3][j/3].addID(fmt.Sprintf("box: %d,%d", i/3, j/3))
		}
	}
	for j := 0; j < 9; j++ {
		sc.columnContainer[j].addID(fmt.Sprintf("col: %d", j))
	}
}

func (sc *SudokuBoard) createBoard(board *[][]byte) {
	sc.boxContainer = [3][3]SudokuContainer{}
	sc.columnContainer = [9]SudokuContainer{}
	sc.rowContainer = [9]SudokuContainer{}
	sc.board = board
	sc.addIdToContainers()

	for i, row := range *board {
		for j, ijthValue := range row {
			sc.addToContainers(i, j, string(ijthValue))
		}
	}
}

// func (sc *SudokuBoard) getBoard() *[][]byte {
//     return &sc.board
// }

func (sc *SudokuBoard) isValid() bool {
	for _, boxRow := range sc.boxContainer {
		for _, ijthContainer := range boxRow {
			if !ijthContainer.isValid() {
				return false
			}
		}
	}

	for _, ithContainer := range sc.columnContainer {
		if !ithContainer.isValid() {
			return false
		}
	}

	for _, jthContainer := range sc.rowContainer {
		if !jthContainer.isValid() {
			return false
		}
	}

	return true
}

func (sc *SudokuBoard) updatePossibleValues() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if string(sc.board[i][j]) != "." {
				continue
			}

			iIndexBox := i / 3
			jIndexBox := j / 3

			boxPossibleValues := sc.boxContainer[iIndexBox][jIndexBox].getPossibleValues()
			columnPossibleValues := sc.columnContainer[j].getPossibleValues()
			rowPossibleValues := sc.rowContainer[i].getPossibleValues()

			result := []string{}
		}
	}
}

func (sc *SudokuBoard) addRestrictedToContainer(i, j int, value string) {
	iIndexBox := i / 3
	jIndexBox := j / 3

	sc.boxContainer[iIndexBox][jIndexBox].addRestricted(i, j, value)
	sc.columnContainer[j].addRestricted(i, j, value)
	sc.rowContainer[i].addRestricted(i, j, value)
}

func (sc *SudokuBoard) calculatePossibleValues() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			psv := sc.calculatePossibleValuesInCoordinate(i, j)
			sc.possibleValues[i][j] = *psv

			// update the restrictedValues in each container
			for _, val := range *psv {
				sc.addRestrictedToContainer(i, j, val)
			}
		}
	}
}

func (sc *SudokuBoard) getUniqueRestrictedFromBox(i, j int) map[string]Point {
	iIndexBox := i / 3
	jIndexBox := j / 3

	// s.restrictedValues = map[string]map[Point]bool{}
	//map[string]Point
	return sc.boxContainer[iIndexBox][jIndexBox].getUniqueRestricted()

}

func (sc *SudokuBoard) getUniqueRestrictedFromRow(i int) map[string]Point {
	return sc.rowContainer[i].getUniqueRestricted()
}

func (sc *SudokuBoard) getUniqueRestrictedFromCol(j int) map[string]Point {
	return sc.columnContainer[j].getUniqueRestricted()
}

func (sc *SudokuBoard) calculatePossibleValuesInCoordinate(i, j int) *[]string {
	if string((*sc.board)[i][j]) != "." {
		// fmt.Printf("This place is filled with: %s\n", string(sc.board[i][j]))
		return &[]string{}
	}

	iIndexBox := i / 3
	jIndexBox := j / 3

	boxPossibleValues := sc.boxContainer[iIndexBox][jIndexBox].getPossibleValues()
	columnPossibleValues := sc.columnContainer[j].getPossibleValues()
	rowPossibleValues := sc.rowContainer[i].getPossibleValues()

	result := []string{}
	for value, _ := range allValues {
		if (*boxPossibleValues)[value] && (*columnPossibleValues)[value] && (*rowPossibleValues)[value] {
			result = append(result, value)
		}
	}

	return &result
}

func (sc *SudokuBoard) getPossibleValues(i, j int) []string {
	return sc.possibleValues[i][j]
}

func (sc *SudokuBoard) spacesLeft() int {
	var spacesLeft int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			// fmt.Printf("Place %d, %d, value %s\n", i, j, string(sc.board[i][j]))
			if string((*sc.board)[i][j]) == "." {
				spacesLeft++
			}
		}
	}
	return spacesLeft
}

func (sc *SudokuBoard) GetFirstEmptyPlace() Point {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if string((*sc.board)[i][j]) == "." {
				return Point{i, j}
			}
		}
	}
	// This should not happen
	return Point{-1, -1}
}

// func (sc *SudokuBoard) ApplyTranslations(translations []Fill) {
//     for _, fill := range traslations {
//         (*sc.board)[fill.point.X][fill.point.Y] = byte(string(fill.value)[0])
//     }
// }

// func (sc *SudokuBoard) ReverseTranslations(translations []Fill) {
//     for _, fill := range traslations {
//         (*sc.board)[fill.point.X][fill.point.Y] = byte(".")
//     }
// }

func (sc *SudokuBoard) ApplyTranslation(translation Fill) {
	//fmt.Printf("value to be saved %s\n", strconv.Itoa(translation.value))
	//(*sc.board)[translation.point.X][translation.point.Y] = byte(strconv.Itoa(translation.value)[0])
	sc.simpleAdd(translation.point.X, translation.point.Y, strconv.Itoa(translation.value))
}

func (sc *SudokuBoard) ReverseTranslation(translation Fill) {
	//(*sc.board)[translation.point.X][translation.point.Y] = byte("."[0])
	sc.simpleRm(translation.point.X, translation.point.Y, strconv.Itoa(translation.value))
}

func (sc *SudokuBoard) Backtrack() {
	translationInOrder := []Fill{}

	//currentValue := 0
	currentPos := 0
	BackTracked := false

	for sc.spacesLeft() != 0 || sc.isValid() == false {

		if !BackTracked {
			tempPoint := sc.GetFirstEmptyPlace()
			fill := Fill{value: 1, point: tempPoint}
			translationInOrder = append(translationInOrder, fill)
			sc.ApplyTranslation(fill)
		}

		if sc.isValid() && translationInOrder[currentPos].value < 9 {
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
				sc.ReverseTranslation(translationInOrder[currentPos])
				translationInOrder = translationInOrder[:len(translationInOrder)-1]
				currentPos--

				// increase the value of the previous
			}

			// this needs to be done always
			sc.ReverseTranslation(translationInOrder[currentPos])
			translationInOrder[currentPos].value++
			sc.ApplyTranslation(translationInOrder[currentPos])
		}

		//         fmt.Printf("Filled: %d\n", len(translationInOrder))

		//         fmt.Printf("Backtracing\n")
		//         fmt.Printf("%s\n", sc.String())

		// fill first place that is emtpy
		// check if its valid
		// if valid, repeat
		// if not valid, go back one place (And remove current from map)
		//   if previous value is 9, go back one place (and remove current from map)
		//   if not, increase value, check if valid
	}
}

func (sc *SudokuBoard) String() string {
	var b strings.Builder

	firstRow := true

	for i := 0; i < 9; i++ {
		if firstRow {
			fmt.Fprintf(&b, "  | ")
			for k := 0; k < 9; k++ {
				fmt.Fprintf(&b, "%d ", k)
			}
			fmt.Fprintf(&b, "\n")
			fmt.Fprintf(&b, "  | ")
			for k := 0; k < 9; k++ {
				fmt.Fprintf(&b, "__")
			}
			fmt.Fprintf(&b, "\n")
			firstRow = false
		}

		fmt.Fprintf(&b, "%d | ", i)

		for j := 0; j < 9; j++ {
			fmt.Fprintf(&b, "%s ", string((*sc.board)[i][j]))
		}
		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}
