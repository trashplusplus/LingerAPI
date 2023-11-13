package main

import(
"net/url"
"log"
)

func parseProxy() (*url.URL, error){
     log.Println(mspray("[Linger] Proxy is loading..."))
     proxyURL, proxyUrlError := url.Parse("http://proxyHere")
        if proxyUrlError != nil {
          log.Fatalf("Url Proxy error: ", proxyUrlError)
        }

      return proxyURL, proxyUrlError
}

