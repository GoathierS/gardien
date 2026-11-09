# Architecture

## Vue d'ensemble

```
            ┌──────────── API /api/v1 (stable, versionnée) ─────────────┐
 Appli de bureau (Wails) │ Appli téléphone (Capacitor) │ Navigateur (NAS) │ Ligne de commande
            └──────────────────────────┬─────────────────────────────────┘
                        Agent gardien (Go, tourne en fond)
   ┌──────────────┬───────────────┬────────────────┬─────────────────┐
   Sources        Destinations    Capteurs santé   Notifications    ← extensions
 (dossiers,      (disque, NAS,   (espace, SMART,  (système, mail,
  téléphone)      S3, cloud)      temp., journaux) téléphone)
```

- **L'agent** est le même programme partout : PC, NAS, serveur. Un seul exécutable sans CGO.
- **Les interfaces** ne parlent à l'agent que par l'API. Elles peuvent changer de techno sans
  toucher au cœur.
- **L'interface est écrite une seule fois** (`web/static`) et réutilisée par l'appli de bureau,
  l'appli téléphone et la page servie par un agent.
- **Les extensions** se branchent derrière des interfaces Go simples. `Capteur` existe ; `Source`,
  `Destination` et `Notificateur` seront définies avec l'intégration de Kopia, dont elles dépendent.
- **Plusieurs agents** pourront se parler : un NAS peut recevoir les sauvegardes des PC et du
  téléphone et afficher l'état de tous les appareils.

## Décisions

| Décision | Raison |
|---|---|
| Go pour l'agent | Promesse de compatibilité de Go 1, exécutable unique, compilation croisée, même langage que Kopia et restic |
| Kopia comme bibliothèque | Moteur éprouvé (utilisé aussi par Velero) : on ne réécrit pas la partie la plus risquée |
| Wails v2 pour le bureau | Vraie fenêtre native, interface web réutilisable, version stable (la v3 est encore en préversion) |
| Capacitor pour le téléphone | Réutilise la même interface sur Android et iOS |
| API versionnée `/api/v1` | Les interfaces et les autres agents en dépendent : stabilité dans le temps |
| Écoute sur 127.0.0.1 par défaut | Rien d'exposé au réseau avant l'authentification |
| Interface sans étape de compilation | Simplicité ; à revoir (Svelte, Vue…) si l'interface grossit |
| Apache 2.0 | Compatible avec Kopia |

## Feuille de route

### v0.1 : la base
- [x] Agent, API `/api/v1`, interface commune
- [x] Appli de bureau Wails (Windows, Linux, macOS)
- [x] Interface `Capteur` et capteur d'espace disque
- [x] Squelette de l'appli téléphone Capacitor
- [ ] Configuration (fichier YAML) : dossiers sources, destinations, planification
- [ ] Interfaces `Source` et `Destination`, sauvegarde via Kopia vers un dossier ou un disque externe
- [ ] Liste des instantanés et restauration d'un fichier
- [ ] Test automatique sauvegarde, restauration et comparaison d'empreintes

### v0.2 : la santé
- [ ] Capteur SMART via smartctl (smartmontools)
- [ ] Capteur de températures CPU, GPU, disques
- [ ] Capteur d'erreurs système (journal d'événements Windows, journalctl, journaux macOS)
- [ ] Interface `Notificateur` et notifications système
- [ ] Sauvegarde déclenchée par un risque

### v0.3 : agents et NAS
- [ ] Agent comme service : Windows, launchd, systemd
- [ ] Icône dans la zone de notification (barre près de l'horloge)
- [ ] Destinations NAS : SMB, SFTP, WebDAV, S3
- [ ] Image Docker pour NAS (Synology, TrueNAS)
- [ ] Installeurs et publication avec GoReleaser

### v0.4 : réseau et téléphone
- [ ] Authentification, écoute sur le réseau local, découverte des agents
- [ ] Appli téléphone : réglages de l'agent, état de tous les appareils, notifications
- [ ] Un agent central (NAS) qui reçoit les sauvegardes des autres

### v1.0
- [ ] Cloud, tableau de bord multi-appareils, sauvegarde des photos du téléphone
