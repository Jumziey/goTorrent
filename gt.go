package main

import(
	"fmt"
	"bittor"
	"errors"
	"log"
	"regexp"
	"os"
	"net/url"
	"net/http"
	"io/ioutil"
)


//Returns "" if tracker does not support scrape
//i.e. the announce does not have an announce after
//the last '/'
func announceToScrape(aUrl* url.URL) *url.URL {

	rAnnounceCheck, err := regexp.Compile("(/announce)([^/]*$)")
	if err != nil {
		log.Fatalln("UPSTREAM REGEXP BROKEN - REPORT!")
	}
	
	a := aUrl.String()
	if !rAnnounceCheck.MatchString(a) {
		return nil
	}
	sUrl, err := url.Parse(rAnnounceCheck.ReplaceAllString(a, "/scrape$2"))
	if err != nil {
		return nil
	}	
	return sUrl
 }
 
 //Takes the master dict and returns a suitable announce url as a Url
 func getAnnounce(m map[string]interface{}) (*url.URL, error) {
 	//Should check if the announce tracker works go through the announce list
 	//until we find someone that works, will now only return the main announce 
 	//tracker no matter what.
 	
 	announce := m["announce"].(string)
 	if announce == "" {
 		return nil, errors.New("Found no annonuce in master dict")
 	}
 	
 	aUrl, err := url.Parse(announce)
 	if err != nil {
 		return nil, errors.New(fmt.Sprint("Error url parsing announce: ", err))
 	}
 	
 	return aUrl, nil
 }
 
func main() {
	fmt.Println("start")
	f,err := os.Open("daily.torrent")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	
	s, err := bittor.GetMainDict(f)
	if err != nil {
		log.Fatalln("Error decoding ", "daily.torrent", ":", err)
	}
	info := bittor.GetInfoDict(s)
	if info == nil {
		log.Fatalln("Found no info dict in torrent file")
	}
	
	//Some reading must have been done in bittor.GetMainDict()
	f.Seek(0,0)
	info_hash, err := bittor.InfoHash(f)
	if err != nil {
		log.Fatalln(err)
	}
	
	aUrl, err := getAnnounce(s)
	if aUrl == nil {
		log.Fatalln(err)
	}
	sUrl := announceToScrape(aUrl)
	if sUrl == nil {
		log.Fatalln("It would seem there is no scrape convention supported")
	}
	
	v := url.Values{}
	v.Set("info_hash", url.QueryEscape(string(info_hash.Sum(nil))) )
	sUrl.RawQuery = v.Encode()
	aUrl.RawQuery = v.Encode()
	
	fmt.Println()
	fmt.Println(fmt.Sprintf("%x", info_hash.Sum(nil)))
	fmt.Println()
	fmt.Println(sUrl.String())
	fmt.Println()
	
	sResp, err := http.Get(sUrl.String())
	if err != nil {
		log.Fatalln(err)
	}
	//defer sResp.Body.Close()
	
	
	bleh, _ := ioutil.ReadAll(sResp.Body)
	fmt.Println(string(bleh))
	

	

}