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
}

func (s *container) notify(cellId string) {
	// TODO Handle the notifications from cell updates
	aLog := newLog("test notification arrived")
	aLog.Error()
}

func (s *container) addID(value string) {
	s.id = value
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
