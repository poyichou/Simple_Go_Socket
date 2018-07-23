package main

import (
	"fmt"
	"net"
	"time"
	"log"
	"strings"
	"sync"
	"io"
	"errors"
)
type Info struct {
	online   bool
	connTime time.Time
}
type Clients struct {
	count int // online client number
	infos map[string]Info
	mux   sync.Mutex
}

var clients = Clients{
	count: 0,
	infos: make(map[string]Info),
}

func offlineClient(conn net.Conn) {
	clients.mux.Lock()
	if info, ok := clients.infos[conn.RemoteAddr().String()]; ok {
		// update to offline
		clients.infos[conn.RemoteAddr().String()] = Info{false, info.connTime}
		clients.count--
	}
	clients.mux.Unlock()
}

func addCient(conn net.Conn) error {
	clients.mux.Lock()
	if clients.count >= 10 {
		defer clients.mux.Unlock()
		return errors.New("too much online clients")
	}
	clients.infos[conn.RemoteAddr().String()] = Info{
		true,
		time.Now(),
	}
	clients.count++
	clients.mux.Unlock()
	return nil
}
// print all clients
func displayClients(conn net.Conn) {
	clients.mux.Lock()
	result := "======== Client history ========\n"
	for addr, info := range clients.infos {
		var indicator string
		if info.online {
			if strings.Compare(addr, conn.RemoteAddr().String()) == 0 {
				indicator = "online <- you\n"
			} else {
				indicator = "online\n"
			}
		} else {
			indicator = "offline\n"
		}
		result += fmt.Sprintf("Client: %s | connection time: %v | " + indicator, addr, info.connTime.Format("2006-01-02 15:04:05"))
	}
	result += fmt.Sprintf("======== Online clients: %d ========\n", clients.count)
	conn.Write([]byte(result))
	clients.mux.Unlock()
}

func serverhandler(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// conn closed
				log.Println(conn.RemoteAddr().String(), " connection closed")
				offlineClient(conn)
			} else {
				// other error
				log.Println(conn.RemoteAddr().String(), " connection error: ", err)
			}
			return
		}
		// if contain "list", print all clients
		if strings.Contains(string(buffer[:n]), "list") {
			displayClients(conn)
		} else {
			conn.Write([]byte(fmt.Sprintf("received: %s\n", string(buffer[:n]))))
		}
	}
}

func main() {
	netListen, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Println(conn.RemoteAddr().String(), " tcp connect successfully")
		addCient(conn)
		go serverhandler(conn)
	}
}
