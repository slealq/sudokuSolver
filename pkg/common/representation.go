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

package common

import (
	"fmt"
	"strings"
)

// PrintSudoku takes a boardData as input, and returns a string representation
// of the board
func PrintSudoku(board *[][]byte) string {

	var sb strings.Builder

	for _, row := range *board {
		for _, value := range row {
			fmt.Fprintf(&sb, "%s ", string(value))
		}
		fmt.Fprintf(&sb, "\n")
	}

	return sb.String()
}
