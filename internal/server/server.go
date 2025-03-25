package server

import (
	"net/http"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
)

func StartServer(cfg *config.AppConfig) error {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	return http.ListenAndServe(":"+cfg.Server.Port, nil)
}
