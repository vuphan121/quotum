package main

import (
	"github.com/gorilla/mux"
	"github.com/vuphan121/quotum/service"
	"github.com/vuphan121/quotum/util"
	"net/http"
)

func main() {
	cfg := util.AppConfig
	router := mux.NewRouter()
	util.Log("Server starting on port 8080")

	auth := util.AuthMiddleware(cfg.APIKey)

	router.HandleFunc("/request/{userID}", service.HandleRequest(cfg)).Methods("POST")

	router.Handle("/status/{userID}", auth(service.HandleStatus(cfg))).Methods("GET")
	router.Handle("/banlist", auth(service.HandleBanlist(cfg))).Methods("GET")
	router.Handle("/health", auth(service.HandleHealth(cfg))).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		util.Log("Error starting server")
	}
}
