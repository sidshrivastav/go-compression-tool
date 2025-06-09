package main

import (
	"flag"
	"fmt"
	"os"

	huffman_coding "go-compression-tool/libs"
)

func main() {
	// --- Command Line Argument Parsing ---
	filePtr := flag.String("file", "", "Path to the input file")
	encodingFilePtr := flag.String("encoding", "", "Name for the output compressed file")

	// Parse flags ONLY ONCE
	flag.Parse()

	if *filePtr == "" {
		fmt.Println("Error: Input file not specified. Use -file=<filename>")
		flag.Usage() // Print usage instructions
		return
	}
	if *encodingFilePtr == "" {
		fmt.Println("Error: Output encoding file not specified. Use -encoding=<filename>")
		flag.Usage() // Print usage instructions
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

	// Step 4: Generate file and prepare headers
	encodingFile, err := os.OpenFile(*encodingFilePtr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file for append:", err)
		return
	}
	defer encodingFile.Close()
	encodingFile.WriteString("--- Header-Start ---\n")
	for ch, prefix := range prefixCodeTable {
		entry := fmt.Sprintf("%c|%s|%d\n", ch, prefix, frequencies[ch])
		encodingFile.WriteString(entry)
	}
	encodingFile.WriteString("--- Header-End ---\n")

}
