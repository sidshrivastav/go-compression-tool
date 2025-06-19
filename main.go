package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	huffman_coding "go-compression-tool/libs"
)

func main() {
	// --- Command Line Argument Parsing ---
	filePtr := flag.String("file", "", "Path to the input file")
	encodingFilePtr := flag.String("encoding", "", "Name for the output compressed file")

	flag.Parse()

	if *filePtr == "" {
		fmt.Println("Error: Input file not specified. Use -file=<filename>")
		flag.Usage()
		return
	}
	if *encodingFilePtr == "" {
		fmt.Println("Error: Output encoding file not specified. Use -encoding=<filename>")
		flag.Usage()
		return
	}

	file, err := os.ReadFile(*filePtr)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Step 1: Calculate character frequencies
	frequencies := huffman_coding.GenerateFrequencyMap(string(file))

	// Step 2: Build Huffman Tree
	huffmanTreeRoot := huffman_coding.BuildHuffmanTree(frequencies)

	// Step 3: Generate prefix code table
	prefixCodeTable := huffman_coding.GeneratePrefixCodeTable(huffmanTreeRoot)

	// Step 4: Prepare output file
	encodingFile, err := os.OpenFile(*encodingFilePtr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
	var bitBufferBuilder strings.Builder
	for _, ch := range string(file) {
		code := prefixCodeTable[ch]
		bitBufferBuilder.WriteString(code)
	}

	bitBuffer := bitBufferBuilder.String()

	// Padding to make it byte-aligned
	padding := (8 - len(bitBuffer)%8) % 8
	bitBuffer += strings.Repeat("0", padding)

	// Write padding info before encoding
	encodingFile.WriteString(fmt.Sprintf("Padding:%d\n", padding))
	encodingFile.WriteString("--- Encoding-Start ---\n")

	// Convert bit string to raw bytes
	byteBuffer := bytes.Buffer{}
	for i := 0; i < len(bitBuffer); i += 8 {
		byteChunk := bitBuffer[i : i+8]
		var b byte
		for j := 0; j < 8; j++ {
			if byteChunk[j] == '1' {
				b |= 1 << (7 - j)
			}
		}
		byteBuffer.WriteByte(b)
	}

	// Write encoded binary bytes
	encodingFile.Write(byteBuffer.Bytes())
	encodingFile.WriteString("\n--- Encoding-End ---\n")

	fmt.Printf("✅ Compressed file saved as: %s\n", *encodingFilePtr)
}
