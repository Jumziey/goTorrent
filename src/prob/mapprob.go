package main

import (
	"fmt"
)

type item struct {
	val interface{}
}

func main() {

	m := make(map[item]item)
	key1 := item{"heja"}
	value1 := item{"bleh"}
	m[key1] = value1
	fmt.Println(m)
	
	n := make(map[item]item)
	key2 := item{map[item]item{key1:value1}}
	value2 := item{"mehe"}
	n[key2] = value2
	
	fmt.Println(n)
}
	