// Commande gardien-bureau : l'application de bureau de gardien.
//
// Elle embarque l'agent et affiche son interface dans une fenêtre native (Wails).
// Compilation (voir CLAUDE.md) :
//
//	Windows : go build -tags desktop,production -ldflags "-w -s -H windowsgui" ./cmd/gardien-bureau
//	macOS   : CGO_LDFLAGS="-framework UniformTypeIdentifiers" go build -tags desktop,production ./cmd/gardien-bureau
//	Linux   : go build -tags desktop,production,webkit2_41 ./cmd/gardien-bureau  (paquets libgtk-3-dev et libwebkit2gtk-4.1-dev)
package main

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/GoathierS/gardien/internal/api"
)

func main() {
	err := wails.Run(&options.App{
		Title:     "gardien",
		Width:     1024,
		Height:    720,
		MinWidth:  720,
		MinHeight: 480,
		// L'API de l'agent sert aussi l'interface : la fenêtre appelle /api/v1/…
		// directement, sans ouvrir de port réseau.
		AssetServer: &assetserver.Options{Handler: api.Handler()},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "gardien-bureau :", err)
		os.Exit(1)
	}
}
