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
	"time"
	"strconv"
	"crypto/rand"
	"flag"
)

type torrentInfo struct {
	id hash.Hash
	complete int
	downloaded int
	incomplete int
}

var version string = "0001"

var port = flag.String("p", "6885", "Sets Client P2P Port")

func peerIdCreation() []byte {

	peer_id := make([]byte,0)
	peer_id = append(peer_id, []byte("GT")...)
	peer_id = append(peer_id, []byte(version)...)
	peer_id = append(peer_id, []byte("-")...)
	
	t := uint16(time.Now().Nanosecond())
	if( t<100) { //must have 3 significant time digits 
		t+=100
	}
	ts := strconv.FormatInt(int64(t),10)
	peer_id = append(peer_id, []byte(ts[0:3])...)
	
	//Add 10 more completly random bytes to peer_id
	r := make([]byte,10)
	rand.Read(r)
	peer_id = append(peer_id,r...)
	
	return peer_id
}



//Returns "" if tracker does not support scrape
//i.e. the announce does not have an announce after
//the last '/'
func announceToScrape(aUrl *url.URL) *url.URL {

	rAnnounceCheck, err := regexp.Compile("(/announce)([^/]*$)")
	if err != nil {
		log.Fatalln("GOLANG UPSTREAM REGEXP BROKEN - REPORT!")
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
	
	fmt.Println("hashy")
	fmt.Printf("% x\n", tHash.Sum(nil))
	fmt.Println(info["piece length"])
	fmt.Println("hashy")
	
	peer_id := peerIdCreation()
	fmt.Println(peer_id)
	fmt.Println(string(peer_id))


	
	uploaded := 0;
	downloaded := 0;
	left := info["length"].(int)
	compact := 0
	event := "started"
	
	aUrl,err := url.Parse(string(src["announce"].(string)))
	if err != nil {
		log.Fatalln("Can't parse announce tracker url: ",err)
	}
	v := url.Values{}
	v.Set("info_hash", string(tHash.Sum(nil)))
	v.Set("peer_id", string(peer_id))
	v.Set("uploaded", strconv.Itoa(uploaded))
	v.Set("downloaded", strconv.Itoa(downloaded))
	v.Set("left", strconv.Itoa(left))
	v.Set("compact", strconv.Itoa(compact))
	v.Set("event", event)
	
	aUrl.RawQuery = v.Encode()
	fmt.Println()
	fmt.Println("Sending:")
	fmt.Println("\t"+aUrl.String())
	aResp, err := http.Get(aUrl.String())
	if err != nil {
		log.Fatalln("Error announcing to tracker: ", err)
	}
	text,_ := ioutil.ReadAll(aResp.Body)
	fmt.Println()
	fmt.Println(string(text))
}


























