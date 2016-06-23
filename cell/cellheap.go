package cell

import "container/heap"

// MinHeap of Cell-type
type Heap []*Cell

func (h Heap) Len() int { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].H < h[j].H }
func (h Heap) Swap(i, j int) { 
    h[i], h[j] = h[j], h[i] 
    h[i].Idx = i
    h[j].Idx = j
}
func (h *Heap) Push(x interface{}) { 
    n := len(*h)
    c := x.(*Cell)
    c.Idx = n
    *h = append(*h, c)
}
func (h *Heap) Pop() interface{} { 
    old := *h
    n := len(old)
    c := old[n-1]
    c.Idx = -1 // for safety
    *h = old[0 : n-1]
    return c
}
func (h *Heap) update(c *Cell, newi int, newh float64) {
    c.I = newi
    c.H = newh
    heap.Fix(h, c.Idx)
}

func NewHeap() Heap {
    ch := make(Heap, 0)
    heap.Init(&ch)
    return ch
}
