package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kuipercm/spaces-summit-famous-places/firestore"
	"github.com/kuipercm/spaces-summit-famous-places/handlers"
	"github.com/kuipercm/spaces-summit-famous-places/pubsub"
	"github.com/kuipercm/spaces-summit-famous-places/storage"
	"github.com/kuipercm/spaces-summit-famous-places/vision"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
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

	imageIdentifier := vision.ImageIdentifier{
		ProjectId: gcpCredentials.ProjectID,
	}
	imageRecordsFirestoreBackend := firestore.Firestore{
		ProjectId:      gcpCredentials.ProjectID,
		CollectionName: "spaces-summit-famous-places",
	}

	uploadHandler := handlers.MultipartUploadHandler{
		Storage:          gcpStorage,
		ImageIdentifier:  imageIdentifier,
		FireStoreBackend: imageRecordsFirestoreBackend,
		MaxSize:          2 << 20, // 2MB max
	}
	baseRouter.Handle("/api/upload", uploadHandler)

	spa := handlers.SpaHandler{StaticPath: "web", IndexPath: "index.html"}
	baseRouter.PathPrefix("/").Handler(spa)

	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler:      baseRouter,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
