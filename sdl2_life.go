package main

import "fmt"
import "github.com/veandco/go-sdl2/sdl"
import "math"
import "math/rand"
import "time"


const WINDOWSIZE = 720 // TODO: query for max y and make biggest square possible
const BBSZ = int32(6) // borderbox size: 720/6 = 120, giving 14400 box grid
const CBSZ = int32(5) // cellbox size: <= BBSZ
const BORDER = int32(1) // normally BBSZ-CBSZ, or (BBSZ-CBSZ)*0.5


func paintboxes(surface *sdl.Surface, board []bool) {
    var x int32
    y := int32(0)
    boxcounter := 0

    for ; y+BBSZ <= WINDOWSIZE; {
        x = int32(0)
        for ; x+BBSZ <= WINDOWSIZE; {
            rect := sdl.Rect{x, y, BBSZ, BBSZ}
            surface.FillRect(&rect, 0x00000000)

            // TODO: take these inner rect objects and put them onto a data structure so we
            //       can do things with them later. Like A* search, or G.O.L.
            //       In the meantime, consult externally defined data structure for liveness
            rect = sdl.Rect{x+BORDER, y+BORDER, CBSZ, CBSZ}
            if board[boxcounter] {
                surface.FillRect(&rect, 0xefefef00)
            } else {
                surface.FillRect(&rect, 0x1f1f1f00)
            }

            x += BBSZ
            boxcounter += 1
        }
        y += BBSZ
    }

}

func checkevents(running bool, board []bool) bool {
    mousebutton := false

    var event sdl.Event
    for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        /* switch event.(type) {
        case *sdl.QuitEvent:
            running = false
        */
        switch t := event.(type) {
        case *sdl.QuitEvent:
            running = false
        case *sdl.MouseButtonEvent:
            if t.Button == sdl.BUTTON_LEFT {
                if t.State == sdl.PRESSED {
                    mousebutton = true
                }
                if t.State == sdl.RELEASED {
                    mousebutton = false
                }
            }
        case *sdl.MouseMotionEvent:
            col := t.X / BBSZ
            row := t.Y / BBSZ
            edge_sz := WINDOWSIZE/BBSZ
            life_cell := row*edge_sz+col
            if mousebutton {
                fmt.Printf("clickdrag on %d\n", life_cell)
            }
            if col > 1000000 {    // XXX
                fmt.Printf("%v\n", board) // XXX
            } // XXX
            // more cooked debug output, actually useful
            //fmt.Printf("[%dms] MouseMotion (%d,%d) = box %d\n", t.Timestamp, row, col, row*edge_sz+col)
            // real raw debug output
            //fmt.Printf("[%dms] MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", 
            //    t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
        }
    }
    return running
}

func RandomizeBoard(board []bool) []bool {
    for i := 0; i < len(board); i += 1 {
        board[i] = (rand.Intn(2) == 1)
    }
    return board
}

func get_neighbor_indexes(i int, board []bool) [8]int {
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
    neighbor_indexes := get_neighbor_indexes(i, board)
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
    for i := 0; i < len(board); i += 1 {
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

func main() {
    var window *sdl.Window
    var surface *sdl.Surface
    var running bool
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

    // Game of Life Setup
    rand.Seed(time.Now().Unix()) // insecure seed sufficient in this case
    field_edge := WINDOWSIZE/BBSZ
    field_size := field_edge * field_edge
    life_board := make([]bool, field_size)
    life_board = RandomizeBoard(life_board)

    // Draw stuff and deal with events
    running = true
    for running  {
        running = checkevents(running, life_board)
        paintboxes(surface, life_board)
        window.UpdateSurface()

        sdl.Delay(200) // 1/5 second

        life_board = CalcNextBoard(life_board)
    }

}
