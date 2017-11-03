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
	var movedBoard string

	var n int
	var err error

	for {
		bytesFromClient := make([]byte, 11)
		n, err = conn.Read(bytesFromClient)

		board = string(bytesFromClient)
		fmt.Printf(" RECEIVED: %q\n", board)

		movedBoard, err = tictactoe.MakeRandomMove(board, squares, SERVERSYMBOL)
		if err != nil {
			n, err = conn.Write([]byte("END"))
			if err != nil {
				return n, fmt.Errorf("playTicTacToe error while writing %v", board)
			}
			break
		}

		n, err = conn.Write([]byte(movedBoard))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe error while writing %v", movedBoard)
		}
		fmt.Printf(" SENT: %q\n", movedBoard)
	}

	return 0, nil
}
