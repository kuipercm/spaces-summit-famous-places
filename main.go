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
	"go.opentelemetry.io/otel"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Interrupt)
	defer cancel()

	projectID := os.Getenv("GCP_PROJECT_ID")

	pubsub.CreateTopic(ctx, projectID, "spaces-summit-famous-places")
	bucket.Create(ctx, projectID, "spaces-summit-famous-places", "spaces-summit-famous-places")

	exporter, err := trace.NewExporter(trace.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("trace::NewExporter: %v", err)
	}
	defer exporter.Shutdown(ctx) // flushes any pending spans

	// WithBatcher to batch trace exports
	// WithTraces == export trace after it's collected i.e. add 500ms latency per request..
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

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
	fileHandler := newFileHandler(firestore)                                         // 2MB max

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok!"))
	})
	router.Handle("/api/uploads", uploadHandler).Methods("POST")
	router.Handle("/api/uploads", fileHandler).Methods("GET")
	router.PathPrefix("/").Handler(spa)

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
