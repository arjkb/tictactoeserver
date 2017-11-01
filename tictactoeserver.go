package main

import (
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:7776")
	if err != nil {
		fmt.Println(" main: ", err)
	}
	defer l.Close()
	fmt.Println("Listening on", l.Addr())

	conn, err := l.Accept()
	defer conn.Close()
	if err != nil {
		fmt.Println(" for loop: ", err)
	}

	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	msg, n, err := receiveMsg(conn)
	if err != nil	{
		fmt.Println(" handleConnection: ", err)
	}
	fmt.Printf("%v bytes read\n", n)
	fmt.Println("Client says:")
	fmt.Println(msg)

	n, err = sendMsg(conn, "Greetings from server!")
	conn.Close()
}

func receiveMsg(conn net.Conn) (msg string, n int, err error) {
	var msg_bytes = make([]byte, 100)
	n, err = conn.Read(msg_bytes)
	if err != nil {
		return msg, n, fmt.Errorf(" receiveMsg: %v", err)
	}

	return string(msg_bytes), n, err
}

func sendMsg(conn net.Conn, msg string) (n int, err error) {
	n, err = conn.Write([]byte(msg))
	if err != nil {
		return n, fmt.Errorf(" error while writing: ", err)
	}
	return n, nil
}
