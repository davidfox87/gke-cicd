package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func releasesHandler(w http.ResponseWriter, r *http.Request) {
	const bucket = "cloud-deploy-releases"
	const object = "releases.json"

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Failed to create GCS client", http.StatusInternalServerError)
		log.Printf("storage.NewClient: %v", err)
		return
	}
	defer client.Close()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			http.Error(w, "releases.json not found in bucket", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to read object", http.StatusInternalServerError)
		}
		log.Printf("Object(%q).NewReader: %v", object, err)
		return
	}
	defer rc.Close()

	// Set proper headers for Grafana
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Stream directly from GCS → HTTP response (zero memory copy for large files)
	if _, err := io.Copy(w, rc); err != nil {
		// Client probably disconnected — log but don't send error body
		log.Printf("Error streaming to client: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/releases.json", releasesHandler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Grafana GCS proxy — /releases.json")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
