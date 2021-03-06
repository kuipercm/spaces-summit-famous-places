package vision

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
)

type ImageRecord struct {
	Filename     string
	Landmarks    []string `firestore:",omitempty"`
	CreationDate time.Time
	Image        string
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
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	image, err := vision.NewImageFromReader(tee)
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

	b, err := io.ReadAll(&buf)
	if err != nil {
		return ImageRecord{}, err
	}

	return ImageRecord{
		Filename:     fileName,
		Landmarks:    landmarks,
		CreationDate: time.Now(),
		Image:        "data:image/png;base64," + base64.StdEncoding.EncodeToString(b),
	}, nil
}
