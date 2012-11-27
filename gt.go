package main

import (
	"fmt"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
	"log"
	"strconv"
)

const (
	itemInteger = iota
	itemBytestring
	itemList
	itemDictionary
)

type itemType int

type item struct {
	typ itemType
	val  interface{}
}

type tlex struct {
	data  []byte
	pos   int
	items []item
	fileName string //For error messages
}

//Returns the next rune in the buffer and increment the pos
func (t *tlex) next() rune {
	r, size := utf8.DecodeRune(t.data[t.pos:])
	if r == utf8.RuneError {
		return r
	}
	t.pos = t.pos + size
	return r
}

//Debug print for fmt.Print
func (t *tlex) String() string {
	var text string
	var typePrefix string
	for i := 0; i < len(t.items); i++ {
		switch (*t).items[i].typ {
		case itemInteger:
			typePrefix = "int"
		case itemBytestring:
			typePrefix = "bytestring"
		case itemList:
			typePrefix = "list"
		case itemDictionary:
			typePrefix = "dic"
		default:
			typePrefix = "type_unknowned"
		}
		itVal := t.items[i].val.(string)
		if i != len(t.items)-1 {
			text = text + typePrefix + ": " + itVal + "\n"
		} else {
			text = text + typePrefix + ": " + itVal
		}
	}
	return text
}

func getInt(t *tlex) item {
	var it item
	buf := make([]byte, 10) //max integer size
	for i := 0; ; i++ {
		r := t.next()
		switch {
		case unicode.IsDigit(r):
			size := utf8.EncodeRune(buf[i:], r)
			if size != 1 {
				log.Fatalln(t.fileName, ": Illegal character found while parsing an integer at pos: ", t.pos);
			}
		case r == 'e' && len(buf) != 0:
			it.typ = itemInteger
			it.val = string(buf)
			return it
		default:
			log.Fatalln(t.fileName, ": Out of bonds parsing an integer at pos: ", t.pos)
		}
	}
	return it
}

func lexInt(t* tlex) {
	it := getInt(t)
	t.items = append(t.items, it)
}

func getBytestring(t *tlex) string {
	//We know the first number is thrown away but we need it
	//We also know a number is only one byte so this is ok
	t.pos = t.pos - 1
	
	numBuf := make([]byte, 10)
	for i := 0;;i++{
		r := t.next()
		switch r {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				numBuf[i] = byte(r) //For numbers: ascii bytes and runes the same
			case ':':
				goto done;
			default:
				log.Fatalln(t.fileName, ": Out of bonds parsing a bytestring at pos: ", t.pos)
		}
	}
done:
	n := firstNil(numBuf) //Helper func, hack for string conversion baddieness. 
	woop := string(numBuf[:n])
	size, err := strconv.Atoi(woop)
	if err != nil {
		log.Fatalln(t.fileName, ": Somthing is haywire while parsing bytestring at ", t.pos, " ERROR: ", err)
	}
	buf := make([]byte, size)
	for i,j := 0,0; j<size;j++  {
		r := t.next()
		if n:=utf8.RuneLen(r); n>1 {
			buf = incByteSlice(buf, n) //Helper func
		}
		r_size := utf8.EncodeRune(buf[i:], r)
		i = i+r_size
	}
	return string(buf)
}

func lexBytestring(t* tlex) {
	var it item
	bytestring := getBytestring(t)
	it.typ = itemBytestring
	it.val = bytestring
	t.items = append(t.items, it)
}
	

func main() {
	var t tlex

	data, err := ioutil.ReadFile("bytestring.torrent")
	if err != nil {
		log.Fatal(err)
	}
	t.fileName = "NeedsToBeImplemented.torrent"
	t.data = data
	for {
		r := t.next()
		switch r {
		case 'd':
			fmt.Println("WOOP! Dict not implemented")
			goto done
		case 'i':
			fmt.Println("OH MY")
			lexInt(&t)
		case 'l':
			fmt.Println("Woop! list not implemented")
			goto done
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			lexBytestring(&t)
		case utf8.RuneError:
			fmt.Println("Now its all done!")
			goto done
		case '\n': //ignoring newline
		default:
			fmt.Println("We got something extra: ", r)
		}
	}
done:

	fmt.Println(&t)
}
