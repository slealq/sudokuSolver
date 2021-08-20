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

type ObserverMock struct {
	Notification bool
	cellRef      *cell
	cellValue    byte
}

func (o *ObserverMock) notify() {
	o.Notification = true
	o.cellValue = o.cellRef.get()
}

var observer ObserverMock

// TestObserverUpdate verifies that the observer receives a notification when
// the observed cell changes value
func TestObserverUpdate(t *testing.T) {

	// Create new cell, and register observer
	aCell := newCell()

	// Create new observer, with notification off. Observer has reference
	// to the cell
	observer = ObserverMock{Notification: false, cellRef: &aCell}

	if observer.Notification == true {
		t.Errorf("Observer notification set: %v", observer.Notification)
	}

	aCell.addObserver("a", &observer)
	aCell.update(byte('1'))

	if observer.Notification != true {
		t.Errorf("Notification should have arrived to observer: %v",
			observer.Notification)
	}

	if observer.cellValue != '1' {
		t.Errorf("Observer didn't receive new cell value: %s",
			string(observer.cellValue))
	}

}
