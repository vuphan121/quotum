package quotum

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vuphan121/quotum/limiter"
	"github.com/vuphan121/quotum/util"
	"net/http"
	"time"
)

func main() {
	cfg := util.AppConfig
	router := mux.NewRouter()
	util.Log("Server starting on port 8080")

	router.HandleFunc("/request/{userID}", handleRequest(cfg)).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		util.Log("Error starting server")
	}
}

func handleRequest(cfg util.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["userID"]
		lim := limiter.CreateUserLimiter(userID, cfg.Rate, cfg.Interval, cfg.Store)

		allowed, bannedUntil := lim.Allow()
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
