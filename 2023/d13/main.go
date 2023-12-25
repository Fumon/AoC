package main

import (
	"fmt"
	"fuaoc2023/day13/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var mirrored_row_count int
	var mirrored_column_count int

	var puzzcount int = 1
	// outer:
	for puzzchan := range u.LineDelimiter(lines, "") {
		var rowcount int
		var rows []string

		puzzlines, linewidth := u.PeekAtFirstLineWidth(puzzchan)
		var columns [][]rune = make([][]rune, linewidth)

		var mirror_indicies []int

		for line := range puzzlines {
			rows = append(rows, line)
			if rowcount > 0 {
				if line == rows[rowcount-1] {
					mirror_indicies = append(mirror_indicies, rowcount)
				}
			}

			// Add to columns
			for i, char := range line {
				columns[i] = append(columns[i], char)
			}

			rowcount++
		}

		// Check mirrors
		var row_found_indicies []int
	rowmirrorloop:
		for _, mirror_index := range mirror_indicies {
			var i = mirror_index - 2
			var j = mirror_index + 1

			for i >= 0 && j < len(rows) {
				if rows[i] != rows[j] {
					continue rowmirrorloop
				}
				i--
				j++
			}
			row_found_indicies = append(row_found_indicies, mirror_index)
		}

		var largest_row_index int
		if len(row_found_indicies) > 0 {
			for _, ind := range row_found_indicies {
				if ind > largest_row_index {
					largest_row_index = ind
				}
			}
		}

		// Check columns
		var column_found_indicies []int
	columnmirrorloop:
		for mirror_index, column := range columns {
			if mirror_index > 0 && string(column) == string(columns[mirror_index-1]) {
				for i, j := mirror_index-2, mirror_index+1; i >= 0 && j < len(columns); i, j = i-1, j+1 {
					if string(columns[i]) != string(columns[j]) {
						continue columnmirrorloop
					}
				}
				column_found_indicies = append(column_found_indicies, mirror_index)
			}
		}

		var largest_column_index int
		if len(column_found_indicies) > 0 {
			for _, ind := range column_found_indicies {
				if ind > largest_column_index {
					largest_column_index = ind
				}
			}
		}

		if largest_column_index == 0 && largest_row_index == 0 {
			panic("no mirrors")
		} else {
			if largest_column_index > largest_row_index {
				mirrored_column_count += largest_column_index
			} else {
				mirrored_row_count += largest_row_index
			}
		}

		fmt.Println(puzzcount, ": Rows: ", row_found_indicies, " Columns: ", column_found_indicies)
		puzzcount++
	}

	return mirrored_column_count + (100 * mirrored_row_count)
}

func Part2(lines <-chan string) int {
	var mirrored_row_count int
	var mirrored_column_count int

	var puzzcount int = 1
	// outer:
	for puzzchan := range u.LineDelimiter(lines, "") {
		var rowcount int
		var rows []string

		puzzlines, linewidth := u.PeekAtFirstLineWidth(puzzchan)
		var columns [][]rune = make([][]rune, linewidth)

		var mirror_indicies []int
		for line := range puzzlines {
			rows = append(rows, line)
			if rowcount > 0 {
				if stringSimilarity(line, rows[rowcount-1]) < 2 {
					mirror_indicies = append(mirror_indicies, rowcount)
				}
			}

			// Add to columns
			for i, char := range line {
				columns[i] = append(columns[i], char)
			}

			rowcount++
		}

		// Check mirrors
		var row_found_indicies []int
		var row_found_indicies_smudge []int
	rowmirrorloop:
		for _, mirror_index := range mirror_indicies {
			var i = mirror_index - 1
			var j = mirror_index + 0

			var smudgebudget = 1
			for i >= 0 && j < len(rows) {
				diff := stringSimilarity(rows[i], rows[j])
				if diff > smudgebudget {
					continue rowmirrorloop
				} else {
					smudgebudget -= diff
				}
				i--
				j++
			}
			if smudgebudget == 1 {
				row_found_indicies = append(row_found_indicies, mirror_index)
			} else {
				row_found_indicies_smudge = append(row_found_indicies_smudge, mirror_index)
			}
		}

		// Check columns
		var column_found_indicies []int
		var column_found_indicies_smudge []int
	columnmirrorloop:
		for mirror_index, column := range columns {
			if mirror_index > 0 {
				diff := stringSimilarity(string(column), string(columns[mirror_index-1]))
				if diff > 1 {
					continue columnmirrorloop
				}
				var smudgebudget = 1 - diff
				for i, j := mirror_index-2, mirror_index+1; i >= 0 && j < len(columns); i, j = i-1, j+1 {
					
					diff = stringSimilarity(string(columns[i]), string(columns[j]))
					if diff > smudgebudget {
						continue columnmirrorloop
					} else {
						smudgebudget -= diff
					}
				}
				if smudgebudget > 0 {
					column_found_indicies = append(column_found_indicies, mirror_index)
				} else {
					column_found_indicies_smudge = append(column_found_indicies_smudge, mirror_index)
				}	
			}
		}

		// var largest_row_index int
		// if len(row_found_indicies) > 0 && len(row_found_indicies_smudge) != 1 {
		// 	panic(fmt.Sprint("Wrong row smudge count: ", len(row_found_indicies_smudge)))
		// } else {
		// 	largest_row_index = row_found_indicies_smudge[0]
		// }

		// var largest_column_index int
		// if len(column_found_indicies) > 0 && len(column_found_indicies_smudge) != 1 {
		// 	panic(fmt.Sprint("Wrong column smudge count: ", len(column_found_indicies_smudge)))
		// } else {
		// 	largest_column_index = column_found_indicies_smudge[0]
		// }

		if len(row_found_indicies_smudge) != 1 {
			if len(column_found_indicies_smudge) != 1 {
				panic("smudge error")
			} else {
				mirrored_column_count += column_found_indicies_smudge[0]
			}
		} else {
			mirrored_row_count += row_found_indicies_smudge[0]
		}

		fmt.Println(puzzcount, ": Rows,Smudge: ", row_found_indicies, row_found_indicies_smudge, " Columns, Smudge: ", column_found_indicies, column_found_indicies_smudge)
		puzzcount++
	}

	return mirrored_column_count + (100 * mirrored_row_count)
}

func stringSimilarity(one, two string) (diffs int) {
	for i, ch := range []byte(one) {
		if ch != two[i] {
			diffs++
		}
	}
	return
}
