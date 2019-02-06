package main

import (
	"github.com/smugcloud/numberserver/server"
)

func main() {
	s := server.Conf{
		Address: ":9000",
	}
	s.Listen()
}
