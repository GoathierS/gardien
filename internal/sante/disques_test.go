package sante

import (
	"context"
	"strings"
	"testing"
)

func TestNiveauOccupation(t *testing.T) {
	cas := []struct {
		pourcent float64
		attendu  Niveau
	}{
		{0, NiveauOK},
		{79.9, NiveauOK},
		{80, NiveauAttention},
		{89.9, NiveauAttention},
		{90, NiveauCritique},
		{100, NiveauCritique},
	}
	for _, c := range cas {
		if got := NiveauOccupation(c.pourcent); got != c.attendu {
			t.Errorf("NiveauOccupation(%v) = %q, attendu %q", c.pourcent, got, c.attendu)
		}
	}
}

func TestPire(t *testing.T) {
	cas := []struct{ a, b, attendu Niveau }{
		{NiveauOK, NiveauOK, NiveauOK},
		{NiveauOK, NiveauAttention, NiveauAttention},
		{NiveauCritique, NiveauAttention, NiveauCritique},
		{NiveauOK, NiveauInconnu, NiveauInconnu},
		{NiveauInconnu, NiveauAttention, NiveauAttention},
	}
	for _, c := range cas {
		if got := Pire(c.a, c.b); got != c.attendu {
			t.Errorf("Pire(%q, %q) = %q, attendu %q", c.a, c.b, got, c.attendu)
		}
	}
}

func TestResumeEspace(t *testing.T) {
	cas := []struct {
		nom      string
		vols     []Volume
		contient string
	}{
		{"vide", nil, "Aucun volume"},
		{"ok", []Volume{{Point: "C:", Utilise: 40, Niveau: NiveauOK}}, "Tous les disques"},
		{"attention", []Volume{{Point: "D:", Utilise: 85, Niveau: NiveauAttention}}, "Bien rempli : D: (85 %)"},
		{"critique prioritaire", []Volume{
			{Point: "C:", Utilise: 95, Niveau: NiveauCritique},
			{Point: "D:", Utilise: 85, Niveau: NiveauAttention},
		}, "Presque plein : C: (95 %)"},
	}
	for _, c := range cas {
		if got := resumeEspace(c.vols); !strings.Contains(got, c.contient) {
			t.Errorf("%s : %q ne contient pas %q", c.nom, got, c.contient)
		}
	}
}

func TestVolumes(t *testing.T) {
	vols, err := Volumes(context.Background())
	if err != nil {
		t.Fatalf("Volumes : %v", err)
	}
	if len(vols) == 0 {
		t.Fatal("aucun volume détecté")
	}
	for _, v := range vols {
		if v.LibreOctets > v.TotalOctets {
			t.Errorf("%s : libre (%d) > total (%d)", v.Point, v.LibreOctets, v.TotalOctets)
		}
	}
}

func TestTous(t *testing.T) {
	rs := Tous(context.Background())
	if len(rs) != len(Capteurs()) {
		t.Fatalf("%d rapports pour %d capteurs", len(rs), len(Capteurs()))
	}
	for _, r := range rs {
		if r.Capteur == "" || r.Nom == "" || r.Resume == "" {
			t.Errorf("rapport incomplet : %+v", r)
		}
	}
}

type capteurEnPanne struct{}

func (capteurEnPanne) ID() string  { return "panne" }
func (capteurEnPanne) Nom() string { return "Capteur en panne" }
func (capteurEnPanne) Mesurer(context.Context) (Rapport, error) {
	return Rapport{}, context.DeadlineExceeded
}

func TestReleveErreur(t *testing.T) {
	r := Releve(context.Background(), capteurEnPanne{})
	if r.Niveau != NiveauInconnu || r.Erreur == "" || r.Capteur != "panne" {
		t.Errorf("une erreur doit donner un rapport de niveau inconnu : %+v", r)
	}
}
