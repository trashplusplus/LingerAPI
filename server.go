package main

import (
  "net/http"
  "encoding/json"
  "log"
  "math/rand"
  "time"
  //"net/url"
)

type LingerServer struct{
  Host string
  Port string
}

//
type ResponseData struct{
  Username string `json:"username"`
  Followers int `json:"followers"`
  Biolink []string `json:"bio"`
  SocLinks []string `json:"soclinks"`

}



//tt is not found response
type NotFoundData struct{
  Message string `json:"message"`
}

func NewLingerServer() *LingerServer {
	return &LingerServer{
      Host: "localhost",
      Port: "3030",
	}
}

func (s *LingerServer) StartServer(){
    //randomize index of proxy in list
    rand.Seed(time.Now().UnixNano())

    serverIP := s.Host + ":" + s.Port
    log.Println(mspray("[LingerAPI]: Started at: "), serverIP)


    //proxy list loading
    //proxyList, proxyErr := getProxyFile("proxy.txt")
    /*
    if proxyErr != nil {
        log.Printf("Ошибка: %v\n", proxyErr)
        return
    }else{
      log.Println("ProxyList sucessfully loaded")
    }
    */


    http.HandleFunc("/api/tiktok", func(w http.ResponseWriter, r *http.Request) {
         linger := NewLinger()
         //get request
        if r.Method != http.MethodGet {
            http.Error(w, "Method is not available", http.StatusMethodNotAllowed)
            return
        }

        //randomize
        /*
        randomIndex := rand.Intn(len(proxyList))
        randomProxy := proxyList[randomIndex]
        proxyURL, proxyUrlError := url.Parse("http://" + randomProxy)
        if proxyUrlError != nil {
          log.Fatalf("Url Proxy error: ", proxyUrlError)
        }
        */

        //proxy client
        /*
        client := &http.Client{
          Transport: &http.Transport{
              Proxy: http.ProxyURL(proxyURL),
          },
          Timeout: time.Second * 15, //timeout
        }

         log.Println("Used Proxy: ", randomProxy)
         */

         name := r.URL.Query().Get("username")
         var responseData interface{}

         isFound, followers, pageHTML := linger.ScrapTikTok(name)

         if isFound {
           bio, internalLinks, err := linger.StartScrapping(name, pageHTML)

         if err != nil{
           log.Println("Error:", err)
            return
         }

          if name == "" {
              http.Error(w, "username param is missing", http.StatusBadRequest)
              return
          }

          responseData = ResponseData{
              Username: name,
              Biolink: bio,
              SocLinks: internalLinks,
              Followers: followers,
            }

         }

         if !isFound {
           responseData = NotFoundData{Message: "tt is not found"}
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

    //server starting...
    err := http.ListenAndServe(serverIP, nil)
    if err != nil {
        log.Printf("Starting server error: %s\n", err)
    }
}
