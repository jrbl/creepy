// Package cell gives definition of grid cells (Cell) and common operations to
// perform upon them, as well as associated types, such as cell.Heap.
package cell


import "math"


type Cell struct {
    I int // Cell's position in a 1d array used to calculate cell's contents
    H float64 // Cell's heuristic value, used for informed search (e.g., ManhattanDistance for A* search)
    Idx int // Cell's position in a cell.Heap used to prioritize for informed search
}

// XY converts Cell.I 1d slice index into 2d positions, with (0, 0) on top left of grid
func (c *Cell) XY(rank int) (int, int) {
    col := c.I % rank
    row := c.I / rank
    return row, col
}

// TaxiDistance calculates the Taxicab distance (Manhattan Distance) between two Cells.
func (c *Cell)  TaxiDistance(goal Cell, rank int) float64 {
    goalX, goalY := goal.XY(rank)
    posX, posY := c.XY(rank)
    dX := math.Abs(float64(posX - goalX))
    dY := math.Abs(float64(posY - goalY))
    // D := 1 // D is an edge weighting
    // return D * (dX + dY)
    return dX + dY
}

// TaxiDistance, plus a tiny fudge-factor calculated from deviation from straight line path
func (c *Cell) FudgeTaxiDistance(goal Cell, start Cell, rank int) float64 {
    goalX, goalY := goal.XY(rank)
    startX, startY := start.XY(rank)
    posX, posY := c.XY(rank)
    heuristic := c.TaxiDistance(goal, rank)
    dx1 := posX - goalX
    dy1 := posY - goalY
    dx2 := startX - goalX
    dy2 := startY - goalY
    crossProduct := math.Abs(float64(dx1*dy2) + float64(dy1*dx2))
    return heuristic + crossProduct*0.001
}

// Get a list of the board indices of this Cell; return -1 for indices that fall off the edge of the board
func (c *Cell) Neighbors(rank int) [8]int {
    n := 0
    neighbors := [8]int{}
    row, col := c.XY(rank)
    for i := -1; i <= 1; i += 1 {
        for j := -1; j <= 1; j += 1 {
            if i == 0 && j == 0 { continue }
            target_row := row + i
            target_col := col + j
            if target_row < 0 { target_row = -1 }              // stop at edge
            if target_col < 0 { target_col = -1 }              // stop at edge
            if target_row >= rank { target_row = -1 }          // stop at edge
            if target_col >= rank { target_col = -1 }          // stop at edge
            neighbors[n] = XYtoI(target_col, target_row, rank)
            n += 1
        }
    }
    return neighbors
}

// Convert 2d positions to 1d slice index
// if either input index is -1, return -1
func XYtoI(x int, y int, rank int) int {
    if x == -1 || y == -1 { return -1 } else { return y * rank + x }
}

