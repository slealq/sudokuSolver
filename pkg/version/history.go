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

package version

import (
	"github.com/slealq/sudokuSolver/pkg/logs"
	"github.com/slealq/sudokuSolver/pkg/sudoku"
)

type History struct {
	changes []*Patch
	board   *sudoku.Board
}

// NewHistory returns a new history instance
func NewHistory(board *sudoku.Board) *History {
	return &History{board: board}
}

// Reverse pops a patch and reverts it from all containers and the board.
// Finally, it returns the popped patch
func (h *History) Reverse() *Patch {
	last := h.last()

	aLog := logs.NewLog(
		logs.BacktrackingReverse,
		string(last.Value),
		last.Coordinate,
	)
	aLog.Info()

	// Reverse from the board
	h.board.Set(last.Coordinate.X, last.Coordinate.Y, byte('.'))

	// Remove from the list of changes
	h.changes = h.changes[:len(h.changes)-1]

	return last
}

// Apply takes a patch object and applies it to the abord. Adds the patch to
// the history changes
func (h *History) Apply(patch *Patch) {
	h.changes = append(h.changes, patch)

	aLog := logs.NewLog(
		logs.BacktrackingSet,
		string(patch.Value),
		patch.Coordinate,
	)
	aLog.Info()

	h.board.Set(patch.Coordinate.X, patch.Coordinate.Y, patch.Value)
}

// Len returns how many changes there are in the history
func (h *History) Len() int {
	return len(h.changes)
}

// last returns the last patch object applied
func (h *History) last() *Patch {
	return h.changes[len(h.changes)-1]
}
