package main

import(
"github.com/gocolly/colly/v2"
"strings"
"io/ioutil"
"regexp"
"strconv"
"time"
"log"
"os"
"unicode/utf8"
//"net/http"
)

type Linger struct{
  Version string
  BioFilter []string
  InternalFilter []string
}


func NewLinger() *Linger{
    newLinger := &Linger{
    Version: "0.2",
    BioFilter: readFilterFromFile("filter/bio.txt"),
    InternalFilter: readFilterFromFile("filter/soc.txt"),
    }

    logger := log.New(os.Stdout,"", log.Ldate|log.Ltime)
    log.SetOutput(logger.Writer())
   //log.Printf("[LingerAPI] v%s searching...", newLinger.Version)

   return newLinger
}


func readFilterFromFile(filename string) []string {
    data, err := ioutil.ReadFile(filename)
    if err != nil {

        log.Println(rspray("Error reading file: "), err)
        return nil
    }

    lines := strings.Split(string(data), "\n")
    var filter []string

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" {
            filter = append(filter, line)
        }
    }

    return filter
}


func decodeBio(bio string) (string, error) {
    // Преобразование Unicode Escape в обычный Unicode
    bio = strings.ReplaceAll(bio, `\u002F`, "/")

    // Проверка на корректное URL-кодирование
    if !utf8.ValidString(bio) {
          log.Println("URL decoding error")
          return "", nil
    }

    // Преобразование в нижний регистр
    bio = strings.ToLower(bio)

   if !strings.HasPrefix(bio, "https://") && !strings.HasPrefix(bio, "http://") {
        bio = "https://" + bio
    }

    return bio, nil
}

func (s *Linger) ScrapTikTok(tiktokUrl string) (bool, int, string){
    log.Println(mspray("[Linger]: Scrapping TikTok..."))
    c := colly.NewCollector()
    found := false
    var followers int
    var pageHTML string
    _prefix := "https://tiktok.com/@"

    //c.WithTransport(client.Transport)

    c.OnHTML("html", func(e *colly.HTMLElement) {
        //get html
        pageHTML = e.Text
        //log.Println("PageHTML: " + pageHTML)
        //get all json with flag followerCount
        startIndex := strings.Index(pageHTML, `{"followerCount":`)
        if startIndex < 0 {
            return
        }

        jsonString := pageHTML[startIndex:]
        //log.Println(jsonString)


      //search tag by regexp
       re := regexp.MustCompile(`"followerCount":\s*(\d+)`)
        match := re.FindStringSubmatch(jsonString)
        if len(match) == 2 {
            followers, _ = strconv.Atoi(match[1])
            found = true
        }
         log.Println(mspray("[Linger]: ") + match[1] + " followers")

    })


    c.Visit(_prefix + tiktokUrl)
	  return found, followers, pageHTML

}


func (s *Linger) StartScrapping(username string, pageHTML string) ([]string, []string, error){

    c := colly.NewCollector()

    tiktUrl := "https://tiktok.com/@" + username

    /*
     if client != nil{
      c.WithTransport(client.Transport)
      //log.Println("Proxy: ", client.Transport)
    }
    */

    var BioLinks[] string
    var InternalLinks[] string


       //Extracting BioLink from json start with tag bioLink
       extractedBio := ""
       //find all html
       //find text starts with bioLink
       startIndex := strings.Index(pageHTML, "bioLink")
       //if not exist - return
       if startIndex < 0 {
            return BioLinks, InternalLinks, nil
        }

        //text started with bioLink
        jsonString := pageHTML[startIndex:]
        //log.Println(jsonString)
        //our reg
        reg := regexp.MustCompile(`"link":"(.*?)"`)
        match := reg.FindStringSubmatch(jsonString)


        if len(match) > 1 {
		      	bioLink := match[1]
		    	  //log.Println("bioLink: " + bioLink)
		    	  res, err := decodeBio(bioLink)
		    	  extractedBio = res
		    	  if err != nil{
		    	  }
            log.Println(yspray("Extracted Bio: " + res))
		    }

    //parsing linktree
        for _, filterLink := range s.BioFilter{
            if extractedBio != "" && strings.Contains(strings.ToLower(extractedBio), filterLink){
                   log.Println(yspray("Link from filter: " + filterLink))
                   BioLinks = append(BioLinks, extractedBio)

                     //inside biolink
                      c.OnHTML("a[href]", func(e * colly.HTMLElement){
                          for _, internalLink := range s.InternalFilter{
                              if e.Attr("href") != "" && strings.Contains(e.Attr("href"), internalLink){
                                    socLink := e.Attr("href")
                                    log.Printf(yspray("[%s]: %s"), filterLink, socLink)
                                    InternalLinks = append(InternalLinks, socLink)
                              }
                          }
                      })
                     if extractedBio != ""{
                      c.Visit(extractedBio)
                     }
            }
        }

    c.OnRequest(func(r *colly.Request){
        log.Println(rspray("Request: "), r.URL)
    })


    startTime := time.Now()
    c.Visit(tiktUrl)
    endTime := time.Now()

    elapsedTime := endTime.Sub(startTime)
    log.Printf(bspray("Time: %s"), elapsedTime)

    /*
    for _, v := range BioLinks {
        log.Println("bioLinks: ", v)
    }

    for _, v := range InternalLinks {
        log.Println("internalLinks: ", v)
    }
    */

    return BioLinks, InternalLinks, nil

}

