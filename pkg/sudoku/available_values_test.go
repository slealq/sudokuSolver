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

import "testing"

// TestUnique verifies that the Unique method works as expected
func TestUnique(t *testing.T) {

	av := AvailableValues{
		'1': true,
		'2': false,
		'3': false,
		'4': false,
		'5': false,
		'6': false,
		'7': false,
		'8': false,
		'9': true,
	}

	if _, unique := av.Unique(); unique {
		t.Errorf("Available values unique: %v, unique=%v, expected otherwise", av, unique)
	}

	av[byte('9')] = false

	if _, unique := av.Unique(); !unique {
		t.Errorf("Available values not unique: %v, unique=%v, expected otherwise", av, unique)
	}

}
