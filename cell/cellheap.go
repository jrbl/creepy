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
// TODO: need test for that also

func NewHeap() CHeap {
    ch := make(CHeap, 0)
    heap.Init(&ch)
    return ch
}
