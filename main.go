package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	filePtr := flag.String("file", "", "Path to the input file")
	flag.Parse()
	if *filePtr == "" {
		fmt.Println("Usage: go run main.go -file=<filename>")
		return
	}

	file, err := os.ReadFile(*filePtr)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	// fmt.Println(string(file))
	frequencies := GenerateFrequency(string(file))
	for _, character := range frequencies {
		fmt.Printf("%q: %d\n", character.Character, character.Frequency)
	}
}
