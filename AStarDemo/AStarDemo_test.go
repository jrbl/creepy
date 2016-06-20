package main

import "fmt"
import "testing"

/* func TestItoXY(t *testing.T) {
    if x, y := ItoXY(8, 3); x != 2 || y != 2 {
		t.Errorf("Bad pair calculation, %d->(%d,%d), not (%d,%d)",
            8, x, y, 2, 2)
    }
} */

func TestManhattanDistance(t *testing.T) {
    // grid of 9 squares, bottom right to top left
    // best distance is 4: straight to an opposing corner, then across to the goal.
    if distance := ManhattanDistance(2, 0, 2, 0); distance != 4 {
		t.Errorf("Incorrect Manhattan distance calculation on 3x3 grid, %d != 4", distance)
    }
    if distance := ManhattanDistance(2, 1, 2, 0); distance != 3 {
		t.Errorf("Incorrect Manhattan distance calculation on 3x3 grid, %d != 3", distance)
    }
    if distance := ManhattanDistance(2, 0, 0, 0); distance != 2 {
		t.Errorf("Incorrect Manhattan distance calculation on 3x3 grid, %d != 2", distance)
    }
}

func TestInformedManhattanDistance(t *testing.T) {
    // grid of 9 squares, buttom right to top left
    if distance := InformedManhattanDistance(8, 0, 8, 3); distance != 4.008 {
        // base position is the starting square
		t.Errorf("Incorrect Informed Manhattan distance calculation on 3x3 grid, %d != 4.008", distance)
    }
    if distance := InformedManhattanDistance(5, 0, 8, 3); distance != 3.006 {
        // 1 north of starting square
		t.Errorf("Incorrect Informed Manhattan distance calculation on 3x3 grid, %d != 3.006", distance)
    }
    if distance := InformedManhattanDistance(7, 0, 8, 3); distance != 3.006 {
        // 1 west of starting square
		t.Errorf("Incorrect Informed Manhattan distance calculation on 3x3 grid, %d != 3.006", distance)
    }
    if distance := InformedManhattanDistance(2, 0, 8, 3); distance != 2.004 {
        // 1 west of starting square
		t.Errorf("Incorrect Informed Manhattan distance calculation on 3x3 grid, %d != 2.004", distance)
    } 
} 

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
