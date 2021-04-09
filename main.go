package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kuipercm/spaces-summit-famous-places/handlers"
	"github.com/kuipercm/spaces-summit-famous-places/pubsub"
	"github.com/kuipercm/spaces-summit-famous-places/storage"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"time"
)

func main() {
	baseRouter := mux.NewRouter().StrictSlash(true)
	baseRouter.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	ctx := context.Background()
	gcpCredentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatal(err)
	}

	gcpStorage := storage.GcpStorage{
		ProjectId:  gcpCredentials.ProjectID,
		BucketName: "spaces-summit-famous-places",
	}
	topicDetails := pubsub.TopicDetails{
		ProjectId: gcpCredentials.ProjectID,
		TopicId:   "spaces-summit-famous-places",
	}

	topicDetails.CreateTopic()
	gcpStorage.CreateBucket(topicDetails)

	uploadHandler := handlers.MultipartUploadHandler{
		Storage: gcpStorage,
		MaxSize: 2 << 20, // 2MB max
	}
	baseRouter.Handle("/api/upload", uploadHandler)

	spa := handlers.SpaHandler{StaticPath: "web", IndexPath: "index.html"}
	baseRouter.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      baseRouter,
		Addr:         ":5080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
