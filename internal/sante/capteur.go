package sante

import (
	"context"
	"fmt"
)

// Capteur est une source de mesures de santé : espace disque, SMART, températures,
// erreurs système…
//
// Pour ajouter un capteur : implémenter cette interface dans un fichier du paquet
// (ou un sous-paquet), puis l'ajouter à la liste renvoyée par Capteurs.
// Un capteur indisponible sur un OS renvoie une erreur claire plutôt que de planter.
type Capteur interface {
	// ID est un identifiant stable, utilisé dans l'API : "espace-disque".
	ID() string
	// Nom est le libellé affiché à l'utilisateur : "Espace disque".
	Nom() string
	// Mesurer relève l'état actuel.
	Mesurer(ctx context.Context) (Rapport, error)
}

// Rapport est le résultat d'une mesure, prêt à être affiché.
type Rapport struct {
	Capteur string `json:"capteur"`
	Nom     string `json:"nom"`
	Niveau  Niveau `json:"niveau"`
	// Resume explique en une phrase ce qui se passe et, si besoin, quoi faire.
	Resume  string `json:"resume"`
	Details any    `json:"details,omitempty"`
	Erreur  string `json:"erreur,omitempty"`
}

// Capteurs renvoie les capteurs disponibles.
func Capteurs() []Capteur {
	return []Capteur{EspaceDisque{}}
}

// Trouver renvoie le capteur d'identifiant id.
func Trouver(id string) (Capteur, bool) {
	for _, c := range Capteurs() {
		if c.ID() == id {
			return c, true
		}
	}
	return nil, false
}

// Releve mesure un capteur. Une erreur devient un rapport de niveau inconnu,
// pour que l'interface puisse l'afficher au lieu d'échouer.
func Releve(ctx context.Context, c Capteur) Rapport {
	r, err := c.Mesurer(ctx)
	if err != nil {
		return Rapport{
			Capteur: c.ID(),
			Nom:     c.Nom(),
			Niveau:  NiveauInconnu,
			Resume:  fmt.Sprintf("Impossible de mesurer : %v", err),
			Erreur:  err.Error(),
		}
	}
	r.Capteur, r.Nom = c.ID(), c.Nom()
	return r
}

// Tous mesure tous les capteurs.
func Tous(ctx context.Context) []Rapport {
	cs := Capteurs()
	rs := make([]Rapport, 0, len(cs))
	for _, c := range cs {
		rs = append(rs, Releve(ctx, c))
	}
	return rs
}
