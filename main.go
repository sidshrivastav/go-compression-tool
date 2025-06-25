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
	compressFlag := flag.Bool("compress", false, "Compress the input file")
	decompressFlag := flag.Bool("decompress", false, "Decompress the input file")
	flag.Parse()

	if *filePtr == "" {
		fmt.Println("Error: Input file not specified. Use -file=<filename>")
		flag.Usage()
		return
	}

	if *compressFlag && *decompressFlag {
		fmt.Println("Error: Please specify only one of -compress or -decompress")
		return
	}

	if *compressFlag {
		err := compressFile(*filePtr)
		if err != nil {
			fmt.Println("Compression failed:", err)
		}
		return
	}

	if *decompressFlag {
		fmt.Println("To be implemented!")
		// err := decompressFile(*filePtr)
		// if err != nil {
		// 	fmt.Println("Decompression failed:", err)
		// }
		return
	}

	fmt.Println("Error: You must specify either -compress or -decompress")
	flag.Usage()
}

func compressFile(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	outputFileName := strings.TrimSuffix(baseName, ext) + ".bin"

	// Step 1–3: Frequency, tree, prefix codes
	frequencies := huffman_coding.GenerateFrequencyMap(string(file))
	huffmanTreeRoot := huffman_coding.BuildHuffmanTree(frequencies)
	prefixCodeTable := huffman_coding.GeneratePrefixCodeTable(huffmanTreeRoot)

	// Step 4: Output file
	encodingFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
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

	// Compress content
	err = huffman_coding.CompressFile(string(file), prefixCodeTable, encodingFile)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Compressed file saved as: %s\n", outputFileName)
	return nil
}
