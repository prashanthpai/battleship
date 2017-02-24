package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func getLineAsInt(reader *bufio.Reader) int {
	numberStr, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	number, err := strconv.ParseInt(strings.TrimSpace(numberStr), 10, 0)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	return int(number)
}

func createMatrix(reader *bufio.Reader, size int) [][]byte {

	ships := make([][]byte, size)
	for i := range ships {
		ships[i] = make([]byte, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			ships[i][j] = '_'
		}
	}

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	for _, xyStr := range strings.Split(strings.TrimSpace(line), ",") {
		xySubStr := strings.Split(xyStr, ":")
		x, _ := strconv.ParseInt(xySubStr[0], 10, 0)
		y, _ := strconv.ParseInt(xySubStr[1], 10, 0)
		ships[x][y] = 'B'
	}

	return ships
}

func prettyPrintMatrix(matrix [][]byte, size int, writer io.Writer) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Fprintf(writer, "%c ", matrix[i][j])
		}
		fmt.Fprintf(writer, "\n")
	}
}

func printResult(p1Count, p2Count int, writer io.Writer) {
	switch {
	case p1Count > p2Count:
		fmt.Fprintln(writer, "Player 1 wins")
	case p1Count < p2Count:
		fmt.Fprintln(writer, "Player 2 wins")
	case p1Count == p2Count:
		fmt.Fprintln(writer, "It is a draw")
	}
}

type missileTarget struct {
	x, y int
}

func loadNextMissile(reader *bufio.Reader, out chan<- missileTarget) {
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	for _, xyStr := range strings.Split(strings.TrimSpace(line), ":") {
		xySubStr := strings.Split(xyStr, ",")
		x, _ := strconv.ParseInt(xySubStr[0], 10, 0)
		y, _ := strconv.ParseInt(xySubStr[1], 10, 0)
		out <- missileTarget{x: int(x), y: int(y)}
	}

	close(out)
}

func launchMissile(ship [][]byte, target missileTarget) bool {
	switch ship[target.x][target.y] {
	case '_':
		ship[target.x][target.y] = 'O'
	case 'B':
		ship[target.x][target.y] = 'X'
		return true
	}
	return false
}

func main() {
	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	gridSize := getLineAsInt(reader)
	totalShips := getLineAsInt(reader)

	fmt.Printf("Grid Size: %d\n", gridSize)
	fmt.Printf("Total Ships : %d\n", totalShips)

	p1Ships := createMatrix(reader, gridSize)
	fmt.Printf("\nP1 Ship Positions:\n")
	prettyPrintMatrix(p1Ships, gridSize, os.Stdout)

	p2Ships := createMatrix(reader, gridSize)
	fmt.Printf("\nP2 Ship Positions:\n")
	prettyPrintMatrix(p2Ships, gridSize, os.Stdout)

	totalMissiles := getLineAsInt(reader)
	fmt.Printf("\nTotal Missiles: %d\n\n", totalMissiles)

	p1Moves := make(chan missileTarget)
	p2Moves := make(chan missileTarget)

	go loadNextMissile(reader, p1Moves)
	go loadNextMissile(reader, p2Moves)

	var p1Hits, p2Hits int
	for i := 0; i < totalMissiles; i++ {
		m1 := <-p1Moves
		if launchMissile(p1Ships, m1) {
			p2Hits++
		}
		m2 := <-p2Moves
		if launchMissile(p2Ships, m2) {
			p1Hits++
		}
	}

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer outputFile.Close()

	writer := io.MultiWriter(outputFile, os.Stdout)

	fmt.Fprintf(writer, "Player1\n")
	prettyPrintMatrix(p1Ships, gridSize, writer)

	fmt.Fprintf(writer, "\nPlayer2\n")
	prettyPrintMatrix(p2Ships, gridSize, writer)

	fmt.Fprintf(writer, "\nP1:%d\n", p1Hits)
	fmt.Fprintf(writer, "P2:%d\n", p2Hits)
	printResult(p1Hits, p2Hits, writer)
}