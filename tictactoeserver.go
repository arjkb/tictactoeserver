package main

import (
	"fmt"
	"github.com/arjunkrishnababu96/tictactoe"
	"log"
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

	n, err := playTicTacToe(conn)
	if err != nil {
		log.Fatalf(" main() n=%v: %v", n, err)
	}
}

func playTicTacToe(conn net.Conn) (int, error) {
	const SERVERSYMBOL = 'O'
	squares := []int{0, 1, 2, 4, 5, 6, 8, 9, 10}
	var board string

	var n int
	var err error

	for {
		bytesFromClient := make([]byte, 11)
		n, err = conn.Read(bytesFromClient)
		if err != nil	{
			return n, fmt.Errorf("playTicTacToe() Error reading from client %v", err)
		}

		board = string(bytesFromClient)
		fmt.Printf(" R: %q\n", board)

		board, err = tictactoe.MakeRandomMove(board, squares, SERVERSYMBOL)
		if err != nil {
			n, err = conn.Write([]byte("END"))
			if err != nil {
				return n, fmt.Errorf("playTicTacToe error while writing %v", board)
			}
			break
		}

		n, err = conn.Write([]byte(board))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe error while writing %v", board)
		}
		fmt.Printf(" S: %q\n", board)
	}

	return 0, nil
}
