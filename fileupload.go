package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kuipercm/spaces-summit-famous-places/bucket"
	"github.com/kuipercm/spaces-summit-famous-places/firestore"
	"github.com/kuipercm/spaces-summit-famous-places/vision"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type multipartUploadHandler struct {
	bucket    bucket.Store
	vision    vision.ImageIdentifier
	fireStore firestore.Store
	maxSize   int64

	tracer trace.Tracer
}

func newUploadHandler(b bucket.Store, v vision.ImageIdentifier, f firestore.Store, maxSize int64) multipartUploadHandler {
	t := otel.Tracer("ssfp/uploads")

	return multipartUploadHandler{
		bucket:    b,
		vision:    v,
		fireStore: f,
		maxSize:   maxSize,
		tracer:    t,
	}
}

func (m multipartUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(m.maxSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	res := make([]vision.ImageRecord, 0, len(r.MultipartForm.File["photos"]))
	for _, h := range r.MultipartForm.File["photos"] {
		fileId, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		content, err := readFileHeader(h)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fileName := fileId.String() + filepath.Ext(h.Filename)
		err = m.bucket.Put(ctx, fileName, bytes.NewReader(content))
		if err != nil {
			fmt.Printf("bucket::Put %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		landmarks, err := m.vision.FindLandmarks(ctx, bytes.NewReader(content), fileName)
		if err != nil {
			fmt.Printf("findLandmarks %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = m.fireStore.Add(ctx, fileName, landmarks)
		if err != nil {
			fmt.Printf("firestore::Add %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res = append(res, landmarks)
	}

	json.NewEncoder(w).Encode(res)
}

func readFileHeader(h *multipart.FileHeader) ([]byte, error) {
	f, err := h.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
