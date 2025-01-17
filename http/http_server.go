package http

import (
	"database/sql"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func HttpServer (db *sql.DB, chainConfig Chain) {

	// create a mux/router for handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		QueryHandler(w,r,db, chainConfig.SolverAddress)
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the HTTP server in a separate goroutine - to alllow for graceful shutdown
	go func() {
		log.Logger.Info().Str("address", srv.Addr).Msg("starting HTTP server")
		if err := srv.ListenAndServe(); err != nil {
			log.Logger.Error().Err(err).Msg("failed to start HTTP server")
		}
	}()
}
