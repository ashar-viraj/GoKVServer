package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"myserver/cache"
	"myserver/db"
)

var mu sync.RWMutex

var Cache = cache.NewCache()

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   int    `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, err := db.DB.Exec("INSERT INTO kvstore (key, value) VALUES ($1, $2)", req.Key, req.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Insert failed: %v", err), http.StatusConflict)
		return
	}

	Cache.Put(req.Key, req.Value)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created key %d", req.Key)
}

func Read(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	keyStr := r.URL.Query().Get("key")
	key, err := strconv.Atoi(keyStr)
	if err != nil {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	if val, ok := Cache.Get(key); ok {
		resp := map[string]string{"value": val, "source": "cache"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var value string
	err = db.DB.QueryRow("SELECT value FROM kvstore WHERE key = $1", key).Scan(&value)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	Cache.Put(key, value)
	resp := map[string]string{"value": value}
	json.NewEncoder(w).Encode(resp)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   int    `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	res, err := db.DB.Exec("UPDATE kvstore SET value = $1 WHERE key = $2", req.Value, req.Key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Update failed: %v", err), http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	Cache.Put(req.Key, req.Value)
	fmt.Fprintf(w, "Updated key %d", req.Key)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key int `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	res, err := db.DB.Exec("DELETE FROM kvstore WHERE key = $1", req.Key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Delete failed: %v", err), http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	Cache.Delete(req.Key)
	fmt.Fprintf(w, "Deleted key %d", req.Key)
}
