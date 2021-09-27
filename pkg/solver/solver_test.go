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

package solver

import (
	"testing"

	"github.com/slealq/sudokuSolver/pkg/common"
	"github.com/slealq/sudokuSolver/pkg/sudoku"
)

const (
	expectedStr = `5 3 4 6 7 8 9 1 2 
6 7 2 1 9 5 3 4 8 
1 9 8 3 4 2 5 6 7 
8 5 9 7 6 1 4 2 3 
4 2 6 8 5 3 7 9 1 
7 1 3 9 2 4 8 5 6 
9 6 1 5 3 7 2 8 4 
2 8 7 4 1 9 6 3 5 
3 4 5 2 8 6 1 7 9 
`
)

func newBoard() *sudoku.Board {
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

	board := sudoku.NewBoard(&data)

	return board
}

// BenchmarkDeterministic performs a benchmark on the deterministic algorithm
// baseline:
// BenchmarkDeterministic-8   	      68	  17646803 ns/op	 3154438 B/op	   77219 allocs/op
func BenchmarkDeterministic(b *testing.B) {

	for i := 0; i < b.N; i++ {
		board := newBoard()
		solver := NewSolver(board)
		solver.Deterministic()
	}

}

// TestDeterministic verifies that the Deterministic algorithm works as
// expected
func TestDeterministic(t *testing.T) {

	board := newBoard()

	if board.IsValid() == false {
		t.Errorf("Board expected to be valid")
	}

	solver := NewSolver(board)
	solver.Deterministic()

	if boardStr := common.PrintSudoku(board.Data()); boardStr != expectedStr {
		t.Errorf("board:\n%sis not the expected\n%s", boardStr, expectedStr)
	}
}

// TestBacktrack verifies that the Backtrack algorithm works as expected
func TestBacktrack(t *testing.T) {

	board := newBoard()
	board.SetDebug(true)

	if board.IsValid() == false {
		t.Errorf("Board expected to be valid")
	}

	solver := NewSolver(board)
	solver.Backtrack()

	if boardStr := common.PrintSudoku(board.Data()); boardStr != expectedStr {
		t.Errorf("board:\n%sis not the expected\n%s", boardStr, expectedStr)
	}
}
