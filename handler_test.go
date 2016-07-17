package pixie_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/gouthamve/pixie"
)

func TestAccept(t *testing.T) {
	content, err := ioutil.ReadFile("config.json")
	ok(t, err)

	cfg := pixie.Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		t.Error(err)
	}
	px, err := pixie.NewPixie(cfg)
	ok(t, err)

	urls := []string{
		"https://www.google.com/",
		"https://m.facebook.com/",
		"http://nvision.org.in/",
	}

	for _, u := range urls {
		req, err := http.NewRequest("GET", u, nil)
		ok(t, err)

		w := httptest.NewRecorder()
		px.Forward(w, req)
		if w.Code < 200 || w.Code >= 400 {
			t.Error("Expected:", 200, "Got:", w.Code)
		}
	}

	return
}

func TestDeny(t *testing.T) {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Error(err)
	}

	cfg := pixie.Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		t.Error(err)
	}
	px, err := pixie.NewPixie(cfg)
	ok(t, err)

	urls := []string{
		"https://www.facebook.com/",
		"yahoo.com",
	}

	for _, u := range urls {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		px.Forward(w, req)
		if w.Code != 403 {
			t.Error("Expected:", 403, "Got:", w.Code)
		}
	}

	return
}

func TestIntegrationAccept(t *testing.T) {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Error(err)
	}

	cfg := pixie.Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		t.Error(err)
	}
	px, err := pixie.NewPixie(cfg)
	ok(t, err)

	s := httptest.NewServer(http.HandlerFunc(px.Forward))

	urls := []string{
		"https://www.google.com/",
		"https://m.facebook.com/",
		"http://nvision.org.in/",
	}

	for _, u := range urls {
		c := getProxyClient(s)

		resp, err := c.Get(u)
		ok(t, err)

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			t.Error("Expected: 200-399", "Got:", resp.StatusCode)
		}
	}
}

func TestIntegrationReject(t *testing.T) {
	content, err := ioutil.ReadFile("config.json")
	ok(t, err)

	cfg := pixie.Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		t.Error(err)
	}
	px, err := pixie.NewPixie(cfg)
	ok(t, err)

	s := httptest.NewServer(http.HandlerFunc(px.Forward))

	urls := []string{
		"https://www.facebook.com",
		"http://yahoo.com",
	}

	for _, u := range urls {
		c := getProxyClient(s)

		resp, err := c.Get(u)
		if err != nil {
			equals(t, (*http.Response)(nil), resp)
			if urlError, ok := err.(*url.Error); !ok || urlError.Err.Error() != "Forbidden" {
				t.Errorf("Expected a *url.Error with our 'Forbidden' error inside; got %#v (%q)", err, err)
			}
			continue
		}

		if resp.StatusCode != 403 {
			t.Error("Expected:", 403, "Got:", resp.StatusCode)
		}
	}
}

func getProxyClient(s *httptest.Server) *http.Client {
	u, err := url.Parse(s.URL)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}

	return client
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
