package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/kuipercm/spaces-summit-famous-places/vision"
	"log"
)

type Firestore struct {
	ProjectId      string
	CollectionName string
}

func (f Firestore) AddImage(record *vision.ImageRecord) error {
	ctx := context.Background()
	client := f.createClient(ctx)

	_, err := client.Collection(f.CollectionName).Doc(record.Filename).Set(ctx, record)
	return err
}

func (f Firestore) createClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, f.ProjectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}
