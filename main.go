package main

import (
	"flag"
	"fmt"
	"os"

	huffman_coding "go-compression-tool/libs"
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

	// Step 1: Calculate each character frequencies
	frequencies := huffman_coding.GenerateFrequencyMap(string(file))
	
	// Step 2: Generate binary tree based on frequencies
	huffmanTreeRoot := huffman_coding.BuildHuffmanTree(frequencies)
	
	// Step 3: Generate prefix code table from huffman tree
	prefixCodeTable := huffman_coding.GeneratePrefixCodeTable(huffmanTreeRoot)

	for ch, prefix := range prefixCodeTable {
		fmt.Printf("%c %s\n", ch, prefix)
	}
}
