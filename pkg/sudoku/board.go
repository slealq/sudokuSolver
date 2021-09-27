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
	"strings"

	"github.com/slealq/sudokuSolver/pkg/common"
	"github.com/slealq/sudokuSolver/pkg/logs"
)

// Complete sudoku Board
type Board struct {
	boxContainer    [3][3]*container
	columnContainer [9]*container
	rowContainer    [9]*container
	data            *[][]byte
	possibleValues  [9][9][]string

	cells   [9][9]*cell
	debug   bool
	History logs.History
}

// NewBoard creates a new board, instanciates containers, adds ids to them
// and finally add values to each container
func NewBoard(data *[][]byte) *Board {

	b := &Board{}

	// TODO: Maybe just enable with debug info
	b.History = logs.History{Capacity: 20}
	b.data = data

	b.newContainers()
	b.initCells()

	for i, row := range *data {
		for j, ijthValue := range row {
			b.addToContainers(i, j, string(ijthValue))
		}
	}

	return b
}

// Debug returns the debug Value
func (b *Board) Debug() bool {
	return b.debug
}

// Set debug value of the board
func (b *Board) SetDebug(debug bool) {
	b.debug = debug

	aLog := logs.NewLog(logs.DebuggingStateChanged, b.debug)
	aLog.Info()
}

// String satifies the stringer interface
func (b *Board) String() string {
	var sb strings.Builder

	firstRow := true

	for i := 0; i < common.ROW_LENGTH; i++ {
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

		for j := 0; j < common.COLUMN_LENGTH; j++ {
			fmt.Fprintf(&sb, "%s ", string((*b.data)[i][j]))
		}
		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}

// GetAvailableValues returns the available values of a given cell position
func (b *Board) GetAvailableValues(coord common.Coordinate) *AvailableValues {
	return b.cells[coord.X][coord.Y].availableValues
}

// FirstEmptyPlace returns the first coordinate that is available and true if
// it's available, false otherwise
func (b *Board) FirstEmptyPlace() (bool, common.Coordinate) {

	for i := 0; i < common.ROW_LENGTH; i++ {
		for j := 0; j < common.COLUMN_LENGTH; j++ {
			if b.cells[i][j].get() == byte('.') {
				return true, common.Coordinate{X: i, Y: j}
			}
		}
	}
	return false, common.Coordinate{}
}

// Set sets a value in the board, given the coordinate and the value
func (b *Board) Set(i, j int, value byte) {

	(*b.data)[i][j] = value
	b.cells[i][j].set(value)

	b.updateHistory()
}

// SpacesLeft returns the amount of positions left empty in the board
func (b *Board) SpacesLeft() int {
	var spacesLeft int
	for i := 0; i < common.ROW_LENGTH; i++ {
		for j := 0; j < common.COLUMN_LENGTH; j++ {
			if string((*b.data)[i][j]) == "." {
				spacesLeft++
			}
		}
	}
	return spacesLeft
}

// Data returns the data ptr from the board
func (b *Board) Data() *[][]byte { return b.data }

// newContainers creates all the containers and initializes them with the
// corresponding ids
func (b *Board) newContainers() {
	for i := 0; i < common.ROW_LENGTH; i++ {
		b.rowContainer[i] = newContainer(fmt.Sprintf("row_%d", i))
		for j := 0; j < common.COLUMN_LENGTH; j++ {
			b.boxContainer[i/3][j/3] = newContainer(
				fmt.Sprintf("box_%di_%dj", i/3, j/3),
			)
		}
	}
	for j := 0; j < common.COLUMN_LENGTH; j++ {
		b.columnContainer[j] = newContainer(fmt.Sprintf("col_%d", j))
	}
}

// IsValid returns true if the board is valid, false otherwise
func (b *Board) IsValid() bool { return b.isValid() }

// initCells creates a cell value for each of the positions in data.
func (b *Board) initCells() {

	if b.data == nil {
		aLog := logs.NewLog(logs.FailedToInitCells)
		aLog.Error()
		panic(aLog.Msg())
	}

	// register observers to all cells,
	for i, row := range *b.data {
		for j := range row {
			aCell := newCell(i, j)
			b.registerObservers(aCell)
			b.cells[i][j] = aCell
		}
	}

	// set the initial layout to all cells
	for i, row := range *b.data {
		for j, value := range row {
			b.cells[i][j].set(value)
		}
	}
}

// registerObservers Add the three corresponding observer for each a given
// cell, and container
func (b *Board) registerObservers(aCell *cell) {

	iBoxIndex := aCell.i / 3
	jBoxIndex := aCell.j / 3

	aCell.addObserver(b.boxContainer[iBoxIndex][jBoxIndex])
	aCell.addObserver(b.rowContainer[aCell.i])
	aCell.addObserver(b.columnContainer[aCell.j])

	b.boxContainer[iBoxIndex][jBoxIndex].addObserver(aCell)
	b.rowContainer[aCell.i].addObserver(aCell)
	b.columnContainer[aCell.j].addObserver(aCell)
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
	if b.Debug() {
		b.History.Push(*b.data)
	}
}

// ========================= REFACTOR BELOW ===================================

func (b *Board) GetFirstEmptyPlace() common.Coordinate {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if string((*b.data)[i][j]) == "." {
				return common.Coordinate{X: i, Y: j}
			}
		}
	}
	// This should not happen
	return common.Coordinate{X: -1, Y: -1}
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

func (b *Board) getUniqueRestrictedFromBox(
	i, j int,
) map[string]common.Coordinate {
	iIndexBox := i / 3
	jIndexBox := j / 3

	return b.boxContainer[iIndexBox][jIndexBox].getUniqueRestricted()

}

func (b *Board) getUniqueRestrictedFromRow(i int) map[string]common.Coordinate {
	return b.rowContainer[i].getUniqueRestricted()
}

func (b *Board) getUniqueRestrictedFromCol(j int) map[string]common.Coordinate {
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
	for value := range common.AllValues {
		if (*boxPossibleValues)[value] && (*columnPossibleValues)[value] &&
			(*rowPossibleValues)[value] {
			result = append(result, value)
		}
	}

	return &result
}

func (b *Board) getPossibleValues(i, j int) []string {
	return b.possibleValues[i][j]
}
