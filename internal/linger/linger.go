package linger

import (
	"LingerAPI/pkg/spray"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly/v2"
)

type Linger struct {
	Version        string
	BioFilter      []string
	InternalFilter []string
}

func NewLinger(bioFilter, internalFilter []string) *Linger {
	newLinger := &Linger{
		Version:        "0.7",
		BioFilter:      bioFilter,
		InternalFilter: internalFilter,
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	log.SetOutput(logger.Writer())

	return newLinger
}

func decodeBio(bio string) (string, error) {
	// Decoding Unicode Escape to default Unicode
	bio = strings.ReplaceAll(bio, `\u002F`, "/")

	//is url correct
	if !utf8.ValidString(bio) {
		log.Println("URL decoding error")
		return "", nil
	}

	//bio to lowercase
	bio = strings.ToLower(bio)

	if !strings.HasPrefix(bio, "https://") && !strings.HasPrefix(bio, "http://") {
		bio = "https://" + bio
	}

	return bio, nil
}

func (s *Linger) MockScrapTikTok(tiktokUrl string) (bool, int, string, int) {
	log.Println(spray.Mspray("[Linger]: Scrapping "), tiktokUrl)
	return true, 111, "testUID", 1010101010
}

func (s *Linger) ScrapTikTokProxy(tiktokUrl string, client *http.Client, pageHTMLChannel chan<- string) {

	log.Println(spray.Mspray("[Linger Proxy]: Scrapping TikTok for"), tiktokUrl)
	c := colly.NewCollector()
	var pageHTML string
	_prefix := "https://tiktok.com/@"

	//proxy
	if client != nil {
		c.WithTransport(client.Transport)
	}

	c.OnHTML("html", func(e *colly.HTMLElement) {
		//e.Response.Headers.Set("User-Agent", "")
		pageHTML = e.Text
		if pageHTML == "" {
			log.Println(spray.Rspray("[Linger Proxy]: PageHTML is empty"))
		}
		//send to the channel
		pageHTMLChannel <- pageHTML
	})

	//visit
	c.Visit(_prefix + tiktokUrl)
}

// GetFollowers
func (s *Linger) GetFollowers(pageHTML string, followersChannel chan<- int) {
	startIndex := strings.Index(pageHTML, `{"followerCount":`)

	if startIndex < 0 {
		followersChannel <- 0
		close(followersChannel)
		log.Println(spray.Mspray("[Linger Proxy]: No Followers Tag"))
		return
	}

	jsonString := pageHTML[startIndex:]

	var followers int //result

	//search tag by regexp
	re := regexp.MustCompile(`"followerCount":\s*(\d+)`)
	match := re.FindStringSubmatch(jsonString)
	if len(match) == 2 {
		followers, _ = strconv.Atoi(match[1])
		followersChannel <- followers
		log.Println(spray.Mspray("[Linger Proxy]: Followers"), followers)
	}

	close(followersChannel)

}

// GetSecUid
func (s *Linger) GetSecUid(pageHTML string, secUidChannel chan<- string) {
	secUidIndex := strings.Index(pageHTML, `"secUid":"`)

	var secUid string

	if secUidIndex > 0 {
		secUidJson := pageHTML[secUidIndex:]
		secUidRegex := regexp.MustCompile(`"secUid"\s*:\s*"([^"]+)"\,`)
		secUidMatch := secUidRegex.FindStringSubmatch(secUidJson)
		if len(secUidMatch) == 2 {
			secUid = secUidMatch[1]
			secUidChannel <- secUid

		}
	} else {
		secUidChannel <- ""
		close(secUidChannel)
		log.Println(spray.Mspray("[Linger Proxy]: No SecUid Tag"))
		return
	}
}

// GetId
func (s *Linger) GetId(pageHTML string, idChannel chan<- int) {
	idIndex := strings.Index(pageHTML, `{"user":{"id":"`)

	var id int

	if idIndex > 0 {
		idJson := pageHTML[idIndex:]
		idRegex := regexp.MustCompile(`{"user":{"id":"\s*(\d+)`)
		idMatch := idRegex.FindStringSubmatch(idJson)

		if len(idMatch) == 2 {
			id, _ = strconv.Atoi(idMatch[1])
			idChannel <- id

		}
	} else {
		idChannel <- 0
		close(idChannel)
		log.Println(spray.Mspray("[Linger Proxy]: No Id Tag"))
		return
	}

}

// GetRedirectUrl
func (s *Linger) GetRedirectUrl(pageHTML string, redirectUrlChannel chan<- string) {
	//redirectURL, captcha

	redirectUrlIndex := strings.Index(pageHTML, `"redirectUrl":"`)
	if redirectUrlIndex > 0 {
		redirectUrlJson := pageHTML[redirectUrlIndex:]
		redirectUrlRegex := regexp.MustCompile(`"redirectUrl":"([^"]+)"`)
		redirectUrlMatch := redirectUrlRegex.FindStringSubmatch(redirectUrlJson)
		if len(redirectUrlMatch) == 2 {
			//log.Println(yspray("Raw redirectUrl: " + redirectUrlMatch[1]))
			_redirectResult, err := decodeString(redirectUrlMatch[1])
			if err != nil {
				log.Println("decodeString error: ", err)
			}
			redirectUrlChannel <- _redirectResult
			log.Println(spray.Yspray("Found redirectUrl: " + _redirectResult))
		}

	} else {
		redirectUrlChannel <- ""
		close(redirectUrlChannel)
		log.Println(spray.Mspray("[Linger Proxy]: No RedirectURL Tag"))
		return
	}

}

func decodeString(input string) (string, error) {
	//decoding json
	var decodedStr string
	err := json.Unmarshal([]byte(`"`+input+`"`), &decodedStr)
	if err != nil {
		return "", fmt.Errorf("decoding json error: : %v", err)
	}

	//changing escape-sequences to default Unicode
	decodedStr, err = strconv.Unquote(`"` + decodedStr + `"`)
	if err != nil {
		return "", fmt.Errorf("unquoting string error: %v", err)
	}

	return decodedStr, nil
}

func (s *Linger) ScrapTikTok(tiktokUrl string) (bool, int, string, int, string) {

	log.Println(spray.Mspray("[Linger]: Scrapping TikTok..."))
	c := colly.NewCollector()
	found := false
	var followers int
	var secUid string
	var id int
	var pageHTML string
	_prefix := "https://tiktok.com/@"

	c.OnHTML("html", func(e *colly.HTMLElement) {

		//e.Response.Headers.Set("User-Agent", "")
		pageHTML = e.Text
		file, err := os.Create("out.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)

		_, err = writer.WriteString(pageHTML)
		if err != nil {
			log.Fatal(err)
		}

		err = writer.Flush()
		if err != nil {
			log.Fatal(err)
		}
		//redirectUrl checking
		redirectUrlIndex := strings.Index(pageHTML, `redirectUrl`)

		if redirectUrlIndex > 0 {
			log.Println(spray.Yspray("Found redirectUrl"))
		}

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
		log.Println(spray.Mspray("[Linger]: ") + match[1] + " followers")

	})

	startTime := time.Now()
	c.Visit(_prefix + tiktokUrl)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	log.Printf(spray.Bspray("Time: %s"), elapsedTime)
	//extracting secUID

	secUidIndex := strings.Index(pageHTML, `"secUid":"`)
	if secUidIndex > 0 {
		secUidJson := pageHTML[secUidIndex:]
		secUidRegex := regexp.MustCompile(`"secUid"\s*:\s*"([^"]+)"\,`)
		secUidMatch := secUidRegex.FindStringSubmatch(secUidJson)
		if len(secUidMatch) == 2 {
			secUid = secUidMatch[1]
		}
	}

	//extracting id
	idIndex := strings.Index(pageHTML, `{"user":{"id":"`)
	if idIndex > 0 {
		idJson := pageHTML[idIndex:]
		idRegex := regexp.MustCompile(`{"user":{"id":"\s*(\d+)`)
		idMatch := idRegex.FindStringSubmatch(idJson)

		if len(idMatch) == 2 {
			id, _ = strconv.Atoi(idMatch[1])
		}
	}

	return found, followers, secUid, id, pageHTML

}

func (s *Linger) ScrapBioLink(extractedBio string, client *http.Client) ([]string, []string, error) {
	log.Println(spray.Gspray("[Linger]: Scrapping bio for "), extractedBio)
	c := colly.NewCollector()

	if client != nil {
		c.WithTransport(client.Transport)
	}

	var BioLinks []string
	var InternalLinks []string

	for _, filterLink := range s.BioFilter {
		if extractedBio != "" && strings.Contains(strings.ToLower(extractedBio), filterLink) {
			//log.Println(yspray("Link from filter: " + filterLink))
			BioLinks = append(BioLinks, extractedBio)

			//inside biolink
			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				for _, internalLink := range s.InternalFilter {
					if e.Attr("href") != "" && strings.Contains(e.Attr("href"), internalLink) {
						socLink := e.Attr("href")
						log.Printf(spray.Yspray("[%s]: %s"), filterLink, socLink)
						InternalLinks = append(InternalLinks, socLink)
					}
				}
			})
			if extractedBio != "" {
				c.Visit(extractedBio)
			}
		}
	}

	c.OnRequest(func(r *colly.Request) {

		log.Println(spray.Rspray("Request: "), r.URL)

	})

	return BioLinks, InternalLinks, nil
}

