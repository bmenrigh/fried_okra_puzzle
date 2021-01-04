package main

import (
	"fmt"
	"errors"
	"os"
	"sort"
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


func transpose(b *Board) Board {
	nb := copy_board(b)

	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			set_val(&nb, y, x, get_val(b, x, y))
		}
	}

	return nb
}


func is_symmetric(b *Board) bool {

	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := y + 1; x < BOARD_WIDTH; x++ {
			if get_val(b, x, y) != get_val(b, y, x) {
				return false
			}
		}
	}

	return true
}


func is_solved(b *Board) bool {

	return get_val(b, 0, 0) == 4
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


func explore_from_board_list(states *map[Board]int, states_solved *map[Board]bool, blist *[]Board, found *int, depth *int, found_solved *int, stop_at_solved bool) {

	*states = make(map[Board]int) // "empty" the states map
	*states_solved = make(map[Board]bool)

	// It is allowed to start the search at a solved state
	// even if stop_at_solved is true
	for _, b := range *blist {
		(*states)[b] = 0
	}

	*found = 0
	d := 0
	for {
		f := 0

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
								if is_solved(&nb) {
									(*states_solved)[nb] = true
									*found_solved += 1

									if stop_at_solved {
										// Don't add solved states to
										// found state list since we're
										// stopping the search at
										// solved states
										continue
									}
								}

								(*states)[nb] = d + 1
								f++
							}
						}
					}
				}
			}
		}

		if f > 0 {
			d++

			*found += f
		} else {
			//fmt.Printf("Got to depth %d without a solution!\n", d)
			*depth = d
			return
		}

	}
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

								if is_solved(&nb) {
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
	for cb, cd := range *states {
		if cd == d && is_solved(&cb) {
			b = cb
			f = true
			break
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


func enumerate_all_piecesets(pcs *[][4]int,) {

	*pcs = make([][4]int, 0)

	for tpc := 1; tpc < (BOARD_WIDTH * BOARD_HEIGHT) - 1; tpc++ {
		for rc := 1; rc <= tpc; rc++ {
			for gc := 0; gc <= tpc - rc; gc++ {
				bc := tpc - (rc + gc)

				if bc == 1 {
					continue
				}

				if rc + gc + bc != tpc {
					fmt.Println("Should be impossible now!");
					continue
				}

				pc := [4]int{((BOARD_WIDTH * BOARD_HEIGHT) - 1) - tpc, rc, bc, gc}

				*pcs = append(*pcs, pc)
			}
		}
	}
}


func enumerate_solved(solved *[]Board, pc [4]int) {

	*solved = make([]Board, 0) // "empty" the solved states list

	b := new_board()
	set_val(&b, 0, 0, 4) // Put the yellow piece in the solved position
	rec_enum(solved, &b, pc, 1)
}


func rec_enum(solved *[]Board, b *Board, pc [4]int, o int) {

	cx, cy := o % BOARD_WIDTH, o / BOARD_WIDTH

	//fmt.Printf("rec_enum at %d, %d\n", cx, cy);

	if cy >= BOARD_HEIGHT {

		*solved = append(*solved, *b)
		//fmt.Printf("Adding board\n")

		return
	}

	for i := 0; i < 4; i++ {
		if pc[i] > 0 {
			npc := pc
			npc[i]--

			if i == 0 {
				// Now just leave this spot blank and move on
				rec_enum(solved, b, npc, o + 1)
			} else if will_fit(b, cx, cy, uint8(i)) {
				nb := copy_board(b)
				set_val(&nb, cx, cy, uint8(i))

				rec_enum(solved, &nb, npc, o + 1)
			}
		}
	}

}


func board_less(b1, b2 *Board) bool {

	for i := 0; i < (BOARD_HEIGHT * BOARD_WIDTH); i++ {
		if ((*b1).Grid[i] < (*b2).Grid[i]) {
			return true
		} else if ((*b1).Grid[i] > (*b2).Grid[i]) {
			return false
		}
	}

	return false
}



func main() {

	//b := new_board()

	// 90 moves
	// b.Grid = [16]uint8{2, 2, 2, 1, 1, 0, 0, 0, 1, 0, 0, 2, 1, 4, 0, 0}

	// bt := transpose(&b)

	// if board_less(&bt, &b) {
	// 	fmt.Printf("Board bt is less than b:\n%s", board_to_string(&bt))
	// } else {
	// 	fmt.Printf("Board b is greater or equal to bt:\n%s", board_to_string(&b))
	// }

	//fmt.Printf("Board b:\n%s", board_to_string(&b))

	//pc := [4]int{11, 4, 0, 0}
	//var solved []Board
	// enumerate_solved(&solved, pc)

	// fmt.Printf("Found %d solved states.\n", len(solved))

	// states := make(map[Board]int)

	// depth, found := 0, 0
	// explore_from_solved(&states, &solved, &found, &depth)

	// fmt.Printf("Found %d total states to a depth of %d.\n", found, depth)

	// d := find_solution(&states, &b)

	// if d > 0 {
	// 	fmt.Printf("Found solution in %d moves!\n", d)
	// } else {
	// 	fmt.Printf("Failed to find a solution!\n")
	// }

	// err := print_solution_steps(&states, d)

	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	var pcs [][4]int
	enumerate_all_piecesets(&pcs);
	for _, pc := range pcs {

		// Debug: to stop at a single configuration
		//if pc != [4]int{4, 1, 10, 0} {
		//continue
		//}

		var solved []Board
		enumerate_solved(&solved, pc)

		can_scramble := 0
		for _, b := range solved {
			if can_move(&b, 0, 0, 1) || can_move(&b, 0, 0, 3) {
				can_scramble++
			}
		}

		all_states := make(map[Board]int)
		var all_states_solved map[Board]bool
		depth, found, found_solved := 0, 0, 0
		explore_from_board_list(&all_states, &all_states_solved, &solved, &found, &depth, &found_solved, false)

		var furthest []Board
		for k, v := range all_states {
			if v == depth {
				furthest = append(furthest, k)
			}
		}

		sym_b := 0
		for _, b := range furthest {
			if is_symmetric(&b) {
				sym_b++
			}
		}

		sort.Slice(furthest, func(i, j int) bool {
			return board_less(&furthest[i], &furthest[j])
		})

		// Now we'll re-explore all the states until all
		// of them have been found; counting how many
		// disjoint sets there are
		all_found_nostart := make(map[Board]bool)
		var largest_solved_nostart []Board
		disjoint_sets_nostart := 0
		symmetry_free_sets_nostart := 0
		largest_set_nostart := 0
		largest_set_solved_nostart := 0
		largest_set_number_nostart := 0
		disjoint_furthest_nostart := make([]Board, 0)
		// First from each furthest state
		for _, b := range furthest {
			_, ok := all_found_nostart[b]

			if !ok {
				disjoint_sets_nostart++

				// Track this as an example state for this set
				disjoint_furthest_nostart = append(disjoint_furthest_nostart, b)

				set_states := make(map[Board]int)
				set_states_solved := make(map[Board]bool)
				start_state := make([]Board, 0, 1)
				start_state = append(start_state, b)
				set_size := 0
				set_depth := 0
				set_solved := 0
				explore_from_board_list(&set_states, &set_states_solved, &start_state, &set_size, &set_depth, &set_solved, true)

				// Account for the one starting state (which isn't 'found')
				set_size++;

				if set_size > largest_set_nostart {
					largest_set_nostart = set_size
					largest_set_solved_nostart = set_solved
					largest_set_number_nostart = disjoint_sets_nostart

					largest_solved_nostart = make([]Board, 0)
					for k, _ := range set_states_solved {
						if is_solved(&k) {
							largest_solved_nostart = append(largest_solved_nostart, k)
						}
					}
				}

				// Check if this set was symmetry-free by seeing if
				// the start-state's transpose is in the set
				tss := transpose(&b)
				_, ok := set_states[tss]

				if !ok {
					symmetry_free_sets_nostart++
				}

				// Put each of these found states in the found map
				for cb, _ := range set_states {
					all_found_nostart[cb] = true
				}
			}
		}
		// Record the stats now that we've searched all the furthest
		disjoint_furthest_sets_nostart := disjoint_sets_nostart;
		symmetry_free_furthest_sets_nostart := symmetry_free_sets_nostart
		furthest_largest_set_nostart := largest_set_nostart
		furthest_largest_set_solved_nostart := largest_set_solved_nostart

		// Continue search
		for b, _ := range all_states {

			if is_solved(&b) {
				continue
			}

			_, ok := all_found_nostart[b]

			if !ok {
				disjoint_sets_nostart++

				set_states := make(map[Board]int)
				set_states_solved := make(map[Board]bool)
				start_state := make([]Board, 0, 1)
				start_state = append(start_state, b)
				set_size := 0
				set_depth := 0
				set_solved := 0
				explore_from_board_list(&set_states, &set_states_solved, &start_state, &set_size, &set_depth, &set_solved, true)

				// Account for the one starting state (which isn't 'found')
				set_size++;

				if set_size > largest_set_nostart {
					largest_set_nostart = set_size
					largest_set_solved_nostart = set_solved
					largest_set_number_nostart = disjoint_sets_nostart

					largest_solved_nostart = make([]Board, 0)
					for k, _ := range set_states_solved {
						if is_solved(&k) {
							largest_solved_nostart = append(largest_solved_nostart, k)
						}
					}
				}

				// Check if this set was symmetry-free by seeing if
				// the start-state's transpose is in the set
				tss := transpose(&b)
				_, ok := set_states[tss]

				if !ok {
					symmetry_free_sets_nostart++
				}

				// Put each of these found states in the found map
				for cb, _ := range set_states {
					all_found_nostart[cb] = true
				}
			}
		}

		largest_states_nostart := make(map[Board]int)
		largest_states_solved_nostart := make(map[Board]bool)
		largest_depth_nostart, largest_found_nostart, largest_found_solved_nostart := 0, 0, 0
		explore_from_board_list(&largest_states_nostart, &largest_states_solved_nostart, &largest_solved_nostart, &largest_found_nostart, &largest_depth_nostart, &largest_found_solved_nostart, true)

		var largest_furthest_nostart []Board
		for k, v := range largest_states_nostart {
			if v == largest_depth_nostart {
				largest_furthest_nostart = append(largest_furthest_nostart, k)
			}
		}

		sort.Slice(largest_furthest_nostart, func(i, j int) bool {
			return board_less(&largest_furthest_nostart[i], &largest_furthest_nostart[j])
		})

		// Now do the exploration again but pass through the start states
		all_found_wstart := make(map[Board]bool)
		var largest_solved_wstart []Board
		disjoint_sets_wstart := 0
		symmetry_free_sets_wstart := 0
		largest_set_wstart := 0
		largest_set_solved_wstart := 0
		largest_set_number_wstart := 0
		disjoint_furthest_wstart := make([]Board, 0)
		// First from each furthest state
		for _, b := range furthest {
			_, ok := all_found_wstart[b]

			if !ok {
				disjoint_sets_wstart++

				// Track this as an example state for this set
				disjoint_furthest_wstart = append(disjoint_furthest_wstart, b)

				set_states := make(map[Board]int)
				set_states_solved := make(map[Board]bool)
				start_state := make([]Board, 0, 1)
				start_state = append(start_state, b)
				set_size := 0
				set_depth := 0
				set_solved := 0
				explore_from_board_list(&set_states, &set_states_solved, &start_state, &set_size, &set_depth, &set_solved, false)

				// Account for the one starting state (which isn't 'found')
				set_size++;

				if set_size > largest_set_wstart {
					largest_set_wstart = set_size
					largest_set_solved_wstart = set_solved
					largest_set_number_wstart = disjoint_sets_wstart

					largest_solved_wstart = make([]Board, 0)
					for k, _ := range set_states_solved {
						if is_solved(&k) {
							largest_solved_wstart = append(largest_solved_wstart, k)
						}
					}
				}

				// Check if this set was symmetry-free by seeing if
				// the start-state's transpose is in the set
				tss := transpose(&b)
				_, ok := set_states[tss]

				if !ok {
					symmetry_free_sets_wstart++
				}

				// Put each of these found states in the found map
				for cb, _ := range set_states {
					all_found_wstart[cb] = true
				}
			}
		}
		// Record the stats now that we've searched all the furthest
		disjoint_furthest_sets_wstart := disjoint_sets_wstart;
		symmetry_free_furthest_sets_wstart := symmetry_free_sets_wstart
		furthest_largest_set_wstart := largest_set_wstart
		furthest_largest_set_solved_wstart := largest_set_solved_wstart

		// Continue search
		for b, _ := range all_states {

			if is_solved(&b) {
				continue
			}

			_, ok := all_found_wstart[b]

			if !ok {
				disjoint_sets_wstart++

				set_states := make(map[Board]int)
				set_states_solved := make(map[Board]bool)
				start_state := make([]Board, 0, 1)
				start_state = append(start_state, b)
				set_size := 0
				set_depth := 0
				set_solved := 0
				explore_from_board_list(&set_states, &set_states_solved, &start_state, &set_size, &set_depth, &set_solved, false)

				// Account for the one starting state (which isn't 'found')
				set_size++;

				if set_size > largest_set_wstart {
					largest_set_wstart = set_size
					largest_set_solved_wstart = set_solved
					largest_set_number_wstart = disjoint_sets_wstart

					largest_solved_wstart = make([]Board, 0)
					for k, _ := range set_states_solved {
						if is_solved(&k) {
							largest_solved_wstart = append(largest_solved_wstart, k)
						}
					}
				}

				// Check if this set was symmetry-free by seeing if
				// the start-state's transpose is in the set
				tss := transpose(&b)
				_, ok := set_states[tss]

				if !ok {
					symmetry_free_sets_wstart++
				}

				// Put each of these found states in the found map
				for cb, _ := range set_states {
					all_found_wstart[cb] = true
				}
			}
		}

		largest_states_wstart := make(map[Board]int)
		largest_states_solved_wstart := make(map[Board]bool)
		largest_depth_wstart, largest_found_wstart, largest_found_solved_wstart := 0, 0, 0
		explore_from_board_list(&largest_states_wstart, &largest_states_solved_wstart, &largest_solved_wstart, &largest_found_wstart, &largest_depth_wstart, &largest_found_solved_wstart, false)

		var largest_furthest_wstart []Board
		for k, v := range largest_states_wstart {
			if v == largest_depth_wstart {
				largest_furthest_wstart = append(largest_furthest_wstart, k)
			}
		}

		sort.Slice(largest_furthest_wstart, func(i, j int) bool {
			return board_less(&largest_furthest_wstart[i], &largest_furthest_wstart[j])
		})

		fmt.Printf("# %dx%d with %d pieces (R:%d, G:%d, B:%d) (solved:%d, scrambleable:%d, found:%d) has %d (symmetric:%d) furthest requiring %d moves\n",
			BOARD_WIDTH, BOARD_HEIGHT, (BOARD_WIDTH * BOARD_HEIGHT) - pc[0], pc[1], pc[3], pc[2],
			len(solved), can_scramble, found, len(furthest), sym_b, depth)
		if depth > 0 {

			fmt.Printf("csv,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
				BOARD_WIDTH, BOARD_HEIGHT, (BOARD_WIDTH * BOARD_HEIGHT) - pc[0], pc[1], pc[3], pc[2],
				len(solved), can_scramble, found, len(furthest), sym_b, depth,
				disjoint_furthest_sets_nostart, symmetry_free_furthest_sets_nostart, furthest_largest_set_nostart + furthest_largest_set_solved_nostart, furthest_largest_set_solved_nostart,
				disjoint_sets_nostart, symmetry_free_sets_nostart, largest_set_nostart + largest_set_solved_nostart, largest_set_solved_nostart, largest_depth_nostart, len(largest_furthest_nostart),
				disjoint_furthest_sets_wstart, symmetry_free_furthest_sets_wstart, furthest_largest_set_wstart, furthest_largest_set_solved_wstart,
				disjoint_sets_wstart, symmetry_free_sets_wstart, largest_set_wstart, largest_set_solved_wstart, largest_depth_wstart, len(largest_furthest_wstart))

			fmt.Printf("The %d furthest states come in %d disjoint set(s) (symmetryfree:%d, largest:%d, solved:%d) out of %d set(s) (symmetryfree:%d, largest:%d, solved:%d, largestdepth:%d, largestfurthest:%d) when exploration is stopped at a solved state\n",
				len(furthest), disjoint_furthest_sets_nostart, symmetry_free_furthest_sets_nostart, furthest_largest_set_nostart + furthest_largest_set_solved_nostart, furthest_largest_set_solved_nostart,
				disjoint_sets_nostart, symmetry_free_sets_nostart, largest_set_nostart + largest_set_solved_nostart, largest_set_solved_nostart, largest_depth_nostart, len(largest_furthest_nostart))

			for n, cb := range disjoint_furthest_nostart {
				fmt.Printf("====== Set %d Example Furthest State (stopped at solved) =====\n", n + 1);
				fmt.Println(board_to_string(&cb))
			}

			if len(largest_furthest_nostart) > 0 {
				fmt.Printf("====== Example Largest Set (%d) Furthest State (stopped at solved) =====\n", largest_set_number_nostart);
				fmt.Println(board_to_string(&largest_furthest_nostart[0]))
			}

			fmt.Printf("The %d furthest states come in %d disjoint set(s) (symmetryfree:%d, largest:%d, solved:%d) out of %d set(s) (symmetryfree:%d, largest:%d, solved:%d, largestdepth:%d, largestfurthest:%d) when exploration continues through solved states\n",
				len(furthest), disjoint_furthest_sets_wstart, symmetry_free_furthest_sets_wstart, furthest_largest_set_wstart, furthest_largest_set_solved_wstart,
				disjoint_sets_wstart, symmetry_free_sets_wstart, largest_set_wstart, largest_set_solved_wstart, largest_depth_wstart, len(largest_furthest_wstart))

			for n, cb := range disjoint_furthest_wstart {
				fmt.Printf("====== Set %d Example Furthest State (continued through solved) =====\n", n + 1);
				fmt.Println(board_to_string(&cb))
			}

			if len(largest_furthest_wstart) > 0 {
				fmt.Printf("====== Example Largest Set (%d) Furthest State (continued through solved) =====\n", largest_set_number_wstart);
				fmt.Println(board_to_string(&largest_furthest_wstart[0]))
			}
		}
	}
}
