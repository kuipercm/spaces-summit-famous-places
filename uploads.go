package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuipercm/spaces-summit-famous-places/firestore"
	"net/http"
	"strconv"
	"time"
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
	qLimit := r.URL.Query().Get("limit")
	qOffset := r.URL.Query().Get("offset")
	qLastCreationDate := r.URL.Query().Get("creationDate")

	if qLastCreationDate != "" {
		m.byLastCreationDate(w, r, qLastCreationDate)
		return
	}

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

	res, err := m.fireStore.List(r.Context(), limit, offset)
	if err != nil {
		fmt.Printf("firestore::list %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (m fileHandler) byLastCreationDate(w http.ResponseWriter, r *http.Request, qLastCreationDate string) {
	layout := "2006-01-02T15:04:05.000Z"
	lastCreationDate, err := time.Parse(layout, qLastCreationDate)
	if err != nil {
		fmt.Printf("parse qLastCreationDate string to date %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := m.fireStore.ListByCreationDate(r.Context(), lastCreationDate)
	if err != nil {
		fmt.Printf("firestore::ListByCreationDate %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
