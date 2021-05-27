package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kuipercm/spaces-summit-famous-places/firestore"
)

type fileHandler struct {
	fireStore firestore.Store
}

func newFileHandler(f firestore.Store) fileHandler {
	return fileHandler{
		fireStore: f,
	}
}

func (m fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	qLastCreationDate := r.URL.Query().Get("creationDate")
	if qLastCreationDate != "" {
		m.byLastCreationDate(ctx, w, r, qLastCreationDate)
		return
	}

	qLimit := r.URL.Query().Get("limit")
	qOffset := r.URL.Query().Get("offset")
	if qLimit == "" || qOffset == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("If not specifing by creationDate, you should specify both limit and offset. offset = " + qOffset + ", qLimit = " + qLimit))
		return
	}

	limit, err := strconv.Atoi(qLimit)
	if err != nil {
		fmt.Printf("parse qLimit string to int %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	offset, err := strconv.Atoi(qOffset)
	if err != nil {
		fmt.Printf("parse qOffset string to int %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := m.fireStore.List(ctx, limit, offset)
	if err != nil {
		fmt.Printf("firestore::list %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (m fileHandler) byLastCreationDate(ctx context.Context, w http.ResponseWriter, r *http.Request, qLastCreationDate string) {
	lastCreationDateMillis, err := strconv.ParseInt(qLastCreationDate, 10, 64)
	if err != nil {
		fmt.Printf("parse qLastCreationDate string to int %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lastCreationDate := time.Unix(0, lastCreationDateMillis*int64(time.Millisecond))
	res, err := m.fireStore.ListByCreationDate(ctx, lastCreationDate)
	if err != nil {
		fmt.Printf("firestore::ListByCreationDate %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
