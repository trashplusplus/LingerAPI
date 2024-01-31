package main

import (
	"log"
	"net/url"
)

func parseProxy() (*url.URL, error) {
	//log.Println(mspray("[Linger] Proxy is loading..."))
	proxyURL, proxyUrlError := url.Parse("proxyhere")
	if proxyUrlError != nil {
		log.Fatalf("Url Proxy error: ", proxyUrlError)
	}
	return proxyURL, proxyUrlError
}
