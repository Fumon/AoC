package main

import (
	"fmt"
	"fuaoc2023/day13/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var mirrored_row_count int
	var mirrored_column_count int

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
				if line == rows[rowcount - 1] {
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
			if mirror_index > 0 && string(column) == string(columns[mirror_index - 1]) {
				for i, j := mirror_index - 2, mirror_index + 1; i > 0 && j < len(columns); i, j = i - 1, j + 1 {
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
	}

	return mirrored_column_count + (100 * mirrored_row_count)
}