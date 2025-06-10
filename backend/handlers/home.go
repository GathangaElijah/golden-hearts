package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DonationsHandler called for path: %s", r.URL.Path)
		fmt.Fprintf(w, "Welcome to the Home Page!")
	})
}
