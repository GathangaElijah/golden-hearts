package handlers

import (
	"fmt"
	"net/http"
)

func CallBackHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Home Page!")
	})
}
