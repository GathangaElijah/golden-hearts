package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type ProjectData struct {
	Id          string `json:"id"`
	Title       string `json:"project-title"`
	Description string `json:"project-description"`
	Amount      string `json:"contribution-target"`
}

func ProjectsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Projects Handler called for path: %s", r.URL.Path)

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Problem opening current working dir\n%v", err)
			http.Error(w, "Internal error: path resolution", http.StatusInternalServerError)
			return
		}
		projectsFilePath := path.Join(cwd, "backend", "data", "projects.json")

		fileContent, err := os.ReadFile(projectsFilePath)
		if err != nil {
			http.Error(w, "Unable to read projects data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(fileContent)
	})
}
