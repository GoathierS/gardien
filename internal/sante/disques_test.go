package sante

import (
	"context"
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
