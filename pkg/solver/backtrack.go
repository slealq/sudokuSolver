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
	"github.com/slealq/sudokuSolver/pkg/logs"
	"github.com/slealq/sudokuSolver/pkg/sudoku"
	"github.com/slealq/sudokuSolver/pkg/version"
)

type BackTracker struct {
	board   *sudoku.Board
	history version.History
}

// newBackTracker returns a new BackTracker struct
func newBackTracker(board *sudoku.Board) *BackTracker {
	aHistory := version.NewHistory(board)
	return &BackTracker{board: board, history: *aHistory}
}

// reachedEnd returns true if there's no more room to backtrack. False otherwise.
func (b *BackTracker) reachedEnd() bool {
	if b.history.Len() == 0 {
		aLog := logs.NewLog(logs.BackTrackWentWrong, b.board.Debug(),
			b.board.History.String())
		aLog.Error()

		return true
	}

	return false
}

// NextStep verifies that it's possible to try more values. It applies the
// nextPatch. Otherwise, returns a bool to inform that the end has been reached
func (b *BackTracker) NextStep(
	areMoreValuesAvailable bool,
	nextPatch *version.Patch,
) (moveForward bool, reachedEnd bool) {

	moveForward = true
	reachedEnd = false

	// If reached the end of the AvailableValues, then go back if possible.
	// If not, the end of backtracking has been reached (without solution)
	if !areMoreValuesAvailable {

		if b.reachedEnd() {
			moveForward = false
			reachedEnd = true
			return
		}

		previousPatch := b.history.Reverse()
		areMoreValuesAvailable = !(previousPatch.Iter.End())
		nextPatch = previousPatch

		return b.NextStep(areMoreValuesAvailable, nextPatch)

	}

	nextPatch.NextValue()
	b.history.Apply(nextPatch)

	return
}
