package proxy

import (
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func ParseProxy() (*url.URL, error) {

	err := godotenv.Load("../configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	proxyString := os.Getenv("PROXY")

	proxyURL, proxyUrlError := url.Parse(proxyString)
	if proxyUrlError != nil {
		log.Fatalf("Url Proxy error: ", proxyUrlError)
	}
	return proxyURL, proxyUrlError
}
