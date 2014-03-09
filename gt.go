package main

import (
	"bittor"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	torrent := torrentInfo{};

	fmt.Println("############start############")
	fmt.Println()
	tName := "crunch.torrent"
	
	tFile, err := os.Open(tName)
	if err != nil {
		log.Fatalln(err)
	}
	defer tFile.Close()
	

	src, err := bittor.GetMainDict(tFile,tName)
	if err != nil {
		log.Fatalln("Error decoding ", tFile, ":", err)
	}
	info := bittor.GetInfoDict(src)
	if info == nil {
		log.Fatalln("Found no info dict in torrent file")
	}

	/*
	tFile.Seek(0, 0)
	info_hash, err := bittor.InfoHash(f)
	if err != nil {
		log.Fatalln(err)
	}
	*/

	aUrl, err := getAnnounce(src);
	if err != nil {
		log.Fatalln(err)
	}
	
	scrape_body, err := scrape(aUrl, info_hash);
	if err != nil {
		log.Fatalln(err)
	}
	
	b_data, err := bittor.GetDictFromByte(scrape_body);
	if err != nil {
		log.Fatalln(err)
	}
	
	actual_torrent := bittor.GetDict(b_data, "files")
	hashy := bittor.GetDict(actual_torrent, string(info_hash.Sum(nil)))
	
	torrent.id = info_hash;
	torrent.complete = hashy["complete"].(int)
	torrent.incomplete = hashy["incomplete"].(int)
	torrent.downloaded = hashy["downloaded"].(int)

	fmt.Println(torrent.complete);
}

