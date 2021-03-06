// A* Demo:

package main

import "fmt"
import "github.com/veandco/go-sdl2/sdl"
import "github.com/jrbl/creepy/cell"
import "math"
import "math/rand"
import "time"


const WINDOWSIZE = 720 // TODO: query for max y and make biggest square possible
//const BBSZ = int32(6) // borderbox size: 720/6 = 120, giving 14400 box grid
//const CBSZ = int32(5) // cellbox size: <= BBSZ
const BBSZ = int32(72) // 240 = 3, 180 = 4, 144 = 5, 120 = 6, 102 = 7, 90 = 8, 80 = 9, 72 = 10, 
const CBSZ = int32(71) // cellbox size: <= BBSZ
const BORDER = int32(1) // normally BBSZ-CBSZ, or (BBSZ-CBSZ)*0.5


// Distance gives the geometric distance between two points on the coordinate plane
func distance(ax, ay, bx, by int) float64 {
    val := math.Sqrt( math.Pow(float64(bx - ax), 2) + math.Pow(float64(by - ay), 2) )
    //fmt.Printf("distance: (%d, %d) -> (%d, %d) = %.4f\n", ax, ay, bx, by, val)
    return val
}

// validateCells confirms that every cell in cs has an .I that fits in board
func validateCells(cs []cell.Cell, rank int) bool {
    val := true
    for _, c := range cs {
        if c.I < 0 || c.I >= rank*rank {
            val = false
        }
    }
    return val
}

// MoveCost uses the calculation of deviation from a straight-line path to calculate  movement cost 
// from a particular cell c to goal. Makes cells which evaluate to True on board effectively infinite in cost.
func moveCostCrossproduct(c, goal, start cell.Cell, rank int, board []bool) float64 {
    posX, posY := c.XY(rank)
    goalX, goalY := goal.XY(rank)
    startX, startY := start.XY(rank)
    dx1 := posX - goalX
    dy1 := posY - goalY
    dx2 := startX - goalX
    dy2 := startY - goalY
    crossProduct := math.Abs(float64(dx1*dy2) + float64(dy1*dx2))
    // empirical adjustment to keep score admissible
    return 0.9 * math.Sqrt(crossProduct)
}


// MoveCost uses chooses between the crossproduct and a straight-up geometric distance calculation for easy comparison.
// Panics on bad cell bounds.
// Short circuits to make cells which evaluate to True extraordinarily high cost.
func MoveCost(c, goal, start cell.Cell, rank int, board []bool) float64 {
    if !validateCells([]cell.Cell{c, goal, start}, rank)  { panic("bad cell bounds") }
    if board[c.I] { return float64(9e9) }

    selector := 1
    if selector == 0 {
        return moveCostCrossproduct(c, goal, start, rank, board)
    } else { 
        ax, ay := c.XY(rank)
        bx, by := goal.XY(rank)
        return distance(ax, ay, bx, by)
    }
}


