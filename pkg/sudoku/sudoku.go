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
	"time"
)

func isValidSudoku(board [][]byte) bool {

	sudoku := SudokuBoard{}
	sudoku.createBoard(board)

	return sudoku.isValid()

	// boxContainer := [3][3]SudokuContainer{}
	// columnContainer := [9]SudokuContainer{}
	// rowContainer := [9]SudokuContainer{}

	// for i, row := range board {
	//     for j, ijthValue := range row {
	// iIndexBox := i / 3
	// jIndexBox := j / 3

	// fmt.Printf("pos %d, %d, val %s, to box %d, %d\n", i, j, string(ijthValue), iIndexBox, jIndexBox)

	// fill box container
	// boxContainer[iIndexBox][jIndexBox].add(string(ijthValue))
	// columnContainer[j].add(string(ijthValue))
	// rowContainer[i].add(string(ijthValue))
	//     }
	// }

	// fmt.Printf("\t->boxcontainer:\n %#v\n", boxContainer)
	// fmt.Printf("\t->columnContainer:\n %#v\n", columnContainer)
	// fmt.Printf("\t->rowContainer:\n %#v\n", rowContainer)

	//     for _, boxRow := range boxContainer {
	//         for _, ijthContainer := range boxRow {
	//             if !ijthContainer.isValid() {
	//                 return false
	//             }
	//         }
	//     }

	//     for _, ithContainer := range columnContainer {
	//         if !ithContainer.isValid() {
	//             return false
	//         }
	//     }

	//     for _, jthContainer := range rowContainer {
	//         if !jthContainer.isValid() {
	//             return false
	//         }
	//     }

	//     return true
}

//
// [["5","3",".",".","7",".",".",".","."]
// ,["6",".",".","1","9","5",".",".","."]
// ,[".","9","8",".",".",".",".","6","."]
// ,["8",".",".",".","6",".",".",".","3"]
// ,["4",".",".","8",".","3",".",".","1"]
// ,["7",".",".",".","2",".",".",".","6"]
// ,[".","6",".",".",".",".","2","8","."]
// ,[".",".",".","4","1","9",".",".","5"]
// ,[".",".",".",".","8",".",".","7","9"]]

// iterate over the row
// iterate over the columns
// is row between 0 to 2 ? -> Upper three containers
// is row between 3 to 5 ? -> middle three containers
// is row between 6 to 8 ? -? lower three containers

// is column between 0 to 2 ? -> map container
// is column between 3 to 5 ? -> map other container
// is column between 6 to 8 ? -> map last container

// in total, theres a [3][3]container matrix, which can be indexed in the inner array

// For the columns and rows, Theres one container for each

// [9]container for column

// [9]container for row

// Which can be indexed in the inner loop.

// append [i]container and [j]container

func getKeysFromRestricted(restricted map[string]Point) []string {
	keys := make([]string, 0, len(restricted))
	for k := range restricted {
		keys = append(keys, k)
	}
	return keys
}

func solveSudoku(board [][]byte) {
	start := time.Now()

	sudoku := SudokuBoard{}
	sudoku.createBoard(&board)

	// fmt.Printf("Spaces left: %d\n", sudoku.spacesLeft())
	filled := false
	for sudoku.spacesLeft() != 0 {
		sudoku.calculatePossibleValues()
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				psv := sudoku.getPossibleValues(i, j)
				if len(psv) != 0 {
					// fmt.Printf("Posible values of %d,%d are: %v\n", i, j, psv)
				}
				if len(psv) == 1 {
					filled = true
					sudoku.add(i, j, psv[0])
					// fmt.Printf("Filling place %d, %d with value: %s\n", i, j, psv[0])
					// fmt.Printf("Results is: %v\n", sudoku.getBoard())
				}

				boxRestricted := sudoku.getUniqueRestrictedFromBox(i, j)
				if len(boxRestricted) == 1 {
					// fmt.Printf("Box restricted values : %d, %d, %v\n", i/3, j/3, boxRestricted)
					value := getKeysFromRestricted(boxRestricted)[0]
					point := boxRestricted[value]

					// fmt.Printf("B Want to fill: %s, in pos %d, %d\n", value, point.X, point.Y)

					sudoku.add(point.X, point.Y, value)
					filled = true
				}

				rowRestricted := sudoku.getUniqueRestrictedFromRow(i)
				if len(rowRestricted) == 1 {
					// fmt.Printf("Row restricted values : %d, %v\n", i, rowRestricted)
					value := getKeysFromRestricted(rowRestricted)[0]
					point := rowRestricted[value]

					// fmt.Printf("R Want to fill: %s, in pos %d, %d\n", value, point.X, point.Y)

					sudoku.add(point.X, point.Y, value)
					filled = true
				}

				colRestricted := sudoku.getUniqueRestrictedFromCol(j)
				if len(colRestricted) == 1 {
					// fmt.Printf("Col restricted values : %d, %v\n", j, colRestricted)
					value := getKeysFromRestricted(colRestricted)[0]
					point := colRestricted[value]

					// fmt.Printf("C Want to fill: %s, in pos %d, %d\n", value, point.X, point.Y)

					sudoku.add(point.X, point.Y, value)
					filled = true
				}
			}
		}

		//         fmt.Printf("%v\n", sudoku.String())

		//         fmt.Printf("Is valid sudoku? :%t\n", sudoku.isValid())

		if !filled {
			// fmt.Println("Don't know how to proceed")
			// fmt.Printf("%v\n", sudoku.String())
			sudoku.Backtrack()
			break
		}
		filled = false
	}

	// board = (*sudoku.getBoard())

	// Code to measure
	duration := time.Since(start)

	// Formatted string, such as "2h3m0.5s" or "4.503μs"
	fmt.Println(duration)
}