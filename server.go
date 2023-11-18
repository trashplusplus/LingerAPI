package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type LingerServer struct {
	Host string
	Port string
}

//
type ResponseData struct {
	UniqueId  string   `json:"uniqueId"`
	Followers int      `json:"followers"`
	SecUid    string   `json:"secUid"`
	Id        int      `json:"id"`
	Biolink   []string `json:"bio"`
	SocLinks  []string `json:"soclinks"`
}

type LinktreeData struct {
	Biolink  []string `json:"bio"`
	SocLinks []string `json:"soclinks"`
}

//tt is not found response
type NotFoundData struct {
	Message string `json:"message"`
}

func NewLingerServer() *LingerServer {
	return &LingerServer{
		Host: "localhost",
		Port: "3030",
	}
}

func (s *LingerServer) StartServer() {
	//randomize index of proxy in list
	rand.Seed(time.Now().UnixNano())

	serverIP := s.Host + ":" + s.Port
	log.Println(mspray("[LingerAPI]: Started at: "), serverIP)

	go func() {
		http.HandleFunc("/api/bio", func(w http.ResponseWriter, r *http.Request) {

			linger := NewLinger()

			name := r.URL.Query().Get("username")
			var responseData interface{}
			bio, internalLinks, err := linger.ScrapBioLink(name)

			if internalLinks == nil {
				log.Println(rspray("[Linger]: No Soclinks for"), name)
			} else {
				responseData = LinktreeData{
					Biolink:  bio,
					SocLinks: internalLinks,
				}
			}

			jsonData, err := json.Marshal(responseData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//response

			w.Header().Set("Content-Type", "application/json")
			if internalLinks == nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			w.Write(jsonData)

		})
	}()

	go func() {
		http.HandleFunc("/api/tiktok", func(w http.ResponseWriter, r *http.Request) {

			linger := NewLinger()
			//get request
			if r.Method != http.MethodGet {
				http.Error(w, "Method is not available", http.StatusMethodNotAllowed)
				return
			}
			/*
			   proxyURL, proxyUrlError := parseProxy()
			   if proxyUrlError != nil {
			     log.Fatalf("Url Proxy error: ", proxyUrlError)
			   }

			   //proxy client
			   client := &http.Client{
			     Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
			   }
			*/

			name := r.URL.Query().Get("username")
			var responseData interface{}

			isFound, followers, secUid, id, pageHTML := linger.ScrapTikTok(name)
			//isFound, followers, secUid, id, pageHTML := linger.MockScrapTikTok(name)

			if isFound {

				bio, internalLinks, err := linger.StartScrapping(name, pageHTML)
				if err != nil {
					log.Println(rspray("Scrapping error: "), err)
				}

				//log.Println(rspray("PageHTML: "), pageHTML)
				responseData = ResponseData{
					UniqueId:  name,
					Biolink:   bio,
					SocLinks:  internalLinks,
					Followers: followers,
					SecUid:    secUid,
					Id:        id,
				}

			}

			if !isFound {
				responseData = NotFoundData{Message: "tt is not found"}
				log.Println(rspray("[Linger]: tt is not found"))
			}

			jsonData, err := json.Marshal(responseData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)

		})
	}()
	//server starting...
	err := http.ListenAndServe(serverIP, nil)
	if err != nil {
		log.Printf("Starting server error: %s\n", err)
	}
}
