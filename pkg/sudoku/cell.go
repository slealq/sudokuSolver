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
	"fmt"
)

// cell represents a single cell in the sudoku board
type cell struct {
	preValue, value byte
	observers       map[string]cellObserver
	id              string
	i, j            int
	availableValues *availableValues
}

// getCellId returns an ID from coordinates given
func getCellId(i, j int) string {
	return fmt.Sprintf("%di,%dj", i, j)
}

// newCell returns a new cell, with the observers map initialized
func newCell(i, j int) *cell {

	aCell := &cell{
		id:       getCellId(i, j),
		i:        i,
		j:        j,
		value:    byte('.'),
		preValue: byte('.'),
	}

	aCell.observers = make(map[string]cellObserver)

	// start with all available values set to true
	aCell.availableValues = newAvailableValues()

	return aCell
}

// Id returns the ID of this cell
func (c *cell) Id() string {
	return c.id
}

// addObserver adds am observer to the observers map. ID is a string that
// should uniquely identify the `newObserver`
func (c *cell) addObserver(newObserver cellObserver) error {

	if _, ok := c.observers[newObserver.Id()]; ok {
		return fmt.Errorf(cellObserverAlreadyRegistered, newObserver.Id())
	}

	c.observers[newObserver.Id()] = newObserver

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

// update receives a notification from the observed containers. Satifies the
// containerObserver interface.
func (c *cell) update(aContainer *container) {

	// ignore notifications if
	aLog := newLog(containerNotificationArrived, c.id, string(c.value), aContainer.id)
	aLog.Info()

	// log the result at the end
	defer func() {
		var availValStr string
		if c.availableValues != nil {
			availValStr = c.availableValues.String()
		} else {
			availValStr = "nil"
		}

		aLog = newLog(cellAvailableValues, c.id, availValStr)
		aLog.Info()
	}()

	// safely ignore update if current cell has a value
	if c.value != byte('.') {
		c.availableValues = nil
		return
	}

	if c.availableValues == nil {
		c.availableValues = newAvailableValues()
	}

	// cell should have reference to three observers (containers), in order
	// to update self available values
	if len(c.observers) != CONTAINERS_PER_CELL {
		aLog := newLog(notAllContainersAvailable, c.id)
		aLog.Error()
		panic(aLog.logMsg)
	}

	// go through all values, verify if they are available or not by polling
	// all containers. Containers can be accessed through observers
	for value := range *c.availableValues {
		availableInContainers := true

		// if any container has this value unavailable, set it as that for
		// this cell. Iterate through all containers to find any with the
		// value unavailable
		for _, obs := range c.observers {

			container := obs.(*container)
			if (*container.availableValues)[value] == false {
				(*c.availableValues)[value] = false
				availableInContainers = false
			}
		}

		// if all containers have it available, set this cell to true for
		// that value
		if availableInContainers {
			(*c.availableValues)[value] = true
		}
	}

}

// notifyObservers sends a notification of update to all observers.
func (c *cell) notifyObservers() {

	for _, obs := range c.observers {
		obs.update(c)
	}
}

// set the value of the cell, and notify all observers about the change
func (c *cell) set(newValue byte) {

	allowedValues := map[byte]bool{
		'.': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}

	// panic if the cell is being updated with an invalid value
	if _, ok := allowedValues[newValue]; !ok {
		aLog := newLog(invalidUpdateValue, string(newValue))
		aLog.Error()
		panic(aLog.logMsg)
	}

	// store value history
	c.preValue = c.value
	c.value = newValue

	c.notifyObservers()
}

// get returns the current value of the cell
func (c *cell) get() byte {
	return c.value
}

// String satifies the stringer interface
func (c *cell) String() string {
	return string(c.value)
}
