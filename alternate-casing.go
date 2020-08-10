package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	text_joined := strings.Join(os.Args[1:], " ")
	for index, letter := range text_joined {
		char := fmt.Sprintf("%c", letter)
		if index%2 == 0 {
			fmt.Print(strings.ToLower(char))
		} else {
			fmt.Print(strings.ToUpper(char))
		}
	}
	fmt.Println()
}
