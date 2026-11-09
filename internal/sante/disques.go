// Package sante surveille l'état du PC au moyen de capteurs (voir Capteur).
package sante

import (
	"context"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

// Niveau résume un état pour l'utilisateur.
type Niveau string

const (
	NiveauOK        Niveau = "ok"
	NiveauAttention Niveau = "attention"
	NiveauCritique  Niveau = "critique"
	NiveauInconnu   Niveau = "inconnu"
)

// gravite ordonne les niveaux, du moins au plus grave.
var gravite = map[Niveau]int{NiveauOK: 0, NiveauInconnu: 1, NiveauAttention: 2, NiveauCritique: 3}

// Pire renvoie le plus grave des deux niveaux.
func Pire(a, b Niveau) Niveau {
	if gravite[b] > gravite[a] {
		return b
	}
	return a
}

// Seuils d'occupation d'un volume, en pourcentage.
const (
	SeuilAttention = 80.0
	SeuilCritique  = 90.0
)

// Volume décrit l'occupation d'un volume monté.
type Volume struct {
	Point       string  `json:"point"`
	Systeme     string  `json:"systeme"`
	TotalOctets uint64  `json:"total_octets"`
	LibreOctets uint64  `json:"libre_octets"`
	Utilise     float64 `json:"utilise_pourcent"`
	Niveau      Niveau  `json:"niveau"`
}

// NiveauOccupation classe un pourcentage d'occupation.
func NiveauOccupation(pourcent float64) Niveau {
	switch {
	case pourcent >= SeuilCritique:
		return NiveauCritique
	case pourcent >= SeuilAttention:
		return NiveauAttention
	default:
		return NiveauOK
	}
}

// Volumes renvoie l'occupation des volumes physiques montés.
// Un volume illisible (lecteur vide, disque débranché) est ignoré.
func Volumes(ctx context.Context) ([]Volume, error) {
	parts, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("liste des volumes : %w", err)
	}

	vols := make([]Volume, 0, len(parts))
	for _, p := range parts {
		u, err := disk.UsageWithContext(ctx, p.Mountpoint)
		if err != nil || u.Total == 0 {
			continue
		}
		vols = append(vols, Volume{
			Point:       p.Mountpoint,
			Systeme:     p.Fstype,
			TotalOctets: u.Total,
			LibreOctets: u.Free,
			Utilise:     u.UsedPercent,
			Niveau:      NiveauOccupation(u.UsedPercent),
		})
	}
	return vols, nil
}

// EspaceDisque est le capteur d'occupation des volumes.
type EspaceDisque struct{}

func (EspaceDisque) ID() string  { return "espace-disque" }
func (EspaceDisque) Nom() string { return "Espace disque" }

func (EspaceDisque) Mesurer(ctx context.Context) (Rapport, error) {
	vols, err := Volumes(ctx)
	if err != nil {
		return Rapport{}, err
	}
	niveau := NiveauOK
	for _, v := range vols {
		niveau = Pire(niveau, v.Niveau)
	}
	return Rapport{Niveau: niveau, Resume: resumeEspace(vols), Details: vols}, nil
}

// resumeEspace écrit la phrase affichée à l'utilisateur.
func resumeEspace(vols []Volume) string {
	if len(vols) == 0 {
		return "Aucun volume détecté."
	}
	var critiques, attention []string
	for _, v := range vols {
		texte := fmt.Sprintf("%s (%.0f %%)", v.Point, v.Utilise)
		switch v.Niveau {
		case NiveauCritique:
			critiques = append(critiques, texte)
		case NiveauAttention:
			attention = append(attention, texte)
		}
	}
	switch {
	case len(critiques) > 0:
		return "Presque plein : " + strings.Join(critiques, ", ") + ". Libère de la place rapidement."
	case len(attention) > 0:
		return "Bien rempli : " + strings.Join(attention, ", ") + ". Pense à faire du tri."
	default:
		return "Tous les disques ont de la place."
	}
}
