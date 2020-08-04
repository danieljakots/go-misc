package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strings"
)

func story() {
	fmt.Println("Pète et Répète sont dans un bateau.")
	fmt.Println("Pète tombe à l'eau, qui reste?")
}

func askAnswer() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func main() {
	text := ""
	for {
		story()
		text = askAnswer()
		if strings.ToLower(text) != "répète" {
			fmt.Println("Non, fais mieux attention")
		} else {
			fmt.Println("Ok je répète")
		}
	}

}
