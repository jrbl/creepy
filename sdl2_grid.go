package main

// import "fmt"
import "github.com/veandco/go-sdl2/sdl"
import "math/rand"
import "time"


const WINDOWSIZE = 720 // TODO: query for max y and make biggest square possible
const BBSZ = int32(6) // borderbox size: 720/6 = 120, giving 14400 box grid
const CBSZ = int32(5) // cellbox size: <= BBSZ
const BORDER = int32(1) // normally BBSZ-CBSZ, or (BBSZ-CBSZ)*0.5


func paintboxes(surface *sdl.Surface) {
    var x int32
    y := int32(0)
    boxcounter := 0
    bbcount := int((WINDOWSIZE/BBSZ) * (WINDOWSIZE/BBSZ))

    // TODO: use this as display for Game of Life. In the meantime, here's a lit box
    rand.Seed(time.Now().Unix()) // time-seeded randomness is insecure, but sufficient
    livebox := rand.Intn(bbcount) // just one of them

    for ; y+BBSZ <= WINDOWSIZE; {
        x = int32(0)
        for ; x+BBSZ <= WINDOWSIZE; {
            rect := sdl.Rect{x, y, BBSZ, BBSZ}
            surface.FillRect(&rect, 0x00000000)

            // TODO: take these inner rect objects and put them onto a data structure so we
            //       can do things with them later. Like A* search, or G.O.L.
            //       In the meantime, paint a randomly selected box yellow.
            rect = sdl.Rect{x+BORDER, y+BORDER, CBSZ, CBSZ}
            if boxcounter == livebox {
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

func checkevents() bool {
    var event sdl.Event
    running := true

    for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        // switch t := event.(type) {
        switch event.(type) {
        case *sdl.QuitEvent:
            running = false
        //case *sdl.MouseMotionEvent:
        //    fmt.Printf("[%d ms] MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
        //        t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
        }
    }
    return running
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
    // SDL Boilerplate End

    // Draw stuff and deal with events
    background := sdl.Rect{0, 0, WINDOWSIZE, WINDOWSIZE}
    surface.FillRect(&background, 0x00000000)

    for running = checkevents(); running; running = checkevents() {
        paintboxes(surface)
        window.UpdateSurface()

        sdl.Delay(200) // 1/5 second
    }

}
