package main

import (
	"github.com/cheynewallace/tabby"
	"go-cesi/cesi"
	"os"
	"strings"
	"time"
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

	firstDay := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	lastDay := firstDay.AddDate(0, 1, -1)

	calendar, err := c.User.GetCalendar(firstDay, lastDay)
	if err != nil {
		panic(err)
	}

	t := tabby.New()
	t.AddHeader("Date", "Commence à", "Termine à", "Intitulé", "Salle", "Intervenant")
	for _, event := range calendar.Events {
		var rooms []string
		for _, room := range event.Salles {
			rooms = append(rooms, room.NomSalle)
		}

		var teachers []string
		for _, teacher := range event.Intervenants {
			teachers = append(teachers, teacher.Nom+" "+teacher.Prenom)
		}

		if len(teachers) == 0 && event.Title != "CESI" {
			teachers = append(teachers, "Autonomie")
		}

		if event.Title == "CESI" {
			event.Title = "Non défini (CESI)"
		}

		t.AddLine(event.Start.Format("02/01/2006"), event.Start.Format("15:04"), event.End.Format("15:04"), event.Title, strings.Join(rooms, ", "), strings.Join(teachers, ", "))
	}

	t.Print()
}
