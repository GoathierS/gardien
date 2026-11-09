---
name: test-restauration
description: Recette pour tester qu'une fonction de sauvegarde de gardien ne perd ni ne modifie aucune donnée. À utiliser pour tout code qui sauvegarde, restaure, déplace ou supprime des instantanés.
---

# Test sauvegarde, restauration, comparaison

Chaque fonction de sauvegarde ou de restauration a un test qui suit ces étapes.

## 1. Préparer un jeu de fichiers

Dans `t.TempDir()`, jamais dans un vrai dossier :
- un fichier texte, un fichier vide, un fichier binaire aléatoire de quelques Mo ;
- un nom Unicode (`été-ñ-日本.txt`) et un sous-dossier imbriqué sur 3 niveaux ;
- si le test le permet : un lien symbolique et un fichier caché.

## 2. Noter l'état des sources

Pour chaque fichier : chemin relatif, taille, empreinte SHA-256, date de modification.

## 3. Sauvegarder

Vers un second `t.TempDir()` servant de destination.

## 4. Vérifier que les sources n'ont pas bougé

Recalcule l'état des sources et compare-le à l'étape 2 : il doit être identique.

## 5. Restaurer et comparer

Restaure dans un troisième `t.TempDir()`. Compare chaque fichier restauré à l'état de l'étape 2 : même liste, mêmes tailles, mêmes empreintes. Les dates sont comparées seulement si la fonction promet de les conserver.

## 6. Tester l'incrémental

Modifie un fichier, supprimes-en un autre, ajoutes-en un troisième, sauvegarde à nouveau. Vérifie que la restauration de l'instantané le plus récent reflète ces changements, et que celle de l'ancien instantané donne toujours l'état de départ.

## Aide

Si elle n'existe pas encore, crée une fonction de test partagée `etatDossier(t, racine) map[string]Empreinte` dans un paquet `internal/testutil`, et réutilise-la partout.
