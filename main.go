package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kuipercm/spaces-summit-famous-places/bucket"
	"github.com/kuipercm/spaces-summit-famous-places/firestore"
	"github.com/kuipercm/spaces-summit-famous-places/pubsub"
	"github.com/kuipercm/spaces-summit-famous-places/vision"
	"github.com/kuipercm/spaces-summit-famous-places/web"
	"golang.org/x/oauth2/google"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Interrupt)
	defer cancel()

	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatal(err)
	}
	projectID := creds.ProjectID

	pubsub.CreateTopic(ctx, projectID, "spaces-summit-famous-places")
	bucket.Create(ctx, projectID, "spaces-summit-famous-places", "spaces-summit-famous-places")

	gcpStorage, err := bucket.New(ctx, projectID, "spaces-summit-famous-places")
	if err != nil {
		log.Fatal(err)
	}

	imageIdentifier, err := vision.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	firestore, err := firestore.New(ctx, projectID, "spaces-summit-famous-places")
	if err != nil {
		log.Fatal(err)
	}

	spa, err := web.NewHandler()
	if err != nil {
		log.Fatal(err)
	}
	uploadHandler := newUploadHandler(gcpStorage, imageIdentifier, firestore, 2<<20) // 2MB max

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").Handler(spa)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok!"))
	})
	router.Handle("/api/upload", uploadHandler)
	port := os.Getenv("PORT")
	srv := http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		fmt.Println("Received shutdown signal, shutting down..")
		srv.Shutdown(context.Background())
	}()

	fmt.Println("Listening on " + port)
	log.Fatal(srv.ListenAndServe())
}
