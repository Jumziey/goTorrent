package main

import (
	"bittor"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"hash"
)

type torrentInfo struct {
	id hash.Hash
	complete int
	downloaded int
	incomplete int
}

//Returns "" if tracker does not support scrape
//i.e. the announce does not have an announce after
//the last '/'
func announceToScrape(aUrl *url.URL) *url.URL {

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

//Takes the master dict and returns a suitable announce url as a url.Url
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

func scrape(aUrl *url.URL,info_hash hash.Hash) ([]byte,error) {
	sUrl := announceToScrape(aUrl)
	if sUrl == nil {
		return nil, errors.New("It would seem there is no scrape convention supported")
	}

	v := url.Values{}
	v.Set("info_hash", string(info_hash.Sum(nil)))
	sUrl.RawQuery = v.Encode()
	aUrl.RawQuery = v.Encode()

	sResp, err := http.Get(sUrl.String())
	if err != nil {
		return nil, err
	}
	defer sResp.Body.Close()

	scrape_body, err := ioutil.ReadAll(sResp.Body)
	if err != nil {
		return nil, err
	}
	return scrape_body, nil
}



func main() {
	tName := "crunch.torrent"
	fmt.Println("############start############")
	
	tData,err := ioutil.ReadFile(tName)
	if err != nil {
		log.Fatalln(err)
	}
	
	
	src, err := bittor.GetMainDict(tData)
	if err != nil {
		log.Fatalln("Error decoding ", tName, ":", err)
	}
	info, err := bittor.GetInfoDict(src)
	if err != nil {
		log.Fatalln("Found no info dict in torrent file")
	}
	
	tHash, err := bittor.GetInfoHash(tData)
	if err != nil {
		log.Fatalln("Problems with hashing the info value of torrent: ", err)
	}
	
	fmt.Printf("% x\n", tHash.Sum(nil))
	fmt.Println(info["piece length"])
}
