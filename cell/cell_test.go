package cell

//import "fmt"
import "testing"

func TestXY(t *testing.T) {
    testBoardRank := 7
    testCell := Cell{I: 0, H: 0.25}
    if x, y := testCell.XY(testBoardRank); x != 0 && y != 0 {
        t.Errorf("Expected position of cell 0 on rank 7 board is (0,0), not (%d,%d)", x, y)
    }
    testCell.I = 6
    if x, y := testCell.XY(testBoardRank); x != 0 && y != 6 {
        t.Errorf("Expected position of cell 6 on rank 7 board is (0,6), not (%d,%d)", x, y)
    }
    testCell.I = 7
    if x, y := testCell.XY(testBoardRank); x != 1 && y != 0 {
        t.Errorf("Expected position of cell 7 on rank 7 board is (1,0), not (%d,%d)", x, y)
    }
}

func TestTaxiDistance1x1(t *testing.T) {
    testBoardRank := 1
    a := Cell{I: 0}
    b := Cell{I: 0}
    if x := a.TaxiDistance(b, testBoardRank); x != 0.0 {
        t.Errorf("Degenerate board 1x1 calculating TaxiDistance incorrectly. %.4f != 0", x)
    }
}

func TestTaxiDistance2x2(t *testing.T) {
    testBoardRank := 2
    a := Cell{I: 0}
    b := Cell{I: 3}
    if x := a.TaxiDistance(b, testBoardRank); x != 2.0 {
        t.Errorf("Degenerate board 2x2 calculating TaxiDistance incorrectly. %.4f != 2", x)
    }
    a.I = 1
    if x := a.TaxiDistance(b, testBoardRank); x != 1.0 {
        t.Errorf("Degenerate board 2x2 calculating TaxiDistance incorrectly. %.4f != 1", x)
    }
    a.I = 2
    if x := a.TaxiDistance(b, testBoardRank); x != 1.0 {
        t.Errorf("Degenerate board 2x2 calculating TaxiDistance incorrectly. %.4f != 1", x)
    }
}

func TestTaxiDistance(t *testing.T) {
    testBoardRank := 7
    a := Cell{I: 0}
    b := Cell{I: 6}
    c := Cell{I: 7}
    if x := a.TaxiDistance(c, testBoardRank); x != 1.0 {
        t.Errorf("Cell %d and %d expected to be north-south adjacent, but x is %.4f", a.I, b.I, x)
    }
    if x := a.TaxiDistance(b, testBoardRank); x != 6.0 {
        t.Errorf("Cell %d and %d expected either end of same row, but x is %.4f", a.I, c.I, x)
    }
    if x := b.TaxiDistance(c, testBoardRank); x != 7.0 {
        t.Errorf("Cell %d and %d expected to adject rows furthes columns, but x is %.4f", b.I, c.I, x)
    }
}

func TestFudgeTaxiDistance(t *testing.T) {
    testBoardRank := 7
    a := Cell{I: 0}
    b := Cell{I: 6}
    c := Cell{I: 7}
    expected := 1.000
    if x := a.FudgeTaxiDistance(c, a, testBoardRank); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", a.I, b.I, x, expected)
    }
    expected = 6.000
    if x := a.FudgeTaxiDistance(b, a, testBoardRank); x != expected {
        t.Errorf("Cell %d and %d same row, but x is %.4f, not %.4f", a.I, c.I, x, expected)
    }
    expected = 7.006
    if x := b.FudgeTaxiDistance(c, a, testBoardRank); x != expected {
        t.Errorf("Cell %d and %d expected to have adjustment, but x is %.4f not %.4f", b.I, c.I, x, expected)
    }
    // TODO(jrbl): establish 3x3 grid, do real line-of-site comparison for various
    //             straight-line paths. (Work them out by hand first.)
}

func TestNeighbors(t *testing.T) {
    tl := Cell{I: 0} // tl is top left of 3x3 grid
    c := Cell{I: 4}  // c is center of 3x3 grid
    br := Cell{I: 8} // br is bottom right of 3x3 grid
    rank := 3
    ns := c.Neighbors(rank)
    expect := [8]int{0, 1, 2, 3, 5, 6, 7, 8}
    if ns != expect {
        t.Errorf("Incorrect neighbor array %v for cell %d, expected %v", ns, c.I, expect)
    }
    ns = tl.Neighbors(rank)
    expect = [8]int{-1, -1, -1, -1, 1, -1, 3, 4}
    if ns != expect {
        t.Errorf("Incorrect neighbor array %v for cell %d, expected %v", ns, tl.I, expect)
    }
    ns = br.Neighbors(rank)
    expect = [8]int{4, 5, -1, 7, -1, -1, -1, -1}
    if ns != expect {
        t.Errorf("Incorrect neighbor array %v for cell %d, expected %v", ns, br.I, expect)
    }
}

func TestXYtoI(t *testing.T) {
    exI, x, y := 0, 0, 0 // top left
    if i := XYtoI(x, y, 3); i != exI {
        t.Errorf("Bad index calculation, (%d, %d)->%d, not %d", x, y, i, exI)
    }
    exI, x, y = 8, 2, 2 // bottom right
    if i := XYtoI(x, y, 3); i != exI {
        t.Errorf("Bad index calculation, (%d, %d)->%d, not %d", x, y, i, exI)
    }
    exI, x, y = 6, 0, 2 // bottom left
    if i := XYtoI(x, y, 3); i != exI {
        t.Errorf("Bad index calculation, (%d, %d)->%d, not %d", x, y, i, exI)
    }
    exI, x, y = 2, 2, 0 // top right
    if i := XYtoI(x, y, 3); i != exI {
        t.Errorf("Bad index calculation, (%d, %d)->%d, not %d", x, y, i, exI)
    }
}
