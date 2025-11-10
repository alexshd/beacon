package main

import (
	"fmt"
	"log"
	"os"

	sudokuexample "github.com/alexshd/beacon/sudoku-example"
)

func main() {
	version := "v1.0"
	port := "9000"

	if len(os.Args) > 1 {
		version = os.Args[1]
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	addr := fmt.Sprintf(":%s", port)
	server := sudokuexample.NewServer(version)
	log.Fatal(server.Start(addr))
}
