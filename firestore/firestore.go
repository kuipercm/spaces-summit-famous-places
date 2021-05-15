package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type Store struct {
	collection *firestore.CollectionRef
}

func (f Store) Add(ctx context.Context, filename string, content interface{}) error {
	_, err := f.collection.Doc(filename).Set(ctx, content)
	return err
}

func New(ctx context.Context, projectID string, collectionName string) (Store, error) {
	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return Store{}, fmt.Errorf("%w: failed to create firestore client", err)
	}

	return Store{
		collection: c.Collection(collectionName),
	}, nil
}
