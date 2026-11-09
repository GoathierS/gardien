# gardien

Sauvegarde et surveillance de la santé du PC, pour Windows, Linux et macOS, avec agents sur NAS
et appli téléphone. Moteur de sauvegarde prévu : Kopia, utilisé comme bibliothèque.

## Les trois briques

1. **L'agent** (`cmd/gardien`) : le cœur, sans interface. Surveille, sauvegarde, expose l'API
   `/api/v1`. Tourne sur PC, NAS, serveur. **Sans CGO.**
2. **L'appli de bureau** (`cmd/gardien-bureau`) : embarque l'agent et affiche l'interface dans une
   fenêtre native avec Wails v2. Besoin de CGO sur Linux et macOS.
3. **L'appli téléphone** (`mobile/`) : Capacitor, qui reprend la même interface et parle à un agent
   du réseau.

**L'interface est écrite une seule fois**, dans `web/static`, et sert aux trois.

## Commandes

```sh
go test ./...                   # tests
go vet ./...                    # analyse statique
gofmt -l .                      # fichiers mal formatés (doit être vide)

go run ./cmd/gardien serve      # agent + interface dans le navigateur : http://127.0.0.1:8484
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build ./cmd/gardien   # agent pour un NAS ARM

# Appli de bureau
go build -tags desktop,production -ldflags "-w -s -H windowsgui" ./cmd/gardien-bureau   # Windows
CGO_LDFLAGS="-framework UniformTypeIdentifiers" go build -tags desktop,production ./cmd/gardien-bureau   # macOS
go build -tags desktop,production,webkit2_41 ./cmd/gardien-bureau   # Linux (libgtk-3-dev, libwebkit2gtk-4.1-dev)

# Téléphone
cd mobile && npm install && npm run ui   # puis npx cap add android / npm run android
```

Sous PowerShell : `$env:GOOS='linux'; $env:CGO_ENABLED='0'; go build ./cmd/gardien`.
L'appli de bureau ne se compile pas en croisé : chaque OS est compilé sur sa machine (voir la CI).

## Structure

- `cmd/gardien` : agent, sous-commandes `serve` et `version`.
- `cmd/gardien-bureau` : appli de bureau Wails. L'API de l'agent sert d'`AssetServer.Handler`.
- `internal/api` : routes HTTP. JSON sous `/api/v1/`, le reste sert l'interface.
- `internal/sante` : surveillance par **capteurs** (interface `Capteur`). Pour en ajouter un :
  implémenter `Capteur`, l'ajouter dans `Capteurs()`, le tester, et ajouter si besoin un rendu
  dédié dans `rendus` de `web/static/app.js`.
- `internal/version` : version injectée par `-ldflags`.
- `web/static` : interface commune (HTML, CSS, JS sans étape de compilation), embarquée avec `embed`.
- `mobile/` : projet Capacitor. `www/` est une copie générée de `web/static` : ne pas la modifier.
- `docs/architecture.md` : architecture et feuille de route.

## Règles absolues

1. **Ne jamais écrire, déplacer ni supprimer dans les dossiers sources de l'utilisateur.** La sauvegarde lit, elle ne modifie pas.
2. **Toute fonction de sauvegarde ou de restauration a un test « sauvegarder puis restaurer, puis comparer les empreintes ».** Voir la skill `test-restauration`.
3. **L'agent n'écoute que sur 127.0.0.1** tant que l'authentification n'existe pas.
4. **Aucun secret dans le dépôt** : mots de passe de dépôt, clés, jetons.
5. **L'agent reste sans CGO** pour que `GOOS=… CGO_ENABLED=0 go build ./cmd/gardien` marche partout. Seule l'appli de bureau peut utiliser CGO.
6. **L'API `/api/v1` ne change pas de façon incompatible** : ajouter des champs, oui ; en retirer ou en renommer, non (passer à `/api/v2`).
7. Code propre à un OS dans des fichiers suffixés (`_windows.go`, `_linux.go`, `_darwin.go`), avec une version de repli.
8. Dans l'interface, aucune donnée de l'API insérée avec `innerHTML` : utiliser `textContent`.

## Conventions

- Identifiants, commentaires, messages et interface en **français**.
- Les erreurs sont enveloppées avec du contexte : `fmt.Errorf("lecture de %s : %w", chemin, err)`.
- Journalisation avec `log/slog`.
- Tests à côté du code (`x_test.go`), en table quand il y a plusieurs cas.
- Un commit par changement cohérent. Messages en français, à l'impératif : « Ajoute la lecture SMART ».

## Avant de terminer une tâche

`gofmt -l .` vide, `go vet ./...` et `go test ./...` passent. Si l'interface a changé, vérifier
aussi l'appli de bureau.
