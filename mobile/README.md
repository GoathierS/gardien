# gardien sur téléphone

Application Android et iOS construite avec [Capacitor](https://capacitorjs.com/).
Elle reprend l'interface commune de `web/static` : on ne la modifie pas ici.

**État :** squelette. Le téléphone affichera l'état de tes PC et de ton NAS quand les agents
pourront être joints sur le réseau local, avec authentification (v0.4).

## Prérequis

- Node.js 22 ou plus récent
- Android : [Android Studio](https://developer.android.com/studio)
- iOS : un Mac avec Xcode

## Commandes

```sh
cd mobile
npm install
npx cap add android      # une seule fois (ou: npx cap add ios)
npm run android          # copie l'interface, synchronise et ouvre Android Studio
```

`npm run ui` copie seulement `web/static` vers `www`.

## À savoir

- Sur iOS, le système limite fortement le travail en arrière-plan : la sauvegarde
  automatique des photos y sera plus restreinte que sur Android.
- L'interface détecte qu'elle tourne sur téléphone (`window.Capacitor`) et appellera
  l'agent dont l'adresse sera enregistrée dans les réglages.
