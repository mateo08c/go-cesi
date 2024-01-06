package cesi

import (
	"crypto/tls"
	"fmt"
	"github.com/kataras/golog"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Options struct {
	Email     string
	Password  string
	UserAgent string
	Debug     bool
}

type Session struct {
	ID string
}

type Cesi struct {
	client    *http.Client
	username  string
	password  string
	useragent string
	debug     bool
	User      *User
}

func New(o *Options) *Cesi {
	if o.Debug {
		golog.SetLevel("debug")
	}

	if o.UserAgent == "" {
		o.UserAgent = "go-cesi v0.1"
	}
	c := newHttpClient()

	ces := &Cesi{
		useragent: o.UserAgent,
		username:  o.Email,
		password:  o.Password,
		client:    c,
		debug:     o.Debug,
	}

	ces.User = &User{
		cesi:           ces,
		Session:        &Session{},
		Establishments: []*Establishment{},
	}

	return ces
}

func (c *Cesi) Get(url string) (*http.Response, error) {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Cesi) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Cesi) PostForm(url string, data url.Values) (*http.Response, error) {
	req, err := c.newRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.client.Do(req)
}

func (c *Cesi) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	golog.Debug("[NewRequest] " + method + " " + url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.useragent)

	return req, nil
}

func newHttpClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			golog.Debug("Redirect: " + req.URL.String())
			return nil
		},
		Jar: newCookieJar(),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func newCookieJar() http.CookieJar {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	return cookieJar
}

func saveResponse(r io.ReadCloser) {
	if _, err := os.Stat("debug"); os.IsNotExist(err) {
		err := os.Mkdir("debug", 0755)
		if err != nil {
			panic(err)
		}
	}
	name := fmt.Sprintf("debug-%d.html", time.Now().Unix())

	out, err := os.Create("debug/" + name)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	p := path.Join(wd, "debug", name)

	golog.Debug("Debug file saved: " + p)
}
