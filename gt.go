package main

import(
	"fmt"
	"github.org/Jumziey/bittor"
	"io/ioutil"
	"log"
)

func main() {
	var t bittor.TorData
	var err error
	
	t.Data, err = ioutil.ReadFile("piece.torrent")
	if err != nil {
		log.Fatalln("Error in main(): ", err)
	}
	
	s := bittor.GetMainDict(&t)
	info := bittor.GetInfoDict(s)
	if info == nil {
		log.Fatalln("Found no info dict in torrent file")
	}
	fmt.Println(info)
}