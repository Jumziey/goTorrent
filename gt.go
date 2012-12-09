package main

import(
	"fmt"
	"github.com/Jumziey/bittor"
	"log"
)

func main() {
	
	s, err := bittor.GetMainDict("daily.torrent")
	if err != nil {
		log.Fatalln("Error decoding torrentfile: ", err)
	}
	info := bittor.GetInfoDict(s)
	if info == nil {
		log.Fatalln("Found no info dict in torrent file")
	}
	fmt.Println(info)
}