package main

import (
	"fmt"
	"errors"
	"os"
)

var (
	val_string = [...]string {".", "R", "B", "G", "Y"}
	val_name = [...]string {"!!BLANK!!", "Red", "Blue", "Green", "Yellow"}
	val_dia = [...]uint8 {0, 1, 2, 3, 3}
	dir_dx = [...]int {0, 0, -1, 1} // Up, down, left, right
	dir_dy = [...]int {-1, 1, 0, 0} // Up, down, left, right
	dir_opp_str = [...]string {"down", "up", "right", "left"}
)

const (
	DIA_LIMIT = 4
	BOARD_WIDTH = 4
	BOARD_HEIGHT = 4
)

type Board struct {
	Grid [BOARD_WIDTH * BOARD_HEIGHT]uint8
}

func new_board() Board {
	b := Board{}

	//b.Grid = make([]uint8, c.Width * c.Height)

	return b
}


func copy_board(b *Board) Board {
	nb := new_board()

	//copy(nb.Grid, b.Grid)
	nb.Grid = b.Grid

	return nb
}

func set_val(b *Board, x int, y int, v uint8) {
	b.Grid[(y * BOARD_WIDTH) + x] = v
}


func get_val(b *Board, x int, y int) uint8 {
	return b.Grid[(y * BOARD_WIDTH) + x]
}


func board_to_string(b *Board) string {

	bstr := ""
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			bstr += val_string[get_val(b, x, y)]
		}
		bstr += "\n"
	}

	return bstr
}


func will_fit(b *Board, x int, y int, v uint8) bool {

	if get_val(b, x, y) != 0 {
		return false
	}

	dia := val_dia[v]

	// Check above
	if y > 0 && val_dia[get_val(b, x, y - 1)] + dia > DIA_LIMIT {
		return false
	}

	// Check below
	if y + 1 < BOARD_HEIGHT && val_dia[get_val(b, x, y + 1)] + dia > DIA_LIMIT {
		return false
	}

	// Check left
	if x > 0 && val_dia[get_val(b, x - 1, y)] + dia > DIA_LIMIT {
		return false
	}

	// Check right
	if x + 1 < BOARD_WIDTH && val_dia[get_val(b, x + 1, y)] + dia > DIA_LIMIT {
		return false
	}

	return true
}


func find_move_target(b *Board, x int, y int, d int) (int, int, error) {

	v := get_val(b, x, y)

	if v == 0 {
		return -1, -1, errors.New("No piece at that location!")
	}

	dx, dy := dir_dx[d], dir_dy[d]
	cx, cy := x, y

	found_occupied := 0
	found_blank := false
	for {
		cx, cy = cx + dx, cy + dy

		if cx < 0 || cy < 0 || cx >= BOARD_WIDTH || cy >= BOARD_HEIGHT {
			return -1, -1, errors.New("Walked off board looking for move target!")
		}

		tv := get_val(b, cx, cy)

		if tv == 0 {
			found_blank = true
			break
		} else {
			found_occupied++
		}
	}

	if found_occupied > 0 && found_blank {
		if will_fit(b, cx, cy, v) {
			return cx, cy, nil
		} else {
			return -1, -1, errors.New("Viable target found but piece won't fit")
		}
	} else {
		return -1, -1, errors.New("No viable target found")
	}

}


func can_move(b *Board, x int, y int, d int) bool {
	_, _, err := find_move_target(b, x, y, d);

	return err == nil
}


func move(b *Board, x int, y int, d int) {
	nx, ny, err := find_move_target(b, x, y, d);

	if err != nil {
		fmt.Fprintf(os.Stderr, "Tried to make illegal move: %s\n", err.Error())
		os.Exit(1)
	}

	v := get_val(b, x, y)
	set_val(b, x, y, 0)
	set_val(b, nx, ny, v)
}


func find_solution(states *map[Board]int, b *Board) int {

	*states = make(map[Board]int) // "empty" the states map

	(*states)[*b] = 0
	d := 0
	f := 0
	for {
		var layerd []Board
		for k, v := range *states {
			if v == d {
				layerd = append(layerd, k)
			}
		}

		for _, cb := range layerd {
			for x := 0; x < BOARD_WIDTH; x++ {
				for y := 0; y < BOARD_HEIGHT; y++ {
					for dir := 0; dir < 4; dir++ {
						if can_move(&cb, x, y, dir) {
							nb := copy_board(&cb)

							move(&nb, x, y, dir)

							_, ok := (*states)[nb]

							if !ok {
								(*states)[nb] = d + 1
								f++

								if get_val(&nb, 0, 0) == 4 {
									//fmt.Printf("Found solution in %d moves!\n", d + 1)
									return d + 1
								}
							}
						}
					}
				}
			}
		}

		if f > 0 {
			d++
		} else {
			//fmt.Printf("Got to depth %d without a solution!\n", d)
			return -1
		}

	}
}


func print_solution_steps(states *map[Board]int, d int) error {

	if d <= 0 {
		return errors.New("Depth must be positive!")
	}

	boards := make([]Board, 0, d + 1)
	steps := make([]string, 0, d + 1)

	var b Board
	f := false
	for k, v := range *states {
		if v == d {
			if get_val(&k, 0, 0) == 4 {
				b = k
				f = true
				break
			}
		}
	}

	if !f {
		return errors.New("Could not find solution board at specified depth")
	}

	boards = append(boards, b)
	for d > 0 {
		for x := 0; x < BOARD_WIDTH; x++ {
			for y := 0; y < BOARD_HEIGHT; y++ {
				for dir := 0; dir < 4; dir++ {
					if can_move(&b, x, y, dir) {
						nb := copy_board(&b)
						move(&nb, x, y, dir)

						nd, ok := (*states)[nb]

						if !ok {
							continue
						}

						if nd >= d {
							continue
						}

						nx, ny, err := find_move_target(&b, x, y, dir)

						if err != nil {
							return err
						}

						v := get_val(&b, x, y)

						step := fmt.Sprintf("Step %d: Move %s at (%d,%d) %s to (%d,%d)",
							d, val_name[v], nx + 1, ny + 1, dir_opp_str[dir], x + 1, y + 1)

						steps = append(steps, step)
						boards = append(boards, nb)
						b = nb
						d = nd
					}

				}
			}
		}
	}

	steps = append(steps, "Step 0: Start with this board")

	for i := len(boards) - 1; i >= 0; i-- {
		fmt.Println(steps[i])
		fmt.Println(board_to_string(&(boards[i])))
	}

	return nil
}



func main() {

	b := new_board()

	// 90 moves
	b.Grid = [16]uint8{2, 2, 2, 1, 1, 0, 0, 0, 1, 0, 0, 2, 1, 4, 0, 0}

	//fmt.Printf("Board b:\n%s", board_to_string(&b))

	states := make(map[Board]int)
	d := find_solution(&states, &b)

	if d > 0 {
		fmt.Printf("Found solution in %d moves!\n", d)
	} else {
		fmt.Printf("Failed to find a solution!\n")
	}

	err := print_solution_steps(&states, d)

	if err != nil {
		fmt.Printf(err.Error())
	}
}
