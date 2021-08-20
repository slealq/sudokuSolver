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

import "fmt"

// cell represents a single cell in the sudoku board
type cell struct {
	value     byte
	observers map[string]cellObserver
}

// newCell returns a new cell, with the observers map initialized
func newCell() cell {
	aCell := cell{}
	aCell.observers = make(map[string]cellObserver)

	return aCell
}

// addObserver adds am observer to the observers map. ID is a string that
// should uniquely identify the `newObserver`
func (c *cell) addObserver(id string, newObserver cellObserver) error {

	if _, ok := c.observers[id]; ok {
		return fmt.Errorf(cellObserverAlreadyRegistered, id)
	}

	c.observers[id] = newObserver

	return nil
}

// rmObserver removes an observer from the observers map. The id is used
// to identify the target observer
func (c *cell) rmObserver(id string) error {

	if _, ok := c.observers[id]; !ok {
		return fmt.Errorf(cellObserverNotFound, id)
	}

	delete(c.observers, id)

	return nil
}

// notifyAll sends a notification of update to all observers.
func (c *cell) notifyAll() {

	for _, obs := range c.observers {
		obs.notify()
	}
}

// update the value of the cell, and notify all observers about the change
func (c *cell) update(newValue byte) {
	c.value = newValue

	c.notifyAll()
}

// get returns the current value of the cell
func (c *cell) get() byte {
	return c.value
}
