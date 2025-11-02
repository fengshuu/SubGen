package server

import (
	"fmt"
	"log"
	"net/http"

	"SubGen/internal/config"
	"SubGen/internal/fetch"
	"SubGen/internal/generator"
)

type Server struct {
	cfg *config.AppConfig
}

func New(cfg *config.AppConfig) *Server { return &Server{cfg: cfg} }

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "SubGen backend is running. Use /config to fetch YAML.")
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if s.cfg == nil {
			log.Println("config not loaded")
			http.Error(w, "config not loaded", http.StatusInternalServerError)
			return
		}
		base, err := fetch.BaseConfig(s.cfg.BaseConfigURL)
		if err != nil {
			log.Printf("failed to fetch base config: %v", err)
			http.Error(w, fmt.Sprintf("failed to fetch base config: %v", err), http.StatusBadGateway)
			return
		}
		final, err := generator.ReplaceProxyProvidersAndEncodeBase64(base, s.cfg.Subscriptions)
		if err != nil {
			log.Printf("failed to generate proxy-providers: %v", err)
			http.Error(w, fmt.Sprintf("failed to generate proxy-providers: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		fmt.Fprint(w, final)
	})

	return mux
}
