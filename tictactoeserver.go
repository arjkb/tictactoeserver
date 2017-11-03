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
	var rBoard, sBoard string

	var n int
	var err error

	for {
		bytesFromClient := make([]byte, 11)
		n, err = conn.Read(bytesFromClient)
		if err != nil	{
			return n, fmt.Errorf("playTicTacToe() Error reading from client %v", err)
		}

		rBoard = string(bytesFromClient)
		if !tictactoe.IsValidBoard(rBoard)	{
			return n, fmt.Errorf("playTicTacToe() client sent invalid board %v", rBoard)
		}
		fmt.Printf(" R: %q\n", rBoard)

		sBoard, err = tictactoe.MakeRandomMove(rBoard, squares, SERVERSYMBOL)
		if err != nil {
			n, err = conn.Write([]byte("END"))
			if err != nil {
				return n, fmt.Errorf("playTicTacToe error while writing %v", sBoard)
			}
			break
		}


		n, err = conn.Write([]byte(sBoard))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe error while writing %v", sBoard)
		}
		fmt.Printf(" S: %q\n", sBoard)
	}

	return 0, nil
}
