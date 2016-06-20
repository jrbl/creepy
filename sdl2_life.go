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


func paintBoxes(surface *sdl.Surface, board []bool) {
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

func lightUpBoxIn(box int, board []bool) {
    board[box] = true
}

func checkEvents(running int, board []bool) int {
    var event sdl.Event
    for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        // TODO: Wouldn't be cool if instead of these printf's we could output to a transparent overlay, quake-console-style?
        switch t := event.(type) {
        case *sdl.QuitEvent:
            running = -1 // code -1 is special, says to die.
        /*case *sdl.MouseButtonEvent:
            fmt.Printf("MouseButtonEvent (%d,%d), %d, %d\n", t.X, t.Y, t.Button, t.State)
            // TODO: Are individual clicks like this interesting? maybe if I implement button controls?
            if t.Button == sdl.BUTTON_LEFT {
                if t.State == sdl.PRESSED {
                    mousebutton = true
                    fmt.Printf("Click!")
                }
                if t.State == sdl.RELEASED {
                    mousebutton = false
                    fmt.Printf("Release!")
                }
            }*/
        case *sdl.MouseMotionEvent:
            col := t.X / BBSZ
            row := t.Y / BBSZ
            edge_sz := WINDOWSIZE/BBSZ
            life_cell := row*edge_sz+col
            lastbox := int32(-1)
            thisbox := int32(-1)
            if t.State & sdl.Button(sdl.ButtonLMask()) == 1 {
                // fmt.Printf("left drag (%d,%d) = #%d // %dx%d\n", row, col, life_cell, t.X, t.Y)
                thisbox = life_cell
                if thisbox != lastbox {
                    lightUpBoxIn(int(life_cell), board)
                    lastbox = thisbox
                }
            }
        /* case *sdl.KeyDownEvent:
            fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
                t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat) */
        case *sdl.KeyUpEvent:
            // Space toggles run state between start and stop
            if t.Keysym.Sym == ' ' && t.Keysym.Mod == 0 {
                if running == 0 { running = 1 } else { running = 0 }  // toggle run state 1/0
            }
            //fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
            //    t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
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

func neighborIndices(i int, board []bool) [8]int {
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

    // Game of Life Setup
    rand.Seed(time.Now().Unix()) // insecure seed is sufficient in this case
    field_edge := WINDOWSIZE/BBSZ
    field_size := field_edge * field_edge
    life_board := make([]bool, field_size)
    life_board = RandomizeBoard(life_board)

    // Draw stuff and deal with events
    // TODO: running states -1,0,1 should probably be some kind of enum. XXX TODO: research enums
    running = 1
    for running != -1 {
        running = checkEvents(running, life_board)
        paintBoxes(surface, life_board)
        window.UpdateSurface()

        sdl.Delay(200) // 1/5 second
        if running == 1 {
            life_board = CalcNextBoard(life_board)
        }
    }

}
