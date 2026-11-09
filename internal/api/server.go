// Package api sert l'interface web et l'API locale de gardien.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/GoathierS/gardien/internal/sante"
	"github.com/GoathierS/gardien/internal/version"
	"github.com/GoathierS/gardien/web"
)

// Handler construit les routes de l'API et de l'interface.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/version", func(w http.ResponseWriter, r *http.Request) {
		ecrireJSON(w, http.StatusOK, map[string]string{"version": version.String()})
	})
	mux.HandleFunc("GET /api/sante/volumes", func(w http.ResponseWriter, r *http.Request) {
		vols, err := sante.Volumes(r.Context())
		if err != nil {
			ecrireJSON(w, http.StatusInternalServerError, map[string]string{"erreur": err.Error()})
			return
		}
		ecrireJSON(w, http.StatusOK, vols)
	})
	mux.Handle("GET /", http.FileServerFS(web.Fichiers()))
	return mux
}

// Serve écoute sur addr jusqu'à l'annulation de ctx.
// Par défaut gardien n'écoute que sur 127.0.0.1 : rien n'est exposé au réseau
// tant qu'il n'y a pas d'authentification.
func Serve(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		<-ctx.Done()
		arret, annuler := context.WithTimeout(context.Background(), 5*time.Second)
		defer annuler()
		_ = srv.Shutdown(arret)
	}()

	slog.Info("interface disponible", "url", "http://"+addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func ecrireJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
