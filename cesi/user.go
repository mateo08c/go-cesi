package cesi

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

type User struct {
	cesi           *Cesi
	ID             int
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	Promotion      string
	Session        *Session
	Establishments []*Establishment
}

type Establishment struct {
	Name string
}

func (u *User) GetIdentifier() string {
	return strings.ToLower(removeAccents(fmt.Sprintf("%s-%s", u.LastName, u.FirstName)))
}

func (c *Cesi) GetCurrentUser() (*User, error) {
	u := fmt.Sprintf(ProfileURL, c.User.Session.ID, c.User.GetIdentifier())

	get, err := c.Get(u)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(get.Body)
	if err != nil {
		return nil, err
	}

	email := doc.Find("#sans_encadres > div:nth-child(1) > div.personne__entete__informations > dl > dd:nth-child(2) > a").First().Text()
	phone := doc.Find("#sans_encadres > div:nth-child(1) > div.personne__entete__informations > dl > dd:nth-child(4) > a").First().Text()
	name := doc.Find("#sans_encadres > div:nth-child(1) > div.personne__entete__carte > div.personne__entete__info > p > strong").First().Text()
	firstName := name[:len(name)/2]
	lastName := name[len(name)/2:]

	userID, _ := doc.Find("#sans_encadres > div.personne__entete.accordeon > div > div.paragraphe--1 > div > div").First().Attr("data-code-personne")
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	promo := doc.Find("#sans_encadres > div.personne__entete.accordeon > h2").First().Text()
	promo = strings.TrimSpace(promo)

	hr := doc.Find("#menu > ul > li.submenu.item.item-session.active > a.sidebar-lien").First().AttrOr("href", c.User.Session.ID)
	ur, err := url.Parse(hr)
	if err != nil {
		return nil, err
	}

	var establishments []*Establishment
	doc.Find("#sans_encadres > div:nth-child(1) > div.personne__entete__groupes > dl > dd").Each(func(i int, s *goquery.Selection) {
		establishments = append(establishments, &Establishment{
			Name: strings.TrimSpace(s.Text()),
		})
	})

	return &User{
		cesi:      c,
		ID:        uid,
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Email:     email,
		Phone:     phone,
		Promotion: promo,
		Session: &Session{
			ID: ur.Path[1:],
		},
		Establishments: establishments,
	}, nil
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
