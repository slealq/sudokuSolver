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

package logs

import (
	"fmt"
	"strings"

	"github.com/slealq/sudokuSolver/pkg/common"
)

// History holds a finite amount of representations of the sudoku board, in
// order to be print for logging
type History struct {
	buffer   [][][]byte
	size     int
	Capacity int
}

// String returns a string representation of the last elements defined in the
// size value
func (h *History) String() string {
	var sb strings.Builder

	for _, entry := range h.buffer {
		fmt.Fprintf(&sb, "%s\n", common.PrintSudoku(&entry))
	}

	return sb.String()
}

// copyData copies the input data to a result variable, copiying each row
// individually
func (h *History) copyData(result *[][]byte, input [][]byte) {

	*result = make([][]byte, 0, common.ROW_LENGTH)

	// copy each row to avoid sharing the same underlying information
	for _, row := range input {
		newRow := make([]byte, common.COLUMN_LENGTH)

		copy(newRow, row)

		*result = append(*result, newRow)
	}
}

// Push adds a new element to the History
func (h *History) Push(data [][]byte) {

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
