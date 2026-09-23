Dawnless

Jeu de rôle en ligne de commande développé en Go, dans le cadre du Projet RED. Le joueur incarne un
voyageur pris dans une boucle temporelle à Emberhollow ; l'histoire complète se trouve dans
[docs/histoire.md](docs/histoire.md).

Prérequis

Go 1.21 ou supérieur (`go version` pour vérifier). Un terminal moderne avec support des couleurs
ANSI : Windows Terminal ou PowerShell sur Windows (pas l'ancien cmd.exe), n'importe quel terminal
sur macOS/Linux.

Lancer le jeu

Depuis la racine du dépôt (le dossier contenant go.mod) :

    go run ./src

Ou pour construire un exécutable réutilisable :

    go build -o dawnless ./src
    ./dawnless

Sur Windows :

    go build -o dawnless.exe ./src
    .\dawnless.exe

Lancer les tests

    go test ./src/...

Contrôles

Écrans de sélection (titre, classe, nom, combat) : flèches ↑↓←→ pour naviguer, Entrée pour valider.

Sur la carte : flèches pour se déplacer, Entrée ou M pour ouvrir le menu à tout moment (le jeu
ne quitte jamais la carte, le menu s'affiche par-dessus).

Menu : touches numériques (1 à 7) pour choisir une option, Entrée pour valider.

Structure du dépôt

    src/    code source Go (un seul package main)
    docs/   documentation : histoire du jeu, carte ASCII
