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
	const CLIENTWON = "client won"
	const SERVERWON = "server won"
	const SERVERSYMBOL = 'O'
	const CLIENTSYMBOL = 'X'

	squares := []int{0, 1, 2, 4, 5, 6, 8, 9, 10}

	var rBoard string
	var sBoard string = tictactoe.GetEmptyBoard()

	var n int
	var err error

	var clientWon, serverWon bool

	for {
		bytesFromClient := make([]byte, 11)
		n, err = conn.Read(bytesFromClient)
		if err != nil {
			return n, fmt.Errorf("playTicTacToe() Error reading from client %v", err)
		}

		rBoard = string(bytesFromClient)
		if !tictactoe.IsValidBoard(rBoard) {
			return n, fmt.Errorf("playTicTacToe() client sent invalid board %v", rBoard)
		}
		fmt.Printf(" R: %q\n", rBoard)

		movCnt, _ := tictactoe.GetMoveDifference(sBoard, rBoard)
		if movCnt != 1 {
			return n, fmt.Errorf("playTicTacToe() client made %d moves", movCnt)
		}

		if tictactoe.HasWon(rBoard, CLIENTSYMBOL) {
			// check if the opponent has won
			sBoard = CLIENTWON
			clientWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rBoard, SERVERSYMBOL); win {
			// check if we can win in one move; make that move
			sBoard, _ = tictactoe.MakeWinMove(rBoard, ptrn, SERVERSYMBOL)
			serverWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rBoard, CLIENTSYMBOL); win {
			// check if opponent can win in one move; block that move
			sBoard, _ = tictactoe.BlockWinMove(rBoard, ptrn, SERVERSYMBOL)
		} else {
			// make a random move
			sBoard, err = tictactoe.MakeRandomMove(rBoard, squares, SERVERSYMBOL)
			if err != nil {
				// error indicates there are no more free positions; send END signal
				sBoard = "END"
			}
		}

		n, err = conn.Write([]byte(sBoard))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe error while writing %v", sBoard)
		}
		fmt.Printf(" S: %q\n", sBoard)

		if sBoard == "END" || serverWon || clientWon {
			break
		}
	}

	return 0, nil
}
