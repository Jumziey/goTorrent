package main

import(
	"log"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	itemInt = iota
	itemString
	itemList
	itemDict
	itemNil
)

type itemType int

type item struct {
	typ	itemType
	val	interface{}
}

func (i item) String() string {
	switch i.typ {
		case itemInt:
			return fmt.Sprint("Int: ", i.val.(int))
		case itemString:
			return fmt.Sprint("string: ", i.val.(string))
		case itemList:
			return fmt.Sprint("LIST: ", i.val.([]item))
		case itemDict:
			return fmt.Sprint("DICT: ", i.val.(map[item]item))
		default:
			return fmt.Sprint("Can't print this item")
	}
	return fmt.Sprint("Unreachable")  //Odd compiler req.
}


type tdata struct {
	data []byte
	pos int
}

func (t *tdata)next() byte {
	b := t.data[t.pos]
	t.pos = t.pos+1
	return b
}

func (t *tdata)peek() byte {
	return t.data[t.pos]
}

func (t *tdata)prev() {
	t.pos = t.pos-1
}

func intParse(t *tdata) int {
	intStr := ""
	var b byte
	for b= t.next(); b != 'e'; b = t.next() {
		intStr = intStr+string(b)
	}
	integ, err := strconv.Atoi(intStr)
	if err != nil {
		log.Fatalln("Error in intParse: ", err)
	}
	return integ
}

func stringParse(t *tdata) string {
	t.prev()
	
	stringSize := ""
	for s:=t.next(); s != ':'; s=t.next() {
		stringSize = stringSize+string(s)
	}
	s_size, err := strconv.Atoi(stringSize)
	if err != nil {
		log.Fatalln("Error in stringParse: ", err)
	}
	
	bstring := make([]byte, s_size)
	for i:=0; i<s_size; i++ {
		bstring[i] = t.next()
	}
	return string(bstring)
}

func listParse(t *tdata) []item {
	var itemSlice []item

	//We read until we reach the end 'e' of the list and make this
	//a list item. we peek so we don't fuck it up for nextItem(*tdata)
	for t.peek() != 'e' {
		it := nextItem(t)
		itemSlice = append(itemSlice, it)
	}
	t.next() //Throw away the 'e'
	
	return itemSlice
}

func dictParse(t *tdata) map[item]item {
	dictMap := make(map[item]item)
	
	//We read until we reach the end 'e' of the dictionary and make this
	//a dictionary item. We peek so we don't fuck it up for nextItem(*tdata).
	//We must be able to read two items at a time, otherwise the torrent is faulty
	//formatted
	for t.peek() != 'e' {
		key := nextItem(t)
		value := nextItem(t)
		dictMap[key] = value
	}
	t.next() //Throw away the 'e'
	
	return dictMap
}

func nextItem(t *tdata) item {
	switch t.next() {
		case 'd':
			return item{itemDict, dictParse(t)}
		case 'l':
			return item{itemList, listParse(t)}
		case 'i':
			return item{itemInt, intParse(t)}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return item{itemString, stringParse(t)}
		default:
			log.Fatalln("Out of bonds in nextItem")
	}
	//Unreachable, but needed due to weird controls in go-compiler
	return item{itemNil, ""}
}

func getItems(t *tdata) []item {
	var itemSlice []item
	
	//Aslong there's been less data read then then data availble, continue fetching items!
	for len(t.data) > t.pos {
		it := nextItem(t)
		itemSlice = append(itemSlice, it)
	}
	return itemSlice
}

func main() {
	var t tdata
	var err error
	
	t.data, err = ioutil.ReadFile("daily.torrent")
	if err != nil {
		log.Fatalln("Error in main(): ", err)
	}
	
	s := getItems(&t)
	fmt.Println(s)
}
