package main

import (
	"io/ioutil"
	"log"
	"net"
)

func main() {
	f, _ := ioutil.ReadFile("../client2.txt")
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		log.Printf("Unable to establish connection.")
	}
	log.Printf("Size of file: %v\n", len(f))

	conn.Write(f)

}
