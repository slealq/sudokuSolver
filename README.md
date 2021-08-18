# sudokuSolver
 A sudoku solver written in Golang, for fun

## Table of contents

---
- [sudokuSolver](#sudokusolver)
  - [Table of contents](#table-of-contents)
  - [Design considerations](#design-considerations)
    - [Deterministic approach](#deterministic-approach)
    - [Backtracking](#backtracking)
  - [Components](#components)
    - [Board representation](#board-representation)
  - [Responsabilities](#responsabilities)
  - [Testing](#testing)


## Design considerations

---

### Deterministic approach

Sudokos might have a deterministic approach to be solved (which is the fastest) 
when for example, there's always a constraint that obligates a particular cell 
to have only one value. 

Example:

    0 1 3  4 5 6  7 8 _
 
By definition of sudoku, the missing value in the above would be a `9`. This
case is relatively straight-forward, as one single `container` (the row) has
all information required to force the value.

Other cases might be more difficult, as every cell belongs to three containers,
always. The intersection of the possible values from the three containers 
determines the possible values for the cell.

Finally, if the intersection only has one possible value, conformtably we could
place that value into the cell. This is the so called `deterministic approach`.

### Backtracking

Although most `easy` and `medium` sudokus can be solved by using the 
`deterministic approach`, not all of them can. When there's no deterministic
solution for any of the cells in the board, `backtracking` must be used to
find the right answer.

This solution will be used as `last-resort`, since any particular board could
have deterministic values up-to some point, where backtracking is required
to finish solving the puzzle. 

And it's assumed that the reverse is also possible. That some level of 
backtracking might be required, until a point where deterministic solutions
can be found. Either way, deterministic solutions are much faster, and so
preferred.

## Components

### Board representation
There are three main components:
- Board
- Containers
- Cells

Whose relationships are as follows:

                        ,-----.                  
                        |Board|                  
                        |-----|                  
                        `-----'                  
                            |                     
                            |                     
        ,-----------------------------------------.
        |Containers: Vertical, Horizontal, Squares|
        |-----------------------------------------|
        `-----------------------------------------'
                            |                     
                        ,-----.                  
                        |Cells|                  
                        |-----|                  
                        `-----'                  


The board is in charge to managing all containers. For each row, there is one 
container. For each column the same, and for each box the same. Each container
has a total of 9 cells. And there are 9 of each type of container.

In total there are 27 containers (9+9+9), which are responsible for yielding
possible values for each cell. 

Cells are shared across containers, so each container doesn't have triplicate
memory.

## Responsabilities

The board is responsible of:

 - Managing containers. Every update of any cell is notified to the board,
which in turn delegates the update to the proper cell.
 - Yielding possible values for each cell.
 - Determining if current board is valid or not (A valid board has no rule
violation in place, although might not be complete).

Containers are responsible of:
 - Managing possible values of each container. The recalculation of possible
values should be triggered by an observer of each cell.

Cells are responsible of:
 - Storing the value of each cell. An observer observes this value, and notifies
containers of possible changes in the values.


## Testing

UTs 