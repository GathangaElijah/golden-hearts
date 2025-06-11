package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golden-hearts/backend/router"
)

func main() {
	mux := router.MultiPlexer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001" // for local development
	}

	serverURL := fmt.Sprintf("http://localhost%s", port)

	// Print the URL to the console
	fmt.Printf("Server starting on %s\n", serverURL)

	err := http.ListenAndServe(":" + port, mux)
	if err != nil {
		log.Fatalf("Error starting the server, %v", err)
		return
	}
}
