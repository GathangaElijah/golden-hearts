package handlers

import (
	"fmt"
	"net/http"
)

func DonationsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Home Page!")
	})
}
