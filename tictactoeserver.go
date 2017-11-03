package main

import (
	"fmt"
	"net"
	"log"
	"github.com/arjunkrishnababu96/tictactoe"
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

	n, err := playTicTacToe(conn)
	if err != nil	{
		log.Fatalf(" main() n=%v: %v", n, err)
	}

	// handleConnection(conn)
}

func playTicTacToe(conn net.Conn) (int, error){
	const SERVERSYMBOL = 'O'
	squares := []int{0,1,2,4,5,6,8,9,10}
	var board string

	var n int
	var err error

	for {
		bytesFromClient := make([]byte, 11)
		n, err = conn.Read(bytesFromClient)

		board = string(bytesFromClient)
		fmt.Printf(" RECEIVED: %q\n", board)

		board, err = tictactoe.MakeRandomMove(board, squares, SERVERSYMBOL)
		if err != nil	{
			n, err = conn.Write([]byte("END"))
			if err != nil	{
				return n, fmt.Errorf("playTicTacToe error while writing %v", board)
			}
		}

		n, err = conn.Write([]byte(board))
		if err != nil	{
			return n, fmt.Errorf("playTicTacToe error while writing %v", board)
		}
		fmt.Printf(" SENT: %q\n", board)
	}
}

func handleConnection(conn net.Conn) {
	msg, n, err := receiveMsg(conn)
	if err != nil {
		fmt.Println(" handleConnection: ", err)
	}
	fmt.Printf("%v bytes read\n", n)
	fmt.Println("Client says:")
	fmt.Println(msg)

	n, err = sendMsg(conn, "END")
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
