package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alexshd/beacon/httpserver-example"
)

func main() {
	// Get port from command line or default to 8080
	port := "8080"
	idMult := 1

	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// ID multiplier: port 8080 -> mult 1 (IDs 10,11,12...)
	//                port 8081 -> mult 2 (IDs 20,21,22...)
	if port == "8080" {
		idMult = 1
	} else if port == "8081" {
		idMult = 2
	} else if len(os.Args) > 2 {
		m, err := strconv.Atoi(os.Args[2])
		if err == nil {
			idMult = m
		}
	}

	addr := fmt.Sprintf(":%s", port)

	// Create server with Law I immutable state and unique ID range
	server := httpserver.NewServerWithIDMultiplier(idMult)

	log.Printf("Server starting with ID multiplier: %d (IDs start at %d)", idMult, idMult*100)

	// Start server
	log.Fatal(server.Start(addr))
}
