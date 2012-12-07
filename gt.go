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
	itemString
	itemList
	itemDictionary
	itemError
	itemNil
)

type itemType int

type item struct {
	typ itemType
	//int.(string), string.(string), list.(listItem), dict.(dictItem) 
	val  interface{}
}

type listItem struct {
	it item
	next *listItem
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

//Debug print for fmt.Print#################
func printList(it *item) string {
	var il* listItem
	var val string
	
	*il = it.val.(listItem)
	val = il.it.val.(string)
	for il = il.next; il != nil; {
		val = val+" "+il.it.val.(string)
	}
	return val
}
	
func (t *tlex) String() string {
	var text, typePrefix, itVal string
	
	for i := 0; i < len(t.items); i++ {
		switch (*t).items[i].typ {
		case itemInteger:
			typePrefix = "int"
			itVal = t.items[i].val.(string)
		case itemString:
			typePrefix = "bytestring"
			itVal = t.items[i].val.(string)
		case itemList:
			typePrefix = "list"
			itVal = printList(&t.items[i])
	//	case itemDictionary:
	//		typePrefix = "dic"
		default:
			typePrefix = "type_unknowned"
			itVal = "WHATTA?"
		}
		
		if i != len(t.items)-1 {
			text = text + typePrefix + ": " + itVal + "\n"
		} else {
			text = text + typePrefix + ": " + itVal
		}
	}
	return text
}

//#################################

func getInt(t *tlex) string {
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
			 return string(buf)
		default:
			log.Fatalln(t.fileName, ": Out of bonds parsing an integer at pos: ", t.pos)
		}
	}
	return "" //Never reach
}

func lexInt(t* tlex) {
	val := getInt(t)
	t.items = append(t.items, item{itemInteger, val})
}

func getString(t *tlex) string {
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

func lexString(t* tlex) {
	var it item
	str := getString(t)
	it.typ = itemString
	it.val = str
	t.items = append(t.items, it)
}

func listParse (t* tlex) (string, itemType) {
	r := t.next()
	switch r {
	case 'i':
		return string(getInt(t)), itemInteger
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return getString(t), itemString
	case 'e':
		return "", itemError
	default:
		log.Fatalln(t.fileName, ": Somthing is haywire while parsing list at ", t.pos)
	}
	return "", itemError //Never reach
}

func getList(t *tlex) listItem {
	var itCur, itPrev *listItem
	var itFirst listItem
	
	val, typ := listParse(t)
	itFirst.it = item{typ, val}
	itPrev = &itFirst
	for val, typ := listParse(t); typ != itemNil;{
		fmt.Println(typ)
		itCur = new(listItem)
		itCur.it = item{typ, val}
		itPrev.next = itCur
		itPrev = itCur
	}
	return itFirst
	
}

func lexList(t *tlex) {
	var it item
	li := getList(t)
	it.typ = itemList
	it.val = li
	t.items = append(t.items, it)
}

func main() {
	var t tlex

	data, err := ioutil.ReadFile("list.torrent")
	if err != nil {
		log.Fatal(err)
	}
	t.fileName = "NeedsToBeImplemented.torrent"
	t.data = data
	for {
		r := t.next()
		switch r {
		case 'd':
			fmt.Println("Woop! dict is not implemented")
			goto done
		case 'i':
			lexInt(&t)
		case 'l':
			lexList(&t)
			goto done
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			lexString(&t)
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