// TODO(jrbl): unit tests, separate module
func PlotRoute(board []bool, creep int, goal int) []int {
    startPos := len(board) - 1 // hardcoded start is end of board
    rank := int(WINDOWSIZE/BBSZ)
    goalCell := cell.Cell{I: 0}  // placeholder, not entered into heap until we find it (hardcoded end is start of board)

    openSet := cell.NewHeap()
    startCell := cell.Cell{I: startPos}
    startCell.H = MoveCost(startCell, goalCell, startCell, rank, board)
    openSet.Add(&startCell)

    closedSet := cell.NewHeap() // cells already evaluated (starts empty)

    parents := make([]int, len(board)) // track how we get to each cell on our path
    parents[startPos] = -1  // special flag value so we know where to stop traversal

    fmt.Println("start h:", startCell.H) // DEBUG
    fmt.Printf("parents: %v\n", parents) // DEBUG

    moveCost := func (c cell.Cell) float64 { return MoveCost(c, goalCell, startCell, rank, board) }

    current := openSet.Remove()               // startCell, because nothing else present
    for current.I != goalCell.I {             // while lowest rank in OPEN is not the GOAL:
        closedSet.Add(current)
        neighbors := current.Neighbors(rank)
        if true { // DEBUG
            fmt.Printf("%di;%.4fh: %v o%d/c%d\n", current.I, current.H, neighbors, len(openSet), len(closedSet))
        }
        for _, pos := range neighbors {          //  for neighbors of current:
            if pos == -1 {
                continue  // Don't fall off the edge of the world.
            }
            // figure out if neighbor is in open, closed, or neither
            _, openCellref := openSet.Search(pos)
            _, closedCellref := closedSet.Search(pos)

            if openCellref != nil {  // if it's in open set, recalculate cost estimate if necessary
                costEstimate := float64(current.H) + moveCost(*openCellref)
                if true { // DEBUG
                    fmt.Printf("%d Neighbor at %d in openCellref, old cost = %.4f", current.I, openCellref.I, openCellref.H)
                }
                if openCellref.H > costEstimate {
                    fmt.Printf(", new cost = %.4f", costEstimate) // DEBUG
                    openSet.Revalue(openCellref, costEstimate)
                    parents[openCellref.I] = current.I // update parent of openCellref to current
                }
                fmt.Printf("\n") // DEBUG
            } else if closedCellref != nil { // if it's in the closed set, move it to open set if necessary
                costEstimate := float64(current.H) + moveCost(*closedCellref)
                if true { // DEBUG
                    fmt.Printf("%d Neighbor at %d in closedCellref, old cost = %.4f", current.I, closedCellref.I, closedCellref.H)
                }
                if closedCellref.H > costEstimate {
                    fmt.Printf(", new cost = %.4f", costEstimate) // DEBUG
                    closedCellref = closedSet.Unlink(closedCellref) // remove from closedSet
                    closedCellref.H = costEstimate  // update H to costEstimate, which is better
                    openSet.Add(closedCellref)  // move it back over to the openSet with new estimate
                    parents[closedCellref.I] = current.I  // update parent of closedCellref to current
                }
                fmt.Printf("\n") // DEBUG
            } else { // neither in open set or closed set, so add it as new
                pCell := cell.Cell{I: pos}
                // cost = g(current) + movementcost(current, neighbor)
                cost := float64(current.H) + moveCost(pCell)
                pCell.H = cost
                if true { // DEBUG
                    fmt.Printf("%d Neighbor at %d in neither set, new cost = %.4f\n", current.I, pos, cost)
                }
                openSet.Add(&pCell)
                parents[pCell.I] = current.I  // set parent of pCell to current
            }
        }
        if true { // DEBUG
            fmt.Printf("%di;%.4fh: %v o%d/c%d\n", current.I, current.H, neighbors, len(openSet), len(closedSet))
        }
        current = openSet.Remove()           // current = remove lowest rank item from OPEN
        fmt.Printf("parents: %v\n", parents)
    }

    /*reconstruct reverse path from goal to start by following parent pointers */
    path := make([]int, 0, len(board))
    path = append(path, 0)
    fmt.Println(parents)
    for i := 0; parents[i] != -1; i = parents[i] {
        path = append(path, parents[i])
    }
    return path
}


// TODO(jrbl): all this drawing stuff to separate module
func getPallette(s *sdl.Surface) (m map[string]uint32) {
    m = make(map[string]uint32)
    // colors by http://tools.medialab.sciences-po.fr/iwanthue/
    m["RED"] = sdl.MapRGB(s.Format, 185, 78, 69) 
    m["GREEN"] = sdl.MapRGB(s.Format, 86, 174, 108)
    m["VIOLET"] = sdl.MapRGB(s.Format, 181, 79, 144)
    m["YELLOW"] = sdl.MapRGB(s.Format, 173, 153, 60)
    m["BLUE"] = sdl.MapRGB(s.Format, 112, 102, 188)
    m["GRAY"] = 0x1f1f1f00
    return
}

func paintBoxes(surface *sdl.Surface, board []bool, creep int, goal int, creep_path []int) {

    board[creep] = false
    board[goal] = false

    rank := WINDOWSIZE/BBSZ
    type coords struct { x, y, w, h int32 }
    coordsFunc := func (i int) coords {
        return coords{(int32(i) % rank) * BBSZ, (int32(i) / rank) * BBSZ, BBSZ, BBSZ}
    }

    colors := getPallette(surface)
    colorFunc := func (i int) uint32 {
        if i == creep {
            return colors["RED"]
        } else if  i == goal {
            return colors["GREEN"]
        } else if board[i] {
            return colors["YELLOW"]
        } else {
            return colors["GRAY"]
        }
    }

    // Draw the field, start and end, and live cells
    for i := 0; i < len(board); i += 1 {
        cs := coordsFunc(i)
        /* draw the black bounding box */
        rect := sdl.Rect{cs.x, cs.y, cs.w, cs.h}
        surface.FillRect(&rect, 0x00000000)
        /* draw the cell box */
        rect = sdl.Rect{cs.x+BORDER, cs.y+BORDER, CBSZ, CBSZ}
        surface.FillRect(&rect, colorFunc(i))
    }

    // Draw the creep path
    for _, j := range creep_path {
        cs := coordsFunc(j)
        rect := sdl.Rect{cs.x+BORDER, cs.y+BORDER, CBSZ, CBSZ}
        if board[j] {
            surface.FillRect(&rect, colors["VIOLET"])
        } else {
            surface.FillRect(&rect, colors["BLUE"])
        }
    }
}

func lightUpBoxIn(box int, board []bool) {
    board[box] = true
}

