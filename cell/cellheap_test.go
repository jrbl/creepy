package cell

import "fmt"
import "testing"


func setUpMinHeap() CHeap {
    // Create a minheap, put 3 items in it,
	ch := NewHeap()
    ch.Add(&Cell{I: 0, H: 0.25})
    ch.Add(&Cell{I: 1, H: 0.001})
    ch.Add(&Cell{I: 2, H: 0.35})
    return ch
}

func TestHeapAdds(t *testing.T) {
    ch := setUpMinHeap()
    
    j := 0    // heap position, only differs from .Idx if heap needs a Fix()
    exp := 1  // insertion ordering recorded as the "value" field I
    if i := ch[j].I; i != exp {
        fmt.Printf("[%d] I: %d, H: %.5f, Idx: %d\n", j, ch[j].I, ch[j].H, ch[j].Idx)
        t.Errorf("Incorrect heap ordering insertion %d should be %d, is %d", j, exp, i)
    }
    j = 1
    exp = 0
    if i := ch[j].I; i != exp {
        fmt.Printf("[%d] I: %d, H: %.5f, Idx: %d\n", j, ch[j].I, ch[j].H, ch[j].Idx)
        t.Errorf("Incorrect heap ordering insertion %d should be %d, is %d", j, exp, i)
    }
    j = 2
    exp = 2
    if i := ch[j].I; i != exp {
        fmt.Printf("[%d] I: %d, H: %.5f, Idx: %d\n", j, ch[j].I, ch[j].H, ch[j].Idx)
        t.Errorf("Incorrect heap ordering insertion %d should be %d, is %d", j, exp, i)
    }
}

func TestHeapRemoves(t *testing.T) {
    ch := setUpMinHeap()
    lastH := -1.00
    for i := 0; i < len(ch); i+= 1 {
        c := ch.Remove()
        // Make sure that when we pop, we get them smallest prio to largest
        if lastH != -1.00 && lastH > c.H {
            t.Errorf("Incorrect heap ordering removal %.4f came after %.4f at position %d", c.H, lastH, i)
        }
    }
}

func TestHeapRevalue(t *testing.T) {
    ch := setUpMinHeap()
    i := int(0)
    min := ch[0].H
    max := ch[0].H
    var lastCell *Cell
    for i = 0; i < len(ch); i += 1 {
        if ch[i].H <= min {
            min = ch[i].H
        }
        if ch[i].H >= max {
            max = ch[i].H
        }
        lastCell = ch[i]
    }
    i -= 1

    // confirm 3rd item has highest value
    if lastCell.H != max {
        t.Errorf("Expected last item inserted to have highest value, but %d has %.4f, > %.4f", i, lastCell.H, max)
    }

    // update 3rd item to lowest value
    // confirm 3rd item has lowest value
    ch.Revalue(lastCell, (min * 0.5))
    for j := 0; j < len(ch); j += 1 {
        if ch[j].I == lastCell.I {  // These must be the same value from the test cell
            if lastCell.H != ch[j].H {
                t.Errorf("After update, H values between cell handle and cell in heap don't match.")
            } else {
                if j != 0 {
                    t.Errorf("By calculation, lowest H value should have lowest index, but index is %d", j)
                }
            }
        }
    }

    // pop lowest-valued item, confirm it's the value we set
    val := ch.Remove()
    if val.H != (min * 0.5) {
        t.Errorf("After Remove, got unexpected cell: Idx: %d, I: %d, H: %d", val.Idx, val.I, val.H)
    }
}
