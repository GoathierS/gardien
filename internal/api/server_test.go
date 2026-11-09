package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoathierS/gardien/internal/sante"
)

func TestVolumes(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/sante/volumes", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("code %d, attendu 200 : %s", rec.Code, rec.Body)
	}
	var vols []sante.Volume
	if err := json.Unmarshal(rec.Body.Bytes(), &vols); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if len(vols) == 0 {
		t.Fatal("aucun volume renvoyé")
	}
}

func TestInterface(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("code %d, attendu 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<title>gardien</title>") {
		t.Error("la page d'accueil n'est pas servie")
	}
}

func TestMethodeRefusee(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/sante/volumes", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("code %d, attendu 405", rec.Code)
	}
}
