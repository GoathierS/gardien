# Architecture

## Vue d'ensemble

```
┌───────────────┐    API locale     ┌──────────────────────────────┐
│ Interface web │ ◄───────────────► │ Démon gardien                │
│ + icône barre │                   │ - planification              │
└───────────────┘                   │ - moteur Kopia (bibliothèque)│
┌───────────────┐   réseau local    │ - santé : SMART, temp., logs │
│ Téléphone     │ ◄───────────────► │ - notifications              │
│ (PWA puis app)│   (avec auth)     └──────────────┬───────────────┘
└───────────────┘                                  │
                          Destinations : disque, NAS (SMB/SFTP/WebDAV), S3, cloud
```

- **Un seul exécutable** par OS, sans dépendance : démon, API et interface embarquée.
- **Interface web** servie par le démon, utilisable sur PC et téléphone.
- **Moteur Kopia** (Apache 2.0) pour le chiffrement, la déduplication, l'historique et les destinations distantes. On ne réécrit pas de moteur de sauvegarde.

## Décisions

| Décision | Raison |
|---|---|
| Go | Compilation croisée simple, exécutable unique, adapté aux services, même langage que Kopia et restic |
| Kopia comme bibliothèque | Moteur éprouvé (utilisé aussi par Velero), évite la partie la plus risquée |
| Interface web plutôt que native | Une seule interface pour les 3 OS et le téléphone |
| Écoute sur 127.0.0.1 par défaut | Rien d'exposé au réseau avant l'authentification |
| Apache 2.0 | Compatible avec Kopia |

## Feuille de route

### v0.1 : la base
- [x] Démon et interface web locale
- [x] Espace disque
- [ ] Configuration (fichier YAML) : dossiers sources, destinations, planification
- [ ] Sauvegarde via Kopia vers un dossier local ou un disque externe
- [ ] Liste des instantanés et restauration d'un fichier
- [ ] Test automatique sauvegarde, restauration et comparaison d'empreintes

### v0.2 : la santé
- [ ] SMART via smartctl (smartmontools)
- [ ] Températures CPU, GPU, disques
- [ ] Erreurs système (journal d'événements Windows, journalctl, journaux macOS)
- [ ] Rapport en français et sauvegarde déclenchée par un risque

### v0.3 : multiplateforme
- [ ] Service Windows, launchd, systemd
- [ ] Destinations NAS : SMB, SFTP, WebDAV, S3
- [ ] Image Docker pour NAS (Synology, TrueNAS)
- [ ] Installeurs et publication avec GoReleaser

### v0.4 : téléphone
- [ ] Authentification et écoute sur le réseau local
- [ ] PWA installable, notifications

### v1.0
- [ ] Cloud, tableau de bord multi-PC, application mobile (sauvegarde des photos)
