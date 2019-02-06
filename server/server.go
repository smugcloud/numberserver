package server

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Conf struct {
	Address        string
	maxConnections uint64
	mut            sync.Mutex
}

type numbers struct {
	m map[string]int
}

func (m *numbers) printMap(t *time.Ticker) {
	for {
		select {
		case <-t.C:
			log.Printf("Current map length: %+v\n", len(m.m))
			log.Printf("Current map: %+v\n", m.m)

		}
	}
}

func (c *Conf) Listen() {
	cn := numbers{
		m: make(map[string]int),
	}
	var counter uint64
	c.maxConnections = 1
	t := time.NewTicker(time.Second * 5)
	go cn.printMap(t)
	l, err := net.Listen("tcp", c.Address)
	if err != nil {
		log.Println("Unable to establish a listener.")
	}
	log.Printf("Listening on localhost%v", c.Address)
	for {
		conn, err := l.Accept()
		atomic.AddUint64(&counter, 1)
		log.Printf("Counter: %v\n", counter)
		if counter > c.maxConnections {
			c.mut.Lock()

			log.Println("Too many connections.")
			conn.Write([]byte(`Too many connections.`))
			atomic.AddUint64(&counter, ^uint64(counter-1))
			c.mut.Unlock()
			conn.Close()

			continue
		}
		if err != nil {
			log.Printf("Unable to accept connection.")
		}
		cn.countNumbers(conn)
		atomic.AddUint64(&counter, ^uint64(counter-1))

	}
}

func (m *numbers) countNumbers(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 16*1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Error reading from the connection: %v\n", err)
	}
	r := bytes.NewReader(buf)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		if _, ok := m.m[scanner.Text()]; !ok {
			m.m[scanner.Text()]++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning %v.\n", err)
	}

}
