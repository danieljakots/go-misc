package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func createHash() {
	fmt.Println("Enter password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	hash, err := bcrypt.GenerateFromPassword(password, 8)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(hash))
}

func checkHash() {
	fmt.Println("Enter password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	hash := os.Args[1]
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) == 1 {
		createHash()
	}

	if len(os.Args) == 2 {
		checkHash()
	}
}
