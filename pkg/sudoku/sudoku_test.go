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
	"testing"

	"github.com/slealq/sudokuSolver/pkg/common"
)

func newData() *[][]byte {
	return &[][]byte{
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
}

func TestSudokuValid(t *testing.T) {

	board := [][]byte{
		{'5', '4', '3', '1', '2', '7', '9', '6', '8'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '.', '.', '.', '.'},
	}

	if isValidSudoku(board) == false {
		t.Error("Board expected to be valid")
	}

}

// BechmarkSolver performs a benchmark in the brute-force algorithm
// baseline:
// BenchmarkSolver-8                 129	   9540439 ns/op	 1710449 B/op	   42574 allocs/op

func BenchmarkSolver(b *testing.B) {

	for i := 0; i < b.N; i++ {
		data := newData()
		solveSudoku(*data)
	}

}

// TestSudokuSolver verifies that the brute-force sudoku solver works as
// expected
func TestSudokuSolver(t *testing.T) {

	data := newData()

	solveSudoku(*data)

	if isValidSudoku(*data) == false {
		t.Errorf("Board expected to be valid")
	}

	expectedStr := `5 3 4 6 7 8 9 1 2 
6 7 2 1 9 5 3 4 8 
1 9 8 3 4 2 5 6 7 
8 5 9 7 6 1 4 2 3 
4 2 6 8 5 3 7 9 1 
7 1 3 9 2 4 8 5 6 
9 6 1 5 3 7 2 8 4 
2 8 7 4 1 9 6 3 5 
3 4 5 2 8 6 1 7 9 
`

	if boardStr := common.PrintSudoku(data); boardStr != expectedStr {
		t.Errorf("board:\n%sis not the expected\n%s", boardStr, expectedStr)
	}
}
