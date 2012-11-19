package main

import(
	"fmt"
	"oi/ioutil"
)

func main() {
	buf := ioutil.Readfile("test.torrent")
	
	// byte strings '4:spam' - ascii
	// integers i<NUMBER>e i.e i24e THIS IS THE ENCODING FOR NUMBERS IN STRINGS
	// list	l<bytestring><bytestring>...e i.e l4:spam5:eggsye -> {"spam", "eggsy"}
	// dictionaries (maps) d<bytestringKEY><bytestringVALUE><bytestringKEY><bytestringVALUE>...e
	//	i.e. d5:spamy3:egg2:he3:keye -> {"spamy":"egg", "he":"key}
	
	//Check type
	//Read in type and return it
	//til error occurs. end of buffer.
	