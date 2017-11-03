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
	const CLIENTSYMBOL = 'X'
	squares := []int{0, 1, 2, 4, 5, 6, 8, 9, 10}
	var rBoard string
	var sBoard string = tictactoe.GetEmptyBoard()

	var n int
	var err error

	var moved bool

	for {
		moved = false
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

		movCnt, _ := tictactoe.GetMoveDifference(sBoard, rBoard)
		if movCnt != 1	{
			return n, fmt.Errorf("playTicTacToe() client made %d moves", movCnt)
		}

		if tictactoe.HasWon(rBoard, CLIENTSYMBOL)	{
			// check is the opponent has won
			fmt.Println("Client won!")
			n, err = conn.Write([]byte("END"))
			if err != nil {
				return n, fmt.Errorf("playTicTacToe error while writing %v", sBoard)
			}
			break
		}

		// check if client can win in the next move
		var patternArray [3]int
		for _, pattern := range(tictactoe.WinPatterns)	{
			copy(patternArray[:], pattern)
			winnable, winMove, _ := tictactoe.IsWinnable(rBoard, CLIENTSYMBOL, patternArray)
			if winnable	{
				fmt.Println("winnable! ", winMove, rBoard)
			}
		}

		if !moved	{
			sBoard, err = tictactoe.MakeRandomMove(rBoard, squares, SERVERSYMBOL)
			moved = false
		}

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
