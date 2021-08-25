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

type availableValues map[byte]bool

// newAvailableValues create a new available values with all the values
// available by default
func newAvailableValues() *availableValues {

	return &availableValues{
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}
}

// String returns the string representation of the available values in the
// container. Satisfies the Stringer interface
func (a *availableValues) String() string {

	var sb strings.Builder

	fmt.Fprint(&sb, "{")
	for value, available := range *a {
		fmt.Fprintf(&sb, "%s: %v, ", string(value), available)
	}
	fmt.Fprint(&sb, "}")

	return sb.String()
}
