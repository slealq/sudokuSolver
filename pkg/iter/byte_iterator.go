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

package iter

import (
	"encoding/json"
)

type ByteIterator struct {
	index  int
	aSlice []byte
}

// Constructor for a new ByteIterator
func NewByteIterator(aSlice []byte) *ByteIterator {
	return &ByteIterator{aSlice: aSlice}
}

// End returns true if there's no more items left. False otherwise
func (b *ByteIterator) End() bool {
	if b.index >= len(b.aSlice) {
		return true
	}
	return false
}

// Next returns true, and the next value if there's a next value. Else, it
// returns false, and any byte value
func (b *ByteIterator) Next() (bool, byte) {

	if !b.End() {
		response := b.aSlice[b.index]
		b.index++

		return true, response
	}

	return false, byte('.')
}

// MarshalJSON satisfies the Marshal interface for encoding in json
func (b ByteIterator) MarshalJSON() ([]byte, error) {
	values := []string{}
	for _, byteValue := range b.aSlice {
		values = append(values, string(byteValue))
	}

	byteIterator := struct {
		Index  int      `json:"index"`
		Values []string `json:"values"`
	}{
		Index:  b.index,
		Values: values,
	}

	return json.Marshal(byteIterator)
}
