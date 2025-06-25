package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	huffman_coding "go-compression-tool/libs"
)

func main() {
	// --- Command Line Argument Parsing ---
	filePtr := flag.String("file", "", "Path to the input file")
	flag.Parse()

	if *filePtr == "" {
		fmt.Println("Error: Input file not specified. Use -file=<filename>")
		flag.Usage()
		return
	}

	file, err := os.ReadFile(*filePtr)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Automatically generate output filename: same name with .bin extension
	baseName := filepath.Base(*filePtr)
	ext := filepath.Ext(baseName)
	outputFileName := strings.TrimSuffix(baseName, ext) + ".bin"

	// Step 1: Calculate character frequencies
	frequencies := huffman_coding.GenerateFrequencyMap(string(file))

	// Step 2: Build Huffman Tree
	huffmanTreeRoot := huffman_coding.BuildHuffmanTree(frequencies)

	// Step 3: Generate prefix code table
	prefixCodeTable := huffman_coding.GeneratePrefixCodeTable(huffmanTreeRoot)

	// Step 4: Prepare output file
	encodingFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return
	}
	defer encodingFile.Close()

	// Write Header
	encodingFile.WriteString("--- Header-Start ---\n")
	for ch, prefix := range prefixCodeTable {
		escapedChar := strings.ReplaceAll(string(ch), "\n", "\\n")
		escapedChar = strings.ReplaceAll(escapedChar, "|", "\\|")
		entry := fmt.Sprintf("%s|%s|%d\n", escapedChar, prefix, frequencies[ch])
		encodingFile.WriteString(entry)
	}
	encodingFile.WriteString("--- Header-End ---\n")

	// Step 5: Encode file content using bit-level Huffman encoding
	err = huffman_coding.CompressFile(string(file), prefixCodeTable, encodingFile)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("✅ Compressed file saved as: %s\n", outputFileName)
}