func checkEvents(running int, board []bool) int {
    var event sdl.Event
    for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        switch t := event.(type) {
        case *sdl.QuitEvent:
            running = -1 // code -1 is special, says to die.
        case *sdl.MouseMotionEvent:
            col := t.X / BBSZ
            row := t.Y / BBSZ
            edge_sz := WINDOWSIZE/BBSZ
            life_cell := int32(cell.XYtoI(int(col), int(row), int(edge_sz)))
            lastbox := int32(-1)
            thisbox := int32(-1)
            if t.State & sdl.Button(sdl.ButtonLMask()) == 1 {
                thisbox = life_cell
                if thisbox != lastbox {
                    lightUpBoxIn(int(life_cell), board)
                    lastbox = thisbox
                }
            }
        }
    }
    return running
}

func RandomizeBoard(board []bool, fill int) []bool {
    if fill < 0 || fill > 100 {
        fill = 15  // FIXME: implement real error throwing here
    }
    for i := range board {
        board[i] = (rand.Intn(100) <= fill)
    }
    return board
}

/* func neighborIndices(i int, board []bool) [8]int {
    var prev_row int
    var next_row int
    var prev_col int
    var next_col int

    neighbors := [8]int{}
    field_size := len(board)
    edge_length := int(math.Sqrt(float64(field_size)))
    row := i / edge_length
    col := i % edge_length
    prev_row = row - 1
    next_row = row + 1
    if prev_row < 0 { // vertical wrap up to bottom
        prev_row = edge_length + prev_row
    }
    if next_row >= edge_length { // vertical wrap down to top
        next_row = next_row - edge_length
    }
    prev_col = col - 1
    next_col = col + 1
    if prev_col < 0 { // horizontal wrap left to right edge
        prev_col = edge_length + prev_col
    }
    if next_col >= edge_length {
        next_col = next_col - edge_length
    }
    neighbors = [8]int{
        prev_row*edge_length+prev_col, 
        prev_row*edge_length+col,
        prev_row*edge_length+next_col,
        row*edge_length+prev_col,
        row*edge_length+next_col,
        next_row*edge_length+prev_col,
        next_row*edge_length+col,
        next_row*edge_length+next_col }
    return neighbors
}

func CountNeighbors(i int, board []bool) int {
    var neighbors int
    neighbor_indexes := neighborIndices(i, board)
    for n := 0; n < 8; n += 1 {
        n_idx := neighbor_indexes[n]
        if board[n_idx] {
            neighbors += 1
        }
    }
    return neighbors
}

func CalcNextBoard(board []bool) []bool {
    new_board := make([]bool, len(board))
    for i := range board {
        neighbors := CountNeighbors(i, board)
        cell := false
        switch {
        case neighbors < 2:
            cell = false
        case neighbors == 2:
            cell = board[i]
        case neighbors == 3:
            cell = true
        case neighbors > 3:
            cell = false
        }
        new_board[i] = cell
    }
    return new_board
}

*/

/* TODO: 
 * Animate it running to prove to yourself how clever you are
 * Try it with different values to RandomizeBoard's fillchance
 * Try it with start and goal in different places, randomize them?
 * Try it with wrapped edges enabled and disabled for the creep
 * (when you want to integrate life features, try each pairing of creepwrap/nocreepwrap & cellwrap/nocellwrap
 */

func main() {
    var window *sdl.Window
    var surface *sdl.Surface
    var running int
    var err error

    // Game Setup
    rand.Seed(time.Now().Unix()) // insecure seed is sufficient in this case
    field_edge := WINDOWSIZE/BBSZ
    fmt.Println(field_edge) // DEBUG
    field_size := field_edge * field_edge
    life_board := make([]bool, field_size)
    life_board = RandomizeBoard(life_board, 25) // normal value 15

    creep := int(field_size-1) // TODO: tunable creep entry?
    goal := 0 // TODO: tunable creep exit?
    creep_path := PlotRoute(life_board, creep, goal)
    fmt.Printf("%v\n", creep_path)

    for i, itm := range life_board { // DEBUG PRINTOUT
        skip := false
        if i % int(field_edge) == 0 {
            fmt.Printf("\n")
        }
        for _, j := range creep_path { if j == i { fmt.Printf("+"); skip = true } }
        if skip { continue }
        if itm {
            fmt.Printf("#")
        } else {
            fmt.Printf(".")
        }
    } ; fmt.Printf("\n") // DEBUG PRINTOUT

    // SDL Boilerplate Start
    if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        panic(err)
    }
    defer sdl.Quit()

    window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
        WINDOWSIZE, WINDOWSIZE, sdl.WINDOW_SHOWN)
    if err != nil {
        panic(err)
    }
    defer window.Destroy()

    surface, err = window.GetSurface()
    if err != nil {
        panic(err)
    }

    // Draw stuff and deal with events
    // TODO: running states -1,0,1 should probably be some kind of enum. XXX TODO: research enums
    running = 1 
    //running = -1 // DEBUG
    for running != -1 {
        running = checkEvents(running, life_board)
        paintBoxes(surface, life_board, creep, goal, creep_path)
        window.UpdateSurface()

        //sdl.Delay(200) // 1/5 second
        //if running == 1 {
        //    life_board = CalcNextBoard(life_board)
        //}
    }

}
