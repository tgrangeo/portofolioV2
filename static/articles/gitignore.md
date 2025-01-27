# Créer et Utiliser un `.gitignore` Global

## Qu'est-ce qu'un fichier `.gitignore` ?

Un fichier `.gitignore` est un fichier spécial utilisé par Git pour indiquer quels fichiers ou dossiers ne doivent pas être suivis dans un dépôt. Cela permet de ne pas inclure des fichiers temporaires, sensibles ou spécifiques à l'environnement dans votre versionnage.

## Pourquoi un `.gitignore` global ?

Un fichier `.gitignore` global est utile lorsque vous travaillez sur plusieurs projets et que vous souhaitez ignorer certains fichiers communs à tous vos dépôts Git. Par exemple :

- Les fichiers de sauvegarde générés par les éditeurs (ex. : `*.swp`, `*.bak`).
- Les fichiers spécifiques à votre système d'exploitation (ex. : `.DS_Store` sur macOS, `Thumbs.db` sur Windows).
- Les configurations locales d'outils comme VSCode (`.vscode/`).

Cela vous évite de devoir ajouter ces règles dans chaque projet individuellement.

## Étapes pour créer un `.gitignore` global

### 1. Créer le fichier `.gitignore_global`
Vous pouvez le créer à un emplacement de votre choix, par exemple dans votre répertoire utilisateur :

```bash
touch ~/.gitignore_global
