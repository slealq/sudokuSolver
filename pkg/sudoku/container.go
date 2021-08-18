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
	numberSet      map[string]int
	possibleValues map[string]bool
	// Restricted values are the ones that
	// are posible and concruent with the posibilities
	// of near containers. Map the posibility with
	// the coordinate where it's posible
	restrictedValues map[string]map[Point]bool
	id               string
}

func (s *container) addID(value string) {
	s.id = value
}

func (s *container) create() {
	s.numberSet = map[string]int{}
	s.createPossibleValues()
	s.restrictedValues = map[string]map[Point]bool{}
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

// should only be used for backtracking
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

// remove should only be used for backtracking, since it
// will irreversible damage restrictedValues
func (s *container) simpleRm(i, j int, value string) {
	if s.numberSet == nil {
		return
	}

	s.numberSet[value]--

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

func (s *container) updatePossibleValues(value *string) {
	// set this value as `not` possible
	s.possibleValues[*value] = false
}

func (s *container) createPossibleValues() {
	result := map[string]bool{}

	for value, _ := range allValues {
		if s.numberSet[value] == 0 {
			result[value] = true
		}
	}
	s.possibleValues = result
}

// possibleValues returns the posible values for this
// sudoku container
func (s *container) getPossibleValues() *map[string]bool {
	return &s.possibleValues
	//     result := map[string]bool{}

	//     for value, _ := range allValues {
	//         if s.numberSet[value] == 0 {
	//             result[value] = true
	//         }
	//     }
	//     return result
}

func (s *container) addRestricted(i, j int, value string) {
	if s.restrictedValues[value] == nil {
		//fmt.Printf("->addRestricted failed | id=\"%10s\", i=%d, j=%d, val=%s |\n", s.id, i, j, value)
		s.restrictedValues[value] = map[Point]bool{}
	}
	//fmt.Printf("add restricted value %d,%d val %s\n", i, j, value)

	s.restrictedValues[value][Point{i, j}] = true

	// if s.id == "box: 2,0" {
	//     fmt.Printf("->adding value to box | id=\"%10s\", i=%d, j=%d, val=%s |\n", s.id, i, j, value)
	//     fmt.Printf("restrictedValues : %v\n", s.restrictedValues)
	// }
}

func (s *container) rmRestrictedValue(value string) {
	if s.restrictedValues[value] == nil {
		// fmt.Printf("rmRestricted failed %d, %d, %s\n", i, j, value)
		return
	}
	//s.restrictedValues = nil
	// if _, ok := s.restrictedValues[value][Point{i,j}]; !ok {
	//     fmt.Printf("rmRestricted value %d,%d not found\n", i, j)
	//     return
	// }

	delete(s.restrictedValues, value)
	// delete(s.restrictedValues[value], Point{i,j})

	// Restricted values should be deleted from the map entirely
	// if the the value is already set
	// for key, _ := range(s.restrictedValues) {
	//     delete(s.restrictedValues, key)
	// }
	// sc.restrictedValues[value][Point{i,j}] = false
}

func (s *container) rmRestrictedPoint(i, j int, value string) {
	if s.restrictedValues[value] == nil {
		// fmt.Printf("rmRestricted failed %d, %d, %s\n", i, j, value)
		return
	}
	//s.restrictedValues = nil
	if _, ok := s.restrictedValues[value][Point{i, j}]; !ok {
		// fmt.Printf("rmRestricted value %d,%d not found\n", i, j)
		return
	}

	//delete(s.restrictedValues, value)
	delete(s.restrictedValues[value], Point{i, j})

	// Restricted values should be deleted from the map entirely
	// if the the value is already set
	// for key, _ := range(s.restrictedValues) {
	//     delete(s.restrictedValues, key)
	// }
	// sc.restrictedValues[value][Point{i,j}] = false
}

func (s *container) getUniqueRestricted() map[string]Point {
	// Return a map of value -> point that are unique to this
	// container (That point can ONLY allow that value)

	result := map[string]Point{}

	// restrictedValues map[string]map[Point]bool

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
