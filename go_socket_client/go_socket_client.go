package main

import (
	"fmt"
	"net"
	"log"
	"strings"
	"io"
)
// channels for communication
var (
	rech = make(chan string)//remote channel
	loch = make(chan string)//local channel
)

func receiveData(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// connection closed by server
				log.Println(conn.RemoteAddr().String(), " connection closed by server")
			} else if ! strings.HasSuffix(err.Error(), "use of closed network connection") {
				/*
				conn might be closed by other goroutine, which make Read() return error
				while we try to fmt.Println() the error, we get "use of closed network connection"
				Currently net.errClosing isn't exported, so I have to check it with strings.HasSuffix()
				*/
				// other error
				log.Println(conn.RemoteAddr().String(), "connection error:", err)
			}
			rech<-""
			return
		}
		rech <- string(buffer[:n])
	}
}

func scanf_send() {
	var str string
	for {
		fmt.Scanf("%s\n", &str)
		loch <-str
	}
}

func tcpConnect() (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8000")
	if err != nil {
		log.Println("net.ResolveTCPAddr failed")
		return nil, err
	}
	/*func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
	    If laddr is nil, a local address is automatically chosen.
	    If the IP field of raddr is nil or an unspecified IP address, the local system is assumed.
	*/
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	conn, err := tcpConnect()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("connect success")
	}
	go scanf_send()
	go receiveData(conn)
	var str string
	for quit := 0; quit == 0; {
		select {
		case str = <-rech:
			// receive from socket
			if strings.Compare(str, "") == 0 {
				quit = 1
			} else {
				fmt.Println(str)
			}
		case str = <-loch:
			// receive from scanf
			if strings.Compare(str, "exit") == 0 {
				defer fmt.Println("connectiion closed")
				quit = 1
			} else {
				// send to server
				conn.Write([]byte(str))
			}
		}
	}
	conn.Close()
	close(loch)
	close(rech)
}
