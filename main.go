package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	Name     string
	Status   string
	Restarts int
	Size     string // New field
	Type     string // New field (Image name/ID)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		// FETCH WITH SIZE: This tells Docker to calculate disk usage
		containers, err := cli.ContainerList(ctx, container.ListOptions{All: true, Size: true})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var stats []ContainerInfo
		for _, c := range containers {
			// Convert bytes to Megabytes for readability
			sizeMB := float64(c.SizeRootFs) / 1024 / 1024

			stats = append(stats, ContainerInfo{
				Name:     c.Names[0][1:],
				Status:   c.Status,
				Restarts: 0, // Simplified for this example
				Size:     fmt.Sprintf("%.2f MB", sizeMB),
				Type:     c.Image, // This shows the Image name or SHA
			})
		}

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, stats)
	})

	println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
