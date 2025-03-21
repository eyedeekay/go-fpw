package main

import (
	"flag"
	"log"
	"net/url"

	fcw "github.com/eyedeekay/go-fpw"
)

func main() {
	starturl := flag.String("url", "https://duckduckgo.com", "URL to open in Firefox")
	flag.Parse()
	uri, err := url.Parse(*starturl)
	if err != nil {
		log.Println(err)
	}
	hasPortable := fcw.PortablePath()
	if hasPortable != "" {
		log.Println("This is a portable installation.", hasPortable)
	} else {
		log.Println("This is not a portable installation.", hasPortable)
	}
	ui, err := fcw.WebAppFirefox(uri.Hostname(), false, *starturl)
	if err != nil {
		log.Println(err)
	}
	log.Println(ui.Log())
	<-ui.Done()
}
