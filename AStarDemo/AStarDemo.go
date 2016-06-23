// A* Demo:
// 
// Todo: Plop a purple dot in lower-left corner, and then plop and green dot in upper right.
//       Implement A* search s.t. the purple dot can plot a course to the green dot around all
//       the yellow trash. Draw the course in white.
//
// TODO XXX TODO

package main

import "fmt"
import "github.com/veandco/go-sdl2/sdl"
import "github.com/jrbl/creepy/cell"
import "math/rand"
import "time"


const WINDOWSIZE = 720 // TODO: query for max y and make biggest square possible
const BBSZ = int32(6) // borderbox size: 720/6 = 120, giving 14400 box grid
const CBSZ = int32(5) // cellbox size: <= BBSZ
const BORDER = int32(1) // normally BBSZ-CBSZ, or (BBSZ-CBSZ)*0.5


func XYtoI(x int, y int, rowSize int) int {
    // Convert 2d positions to 1d slice index
    return y * rowSize + x
}

func PlotRoute(board []bool, creep int, goal int) cell.Heap {
//func PlotRoute(board []bool, creep int, goal int) []bool {
    setSize := len(board)
    startPos := setSize - 1
    goalPos := 0
    rank := int(WINDOWSIZE/BBSZ)
    goalCell := cell.Cell{I: goalPos}  // placeholder, not entered into heap until we find it

    // OPEN = priority queue containing START
    openSet := cell.NewHeap()
    startCell := cell.Cell{I: startPos}
    startCell.H = startCell.FudgeTaxiDistance(goalCell, startCell, rank)
    openSet.Push(&startCell)

    // CLOSED = empty set
    closedSet := cell.NewHeap() // cells already evaluated (starts empty)

    current := startCell
    for current.I != goalCell.I {             // while lowest rank in OPEN is not the GOAL:
        current := openSet.Pop().(*cell.Cell) // current = remove lowest rank item from OPEN
        closedSet.Push(&current)              //  add current to CLOSED
        neighbors := current.Neighbors(rank)
        fmt.Printf("%d: ", current.I) // XXX
        for pos := range neighbors {          //  for neighbors of current:
            fmt.Printf("%d ", pos) // XXX
        }
        fmt.Printf("\n") // XXX
    }

//    cost = g(current) + movementcost(current, neighbor)
//    if neighbor in OPEN and cost less than g(neighbor):
//      remove neighbor from OPEN, because new path is better
//    if neighbor in CLOSED and cost less than g(neighbor): ⁽²⁾
//      remove neighbor from CLOSED
//    if neighbor not in OPEN and neighbor not in CLOSED:
//      set g(neighbor) to cost
//      add neighbor to OPEN
//      set priority queue rank to g(neighbor) + h(neighbor)
//      set neighbor's parent to current

/*reconstruct reverse path from goal to start
by following parent pointers
*/

    fmt.Printf("[%d] i: %d, h: %.5f, idx: %d\n", -1, current.I, current.H, current.Idx)
    if len(openSet) == 0 { return closedSet } else { return openSet }  // XXX placeholder to make it compile
}


func paintBoxes(surface *sdl.Surface, board []bool, creep int, goal int, creep_path []uint) {
    var x int32
    y := int32(0)
    boxcounter := 0

    board[creep] = false
    board[goal] = false
    if len(creep_path) > 30000 { // impossible condition to avoid errors from unused creep_path XXX
        board[creep] = false
    } // XXX

    for y+BBSZ <= WINDOWSIZE {
        x = int32(0)
        for ; x+BBSZ <= WINDOWSIZE; {
            rect := sdl.Rect{x, y, BBSZ, BBSZ}
            surface.FillRect(&rect, 0x00000000)

            // TODO: take these inner rect objects and put them onto a data structure so we
            //       can do things with them later. Like A* search, or G.O.L.
            //       In the meantime, consult externally defined data structure for liveness
            rect = sdl.Rect{x+BORDER, y+BORDER, CBSZ, CBSZ}
            // TODO: These colors are fine to work with, but they don't look anything like I expect. Look at the implementation, and at 
            //        http://tools.medialab.sciences-po.fr/iwanthue/ and figure out why the creep's not blue, the goal isn't green, and 
            //        the walls aren't gold.
            if boxcounter == creep {
                surface.FillRect(&rect, 0x3c74fe00)
            } else if boxcounter == goal {
                surface.FillRect(&rect, 0x2edfb400)
            } else if board[boxcounter] {
                surface.FillRect(&rect, 0xb2af0000)
            } else {
                surface.FillRect(&rect, 0x1f1f1f00)
            }

            x += BBSZ
            boxcounter += 1
        }
        y += BBSZ
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
            life_cell := int32(XYtoI(int(col), int(row), int(edge_sz)))
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
 * Implement A* search
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

    // Game Setup
    rand.Seed(time.Now().Unix()) // insecure seed is sufficient in this case
    field_edge := WINDOWSIZE/BBSZ
    field_size := field_edge * field_edge
    life_board := make([]bool, field_size)
    life_board = RandomizeBoard(life_board, 15)

    creep := int(field_size-1) // TODO: tunable creep entry?
    goal := 0 // TODO: tunable creep exit?
    creep_path := make([]uint, field_size)

    // Draw stuff and deal with events
    // TODO: running states -1,0,1 should probably be some kind of enum. XXX TODO: research enums
    running = 1
    for running != -1 {
        running = checkEvents(running, life_board)
        _ = PlotRoute(life_board, creep, goal)
        paintBoxes(surface, life_board, creep, goal, creep_path)
        window.UpdateSurface()

        //sdl.Delay(200) // 1/5 second
        //if running == 1 {
        //    life_board = CalcNextBoard(life_board)
        //}
    }

}
