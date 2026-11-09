package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoathierS/gardien/internal/sante"
)

func requete(t *testing.T, methode, chemin string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, httptest.NewRequest(methode, chemin, nil))
	return rec
}

func TestSante(t *testing.T) {
	rec := requete(t, http.MethodGet, "/api/v1/sante")
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d, attendu 200 : %s", rec.Code, rec.Body)
	}
	var rs []sante.Rapport
	if err := json.Unmarshal(rec.Body.Bytes(), &rs); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if len(rs) == 0 {
		t.Fatal("aucun rapport renvoyé")
	}
}

func TestSanteCapteur(t *testing.T) {
	rec := requete(t, http.MethodGet, "/api/v1/sante/espace-disque")
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d, attendu 200 : %s", rec.Code, rec.Body)
	}
	var r sante.Rapport
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
		t.Fatalf("JSON invalide : %v", err)
	}
	if r.Capteur != "espace-disque" {
		t.Errorf("capteur %q, attendu espace-disque", r.Capteur)
	}
}

func TestCapteurInconnu(t *testing.T) {
	if rec := requete(t, http.MethodGet, "/api/v1/sante/nexiste-pas"); rec.Code != http.StatusNotFound {
		t.Fatalf("code %d, attendu 404", rec.Code)
	}
}

func TestInterface(t *testing.T) {
	rec := requete(t, http.MethodGet, "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("code %d, attendu 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<title>gardien</title>") {
		t.Error("la page d'accueil n'est pas servie")
	}
}

func TestMethodeRefusee(t *testing.T) {
	if rec := requete(t, http.MethodPost, "/api/v1/sante"); rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("code %d, attendu 405", rec.Code)
	}
}
