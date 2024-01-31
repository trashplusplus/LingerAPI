package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"io/ioutil"
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
	//filters init
	//deviceHeader := readFilterFromFile("filter/deviceHeader.txt")
	bioFilter := readFilterFromFile("filter/bio.txt")
	internalFilter := readFilterFromFile("filter/soc.txt")

	log.Println(mspray(" _ _                 _____ _____ _____ "))
	log.Println(mspray("| |_|___ ___ ___ ___|  _  |  _  |     |"))
	log.Println(mspray("| | |   | . | -_|  _|     |   __|-   -|"))
	log.Println(mspray("|_|_|_|_|_  |___|_| |__|__|__|  |_____|"))
	log.Println(mspray("        |___|                          "))
	log.Println(mspray("[LingerAPI]: Started at: "), serverIP)


  //test
  /*
	_r, _ := readFromFile("catia_carla.txt")
	ch := make(chan string)

  l := NewLinger(bioFilter, internalFilter)
  go l.GetRedirectUrl(_r, ch)

  res := <-ch

  log.Println("result: ",res)
  */


		http.HandleFunc("/api/bio", func(w http.ResponseWriter, r *http.Request) {

			linger := NewLinger(bioFilter, internalFilter)

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


		http.HandleFunc("/api/tiktok", func(w http.ResponseWriter, r *http.Request) {

			linger := NewLinger(bioFilter, internalFilter)
			//get request
			if r.Method != http.MethodGet {
				http.Error(w, "Method is not available", http.StatusMethodNotAllowed)
				return
			}

      //proxy
			proxyURL, proxyUrlError := parseProxy()
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

			//isFound, followers, secUid, id := linger.MockScrapTikTok(name)

      go linger.ScrapTikTokProxy(name, client, pageHTMLChannel)


      //ожидаем данних из канала
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

      var internalLinks[] string
      var bio[] string

       if secUid != "" {
         bio, internalLinks, _ = linger.StartScrapping(name, pageHTML)
            //response
            responseData = ResponseData{
                UniqueId: name,
                Biolink: bio,
                SocLinks: internalLinks,
                Followers: followers,
                SecUid: secUid,
                Id: id,
            }


       }else{
          go linger.GetRedirectUrl(pageHTML, redirectUrlChannel)
          redirectUrl := <-redirectUrlChannel


         if redirectUrl != "" {
             responseData = NotFoundData{
              Message: "tt is exist!",
            }
         }else{
            //not found
            responseData = NotFoundData{
              Message: "tt is not found",
            }
              log.Println(rspray("[Linger]: tt is not found | " + name))
          }


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

      //close(pageHTMLChannel)
      //close(followersChannel)
      //close(secUidChannel)
      //close(idChannel)
		})

	//server starting...
	err := http.ListenAndServe(serverIP, nil)
	if err != nil {
		log.Printf("Starting server error: %s\n", err)
	}
}


func readFromFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
