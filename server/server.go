package server

import (
	"bufio"
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
	m          map[string]int
	newUnique  uint32
	uniqueTTL  uint32
	lastSize   uint32
	duplicates uint32
	mut        sync.Mutex
}

func (m *numbers) printMap(t *time.Ticker) {
	for {
		select {
		case <-t.C:
			if m.uniqueTTL != uint32(len(m.m)) {
				m.mut.Lock()
				m.newUnique = uint32(len(m.m)) - m.newUnique
				m.uniqueTTL = uint32(len(m.m))
				m.mut.Unlock()
			}
			log.Printf("Received %v unique numbers, %v duplicates.  Unique total: %v\n", m.newUnique, m.duplicates, m.uniqueTTL)

			m.mut.Lock()
			m.newUnique = 0
			m.mut.Unlock()

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
		go cn.countNumbers(conn)
		atomic.AddUint64(&counter, ^uint64(counter-1))

	}
}

func (m *numbers) countNumbers(conn net.Conn) {
	defer conn.Close()

	loops := 0

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		loops++

		l := scanner.Text()

		m.m[l]++
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning %v.\n", err)
	}

	for _, v := range m.m {
		if v > 1 {
			m.duplicates++
		}
	}
}
