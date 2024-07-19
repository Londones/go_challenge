# Utiliser swaggo/swag

## Information pratique:
Swagger version : v1.16.3

## Utilisation:
1. D'abord mettre les différentes annotations dans les fichiers que l'on souhaite
   utiliser, (pour nous ce sera les fichiers internal/controller ainsi que le
   cmd/api/main.go).

2. Dans le main.go doit apparaître les informations génrales de l'api ainsi que
   tous les modules que l'on va utiliser (internal/controller dans notre cas)

3. Pour connaitre les bonnes annotations:
- https://github.com/swaggo/swag/?tab=readme-ov-file#general-api-info
- https://github.com/swaggo/swag/?tab=readme-ov-file#api-operation

4. Après avoir mis les annotations dans les fichiers et dans le main.go il faut
   initialiser le swagger à partir de la racine du projet avec la commande:
```shell
swag init --parseDependency -d ./internal/handlers/ -g ../../cmd/api/main.go
```

### Explication de la commande:
La commande ci-dessus prends plusieurs paramètres:
- ``--parseDependency`` --> Se base sur le fichier de référence indiqué avec le paramètres
  ``-g`` et parse toutes les dépendences qui y figurent.
- ``-d`` Spécifie le chemin d'accès pour les dossier à parser
- ``-g`` Spécifie le chemin d'accès pour le fichier main.go.

