package service

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/vuphan121/quotum/limiter"
	"github.com/vuphan121/quotum/storage"
	"github.com/vuphan121/quotum/util"
)

var (
	startTime    = time.Now()
	requestCount int64
	requestMux   sync.Mutex
)

func HandleRequest(cfg util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["userID"]
		lim := limiter.CreateUserLimiter(userID, cfg.Rate, cfg.Interval, cfg.Store, cfg.Algorithm)

		allowed, bannedUntil := lim.Allow()
		incrementRequestCount()

		w.Header().Set("Content-Type", "application/json")
		if allowed {
			json.NewEncoder(w).Encode(map[string]bool{"allowed": true})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"allowed":     false,
				"bannedUntil": bannedUntil.UTC().Format(time.RFC3339),
			})
		}
	}
}

func HandleStatus(cfg util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["userID"]
		state, err := cfg.Store.GetState(userID)

		w.Header().Set("Content-Type", "application/json")

		if err != nil || state.BannedUntil == nil || time.Now().After(*state.BannedUntil) {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "banned",
			"bannedUntil": state.BannedUntil.UTC().Format(time.RFC3339),
		})
	}
}

// make sure it supports both mem and redis
func HandleBanlist(cfg util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memStore, ok := cfg.Store.(*storage.MemoryStorage)
		w.Header().Set("Content-Type", "application/json")

		if !ok {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "banlist only supported for memory storage",
			})
			return
		}

		memStore.Lock()
		defer memStore.Unlock()

		now := time.Now()
		banned := []map[string]string{}

		for userID, state := range memStore.Data() {
			if state.BannedUntil != nil && now.Before(*state.BannedUntil) {
				banned = append(banned, map[string]string{
					"userID":      userID,
					"bannedUntil": state.BannedUntil.UTC().Format(time.RFC3339),
				})
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"bannedUsers": banned,
		})
	}
}

func HandleHealth(cfg util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uptime := int(time.Since(startTime).Seconds())
		storageType := "unknown"

		switch cfg.Store.(type) {
		case *storage.MemoryStorage:
			storageType = "memory"
		default:
			storageType = "unknown"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":            "ok",
			"uptimeSeconds":     uptime,
			"requestsProcessed": getRequestCount(),
			"storage":           storageType,
			"storageStatus":     "connected",
		})
	}
}

func incrementRequestCount() {
	requestMux.Lock()
	defer requestMux.Unlock()
	requestCount++
}

func getRequestCount() int64 {
	requestMux.Lock()
	defer requestMux.Unlock()
	return requestCount
}
