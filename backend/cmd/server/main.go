package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"plugins.dev/backend/internal/api"
	"plugins.dev/backend/internal/github"
	"plugins.dev/backend/internal/notify"
	"plugins.dev/backend/internal/registry"
)

func main() {
	gh := github.NewClient(os.Getenv("GITHUB_TOKEN"))
	r := api.NewRouterWithDeps(registry.New(), gh, notify.NewEngine(5*time.Minute))
	log.Println("backend listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
