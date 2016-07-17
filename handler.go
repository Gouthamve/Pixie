package pixie

import (
	"io"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/vulcand/oxy/forward"
)

// Pixie is the encapsulating struct
type Pixie struct {
	fwd         *forward.Forwarder
	acceptRegex []*regexp.Regexp
	denyRegex   []*regexp.Regexp
}

// NewPixie returns a new Pixie :P
func NewPixie(cfg Config) (*Pixie, error) {
	fwd, err := forward.New()
	if err != nil {
		return nil, err
	}

	px := Pixie{
		fwd: fwd,
	}

	px.acceptRegex = make([]*regexp.Regexp, len(cfg.Accept))
	px.denyRegex = make([]*regexp.Regexp, len(cfg.Deny))
	for i, v := range cfg.Accept {
		px.acceptRegex[i] = regexp.MustCompile(v)
	}
	for i, v := range cfg.Deny {
		px.denyRegex[i] = regexp.MustCompile(v)
	}

	return &px, nil
}

// Forward the main handler
func (px *Pixie) Forward(w http.ResponseWriter, req *http.Request) {
	accepted := false
	for i := range px.acceptRegex {
		if px.acceptRegex[i].MatchString(req.URL.String()) {
			accepted = true
			break
		}
	}

	if !accepted {
		for i := range px.denyRegex {
			if px.denyRegex[i].MatchString(req.URL.String()) {
				http.Error(w, "", http.StatusForbidden)
				return
			}
		}
	}

	if req.Method == "CONNECT" {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}

		clientConn, _, err := hj.Hijack()
		if err != nil {
			panic("Cannot hijack connection " + err.Error())
		}

		targetSiteCon, err := net.Dial("tcp", req.URL.Host)
		if err != nil {
			log.Printf("Error establishing remote connection: %s\n", err)
			if _, err := io.WriteString(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n"); err != nil {
				log.Printf("Error responding to client: %s\n", err)
			}
		}

		clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		tTargetSiteCon := targetSiteCon.(*net.TCPConn)
		tClientConn := clientConn.(*net.TCPConn)

		go copyAndClose(tTargetSiteCon, tClientConn)
		go copyAndClose(tClientConn, tTargetSiteCon)

	} else {
		px.fwd.ServeHTTP(w, req)
	}
}

func copyAndClose(src, dst *net.TCPConn) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("Error copying to client: %s\n", err.Error())
	}

	dst.CloseWrite()
	src.CloseRead()
}

// Config is the access deny rules
type Config struct {
	Accept []string `json:"accept,omitempty"`
	Deny   []string `json:"deny,omitempty"`
}
