package router

import (
	"net/http"

	"golden-hearts/backend/handlers"
)

func MultiPlexer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", handlers.HomeHandler())
	mux.Handle("/projects", handlers.ProjectsHandler())
	mux.Handle("/donations", handlers.DonationsHandler())

	return mux
}
