# Project go-challenge

Ce projet est le back office de purrfect-match une application réalisée en flutter d'adoption de chats.

## Getting Started

### Téléchargement depuis git
```sh
git clone https://github.com/Londones/go-challenge
```

### Fichier .env

Les variables d'environnements sont stockés dans un .env à la racine du projet.
```
// .env

PORT=
APP_ENV=

DB_HOST=
DB_PORT=
DB_DATABASE=
DB_USERNAME=
DB_PASSWORD=

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

SERVER_URL=http://localhost:8080
CLIENT_URL=http://localhost:8080
CLIENT_CALLBACK_URL=http://localhost:8080/auth/google/callback
APP_CLIENT_URL=purrfect_match://app_oauth

UPLOAD_CARE_SECRET_KEY=
UPLOAD_CARE_PUBLIC_KEY=

FIREBASE_SDK=

JWT_SECRET=

SESSION_KEY=

SESSION_SECRET=
```

Pour le client googleID qui correspond à la connexion avec google suivre ce lien afin de créer
un google id:
https://support.google.com/workspacemigrate/answer/9222992?hl=fr

### Les dépendances:
Pour récupérer et mettre à jour les dépendance du projet: 
```shell
go mod tidy
```

### Utilisation du swagger. 
Le swagger du projet se lance avec la commande suivante: 
```shell
swag init --parseDependency -d internal/handlers/ -g ../../cmd/api/main.go
```

### Lancement du serveur: 
Pour lancer le serveur go, à la racine du projet exécuter la commande:
```shell
go run cmd/api/main.go
```
Des fixtures peuvent être jouées automatiquement au lancement du serveur, il suffit de décommenter la ligne 9 et les lignes 109-156 du fichier ./internal/database/database.go

Les tests unitaires de l'application se trouvent sur la branche tests du git.

## Description du projet

### Organisation des dossiers

Les dossiers du projets sont organisés comme ceci:
```
./assets/ --> Contient certains assets du projet les autres sont récupérés depuis UPLOAD_CARE.

./cmd/api/main.go --> Fichier principale de l'application, il gère le lancement du serveur ainsi que
    des procédure de lancement de la base de données.
    
./docs/ --> Dossier contenant les fichiers générés et utilisés par le swagger de l'application.

./internal/ --> Dossier principal du projet, il contient toute la logique ainsi que les dossiers 
    relatifs à la base de données.

./internal/api/ --> Contient le fichier permettant d'accéder à UPLOAD_CARE pour récupérer les assets
    du projet.
./internal/auth/ --> Dossier contenant la logique d'authentification et de gestion des tokens.
./internal/config/ --> Dossier contenant la logique d'accès à firebase, qui est utilisé en prod pour
    héberger l'application.
./internal/database/ --> Contient le dossier /queries possédant toute la logique permettant de requéter 
    la base de donnée de l'application. /database contient aussi les fichers permettant de setup la base 
    de donnée et de jouer les fixtures de l'application.
./internal/fixture/ --> Contient l'ensemble des fixtures de l'application.
./internal/handlers/ --> Contient l'ensembles des fichiers permettants de traiter les demandes client et 
    d'appeler les bonnes queries pour la base de donnée.
./internal/models/ --> Contient tous les models de l'application.
./internal/server/ --> Contient les fichiers gérants les routes les middleware pour gérer l'accès aux routes
    selon les rôles ainsi que le fichier permettant la création d'un serveur.
./internal/utils/ --> Contient le logger de l'application.

./tests/ --> Contient le fichiers de test de handlers de l'application.
```

