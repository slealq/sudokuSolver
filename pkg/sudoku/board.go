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
	cells           [9][9]cell
	history         history
	debug           bool
}

// newBoard creates a new board, instanciates containers, adds ids to them
// and finally add values to each container
func newBoard(data *[][]byte) *Board {

	b := &Board{}

	// TODO: Maybe just enable with debug info
	b.history = history{Capacity: 20}

	b.boxContainer = [3][3]container{}
	b.columnContainer = [9]container{}
	b.rowContainer = [9]container{}
	b.data = data
	b.addIdToContainers()

	for i, row := range *data {
		for j, ijthValue := range row {
			b.addToContainers(i, j, string(ijthValue))
		}
	}

	return b
}

// initCells creates a cell value for each of the positions in data.
func (b *Board) initCells() {

	b.cells = [ROW_LENGTH][COLUMN_LENGTH]cell{}

	for i, row := range *b.data {
		for j, value := range row {
			aCell := newCell(i, j)

			aCell.update(value)
			b.cells[i][j] = aCell

			// TODO register observers
		}
	}
}

// addToContainers add each specific cell to all the containers that should
// track it
func (b *Board) addToContainers(i, j int, value string) {
	iIndexBox := i / 3
	jIndexBox := j / 3

	b.boxContainer[iIndexBox][jIndexBox].add(i, j, value)
	b.columnContainer[j].add(i, j, value)
	b.rowContainer[i].add(i, j, value)
}

// updateHistory adds a new entry to the history if debug flag is enabled
func (b *Board) updateHistory() {
	if b.debug {
		b.history.push(*b.data)
	}
}

func (b *Board) add(i, j int, value string) {
	b.addToContainers(i, j, value)
	b.rmRestrictedFromContainers(i, j, value)
	(*b.data)[i][j] = byte(value[0])
}

// simpleAdd adds a value to all the containers and updates the data
// accordingly, but doesn't update  the restricted values
func (b *Board) simpleAdd(i, j int, value string) {
	b.boxContainer[i/3][j/3].simpleAdd(i, j, value)
	b.columnContainer[j].simpleAdd(i, j, value)
	b.rowContainer[i].simpleAdd(i, j, value)

	(*b.data)[i][j] = byte(value[0])

	b.updateHistory()
}

// simpleRm removes a value from all containers and updates the board data
// accordingly, but doesn't update restricted values
func (b *Board) simpleRm(i, j int, value string) {
	b.boxContainer[i/3][j/3].simpleRm(i, j, value)
	b.columnContainer[j].simpleRm(i, j, value)
	b.rowContainer[i].simpleRm(i, j, value)

	(*b.data)[i][j] = byte("."[0])

	b.updateHistory()
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
		return &[]string{}
	}

	iIndexBox := i / 3
	jIndexBox := j / 3

	boxPossibleValues := b.boxContainer[iIndexBox][jIndexBox].getPossibleValues()
	columnPossibleValues := b.columnContainer[j].getPossibleValues()
	rowPossibleValues := b.rowContainer[i].getPossibleValues()

	result := []string{}
	for value := range allValues {
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

// ApplyTranslation takes a translation and applies it to all containers
// and to the board
func (b *Board) ApplyTranslation(translation Fill) {
	b.simpleAdd(translation.point.X, translation.point.Y, strconv.Itoa(translation.value))
}

// ReverseTranslation takes a translation and reverts it from all containers
// and the board
func (b *Board) ReverseTranslation(translation Fill) {
	b.simpleRm(translation.point.X, translation.point.Y, strconv.Itoa(translation.value))
}

// Backtrack performs a backtracking algorithm to the current board values,
// in which it tests values and goes backwards if it reaches a point where the
// values make the board invalid. Backtracking should end when all the cells
// in the board are filled
func (b *Board) Backtrack() {
	// Check board is valid before calling backtracking, otherwise it will
	// never be able to solve
	if !b.isValid() {
		aLog := newLog(cannotBacktrack, b.String())
		aLog.Error()
		return
	}

	// Holds the translations made, so they can be reversed if required
	translationInOrder := []Fill{}
	currentTrans := 0

	// newPos flag is turned on when the current position hasn't begun
	// testing new numbers yet. Meaning we are arriving at this position for
	// the first time.
	newPos := false

	for b.spacesLeft() != 0 || b.isValid() == false {

		aLog := newLog(backTrackingStats, len(translationInOrder), newPos)
		aLog.Info()

		if !newPos {
			// since newPos flag is false, make a new Fill and add it
			// to the translations
			tempPoint := b.GetFirstEmptyPlace()
			fill := Fill{value: 1, point: tempPoint}

			translationInOrder = append(translationInOrder, fill)

			b.ApplyTranslation(fill)
		}

		// When it's not the first backtracking, and the board is currently
		// valid, advance to the next position
		if b.isValid() && translationInOrder[currentTrans].value < 9 {
			// continue back tracking
			currentTrans++
			newPos = false

		} else
		// If board is not valid, then remove current translation and go
		// backwards until a new valid position is reached
		{
			newPos = true

			if translationInOrder[currentTrans].value == 9 {
				if len(translationInOrder) <= 1 {

					aLog := newLog(backTrackWentWrong, b.debug,
						b.history.String())
					aLog.Error()
					break
				}
				// remove this element
				b.ReverseTranslation(translationInOrder[currentTrans])
				translationInOrder = translationInOrder[:len(translationInOrder)-1]
				currentTrans--

				// increase the value of the previous
			}

			// At this point, currentTrans should have already been decremented
			// so it's time to increase previous value IF that value is able
			// to increase. If not, move backwards further

			// We'll reverse the translation, increase the value, and apply it
			// back again
			b.ReverseTranslation(translationInOrder[currentTrans])
			translationInOrder[currentTrans].value++
			b.ApplyTranslation(translationInOrder[currentTrans])
		}
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

	return sb.String()
}
