// Package rest
package rest

import (
	"net/http"

	"horizonx-server/internal/config"
	"horizonx-server/internal/transport/rest/middleware"
	"horizonx-server/internal/transport/websocket"
)

type RouterDeps struct {
	WS      *websocket.Handler
	Auth    *AuthHandler
	Metrics *MetricsHandler
}

func NewRouter(cfg *config.Config, deps *RouterDeps) http.Handler {
	mux := http.NewServeMux()

	globalMw := middleware.New()
	globalMw.Use(middleware.CORS(cfg))
	globalMw.Use(middleware.CSRF(cfg))

	authMw := middleware.New()
	authMw.Use(middleware.JWT(cfg))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /auth/register", deps.Auth.Register)
	mux.HandleFunc("POST /auth/login", deps.Auth.Login)
	mux.HandleFunc("POST /auth/logout", deps.Auth.Logout)

	mux.HandleFunc("/ws", deps.WS.Serve)
	mux.Handle("GET /metrics", authMw.Then(http.HandlerFunc(deps.Metrics.Get)))

	// Placeholder for new feature routes
	// mux.HandleFunc("/ssh", handler.HandleSSH)
	// mux.HandleFunc("/deploy", handler.HandleDeploy)

	return globalMw.Apply(mux)
}
