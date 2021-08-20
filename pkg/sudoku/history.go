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
)

// history holds a finite amount of representations of the sudoku board, in
// order to be print for logging
type history struct {
	buffer   [][][]byte
	size     int
	Capacity int
}

// String returns a string representation of the last elements defined in the
// size value
func (h *history) String() string {
	var sb strings.Builder

	for _, entry := range h.buffer {
		fmt.Fprintf(&sb, "%s\n", printSudoku(&entry))
	}

	return sb.String()
}

// copyData copies the input data to a result variable, copiying each row
// individually
func (h *history) copyData(result *[][]byte, input [][]byte) {

	*result = make([][]byte, 0, ROW_LENGTH)

	// copy each row to avoid sharing the same underlying information
	for _, row := range input {
		newRow := make([]byte, COLUMN_LENGTH)

		copy(newRow, row)

		*result = append(*result, newRow)
	}
}

// push adds a new element to the history
func (h *history) push(data [][]byte) {

	// create copy of data for storage
	var newData [][]byte
	h.copyData(&newData, data)

	if h.Capacity > h.size {
		h.buffer = append(h.buffer, newData)
	} else {
		// remove first element
		h.buffer = h.buffer[1:]
		// add to the end
		h.buffer = append(h.buffer, newData)
	}

	h.size++
}
