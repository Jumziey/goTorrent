package main

/* First iteration just test for correct bytestreams in utf8 encoding, no sha1's
*/
import(
	"unicode/utf8"
	"log"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	itemInteger = iota
	itemString
	itemList
	itemDictionary
	itemError
	itemNil
)

type itemType int

type tlex struct {
	data []byte
	pos int
}

func (t *tlex)next() rune {
	r, size := utf8.DecodeRune(t.data[t.pos:])
	if r == utf8.RuneError {
		log.Fatalln("(tlex).rune() return utf8.RuneError")
	}
	t.pos = t.pos+size
	return r
}

func (t *tlex)peek() rune {
	r, _ := utf8.DecodeRune(t.data[t.pos:])
	if r == utf8.RuneError {
		log.Fatalln("(tlex).rune() return utf8.RuneError")
	}
	return r
}

func intParse(t *tlex) string {
	intStr := ""
	for r:= t.next(); r != 'e'; r = t.next() {
		intStr = intStr+string(r)
	}
	return intStr
}

func stringParse(t *tlex) string {
	//A "back step" statement only valid when we know the rune is exactly one byte. 
	//This is the case now since this function is only called if lexParse stumble upons 
	//a digit. We need this digit to figure out how long the string is. 
	t.pos = t.pos - 1
	
	stringSize := ""
	for s:=t.next(); s != ':'; s=t.next() {
		stringSize = stringSize+string(s)
	}
	s_size, err := strconv.Atoi(stringSize)
	if err != nil {
		log.Fatalln("Error in stringParse: ", err)
	}
	
	rString := ""
	for i:=0; i<s_size; i++ {
		rString = rString+string(t.next())
	}
	return rString
}

func listParse(t *tlex) []string {
	var str []string

	//We read until we reach the end 'e' of the list and make this
	//a list item, we peek so we don't fuck it up for lexParse(*tlex)
	for t.peek() != 'e' {
		s := lexParse(t)
		str = append(str, s...)
	}
	t.next() //Throw away the 'e'
	return str
}

func lexParse(t *tlex) []string {
	switch t.next() {
		case 'l':
			return listParse(t)
		case 'i':
			return []string{intParse(t)}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return []string{stringParse(t)}
		default:
			log.Fatalln("Out of bonds in lexParse")
	}
	//Unreachable
	return []string{""}
}

func lexStart(t *tlex) []string {
	var str []string
	
	//Aslong there's been less data read then then data availble, get a new item!
	for len(t.data) > t.pos {
		s := lexParse(t)
		str = append(str, s...)
	}
	return str
}

func main() {
	var t tlex
	var err error
	
	t.data, err = ioutil.ReadFile("list.torrent")
	if err != nil {
		log.Fatalln("Error in main(): ", err)
	}
	
	s := lexStart(&t)
	fmt.Println(s)
}
	