package main

import (
	"fmt"
	"io/ioutil"
	"unicode"
	"unicode/utf8"
)

type itemType int

const (
	itemInteger = iota
	itemBytestring
	itemList
	itemDictionary
)

type item struct {
	Type itemType
	val  string
}

type tlex struct {
	data  []byte
	pos   int
	items []item
}

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
		switch (*t).items[i].Type {
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
		if i != len(t.items)-1 {
			text = text + typePrefix + ": " + t.items[i].val + "\n"
		} else {
			text = text + typePrefix + ": " + t.items[i].val
		}
	}
	return text
}

func intFunc(t *tlex) {
	buf := make([]byte, 10) //max integer size
	for i := 0; ; i++ {
		r := t.next()
		switch {
		case unicode.IsDigit(r):
			size := utf8.EncodeRune(buf[i:], r)
			if size != 1 {
				panic("Integer is not an INTEGER! WRONGLY ENCODED YOU BASTARD!")
			}
		case r == 'e' && len(buf) != 0:
			var i item
			i.Type = itemInteger
			i.val = string(buf)
			t.items = append(t.items, i)
			goto done
		default:
			panic("SOMETHING WENT HAYWIRE IN INTFUNC!")
		}
	}
done:
}

func bytestringFunc(t *tlex) {
	

func main() {
	var t tlex

	data, err := ioutil.ReadFile("bytestring.torrent")
	if err != nil {
		panic(err)
	}
	t.data = data
	for {
		r := t.next()
		switch r {
		case 'd':
			fmt.Println("WOOP! Dict not implemented")
			goto done
		case 'i':
			intFunc(&t)
		case 'l':
			fmt.Println("Woop! list not implemented")
			goto done
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			fmt.Println("Woop! Bytestring not implemented")
			goto done
		case utf8.RuneError:
			fmt.Println("Now its all done!")
			goto done
		default:
			fmt.Println("We got something extra: ", r)
		}
	}
done:
	fmt.Println(&t)

	/*
		NR:<bytestring> a string
		i---e			a number
		d---e	a map[string]string
		l---e a []string

		announce
		info (1file case)
			length // Length of the whole file. 
			piece length //Length of a piece in Bytes
			pieces // string of multiple 20
			name

		byte strings '4:spam' - ascii
		integers i<NUMBER>e i.e i24e TO BE TRANSLATED TO INTEGERS
		list	l<bytestring><bytestring>...e i.e l4:spam5:eggsye -> {"spam", "eggsy"}
		dictionaries (maps) d<bytestringKEY><bytestringVALUE><bytestringKEY><bytestringVALUE>...e
		i.e. d5:spamy3:egg2:he3:keye -> {"spamy":"egg", "he":"key}

		Check type
		Read in type and return it
		til error occurs. end of buffer.

	*/
}
