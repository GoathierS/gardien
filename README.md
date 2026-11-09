# gardien

**Sauvegarde et surveillance de la santé de ton PC, expliquées simplement.**
Windows, Linux et macOS, avec agents pour NAS et appli téléphone. Projet en tout début de développement.

gardien veille sur tes appareils et protège tes fichiers :

- il **surveille** les disques (SMART), les températures, les erreurs système et l'espace libre ;
- il **explique** en français ce qui ne va pas et quoi faire ;
- il **sauvegarde selon le risque** : si un disque se dégrade, il lance une sauvegarde et te prévient ;
- il **vérifie** que tes sauvegardes se restaurent vraiment.

## Comment c'est fait

| Brique | Rôle | Techno |
|---|---|---|
| **Agent** | Surveille, sauvegarde, expose une API. Tourne sur PC, NAS, serveur | Go, sans dépendance |
| **Appli de bureau** | Fenêtre native qui embarque l'agent | Go + [Wails](https://wails.io) |
| **Appli téléphone** | Suit l'état de tous tes appareils | [Capacitor](https://capacitorjs.com) |

Les trois partagent la même interface (`web/static`). Ton NAS peut recevoir les sauvegardes
de tous tes PC et de ton téléphone.

## État actuel

| Fonction | État |
|---|---|
| Appli de bureau Windows, Linux, macOS | ✅ |
| Agent et API `/api/v1` | ✅ |
| Capteur d'espace disque | ✅ |
| Sauvegarde (moteur Kopia) | 🔜 v0.1 |
| SMART, températures, erreurs système | 🔜 v0.2 |
| Agent comme service, image Docker pour NAS | 🔜 v0.3 |
| Téléphone connecté aux agents | 🔜 v0.4 |

Feuille de route détaillée : [docs/architecture.md](docs/architecture.md).

## Essayer

Il faut [Go](https://go.dev/dl/) 1.27 ou plus récent.

```sh
git clone https://github.com/GoathierS/gardien.git
cd gardien

# Appli de bureau (Windows)
go build -tags desktop,production -ldflags "-H windowsgui" ./cmd/gardien-bureau

# Ou l'agent seul, dans le navigateur : http://127.0.0.1:8484
go run ./cmd/gardien serve
```

Commandes pour macOS, Linux et le téléphone : voir [CLAUDE.md](CLAUDE.md) et [mobile/README.md](mobile/README.md).

## Licence

[Apache 2.0](LICENSE).
