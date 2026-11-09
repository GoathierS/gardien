# gardien

Sauvegarde et surveillance de la santé du PC, pour Windows, Linux et macOS.
Un démon Go sert une interface web locale. Moteur de sauvegarde prévu : Kopia, utilisé comme bibliothèque.

## Commandes

```sh
go build ./cmd/gardien          # compile
go run ./cmd/gardien serve      # lance l'interface : http://127.0.0.1:8484
go test ./...                   # tests
go vet ./...                    # analyse statique
gofmt -l .                      # fichiers mal formatés (doit être vide)
GOOS=linux GOARCH=amd64 go build ./cmd/gardien   # compilation croisée (darwin, windows aussi)
```

Sous PowerShell, la compilation croisée s'écrit `$env:GOOS='linux'; go build ./cmd/gardien`.

## Structure

- `cmd/gardien` : point d'entrée et sous-commandes (`serve`, `version`).
- `internal/api` : routes HTTP. Les routes JSON sont sous `/api/`, le reste sert l'interface.
- `internal/sante` : surveillance (espace disque ; à venir : SMART, températures, journaux système).
- `internal/version` : version injectée par `-ldflags`.
- `web/static` : interface web, embarquée dans l'exécutable avec `embed`.
- `docs/architecture.md` : architecture et feuille de route.

## Règles absolues

1. **Ne jamais écrire, déplacer ni supprimer dans les dossiers sources de l'utilisateur.** La sauvegarde lit, elle ne modifie pas.
2. **Toute fonction de sauvegarde ou de restauration a un test « sauvegarder puis restaurer, puis comparer les empreintes ».** Voir la skill `test-restauration`.
3. **Pas d'écoute réseau hors de 127.0.0.1** tant que l'authentification n'existe pas.
4. **Aucun secret dans le dépôt** : mots de passe de dépôt, clés, jetons.
5. **Pas de dépendance CGO** sans en discuter : la compilation croisée doit rester `GOOS=… go build`.
6. Le code propre à un OS va dans des fichiers suffixés (`_windows.go`, `_linux.go`, `_darwin.go`), avec une version de repli pour les autres OS.

## Conventions

- Identifiants, commentaires, messages et interface en **français**.
- Les erreurs sont enveloppées avec du contexte : `fmt.Errorf("lecture de %s : %w", chemin, err)`.
- Journalisation avec `log/slog`.
- Tests à côté du code (`x_test.go`), en table quand il y a plusieurs cas.
- Un commit par changement cohérent. Messages en français, à l'impératif : « Ajoute la lecture SMART ».

## Avant de terminer une tâche

`gofmt -l .` vide, `go vet ./...` et `go test ./...` passent.
