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
	"encoding/json"
	"fmt"

	"github.com/slealq/sudokuSolver/pkg/common"
	"github.com/slealq/sudokuSolver/pkg/logs"
	"github.com/slealq/sudokuSolver/pkg/sudoku"
	"github.com/slealq/sudokuSolver/pkg/version"
)

type Solver struct {
	board *sudoku.Board
}

// NewSolver returns a new Solver. Receives a ptr to a board, which is the
// target to solve
func NewSolver(board *sudoku.Board) *Solver {

	solver := &Solver{board: board}

	return solver
}

// Deterministic solves the board using the `deterministic approach`, which
// is based on the basic rules of sudoku, assuming that the board following
// steps are all deterministic. Returns an error if it reaches a point where
// there's no more deterministic steps.
func (s *Solver) Deterministic() error {

	var fillFound bool
	var availableValues *sudoku.AvailableValues

	for s.board.SpacesLeft() != 0 {

		fillFound = false

		// Go through the complete board, searching for an cell that has
		// only one possible value. If so, fill it
		for i := 0; i < common.ROW_LENGTH; i++ {
			for j := 0; j < common.COLUMN_LENGTH; j++ {
				availableValues =
					s.board.GetAvailableValues(common.Coordinate{X: i, Y: j})

				if value, unique := availableValues.Unique(); unique {
					fillFound = true
					s.board.Set(i, j, value)
				}
			}
		}

		// If no place to fill was found, there's no more places the
		// deterministic approach can fill
		if !fillFound {
			aLog := logs.NewLog(logs.DeterministicApprNoMoreSteps)
			aLog.Error()
			return fmt.Errorf(aLog.Msg())
		}
	}

	return nil
}

// Backtrack performs a backtracking algorithm to the current board values,
// in which it tests values and goes backwards if it reaches a point where the
// values make the board invalid. Backtracking should end when all the cells
// in the board are filled
func (s *Solver) Backtrack() {
	// Check board is valid before calling backtracking, otherwise it will
	// never be able to solve
	// TODO: This might not be required if the board is never able to be
	// invalid
	if !s.board.IsValid() {
		aLog := logs.NewLog(logs.CannotBacktrack, s.board.String())
		aLog.Error()
		return
	}

	backtracker := newBackTracker(s.board)

	// moveForward flag is turned on when the current position hasn't begun
	// testing new numbers yet. Meaning we are arriving at this position for
	// the first time.
	moveForward := true
	var areMoreValuesAvailable bool
	var nextPatch *version.Patch
	var reachedEnd bool

	for s.board.SpacesLeft() != 0 {

		aLog := logs.NewLog(
			logs.BackTrackingStats,
			backtracker.history.Len(),
			moveForward,
			s.board.String(),
		)
		aLog.Info()

		if moveForward {
			_, newCoord := s.board.FirstEmptyPlace()
			it := s.board.GetAvailableValues(newCoord).NewIterator()

			aLog := logs.NewLog(
				logs.StartNewBackTrackPos,
				newCoord,
				s.board.GetAvailableValues(newCoord).String(),
			)
			aLog.Info()

			patch := version.NewPatch(newCoord, it)

			areMoreValuesAvailable = !(it.End())
			nextPatch = patch

		} else {
			previousPatch := backtracker.history.Reverse()

			areMoreValuesAvailable = !(previousPatch.Iter.End())
			nextPatch = previousPatch
		}

		json.Marshal(nextPatch)
		aLog = logs.NewLog(
			logs.NextPatch,
			nextPatch,
		)
		aLog.Info()

		moveForward, reachedEnd = backtracker.NextStep(
			areMoreValuesAvailable,
			nextPatch,
		)

		if reachedEnd {
			break
		}
	}
}
