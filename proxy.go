package main

import (
	"log"
	"net/url"
)

func parseProxy() (*url.URL, error) {
	proxyURL, proxyUrlError := url.Parse("")
	if proxyUrlError != nil {
		log.Fatalf("Url Proxy error: ", proxyUrlError)
	}
	return proxyURL, proxyUrlError
}
