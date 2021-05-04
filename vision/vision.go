package vision

import (
	"context"
	"io"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
)

type ImageRecord struct {
	Filename     string
	Landmarks    []string
	CreationDate time.Time
}

type ImageIdentifier struct {
	vision *vision.ImageAnnotatorClient
}

func New(ctx context.Context) (ImageIdentifier, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return ImageIdentifier{}, err
	}

	return ImageIdentifier{
		vision: client,
	}, nil
}

func (i ImageIdentifier) FindLandmarks(ctx context.Context, r io.Reader, fileName string) (ImageRecord, error) {
	image, err := vision.NewImageFromReader(r)
	if err != nil {
		return ImageRecord{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	annotations, err := i.vision.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		return ImageRecord{}, err
	}

	landmarks := make([]string, 0, len(annotations))
	for i := range annotations {
		landmarks = append(landmarks, annotations[i].Description)
	}

	return ImageRecord{
		Filename:     fileName,
		Landmarks:    landmarks,
		CreationDate: time.Now(),
	}, nil
}
