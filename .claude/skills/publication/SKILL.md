---
name: publication
description: Prépare une nouvelle version de gardien : vérifications, changelog, tag et compilation pour Windows, Linux et macOS. À utiliser quand on veut publier une version (v0.1.0, etc.).
---

# Publier une version

Demande toujours confirmation avant de pousser un tag ou de créer une release sur GitHub : c'est public.

1. **Vérifier l'état** : branche `main` à jour, aucun changement non commité.
2. **Contrôles** : `gofmt -l .` vide, `go vet ./...`, `go test ./...`.
3. **Compilation croisée** : `windows/amd64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, avec `CGO_ENABLED=0` et
   `-ldflags "-X github.com/GoathierS/gardien/internal/version.Version=<version>"`.
4. **Changelog** : ajoute une section en haut de `CHANGELOG.md` (créé s'il n'existe pas) à partir de `git log <tag précédent>..HEAD`, regroupée en Ajouts, Corrections, Changements. Rédige pour un utilisateur, pas pour un développeur.
5. **Feuille de route** : coche ce qui est fait dans `docs/architecture.md`.
6. **Commit** « Publie <version> », puis tag annoté `<version>`.
7. **Après confirmation** : `git push` et `git push origin <version>`, puis `gh release create <version>` avec le changelog et les exécutables.

Quand GoReleaser sera configuré, les étapes 3 et 7 passeront par `goreleaser release`.
