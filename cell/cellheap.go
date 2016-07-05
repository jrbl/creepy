package cell

import "container/heap"


// MinHeap of Cell-type
type CHeap []*Cell

func (h CHeap) Len() int { return len(h) }
func (h CHeap) Less(i, j int) bool { return h[i].H < h[j].H }
func (h CHeap) Swap(i, j int) { 
    h[i], h[j] = h[j], h[i] 
    h[i].Idx = i
    h[j].Idx = j
}
func (h *CHeap) Push(x interface{}) { 
    n := len(*h)
    c := x.(*Cell)
    c.Idx = n
    *h = append(*h, c)
}
func (h *CHeap) Pop() interface{} { 
    old := *h
    n := len(old)
    c := old[n-1]
    c.Idx = -1 // for safety
    *h = old[0 : n-1]
    return c
}
func (h *CHeap) update(c *Cell, newi int, newh float64) {
    c.I = newi
    c.H = newh
    heap.Fix(h, c.Idx)
}


func (h *CHeap) Add(c *Cell) {
    heap.Push(h, c)
}

func (h *CHeap) Remove() *Cell {
    return heap.Pop(h).(*Cell)
}

func (h *CHeap) Revalue(c *Cell, new_h float64) {
    h.update(c, c.I, new_h)
}

// Unlink uses heap implementation to remove Cell c from CHeap h in log2*N time.
// Returns the cell requested
func (h *CHeap) Unlink(c *Cell) *Cell {
    hValue := c.H
    h.Revalue(c, -99999.99999)
    cPopped := h.Remove()
    cPopped.H = hValue
    cPopped.Idx = -1 // congruent with Pop() above
    return cPopped
}

// Search linearly scans the structure of the heap, looking for the
// first instance of i to appear in some Cell.I.
// Returns j, the Cell's linear index, and the Cell.
// In a properly valued heap, j == Cell.Idx
func (h CHeap) Search(i int) (int, *Cell) {
    for j := 0; j < len(h); j += 1 {
        if h[j].I == i {
            return j, h[j]
        }
    }
    return -1, nil
}

func NewHeap() CHeap {
    ch := make(CHeap, 0)
    heap.Init(&ch)
    return ch
}
