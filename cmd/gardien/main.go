// Commande gardien : sauvegarde et surveillance de la santé du PC.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/GoathierS/gardien/internal/api"
	"github.com/GoathierS/gardien/internal/version"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "serve":
		fs := flag.NewFlagSet("serve", flag.ExitOnError)
		addr := fs.String("addr", "127.0.0.1:8484", "adresse d'écoute de l'interface web")
		_ = fs.Parse(os.Args[2:])

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if err := api.Serve(ctx, *addr); err != nil {
			slog.Error("arrêt du serveur", "err", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println(version.String())
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Utilisation : gardien <commande>

Commandes :
  serve     lance le service et l'interface web (http://127.0.0.1:8484)
  version   affiche la version`)
}
