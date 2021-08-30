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

package logs

var (
	// cell.go
	CellObserverAlreadyRegistered = "Cell observer id '%s' has already been registered"
	CellObserverNotFound          = "Cell observer id '%s' was not found"
	InvalidUpdateValue            = "Invalid value for cell: %s"
	ContainerNotificationArrived  = "Cell '%-5s' with value '%s' received notification from container: '%s'"
	CellAvailableValues           = "Cell '%-5s' available values updated: %v"
	NotAllContainersAvailable     = "Cell '%-5s received an update, but doesn't have complete references to containers"

	// board.go
	BackTrackWentWrong = "Backtracking went wrong; debug=%v; history=\n%s"
	BackTrackingStats  = "Backtracking translations: %d, newPos: %v"
	CannotBacktrack    = "Cannot backtrack. Board is invalid:\n%s"
	FailedToInitCells  = "Failed initializing cells"

	// sudoku.go
	CallingBacktracking = "Calling bactracking, isValid: %v, board:\n%s"

	// container.go
	CellNotificationArrived            = "Container '%-9s' received notification from cell: '%s', value: '%s'"
	CellUpdateInvalidValue             = "Container '%-9s' received notification with an invalid cell value: '%s'"
	ContainerValueNotAvailable         = "Value '%s' not available in container '%-9s'"
	CellPrevValueInvalid               = "Cell '%s' previous value: '%s' is invalid"
	ContainerAvailableValues           = "Container '%-9s' available values updated: %v"
	ContainerObserverAlreadyRegistered = "Container observer '%s' has already been registered"
	ContainerObserverNotFound          = "The container observer id '%s' was not found"
)

var (
	LOG_FILENAME   = "logs.txt"
	INFO_HEADER    = "INFO: "
	WARNING_HEADER = "WARNING: "
	ERROR_HEADER   = "ERROR: "
)
