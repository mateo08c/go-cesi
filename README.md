# go-cesi 📘

[![GoDoc](https://godoc.org/github.com/mateo08c/go-cesi?status.svg)](https://godoc.org/github.com/github.com/mateo08c/go-cesi)
[![Go Report Card](https://goreportcard.com/badge/github.com/mateo08c/go-cesi)](https://goreportcard.com/report/github.com/mateo08c/go-cesi)

Une bibliothèque Go pour se connecter et récupérer des informations depuis l'ENT du CESI.

- [X] 🔑 Authentification
- [X] 👤 Récupération des informations de l'utilisateur
- [X] 🏫 Récupération des informations des établissements
- [X] 📓 Récupération des informations des cours
- [ ] 🔮 Voir mon avenir...




## Installation 💻

Pour installer cette bibliothèque, utilisez la commande \`go get\` :

```bash
go get ggithub.com/mateo08c/go-cesi
```

## Utilisation 🚀

Voici un exemple d'utilisation de cette bibliothèque :

```go
package main

import (
	cesi "github.com/mateo08c/go-cesi"
	"os"
)

func main() {
	c := cesi.New(&cesi.Options{
		Email:    os.Getenv("CESI_MAIL"),
		Password: os.Getenv("CESI_PASSWORD"),
	})

	err := c.Login()
	if err != nil {
		panic(err)
	}

	for _, e := range c.User.Establishments {
		println(e.Name)
	}

	golog.Infof("Session ID: %s", c.User.Session.ID)
	golog.Infof("Firstname: %s", c.User.FirstName)
	golog.Infof("Lastname: %s", c.User.LastName)
	golog.Infof("Email: %s", c.User.Email)
	golog.Infof("Phone: %s", c.User.Phone)
	golog.Infof("Promotion: %s", c.User.Promotion)
}

```

## Contribution 🤝

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou à soumettre une pull request.

## Disclaimer ⚠️
La bibliothèque effectue de nombreuses requêtes sur l'ENT du CESI, ce qui pourrait entraîner un blocage temporaire de votre compte si vous effectuez un grand nombre de requêtes en peu de temps. Il est important de noter que je décline toute responsabilité quant à l'utilisation que vous faites de cette bibliothèque.

De plus, il est essentiel de comprendre que cette bibliothèque n'est pas officielle et que je n'ai aucun lien d'affiliation avec le CESI.

Cette bibliothèque a été créée dans le cadre d'un projet visant à automatiser et synchroniser mon ENT avec mon calendrier Google. 

## Licence ⚖️

Cette bibliothèque est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus d'informations.
