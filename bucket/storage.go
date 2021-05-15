package bucket

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

type Store struct {
	ProjectID string
	bucket    *storage.BucketHandle
}

func New(ctx context.Context, projectID, bucketName string) (Store, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return Store{}, fmt.Errorf("%w: failed to create storage client", err)
	}

	return Store{
		ProjectID: projectID,
		bucket:    client.Bucket(bucketName),
	}, nil
}

func (s Store) Put(ctx context.Context, objectName string, data io.Reader) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	wc := s.bucket.Object(objectName).NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, data); err != nil {
		return err
	}
	return nil
}

func Create(ctx context.Context, projectID, bucketName, topicID string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed to create storage client", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)
	if _, err := bucket.Attrs(ctx); err != nil && !errors.Is(storage.ErrBucketNotExist, err) {
		fmt.Printf("Bucket %s already exists... skipping\n", bucketName)
		return nil
	}

	attrs := storage.BucketAttrs{Location: "europe-west4"}
	if err := bucket.Create(ctx, projectID, &attrs); err != nil {
		return fmt.Errorf("%w: failed to create bucket %s", err, bucketName)
	}

	fmt.Printf("Created bucket %v in %v with storage class %v\n", bucketName, attrs.Location, attrs.StorageClass)

	if _, err := bucket.AddNotification(ctx, &storage.Notification{
		TopicProjectID: projectID,
		TopicID:        topicID,
		PayloadFormat:  storage.JSONPayload,
	}); err != nil {
		return fmt.Errorf("%w: failed to add notification to %s", err, bucketName)
	}

	return nil
}
