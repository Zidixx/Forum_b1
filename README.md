# LE VESTIAIRE - Forum Football

Forum web dédié au football, construit en Go avec SQLite. Pas de framework.

## Prerequis

- **Go 1.21+** : [golang.org/dl](https://golang.org/dl/)
- **GCC** : necessaire pour la compilation de go-sqlite3
  - macOS : `xcode-select --install`
  - Ubuntu/Debian : `sudo apt install gcc`
  - Windows : installer [MinGW-w64](https://www.mingw-w64.org/)

## Installation

```bash
# 1. Cloner le repo
git clone https://github.com/zidixx/Forum_b1.git
cd Forum_b1

# 2. Installer les dependances Go
go mod tidy
```

## Configuration OAuth (Google & GitHub)

La connexion via Google et GitHub necessite des cles OAuth. Sans ces cles, le forum fonctionne normalement mais les boutons "Se connecter avec Google/GitHub" ne marcheront pas.

### Etape 1 : Recuperer le fichier `.env`

Demande le fichier `.env` au proprietaire du projet (par message prive, jamais via git).

### Etape 2 : Placer le fichier `.env`

Copie le fichier `.env` a la racine du projet (au meme niveau que `go.mod`) :

```
Forum_b1/
  .env          <-- ici
  go.mod
  cmd/
  internal/
  ...
```

### Etape 3 : Verifier le contenu

Le fichier `.env` doit contenir :

```env
# OAuth Google
GOOGLE_CLIENT_ID=xxxxx
GOOGLE_CLIENT_SECRET=xxxxx

# OAuth GitHub
GITHUB_CLIENT_ID=xxxxx
GITHUB_CLIENT_SECRET=xxxxx

# Base URL pour les callbacks OAuth
OAUTH_REDIRECT_BASE=https://localhost:8443
```

> Un fichier `.env.example` est fourni dans le repo comme reference (sans les vraies cles).

## Certificat TLS (HTTPS)

Le serveur demarre en HTTPS. Le dossier `tls/` contient deja les certificats auto-signes (`server.crt` et `server.key`). Rien a faire.

> Ton navigateur affichera un avertissement "connexion non securisee" au premier acces — c'est normal pour un certificat auto-signe. Clique sur "Avance" puis "Continuer".

## Lancement

```bash
go run cmd/web/main.go
```

Le serveur demarre sur **https://localhost:8443**.

Pour changer le port :

```bash
PORT=3000 go run cmd/web/main.go
```

## Structure du projet

```
cmd/web/main.go           Point d'entree, routes
internal/
  db/                     Connexion SQLite, migrations
  models/                 Structures de donnees
  repository/             Acces base de donnees
  service/                Logique metier
  handler/                Handlers HTTP
  middleware/             Auth, rate limiter, contexte user
  utils/                  Validation, hash, upload, UUID
templates/                Templates HTML
static/
  css/                    Feuilles de style
  js/                     JavaScript
  img/                    Images statiques
  uploads/                Images uploadees par les users
sql/                      Schema et seed SQL
tls/                      Certificats HTTPS
data/                     Base SQLite (creee automatiquement)
```

## Fonctionnalites

- Inscription / Connexion (email + mot de passe)
- Connexion OAuth (Google, GitHub)
- Creer, modifier, supprimer des posts
- Commenter les posts (avec reponses imbriquees)
- Like / Dislike sur posts et commentaires
- Repost
- Upload d'images (JPEG, PNG, GIF — max 20 Mo)
- Recherche de posts
- Filtrage par ligue (Ligue 1, Premier League, La Liga, Bundesliga, Serie A, Champions League, Europa League)
- Rate limiting anti-DDoS

## Routes principales

| Methode | Route | Description |
|---------|-------|-------------|
| GET | `/` | Accueil — liste des posts |
| GET/POST | `/register` | Inscription |
| GET/POST | `/login` | Connexion |
| POST | `/logout` | Deconnexion |
| GET | `/auth/google/login` | Connexion via Google |
| GET | `/auth/github/login` | Connexion via GitHub |
| GET/POST | `/post/create` | Creer un post |
| GET | `/post/{id}` | Detail d'un post |
| GET/POST | `/post/edit/{id}` | Modifier un post |
| POST | `/post/delete/{id}` | Supprimer un post |
| POST | `/post/react/{id}` | Like/dislike un post |
| POST | `/comment/create` | Creer un commentaire |
| POST | `/comment/react/{id}` | Like/dislike un commentaire |
| POST | `/repost/{id}` | Reposter |
| GET | `/search` | Recherche |
| GET | `/my-posts` | Mes posts |
| GET | `/liked-posts` | Posts aimes |

## Base de donnees

La base SQLite est creee automatiquement au demarrage dans `data/forum.db`. Le schema est dans `sql/schema.sql`, les categories par defaut dans `sql/seed.sql`.

## Docker

```bash
docker build -t forum .
docker run -p 8443:8443 forum
```

## Reactions (like/dislike)

- Cliquer sur like quand on a deja like = annule le like
- Cliquer sur dislike quand on a like = remplace par dislike (et inversement)
- Un seul vote par utilisateur par post/commentaire
