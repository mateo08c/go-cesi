package main

import (
	"github.com/kataras/golog"
	"github.com/mateo08c/go-cesi/cesi"
	"os"
)

func main() {
	c := cesi.New(&cesi.Options{
		Email:    os.Getenv("CESI_MAIL"),
		Password: os.Getenv("CESI_PASSWORD"),
	})

	err := c.Login()
	if err != nil {
		golog.Fatal(err)
	}

	for _, e := range c.User.Establishments {
		golog.Infof("Establishment: %s", e.Name)
	}

	golog.Infof("Session ID: %s", c.User.Session.ID)
	golog.Infof("Firstname: %s", c.User.FirstName)
	golog.Infof("Lastname: %s", c.User.LastName)
	golog.Infof("Email: %s", c.User.Email)
	golog.Infof("Phone: %s", c.User.Phone)
	golog.Infof("Promotion: %s", c.User.Promotion)
}
