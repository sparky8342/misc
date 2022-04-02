package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func valid(grid []int, pos int) bool {
	for i := 0; i < 81; i++ {
		if grid[i] != grid[pos] || i == pos {
			continue
		}
		// same row
		if int(i/9) == int(pos/9) {
			return false
		}
		// same column
		if i%9 == pos%9 {
			return false
		}
		// same square
		if int(i/27) == int(pos/27) && int(i%9/3) == int(pos%9/3) {
			return false
		}
	}
	return true
}

func search(grid []int, n int) bool {
	for i := 0; i < 81; i++ {
		if grid[i] == 0 {
			for digit := 1; digit <= 9; digit++ {
				grid[i] = digit
				if valid(grid, i) {
					if search(grid, i+1) {
						return true
					}
				}
			}
			grid[i] = 0
			return false
		}
	}
	return true
}

func read_puzzle(filename string) []int {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	grid := []int{}
	for _, ch := range data {
		if ch != '\n' {
			grid = append(grid, int(ch-'0'))
		}
	}
	return grid
}

func print_puzzle(puzzle []int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Print(puzzle[i*9+j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: " + os.Args[0] + " [filename]")
		return
	}
	puzzle := read_puzzle(os.Args[1])
	print_puzzle(puzzle)
	search(puzzle, 0)
	print_puzzle(puzzle)
}

/*
005300000
800000020
070010500
400005300
010070006
003200080
060500009
004000030
000009700
*/
