package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gouthamve/pixie"
)

func main() {

	configFile := flag.String("config", "config.json", "The config file path")
	flag.Parse()

	content, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}

	cfg := pixie.Config{}
	if err := json.Unmarshal(content, &cfg); err != nil {
		panic(err)
	}

	px, err := pixie.NewPixie(cfg)
	if err != nil {
		panic(err)
	}

	s := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(px.Forward),
	}
	log.Fatalln(s.ListenAndServe())
}
