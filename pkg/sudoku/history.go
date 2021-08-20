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

// sudokuRepr is a representation of the board in string
type sudokuRepr string

// history holds a finite amount of representations of the sudoku board, in
// order to be print for logging
type history struct {
	buffer   []sudokuRepr
	size     int
	Capacity int
}

// get returns a string representation of the last elements defined in the
// size value
func (h *history) get() string {
	var sb strings.Builder

	for _, entry := range h.buffer {
		fmt.Fprintf(&sb, "%s\n", entry)
	}

	return sb.String()
}

// push adds a new element to the history
func (h *history) push(r sudokuRepr) {
	if h.Capacity > h.size {
		h.buffer = append(h.buffer, r)
	} else {
		// remove first element
		h.buffer = h.buffer[1:]
		// add to the end
		h.buffer = append(h.buffer, r)
	}

	h.size++
}
