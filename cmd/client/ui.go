package main

import "fmt"

func render(board [9]string) {
	fmt.Println()
	for r := range 3 {
		i := r * 3
		fmt.Printf(" %s | %s | %s \n", board[i], board[i+1], board[i+2])
		if r < 2 {
			fmt.Println("---+---+---")
		}
	}
	fmt.Println("\nplay a cell (1-9) or q to quit:")
}
