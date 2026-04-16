# Forum

Forum web en Go avec SQLite, sans framework.

## Prérequis

- Go 1.21+
- GCC (nécessaire pour go-sqlite3)

## Lancement

```bash
# Installer les dépendances
go mod tidy

# Lancer le serveur
go run cmd/web/main.go
```

Le serveur démarre sur `http://localhost:8080`.

Pour changer le port :
```bash
PORT=3000 go run cmd/web/main.go
```

Si le binaire est lancé depuis un autre dossier, passer le chemin racine du projet en argument :
```bash
go run cmd/web/main.go /chemin/vers/le/projet
```

## Structure

```
/cmd/web/main.go          Point d'entrée
/internal/
  /db/                     Connexion SQLite, migrations
  /models/                 Structures de données
  /repository/             Accès base de données
  /service/                Logique métier
  /handler/                Handlers HTTP
  /middleware/             Auth, contexte utilisateur
  /utils/                  Validation, hash, upload, UUID
/templates/                Templates HTML
/static/css/               Feuilles de style
/static/uploads/           Images uploadées
/sql/                      Schema et seed SQL
/data/                     Base SQLite (créée automatiquement)
```

## Base de données

La base SQLite est créée automatiquement au démarrage dans `data/forum.db`.
Le schéma est dans `sql/schema.sql`, les catégories par défaut dans `sql/seed.sql`.

## Images

Les images uploadées sont stockées dans `static/uploads/`.
Formats acceptés : JPEG, PNG, GIF. Taille max : 20 Mo.

## Routes

| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/` | Accueil (liste des posts, filtre par catégorie via `?category=ID`) |
| GET | `/register` | Page d'inscription |
| POST | `/register` | Inscription |
| GET | `/login` | Page de connexion |
| POST | `/login` | Connexion |
| POST | `/logout` | Déconnexion |
| GET | `/post/create` | Formulaire nouveau post |
| POST | `/post/create` | Création de post |
| GET | `/post/{id}` | Détail d'un post |
| GET | `/post/edit/{id}` | Formulaire modification post |
| POST | `/post/edit/{id}` | Modification du post |
| POST | `/post/delete/{id}` | Suppression du post |
| POST | `/post/react/{id}` | Like/dislike un post |
| POST | `/comment/create` | Créer un commentaire |
| GET | `/comment/edit/{id}` | Formulaire modification commentaire |
| POST | `/comment/edit/{id}` | Modification du commentaire |
| POST | `/comment/delete/{id}` | Suppression du commentaire |
| POST | `/comment/react/{id}` | Like/dislike un commentaire |
| GET | `/my-posts` | Mes posts |
| GET | `/liked-posts` | Posts aimés |

## Tests

```bash
go test ./...
```

## Réactions (like/dislike)

- Cliquer sur like quand on a déjà liké annule le like
- Cliquer sur dislike quand on a liké remplace par dislike (et inversement)
- Un seul vote par utilisateur par post/commentaire

## Rôles

- `guest` : lecture seule
- `user` : créer, commenter, liker, modifier/supprimer son contenu
- `admin`/`moderator` : prévus dans le schéma (champ `role` dans `users`)

## Intégration Docker/serveur

Le projet est conçu pour être branché facilement :
- Pas de dépendance à Docker
- Port configurable via `PORT`
- Chemins relatifs ou configurables via argument
- Routes centralisées dans `main.go`
