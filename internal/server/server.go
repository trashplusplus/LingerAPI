package server

import (
	"LingerAPI/internal/linger"
	"LingerAPI/internal/proxy"
	"LingerAPI/pkg/filter"
	"LingerAPI/pkg/spray"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var semaphore = make(chan struct{}, 64)

type LingerServer struct {
	Host string
	Port string
}

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

// tt is not found response
type NotFoundData struct {
	Message string `json:"message"`
}

func NewLingerServer() *LingerServer {

	err := godotenv.Load("../configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &LingerServer{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}

}

func (s *LingerServer) StartServer() {

	serverIP := s.Host + ":" + s.Port
	//filters init
	bioFilter := filter.ReadFilterFromFile("../configs/filter/bio.txt")
	internalFilter := filter.ReadFilterFromFile("../configs/filter/soc.txt")

	log.Println(spray.Mspray(" _ _                 _____ _____ _____ "))
	log.Println(spray.Mspray("| |_|___ ___ ___ ___|  _  |  _  |     |"))
	log.Println(spray.Mspray("| | |   | . | -_|  _|     |   __|-   -|"))
	log.Println(spray.Mspray("|_|_|_|_|_  |___|_| |__|__|__|  |_____|"))
	log.Println(spray.Mspray("        |___|                          "))
	log.Println(spray.Mspray("[LingerAPI]: Started at: "), serverIP)
	log.Println(spray.Mspray("[LingerAPI]: Your proxy: "), os.Getenv("PROXY"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "LingerAPI", http.StatusNotFound)
	})

	//api/bio endpoint handler
	http.HandleFunc("/api/bio", func(w http.ResponseWriter, r *http.Request) {

		//semaphore queue to load balance
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		//proxy
		proxyURL, proxyUrlError := proxy.ParseProxy()
		if proxyUrlError != nil {
			log.Fatalf("Url Proxy error: ", proxyUrlError)
		}

		//proxy client
		client := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		}

		linger := linger.NewLinger(bioFilter, internalFilter)

		name := r.URL.Query().Get("username")
		var responseData interface{}
		bio, internalLinks, _ := linger.ScrapBioLink(name, client)

		if internalLinks == nil {
			log.Println(spray.Rspray("[Linger]: No Soclinks for"), name)
			responseData = NotFoundData{
				Message: "no soclinks for " + name,
			}
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
	//api/tiktok endpoint handler
	http.HandleFunc("/api/tiktok", func(w http.ResponseWriter, r *http.Request) {
		//semaphore queue to load balance
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		if r.URL.Query().Get("username") == "" {
			w.Header().Set("Content-Type", "application/json")
			jsonData, _ := json.Marshal(NotFoundData{
				Message: "username is empty",
			})
			w.WriteHeader(http.StatusNotFound)
			w.Write(jsonData)
			return
		}

		_responseCode := http.StatusOK
		linger := linger.NewLinger(bioFilter, internalFilter)
		//get request
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not available", http.StatusMethodNotAllowed)
			return
		}

		//proxy
		proxyURL, proxyUrlError := proxy.ParseProxy()
		if proxyUrlError != nil {
			log.Fatalf("Url Proxy error: ", proxyUrlError)
		}

		//proxy client
		client := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		}
		//contrller key
		name := r.URL.Query().Get("username")

		var responseData interface{}

		pageHTMLChannel := make(chan string)
		followersChannel := make(chan int)
		secUidChannel := make(chan string)
		idChannel := make(chan int)
		redirectUrlChannel := make(chan string)

		go linger.ScrapTikTokProxy(name, client, pageHTMLChannel)

		//waiting data from pageHTMLChannel
		pageHTML, ok := <-pageHTMLChannel
		if !ok {
			log.Println("Channel is closed, no cakes :'(")
			return
		}

		close(pageHTMLChannel)

		go linger.GetFollowers(pageHTML, followersChannel)
		go linger.GetSecUid(pageHTML, secUidChannel)
		go linger.GetId(pageHTML, idChannel)

		followers := <-followersChannel
		secUid := <-secUidChannel
		id := <-idChannel

		var internalLinks []string
		var bio []string

		if secUid != "" {
			bio, internalLinks, _ = linger.StartScrapping(name, pageHTML)
			//response
			responseData = ResponseData{
				UniqueId:  name,
				Biolink:   bio,
				SocLinks:  internalLinks,
				Followers: followers,
				SecUid:    secUid,
				Id:        id,
			}

		} else {
			go linger.GetRedirectUrl(pageHTML, redirectUrlChannel)
			redirectUrl := <-redirectUrlChannel

			if redirectUrl != "" {
				responseData = NotFoundData{
					Message: "tt is exist!",
				}

			} else {
				//not found
				responseData = NotFoundData{
					Message: "tt is not found",
				}
				_responseCode = http.StatusNotFound
				log.Println(spray.Rspray("[Linger]: tt is not found | " + name))
			}

		}

		jsonData, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(_responseCode)
		w.Write(jsonData)

	})

	//server starting...
	err := http.ListenAndServe(serverIP, nil)
	if err != nil {
		log.Printf("Starting server error: %s\n", err)
	}
}
