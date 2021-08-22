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

func TestSudokuStringer(t *testing.T) {
	data := [][]byte{
		{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}

	board := newBoard(&data)

	expectedStr := `  | 0 1 2 3 4 5 6 7 8 
  | __________________
0 | 5 3 . . 7 . . . . 
1 | 6 . . 1 9 5 . . . 
2 | . 9 8 . . . . 6 . 
3 | 8 . . . 6 . . . 3 
4 | 4 . . 8 . 3 . . 1 
5 | 7 . . . 2 . . . 6 
6 | . 6 . . . . 2 8 . 
7 | . . . 4 1 9 . . 5 
8 | . . . . 8 . . 7 9 
`

	if boardStr := board.String(); boardStr != expectedStr {
		t.Errorf("String : \n%sis not the expected:\n%s", boardStr, expectedStr)
	}
}
