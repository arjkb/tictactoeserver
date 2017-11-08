package main

import (
	"fmt"
	"github.com/arjunkrishnababu96/tictactoe"
	"net"
	"os"
	"strings"
)

func main() {

	address := os.Args[1]
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(" main: ", err)
	}
	defer l.Close()
	fmt.Println("Listening on", l.Addr())

	for {
		conn, err := l.Accept()
		// defer conn.Close()
		if err != nil {
			fmt.Println(" for loop: ", err)
		}
		go playTicTacToe(conn)
	}
}

func playTicTacToe(conn net.Conn) (int, error) {
	defer conn.Close()
	var rBoard string
	var sBoard string = tictactoe.GetEmptyBoard()

	var n int
	var err error

	var clientWon, serverWon bool

InfiniteLoop:
	for {
		bytesFromClient := make([]byte, 11)
		n, err = conn.Read(bytesFromClient)
		if err != nil {
			return n, fmt.Errorf("playTicTacToe() Error reading from client %v", err)
		}

		rBoard = string(bytesFromClient)
		if strings.Contains(rBoard, tictactoe.TIE) {
			fmt.Println("tie")
			break
		}

		if !tictactoe.IsValidBoard(rBoard) {
			return n, fmt.Errorf("playTicTacToe() client sent invalid board %v", rBoard)
		}
		fmt.Printf(" R: %q\n", rBoard)

		movCnt, _ := tictactoe.GetMoveDifference(sBoard, rBoard)
		if movCnt != 1 {
			return n, fmt.Errorf("playTicTacToe() client made %d moves", movCnt)
		}

		if tictactoe.HasWon(rBoard, tictactoe.CLIENTSYMBOL) {
			// check if the opponent has won
			sBoard = tictactoe.CLIENTWON
			clientWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rBoard, tictactoe.SERVERSYMBOL); win {
			// check if we can win in one move; make that move
			sBoard, _ = tictactoe.MakeWinMove(rBoard, ptrn, tictactoe.SERVERSYMBOL)
			serverWon = true
		} else if win, ptrn := tictactoe.CanWinNext(rBoard, tictactoe.CLIENTSYMBOL); win {
			// check if opponent can win in one move; block that move
			sBoard, _ = tictactoe.BlockWinMove(rBoard, ptrn, tictactoe.SERVERSYMBOL)
		} else if tictactoe.IsFree(rBoard, 5) {
			// can play center
			sBoard, _ = tictactoe.MakeMove(rBoard, 5, tictactoe.SERVERSYMBOL)
			fmt.Println("playing center! %v %v", rBoard, sBoard)

			// DOWN: Play opposite corner
		} else if rBoard[0] == tictactoe.CLIENTSYMBOL && tictactoe.IsFree(rBoard, 10) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 10, tictactoe.SERVERSYMBOL)
		} else if rBoard[2] == tictactoe.CLIENTSYMBOL && tictactoe.IsFree(rBoard, 8) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 8, tictactoe.SERVERSYMBOL)
		} else if rBoard[8] == tictactoe.CLIENTSYMBOL && tictactoe.IsFree(rBoard, 2) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 2, tictactoe.SERVERSYMBOL)
		} else if rBoard[10] == tictactoe.CLIENTSYMBOL && tictactoe.IsFree(rBoard, 0) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 0, tictactoe.SERVERSYMBOL)

			// DOWN: Play empty corner
		} else if tictactoe.IsFree(rBoard, 0) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 0, tictactoe.SERVERSYMBOL)
		} else if tictactoe.IsFree(rBoard, 2) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 2, tictactoe.SERVERSYMBOL)
		} else if tictactoe.IsFree(rBoard, 8) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 8, tictactoe.SERVERSYMBOL)
		} else if tictactoe.IsFree(rBoard, 10) {
			sBoard, _ = tictactoe.MakeMove(rBoard, 10, tictactoe.SERVERSYMBOL)

		} else {
			// make a random move
			sBoard, err = tictactoe.MakeRandomMove(rBoard, tictactoe.AllSquares, tictactoe.SERVERSYMBOL)
			if err != nil {
				// error indicates there are no more free positions
				// happens only when there is a tie
				sBoard = tictactoe.TIE
			}
		}

		n, err = conn.Write([]byte(sBoard))
		if err != nil {
			return n, fmt.Errorf("playTicTacToe error while writing %v", sBoard)
		}

		switch {
		case sBoard == tictactoe.TIE:
			fmt.Println(tictactoe.TIE)
			break InfiniteLoop
		case serverWon:
			fmt.Println(tictactoe.SERVERWON)
			break InfiniteLoop
		case clientWon:
			fmt.Println(tictactoe.CLIENTWON)
			break InfiniteLoop
		default:
			fmt.Printf(" S: %q\n", sBoard)
		}
	}

	return 0, nil
}
