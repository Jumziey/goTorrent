package main

import (
)


func firstNil(b []byte) int {
	var i int
	for i = len(b) - 1; ; i-- {
		if b[i] != 0 {
			break
		}
	}
	return i + 1
}

func incByteSlice(b []byte, n int) []byte {
	buf := make([]byte, len(b)+n);
	
	for i:=0; i<len(b); i++ {
		buf[i] = b[i]
	}
	return buf
}

func stringSlicetoString(s []string) string {
	var str string
	for i:=0; i<len(s); i++ {
		str = str+" "
		str = str+s[i]
	}
	return str
}