// main function to scrap bio and internal links
func (s *Linger) StartScrapping(username string, pageHTML string) ([]string, []string, error) {

	c := colly.NewCollector()

	var BioLinks []string
	var InternalLinks []string

	//Extracting BioLink from json begins with tag bioLink
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

	//our regex
	reg := regexp.MustCompile(`"link":"(.*?)"`)
	match := reg.FindStringSubmatch(jsonString)

	if len(match) > 1 {
		bioLink := match[1]
		//log.Println("bioLink: " + bioLink)
		res, err := decodeBio(bioLink)
		extractedBio = res
		if err != nil {
			log.Println("decodeBio error: ", err)
		}
		log.Println(spray.Yspray("Extracted Bio: " + res))
	}

	//parsing biolink
	for _, filterLink := range s.BioFilter {
		if extractedBio != "" && strings.Contains(strings.ToLower(extractedBio), filterLink) {
			log.Println(spray.Yspray("Link from filter: " + filterLink))
			BioLinks = append(BioLinks, extractedBio)

			//inside biolink
			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				for _, internalLink := range s.InternalFilter {
					if e.Attr("href") != "" && strings.Contains(e.Attr("href"), internalLink) {
						socLink := e.Attr("href")
						log.Printf(spray.Yspray("[%s]: %s"), filterLink, socLink)
						InternalLinks = append(InternalLinks, socLink)
					}
				}
			})
			if extractedBio != "" {
				c.Visit(extractedBio)
			}
		}
	}

	c.OnRequest(func(r *colly.Request) {

		log.Println(spray.Rspray("Request: "), r.URL)

	})

	return BioLinks, InternalLinks, nil

}
