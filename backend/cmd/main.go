package main

import (
	"fmt"
	"golden-hearts/backend/router"
	"log"
	"net/http"
)

func main() {

	mux := router.MultiPlexer()

	port := ":5001"
	serverURL := fmt.Sprintf("http://localhost%s", port)

	// Print the URL to the console
	fmt.Printf("Server starting on %s\n", serverURL)

	err := http.ListenAndServe(":5001", mux)

	if err != nil {
		log.Fatalf("Error starting the server, %v", err)
		return
	}

}
