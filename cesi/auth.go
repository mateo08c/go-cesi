package cesi

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/golog"
	"net/url"
	"strings"
)

const (
	optionsFormSelector = "#options"
	errorTextSelector   = "#errorText"
)

type SAMLRequestData struct {
	Action      string
	RelayState  string
	SAMLRequest string
}

// 1. Get the page from https://ent.cesi.fr/
// 2. Send the login page to https://wayf.cesi.fr/login?client_name=ClientIdpViaCesiFr&needs_client_redirection=true&UserName=xxxxxx@viacesi.fr
// 3. parse html and get the form action, RelayState and SAMLRequest
// TODO: Finish to explain the login process

func (c *Cesi) Login() error {
	if c.username == "" || c.password == "" {
		return ErrMissingCredentials
	}

	if err := c.initConnection(); err != nil {
		return err
	}

	action, err := c.initSAML()
	if err != nil {
		return err
	}

	err = c.sendCredentials(action)
	if err != nil {
		return err
	}

	user, err := c.GetCurrentUser()
	if err != nil {
		return err
	}

	c.User = user

	return nil
}

func (c *Cesi) sendCredentials(u string) error {
	golog.Debug("(sendCredentials) Send the login page to ", u)
	resp, err := c.PostForm(u, map[string][]string{
		"UserName":   {c.username},
		"Password":   {c.password},
		"AuthMethod": {"FormsAuthentication"},
	})
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	if errorText := doc.Find(errorTextSelector).Length(); errorText > 0 {
		var errors []string
		doc.Find(errorTextSelector).Each(func(i int, s *goquery.Selection) {
			errors = append(errors, s.Text())
		})

		return Error(strings.Join(errors, " - "))
	}

	form := doc.Find("form[name='hiddenform']")
	if form.Length() == 0 {
		return Error("(sendCredentials) Error: form not found")
	}

	action, _ := form.Attr("action")
	if action == "" {
		return Error("(sendCredentials) Error: action not found")
	}

	relayState, _ := form.Find("input[name='RelayState']").Attr("value")
	if relayState == "" {
		return Error("(sendCredentials) Error: RelayState not found")
	}

	samlResponse, _ := form.Find("input[name='SAMLResponse']").Attr("value")
	if samlResponse == "" {
		return Error("(sendCredentials) Error: SAMLRequest not found")
	}

	resp, err = c.PostForm(action, map[string][]string{
		"RelayState":  {relayState},
		"SAMLRequest": {samlResponse},
	})
	if err != nil {
		return err
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	//get href in .sidebar-lien
	eid := doc.Find(".sidebar-lien").First()
	if eid.Length() == 0 {
		return Error("(sendCredentials) Error: eid not found")
	}

	href, _ := eid.Attr("href")
	if href == "" {
		return Error("(sendCredentials) Error: href not found")
	}

	ur, err := url.Parse(href)
	if err != nil {
		return err
	}

	c.User.Session.ID = ur.Path[1:]

	//get name in #compte.plier-deplier__bouton
	name := doc.Find("#compte .plier-deplier__bouton .session_libelle").First()
	if name.Length() == 0 {
		return Error("(sendCredentials) Error: name not found")
	}

	sl := strings.Split(name.Text(), " ")
	if len(sl) != 2 {
		return Error("(sendCredentials) Error: name not found")
	}

	c.User.FirstName = sl[0]
	c.User.LastName = sl[1]

	return nil
}

func (c *Cesi) initSAML() (string, error) {
	samlData, err := c.requestSAML()
	if err != nil {
		return "", err
	}

	return c.sendSAML(samlData)
}

func (c *Cesi) sendSAML(s *SAMLRequestData) (string, error) {
	golog.Debug("(sendSAML) Send the form to ", s.Action)
	resp, err := c.PostForm(s.Action, map[string][]string{
		"RelayState":  {s.RelayState},
		"SAMLRequest": {s.SAMLRequest},
	})
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	form := doc.Find(optionsFormSelector)
	if form.Length() == 0 {
		return "", Error("unable to find " + optionsFormSelector)
	}

	action, exist := form.Attr("action")
	if !exist || action == "" {
		return "", Error("unable to find action")
	}

	return action, nil
}

func (c *Cesi) requestSAML() (*SAMLRequestData, error) {
	golog.Debug("(requestSAML) Send the login page to ", WayfBaseURL)
	u, err := url.Parse(WayfBaseURL)
	if err != nil {
		return nil, err
	}

	u.Path = "/login"

	uv := url.Values{
		"client_name":              {"ClientIdpViaCesiFr"},
		"needs_client_redirection": {"true"},
		"UserName":                 {c.username},
	}

	u.RawQuery = uv.Encode()

	resp, err := c.PostForm(u.String(), nil)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	form := doc.Find("form")
	if form.Length() == 0 {
		return nil, Error("unable to find form")
	}

	action, _ := form.Attr("action")
	if action == "" {
		return nil, Error("unable to find action")
	}

	relayState, _ := form.Find("input[name='RelayState']").Attr("value")
	if relayState == "" {
		return nil, Error("unable to find RelayState")
	}

	samlResponse, _ := form.Find("input[name='SAMLRequest']").Attr("value")
	if samlResponse == "" {
		return nil, Error("unable to find SAMLRequest")
	}

	return &SAMLRequestData{
		Action:      action,
		RelayState:  relayState,
		SAMLRequest: samlResponse,
	}, nil
}

func (c *Cesi) initConnection() error {
	golog.Debug("(initConnection) Get the page from", BaseURL)
	resp, err := c.Get(BaseURL)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	if doc.Find("#login").Length() == 0 {
		return ErrInitConnection
	}

	return nil
}
