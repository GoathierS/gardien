# gardien

**Sauvegarde et surveillance de la santé de ton PC, expliquées simplement.**
Pour Windows, Linux et macOS. Projet en tout début de développement.

gardien veille sur ton ordinateur et protège tes fichiers :

- il **surveille** les disques (SMART), les températures, les erreurs système et l'espace libre ;
- il **explique** en français ce qui ne va pas et quoi faire ;
- il **sauvegarde selon le risque** : si un disque se dégrade, il lance une sauvegarde et te prévient ;
- il **vérifie** que tes sauvegardes se restaurent vraiment.

Destinations prévues : disque externe, NAS (SMB, SFTP, WebDAV, S3) et cloud.
Suivi depuis le téléphone par l'interface web.

## État actuel

| Fonction | État |
|---|---|
| Interface web locale | ✅ |
| Espace disque | ✅ |
| Sauvegarde (moteur Kopia) | 🔜 v0.1 |
| SMART, températures, erreurs système | 🔜 v0.2 |
| Service Windows, macOS, Linux ; image Docker pour NAS | 🔜 v0.3 |
| Téléphone (PWA, notifications) | 🔜 v0.4 |

Feuille de route détaillée : [docs/architecture.md](docs/architecture.md).

## Essayer

Il faut [Go](https://go.dev/dl/) 1.27 ou plus récent.

```sh
git clone https://github.com/GoathierS/gardien.git
cd gardien
go run ./cmd/gardien serve
```

Puis ouvre <http://127.0.0.1:8484>.

## Compiler

```sh
go build ./cmd/gardien                          # pour ton OS
GOOS=linux   GOARCH=amd64 go build ./cmd/gardien
GOOS=darwin  GOARCH=arm64 go build ./cmd/gardien
GOOS=windows GOARCH=amd64 go build ./cmd/gardien
```

## Licence

[Apache 2.0](LICENSE).
