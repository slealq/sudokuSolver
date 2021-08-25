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
	"strings"
)

// container is a single sudoku container from a sudoku board.
// Nine cell values are stored in a single container. One single cell
// can be shared across several values.
//
// container is responsible of storing the values stored inside the container,
// as well as updating the possibleValues and the restrictedValues as well
type container struct {
	numberSet map[string]int
	// possibleValues
	possibleValues map[string]bool
	// Restricted values are the ones that
	// are posible and concruent with the posibilities
	// of near containers. Map the posibility with
	// the coordinate where it's posible
	restrictedValues map[string]map[Point]bool
	id               string
	observers        map[string]containerObserver
	availableValues  *availableValues
}

// newContainer returns a pointer to a new container initialized
func newContainer(id string) *container {
	aContainer := &container{id: id}

	aContainer.observers = make(map[string]containerObserver)

	// start with all available values set to true
	aContainer.availableValues = newAvailableValues()

	return aContainer
}

// Id returns the ID of this container
func (s *container) Id() string {
	return s.id
}

// update satifies the cellObserver interface, and is used to receive
// notifications from the cells when the value changes.
func (s *container) update(aCell *cell) {

	// cell notification arrived
	aLog := newLog(cellNotificationArrived, s.id, aCell.id, aCell.String())
	aLog.Info()

	// log the result at the end
	defer func() {
		aLog = newLog(containerAvailableValues, s.id, s.availValStr())
		aLog.Info()
	}()

	// case where the update removes a value from the board, in which case
	// enable again the previous value as available
	if aCell.value == byte('.') {
		if _, ok := (*s.availableValues)[aCell.preValue]; !ok {
			aLog := newLog(cellPrevValueInvalid, aCell.id, string(aCell.preValue))
			aLog.Warn()
			return
		}

		// If previous value is not empty, restore it as available
		(*s.availableValues)[aCell.preValue] = true
		s.notifyObservers()
	} else
	// value should be valid. Remove that value from available values
	{
		var available, ok bool
		if available, ok = (*s.availableValues)[aCell.value]; !ok {
			aLog := newLog(cellUpdateInvalidValue, s.id, aCell.String())
			aLog.Error()
			panic(aLog.logMsg)
		}

		// double check that value is still available in the container
		if !available {
			aLog := newLog(containerValueNotAvailable, aCell.String(), s.id)
			aLog.Error()
			panic(aLog.logMsg)
		}

		(*s.availableValues)[aCell.value] = false
		s.notifyObservers()
	}
}

// availValStr returns the string representation of the available values in
// the container
func (s *container) availValStr() string {

	var sb strings.Builder

	fmt.Fprint(&sb, "{")
	for value, available := range *s.availableValues {
		fmt.Fprintf(&sb, "%s: %v, ", string(value), available)
	}
	fmt.Fprint(&sb, "}")

	return sb.String()
}

// addObserver adds am observer to the observers map. The observers all get
// a notification when the availableValues of this container change
func (s *container) addObserver(newObserver containerObserver) error {

	if _, ok := s.observers[newObserver.Id()]; ok {
		return fmt.Errorf(containerObserverAlreadyRegistered, newObserver.Id())
	}

	s.observers[newObserver.Id()] = newObserver

	return nil
}

// rmObserver removes an observer from the observers map. The id is used
// to identify the target observer
func (s *container) rmObserver(obs containerObserver) error {

	if _, ok := s.observers[obs.Id()]; !ok {
		return fmt.Errorf(containerObserverNotFound, obs.Id())
	}

	delete(s.observers, obs.Id())

	return nil
}

// notifyObservers sends a notification of update to all observers.
func (s *container) notifyObservers() {

	for _, obs := range s.observers {
		obs.update(s)
	}
}

func (s *container) create() {
	s.numberSet = map[string]int{}
	s.createPossibleValues()
	s.restrictedValues = map[string]map[Point]bool{}
}

// createPossibleValues stores
func (s *container) createPossibleValues() {
	result := map[string]bool{}

	for value := range allValues {
		if s.numberSet[value] == 0 {
			result[value] = true
		}
	}
	s.possibleValues = result
}

func (s *container) add(i, j int, value string) {
	if s.numberSet == nil {
		s.create()
	}

	if value == "." {
		// skip counting dots
		return
	}

	s.numberSet[value]++
	s.updatePossibleValues(&value)
	// If no restricted values set, don't remove
	if len(s.restrictedValues) != 0 {
		s.rmRestrictedValue(value)
	}
}

// simpleAdd adds a value to the numberSet of the container. Should only be
// used for backtracking
func (s *container) simpleAdd(i, j int, value string) {
	if s.numberSet == nil {
		s.create()
	}

	if value == "." {
		// skip counting dots
		return
	}

	s.numberSet[value]++
}

// simpleRm removes a value from the numberSet of this container.
// remove should only be used for backtracking, since it
// will irreversible damage restrictedValues
func (s *container) simpleRm(i, j int, value string) {
	if s.numberSet == nil {
		return
	}

	s.numberSet[value]--
}

func (s *container) updatePossibleValues(value *string) {
	// set this value as `not` possible
	s.possibleValues[*value] = false
}

func (s *container) addRestricted(i, j int, value string) {
	if s.restrictedValues[value] == nil {
		s.restrictedValues[value] = map[Point]bool{}
	}

	s.restrictedValues[value][Point{i, j}] = true
}

func (s *container) rmRestrictedValue(value string) {
	if s.restrictedValues[value] == nil {
		return
	}

	delete(s.restrictedValues, value)
}

func (s *container) rmRestrictedPoint(i, j int, value string) {
	if s.restrictedValues[value] == nil {
		return
	}
	if _, ok := s.restrictedValues[value][Point{i, j}]; !ok {
		return
	}

	delete(s.restrictedValues[value], Point{i, j})
}

func (s *container) getUniqueRestricted() map[string]Point {
	// Return a map of value -> point that are unique to this
	// container (That point can ONLY allow that value)

	result := map[string]Point{}

	for value, pointSet := range s.restrictedValues {
		if pointSet == nil {
			break
		}
		if len(pointSet) == 1 {
			keys := make([]Point, 0, len(pointSet))
			for k := range pointSet {
				keys = append(keys, k)
			}

			// Add the unique value
			result[value] = keys[0]
		}
	}
	return result
}

// possibleValues returns the posible values for this
// sudoku container
func (s *container) getPossibleValues() *map[string]bool {
	return &s.possibleValues
}

// isValid says if the values stored in this container
// are all different
func (s *container) isValid() bool {
	for _, count := range s.numberSet {
		if count > 1 {
			return false
		}
	}
	return true
}
