package vision

import (
	"bytes"
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"time"
)

type ImageIdentifier struct {
	ProjectId string
}

type ImageRecord struct {
	Filename     string
	Landmarks    []string
	CreationDate time.Time
}

func (i ImageIdentifier) FindLandmarks(imageBytes []byte, fileName string) (*ImageRecord, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	image, err := vision.NewImageFromReader(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	annotations, err := client.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	var landmarks []string
	if len(annotations) == 0 {
		landmarks = make([]string, 1)
	} else {
		landmarks = make([]string, len(annotations))
		for i, annotation := range annotations {
			landmarks[i] = annotation.Description
		}
	}

	record := ImageRecord{
		Filename:     fileName,
		Landmarks:    landmarks,
		CreationDate: time.Now(),
	}

	return &record, nil
}
