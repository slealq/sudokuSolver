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
)

// TestHistoryOrder verifies that when the maximum capacity of the history
// is reached, the first item is removed from the history, and the new item
// is pushed to the end. History is read from top (oldest) to bottom (newest)
func TestHistoryOrder(t *testing.T) {

	var historyStr string
	hist := history{Capacity: 3}

	hist.push("1")
	hist.push("2")
	hist.push("3")

	historyStr = hist.get()
	if expected := "1\n2\n3\n"; historyStr != expected {
		t.Errorf("History should be: %s, was: %s", historyStr, expected)
	}

	hist.push("4")

	historyStr = hist.get()
	if expected := "2\n3\n4\n"; historyStr != expected {
		t.Errorf("History should be: %s, was: %s", historyStr, expected)
	}
}
