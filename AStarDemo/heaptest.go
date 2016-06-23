package main

import "fmt"
import "container/heap"
import "math/rand"


type Cell struct {
    i int
    h float64
    idx int // the position in our priority heap
}

// MinHeap of Cell-type
type CellHeap []*Cell

func (h CellHeap) Len() int { return len(h) }
func (h CellHeap) Less(i, j int) bool { return h[i].h < h[j].h }
func (h CellHeap) Swap(i, j int) { 
    h[i], h[j] = h[j], h[i] 
    h[i].idx = i
    h[j].idx = j
}
func (h *CellHeap) Push(x interface{}) { 
    n := len(*h)
    cell := x.(*Cell)
    cell.idx = n
    *h = append(*h, cell)
}
func (h *CellHeap) Pop() interface{} { 
    old := *h
    n := len(old)
    cell := old[n-1]
    cell.idx = -1 // for safety
    *h = old[0 : n-1]
    return cell
}
func (h *CellHeap) update(c *Cell, newi int, newh float64) {
    c.i = newi
    c.h = newh
    heap.Fix(h, c.idx)
}

func main() {
    fmt.Printf("Hello, World!\n")

    // Create a minheap, put the items in it, and
	// establish the invariants
	//ch := make(CellHeap, rand.Intn(44))
	ch := make(CellHeap, 0)
    heap.Init(&ch)
    for j := 0; j < 9; j += 1 {
        cell := Cell{i: j, h: rand.Float64()}
        heap.Push(&ch, &cell)
    }
    
    // debug printout, raw heap refs
    for j := 0; j < 9; j += 1 {
        fmt.Printf("[%d] i: %d, h: %.5f, idx: %d\n", j, ch[j].i, ch[j].h, ch[j].idx)
    }

    fmt.Printf("----\n")

    // pretty print (destructive)
    counter := 0
    for ch.Len() > 0 {
        cell := heap.Pop(&ch).(*Cell)
        fmt.Printf("[%d] i: %d, h: %.5f, idx: %d\n", counter, cell.i, cell.h, cell.idx)
        counter += 1
    }
}
