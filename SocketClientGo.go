package main

import (
	"bufio"
	"fmt"
	"github.com/go-stomp/stomp"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	connected     bool
	connectedSync sync.Mutex
	buffData      = make(chan []byte)
)

func main() {
	fmt.Println("Client started...")
	//u := url.URL{
	//	Scheme: "ws",
	//	Host:   "localhost:8091",
	//	Path:   "/nexsocket-notif/325/qktbggje/websocket",
	//}
	for {
		connectedSync.Lock()
		alreadyConnected := connected
		connectedSync.Unlock()
		if !alreadyConnected {
			//conn, err := net.Listen("TCP", u.String())
			//if err != nil {
			//	fmt.Println(err.Error())
			//	time.Sleep(time.Duration(5) * time.Second)
			//	continue
			//}

			//conn, resp, err := websocket.DefaultDialer.Dial("ws://localhost:8091/nexsocket-notif", nil)
			//fmt.Println(resp)
			////url := "wss://socket.gromart.club/nexsocket-notif"
			conn, err := websocket.Dial("ws://localhost:8091/nexsocket-notif", "tcp", "http://localhost:8091")
			//
			//conn, err := net.Dial("tcp", "socket.gromart.club:8443")
			if err != nil {
				log.Fatal("gagal connect")
			}

			var options = []func(*stomp.Conn) error{
				stomp.ConnOpt.AcceptVersion(stomp.V10),
				stomp.ConnOpt.AcceptVersion(stomp.V11),
				stomp.ConnOpt.AcceptVersion(stomp.V12),
				stomp.ConnOpt.Header("destination", "/ping"),
				stomp.ConnOpt.Header("X-Token-Nexsoft", "bmV4Y2hpZWYuZGV2LWM1NzkwMTcwYjVmNzQ2ZGViNGYwYmZjMTUyOGU0MjQ0LTE2NTA4NTI0MjI4MTg"),
				stomp.ConnOpt.Header("heart-beat", "25000,25000"),
				stomp.ConnOpt.Header("Cookie", "SOCKETSERVER="),
				stomp.ConnOpt.Header("content-type", "application/json"),
				stomp.ConnOpt.Header("content-length", "0"),
				stomp.ConnOpt.ReadBufferSize(20060),
			}

			//sc, err := stomp.Dial("tcp", "localhost:8091", options...)
			sc, err := stomp.Connect(conn, options...)
			if err != nil {
				log.Fatal("error nya adalah ==> ", err)
			}

			sub, err := sc.Subscribe("/user/pong", stomp.AckClient)
			if err != nil {
				log.Fatal("error subscribe : ", err)
			}

			err = sc.Send(
				"/nexsocket/ping",            // destination
				"text/plain",               // content-type
				[]byte("ping"))
			if err != nil {
				log.Fatal("error kirim : ", err)
			}

			md := <-sub.C
			if md.Err != nil {
				log.Fatal("receive greeting message caught error: %v", md.Err)
			} else {
				fmt.Printf("----> receive new message: %v\n", string(md.Body))
				log.Fatal("berhasil")
			}









			connectedSync.Lock()
			connected = true
			connectedSync.Unlock()
			//go sendingData(conn)
			//go receivingData(conn)
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func receivingData(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(conn.RemoteAddr().String() + ": disconnected")
			conn.Close()
			connectedSync.Lock()
			connected = false
			connectedSync.Unlock()
			fmt.Println(conn.RemoteAddr().String() + ": end receiving data")
			return
		}
		fmt.Print(conn.RemoteAddr().String() + ": received " + message)
	}
}

func sendingData(conn net.Conn) {
	i := 0
	for {
		_, err := fmt.Fprintf(conn, strconv.Itoa(i)+". data from client\n")
		i++
		if err != nil {
			fmt.Println(conn.RemoteAddr().String() + ": end sending data")
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}
