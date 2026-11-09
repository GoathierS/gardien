// Package sante surveille l'état du PC : espace disque pour l'instant,
// puis SMART, températures et erreurs système.
package sante

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

// Niveau résume un état pour l'utilisateur.
type Niveau string

const (
	NiveauOK        Niveau = "ok"
	NiveauAttention Niveau = "attention"
	NiveauCritique  Niveau = "critique"
)

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
		return nil, err
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
