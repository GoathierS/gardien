// Package web embarque l'interface web dans l'exécutable.
package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var static embed.FS

// Fichiers renvoie les fichiers de l'interface, racine = dossier static.
func Fichiers() fs.FS {
	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err) // impossible : le dossier est embarqué à la compilation
	}
	return sub
}
