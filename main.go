package main

import (
	"github.com/smugcloud/numberserver/server"
)

func main() {
	s := server.Conf{
		Address: ":4000",
	}
	s.Listen()
}
