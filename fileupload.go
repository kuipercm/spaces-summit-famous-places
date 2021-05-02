package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kuipercm/spaces-summit-famous-places/bucket"
	"github.com/kuipercm/spaces-summit-famous-places/firestore"
	"github.com/kuipercm/spaces-summit-famous-places/vision"
)

type multipartUploadHandler struct {
	bucket    bucket.Store
	vision    vision.ImageIdentifier
	fireStore firestore.Store
	maxSize   int64
}

func newUploadHandler(b bucket.Store, v vision.ImageIdentifier, f firestore.Store, maxSize int64) multipartUploadHandler {
	return multipartUploadHandler{
		bucket:    b,
		vision:    v,
		fireStore: f,
		maxSize:   maxSize,
	}
}

func (m multipartUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(m.maxSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	for _, h := range r.MultipartForm.File["photo"] {
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
		if err = m.bucket.Put(ctx, fileName, bytes.NewReader(content)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		landmarks, err := m.vision.FindLandmarks(ctx, bytes.NewReader(content), fileName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := m.fireStore.Add(ctx, fileName, landmarks); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func readFileHeader(h *multipart.FileHeader) ([]byte, error) {
	f, err := h.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
