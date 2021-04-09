package storage

import (
	"context"
	"fmt"
	"github.com/kuipercm/spaces-summit-famous-places/pubsub"
	"google.golang.org/api/iterator"
	"io"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type GcpStorage struct {
	ProjectId  string
	BucketName string
}

func (s GcpStorage) CreateBucket(t pubsub.TopicDetails) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	buckets := client.Buckets(ctx, s.ProjectId)
	for {
		bucketAttrs, err := buckets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("client.Buckets(%q).Next: %v", s.ProjectId, err)
		}
		if bucketAttrs.Name == s.BucketName {
			fmt.Fprintf(os.Stdout, "Bucket %v already exists\n", s.BucketName)
			return nil
		}
	}

	storageClassAndLocation := &storage.BucketAttrs{
		Location: "europe-west4",
	}

	bucket := client.Bucket(s.BucketName)
	if err := bucket.Create(ctx, s.ProjectId, storageClassAndLocation); err != nil {
		return fmt.Errorf("Bucket(%q).Create: %v", s.BucketName, err)
	}
	fmt.Fprintf(os.Stdout, "Created bucket %v in %v with storage class %v\n", s.BucketName, storageClassAndLocation.Location, storageClassAndLocation.StorageClass)

	_, err = bucket.AddNotification(ctx, &storage.Notification{
		TopicProjectID: t.ProjectId,
		TopicID:        t.TopicId,
		PayloadFormat:  storage.JSONPayload,
	})
	if err != nil {
		return fmt.Errorf("Bucket(%q).AddNotification: %v", s.BucketName, err)
	}

	return nil
}

func (s GcpStorage) StoreFile(objectName string, h *multipart.FileHeader) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := h.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(s.BucketName).Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(os.Stdout, "Blob %v uploaded.\n", objectName)
	return nil
}
