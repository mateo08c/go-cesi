package cesi

import (
	"encoding/json"
	"fmt"
	"time"
)

type Calendar struct {
	Events []*Event
}

type Event struct {
	Code      string      `json:"code"`
	Title     string      `json:"title"`
	AllDay    bool        `json:"allDay"`
	Nightly   bool        `json:"nightly"`
	Start     time.Time   `json:"start"`
	End       time.Time   `json:"end"`
	Url       string      `json:"url"`
	NomModule string      `json:"nomModule"`
	Matiere   interface{} `json:"matiere"`
	Theme     string      `json:"theme"`
	Salles    []struct {
		NomSalle string `json:"nomSalle"`
	} `json:"salles"`
	Intervenants []struct {
		SousTitre           interface{}   `json:"sousTitre"`
		Profils             interface{}   `json:"profils"`
		GroupesPedagogiques interface{}   `json:"groupesPedagogiques"`
		UrlFiche            string        `json:"urlFiche"`
		UrlPhoto            string        `json:"urlPhoto"`
		Nom                 string        `json:"nom"`
		Prenom              string        `json:"prenom"`
		Code                string        `json:"code"`
		AdresseMail         string        `json:"adresseMail"`
		UrlAgenda           string        `json:"urlAgenda"`
		Sessions            []interface{} `json:"sessions"`
		Inconnu             bool          `json:"inconnu"`
	} `json:"intervenants"`
	ParticipantsPersonne interface{} `json:"participantsPersonne"`
	Participants         []struct {
		LibelleGroupe string `json:"libelleGroupe"`
		CodeGroupe    string `json:"codeGroupe"`
		CodeSession   string `json:"codeSession"`
	} `json:"participants"`
}

func (ct *Event) UnmarshalJSON(b []byte) (err error) {
	type Alias Event
	aux := &struct {
		*Alias
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Alias: (*Alias)(ct),
	}

	if err = json.Unmarshal(b, &aux); err != nil {
		return
	}

	ct.Start, err = time.ParseInLocation("2006-01-02T15:04:05-07", aux.Start, time.Local)
	if err != nil {
		return
	}

	ct.End, err = time.ParseInLocation("2006-01-02T15:04:05-07", aux.End, time.Local)
	if err != nil {
		return
	}

	return nil
}

func (u *User) GetCalendar(start time.Time, end time.Time) (*Calendar, error) {
	s := start.Format("2006-01-02")
	e := end.Format("2006-01-02")

	url := fmt.Sprintf(CalendarURL, s, e, u.ID)
	
	resp, err := u.cesi.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return &Calendar{}, nil
	}

	var events []*Event
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return nil, err
	}

	return &Calendar{
		Events: events,
	}, nil
}
