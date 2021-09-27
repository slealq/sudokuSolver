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
	"encoding/json"
	"fmt"

	"github.com/slealq/sudokuSolver/pkg/common"
	"github.com/slealq/sudokuSolver/pkg/iter"
	"github.com/slealq/sudokuSolver/pkg/logs"
)

// Patch is a structure that holds a coordinate and a value, which represent a
// particular change in that coordinate.
type Patch struct {
	Coordinate common.Coordinate
	Iter       *iter.ByteIterator
	Value      byte
}

// NewPatch creates a new Diff instance
func NewPatch(coord common.Coordinate, iter *iter.ByteIterator) *Patch {
	return &Patch{Iter: iter, Coordinate: coord}
}

// NextValue sets the value for the next available value in the iterator
func (p *Patch) NextValue() error {

	// TODO: Maybe it doesn't make much sense to have this here. Refactor
	// later into another type

	if ok, value := p.Iter.Next(); ok {
		p.Value = value
	} else {
		aLog := logs.NewLog(logs.NoNextValue, *p)
		aLog.Error()

		return fmt.Errorf(aLog.Msg())
	}

	return nil
}

// MarshalJSON satisfies the Marshal interface for encoding in json
func (p *Patch) MarshalJSON() ([]byte, error) {

	var value string
	if p.Value != byte(0) {
		value = string(p.Value)
	} else {
		value = ""
	}

	patch := struct {
		Coordinate common.Coordinate  `json:"coordinate"`
		Iter       *iter.ByteIterator `json:"iterator"`
		Value      string             `json:"value"`
	}{
		Coordinate: p.Coordinate,
		Iter:       p.Iter,
		Value:      value,
	}

	return json.Marshal(patch)
}

// String satisfies the stringer interface
func (p *Patch) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}
