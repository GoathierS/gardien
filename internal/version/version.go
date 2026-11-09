// Package version expose la version de gardien, injectée à la compilation.
package version

import "runtime"

// Version est remplacée à la compilation :
//
//	go build -ldflags "-X github.com/GoathierS/gardien/internal/version.Version=v0.1.0" ./cmd/gardien
var Version = "dev"

// String renvoie la version avec l'OS et l'architecture.
func String() string {
	return "gardien " + Version + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
}
