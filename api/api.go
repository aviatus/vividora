package api

import (
	"aviatus/vividora/internal/config"
	"aviatus/vividora/internal/errors"
	"aviatus/vividora/internal/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func StartServer() error {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/Save", handleSaveRequest)
	http.HandleFunc("/Restore", handleRestoreRequest)

	err := http.ListenAndServe(config.ExternalPort, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func handleRestoreRequest(w http.ResponseWriter, r *http.Request) {
	if config.Mode == config.REPLICA {
		http.Error(w, errors.MethodNotAllowedError, http.StatusMethodNotAllowed)
	}

	err := store.RestoreFromSnapshot(config.SnapshotPath + "snapshot-2023-06-27-22:08:35.vdb")
	if err != nil {
		http.Error(w, errors.SnapshotRestoreError, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func handleSaveRequest(w http.ResponseWriter, r *http.Request) {
	if config.Mode == config.REPLICA {
		http.Error(w, errors.MethodNotAllowedError, http.StatusMethodNotAllowed)
	}

	err := store.TakeSnapshot()
	if err != nil {
		http.Error(w, errors.SnapshotTakeError, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch config.Mode {
	case config.REPLICA:
		switch r.Method {
		case http.MethodGet:
			getValue(w, r)
		default:
			http.Error(w, errors.MethodNotAllowedError, http.StatusMethodNotAllowed)
		}
	case config.MASTER:
		switch r.Method {
		case http.MethodGet:
			getValue(w, r)
		case http.MethodPost:
			setValue(w, r)
		case http.MethodPut:
			updateValue(w, r)
		case http.MethodDelete:
			deleteValue(w, r)
		default:
			http.Error(w, errors.MethodNotAllowedError, http.StatusMethodNotAllowed)
		}
	}
}

func getValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:] // Get the key from the URL path
	// Check if the key exists in the key-value store
	value, ok := store.Get(key)
	if !ok {
		http.Error(w, errors.KeyNotFoundError, http.StatusNotFound)
		return
	}

	// Return the value as the response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}

func setValue(w http.ResponseWriter, r *http.Request) {
	// Read the request body to get the key-value pair
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Extract the key-value pair from the request body
	var keyValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	err = json.Unmarshal(body, &keyValue)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if keyValue.Key == "" {
		http.Error(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	if len(keyValue.Key+keyValue.Value) > config.MaxItemSize {
		http.Error(w, "Key size exceeds maximum key size", http.StatusBadRequest)
		return
	}

	// Store the key-value pair in the key-value store
	err = store.Set(keyValue.Key, keyValue.Value)
	if err != nil {
		http.Error(w, "Failed to set key-value pair", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Value set successfully")
}

func updateValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:] // Get the key from the URL path

	// Check if the key exists in the key-value store
	_, ok := store.Get(key)
	if !ok {
		http.Error(w, errors.KeyNotFoundError, http.StatusNotFound)
		return
	}

	// Read the request body to get the updated value
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Update the value in the key-value store
	store.Set(key, string(body))

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Value updated successfully")
}

func deleteValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:] // Get the key from the URL path

	// Check if the key exists in the key-value store
	_, ok := store.Get(key)
	if !ok {
		http.Error(w, errors.KeyNotFoundError, http.StatusNotFound)
		return
	}

	// Delete the key from the key-value store
	store.Delete(key)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Value deleted successfully")
}
