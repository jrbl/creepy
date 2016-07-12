package main

import "fmt"
import "github.com/jrbl/creepy/cell"
import "testing"


func setUp(size int) (int, []bool) {
    fmt.Println("Setting up", size)
    return size, make([]bool, size*size)
}

// [.] and [*]
func TestMoveCost1x1(t *testing.T) {
    rank, board := setUp(1)
    goal := cell.Cell{I: 0}
    pos := cell.Cell{I: 0}
    start := cell.Cell{I: 0}
    expected := 0.000
    if x := MoveCost(pos, goal, start, rank, board); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", pos.I, goal.I, x, expected)
    }
    board[0] = true
    expected = 9e9
    if x := MoveCost(pos, goal, start, rank, board); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", pos.I, goal.I, x, expected)
    }
}

// [g.] and [*.] and [g*] and [s.]
// [.s]     [.s]     [.s]     [.g]
func TestMoveCost2x2(t *testing.T) {
    rank, board := setUp(2)
    goal := cell.Cell{I: 0}
    pos := cell.Cell{I: 3}
    start := cell.Cell{I: 3}
    expected := 1.000
    if x := MoveCost(pos, goal, start, rank, board); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", pos.I, goal.I, x, expected)
    }
    board[0] = true
    expected = 9e9
    if x := MoveCost(pos, goal, start, rank, board); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", pos.I, goal.I, x, expected)
    }
}

func TestMoveCost3x3(t *testing.T) {
    testBoardRank := 3
    testBoard := make([]bool, testBoardRank*testBoardRank)
    goal := cell.Cell{I: 0}
    pos := cell.Cell{I: 4}
    start := cell.Cell{I: 7} // XXX: WTF?
    expected := 1.000
    if x := MoveCost(pos, start, goal, testBoardRank, testBoard); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", pos.I, goal.I, x, expected)
    }
    pos.I = 8
    expected = 2.000
    if x := MoveCost(pos, start, goal, testBoardRank, testBoard); x != expected {
        t.Errorf("Cell %d and %d close, but x is %.4f not %.4f", pos.I, goal.I, x, expected)
    }
//    expected = 6.000
//    if x := a.FudgeTaxiDistance(b, a, testBoardRank); x != expected {
//        t.Errorf("Cell %d and %d same row, but x is %.4f, not %.4f", a.I, c.I, x, expected)
//    }
//    expected = 7.006
//    if x := b.FudgeTaxiDistance(c, a, testBoardRank); x != expected {
//        t.Errorf("Cell %d and %d expected to have adjustment, but x is %.4f not %.4f", b.I, c.I, x, expected)
//    }
    // TODO(jrbl): establish 3x3 grid, do real line-of-site comparison for various
    //             straight-line paths. (Work them out by hand first.)
}

/*func TestFudgeTaxiDistance(t *testing.T) {
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
} */

/* func TestInitializeBoard(t *testing.T) {
	// TODO: is it good practice to make InitializeBoard return error r/t log and die directly,
	//       so that we can unit test its boundary conditions? What about logging and returning
	//       an error? What's idiomatic?
	test_board := InitializeBoard(3)
	if size := len(test_board); size != 3 {
		t.Errorf("Expected test board to be two rows tall, got %d instead.", size)
	}
}

func Test_get_neighbor_coords(t *testing.T) {
	test_board := InitializeBoard(3)
	var neighbor_coords [8][2]int
	// top left corner
	neighbor_coords = get_neighbor_coords(0, 0, test_board)
	if neighbor_coords != [8][2]int{{2, 2}, {2, 0}, {2, 1}, {0, 2}, {0, 1}, {1, 2}, {1, 0}, {1, 1}} {
		t.Errorf("Unexpected neighbor coordinate set for 0,0")
		fmt.Printf("%v", neighbor_coords)
	}
	// center
	neighbor_coords = get_neighbor_coords(1, 1, test_board)
	if neighbor_coords != [8][2]int{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}} {
		t.Errorf("Unexpected neighbor coordinate set for 1,1")
		fmt.Printf("%v", neighbor_coords)
	}
	// bottom right corner
	neighbor_coords = get_neighbor_coords(2, 2, test_board)
	if neighbor_coords != [8][2]int{{1, 1}, {1, 2}, {1, 0}, {2, 1}, {2, 0}, {0, 1}, {0, 2}, {0, 0}} {
		t.Errorf("Unexpected neighbor coordinate set for 2,2")
		fmt.Printf("%v", neighbor_coords)
	}
}

func Test_count_neighbors(t *testing.T) {
	test_board := InitializeBoard(6)
	test_board[0][0] = 1
	test_board[2][2] = 1
	if c := count_neighbors(0, 0, test_board); c != 0 {
		t.Errorf("Incorrect live neighbor count for 0, 0: %d", c)
	}
	if c := count_neighbors(0, 1, test_board); c != 1 {
		t.Errorf("Incorrect live neighbor count for 0, 1: %d", c)
	}
	if c := count_neighbors(1, 0, test_board); c != 1 {
		t.Errorf("Incorrect live neighbor count for 1, 0: %d", c)
	}
	if c := count_neighbors(1, 1, test_board); c != 2 {
		t.Errorf("Incorrect live neighbor count for 1, 1: %d", c)
	}
}
*/
