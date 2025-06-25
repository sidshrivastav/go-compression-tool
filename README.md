# Huffman Compression Tool

A simple CLI tool to compress and decompress text files using the Huffman coding algorithm. Built with Go.

## 🎯 Challenge Context

This project is a solution to the [Huffman Coding Challenge](https://codingchallenges.fyi/challenges/challenge-huffman) from [CodingChallenges.fyi](https://codingchallenges.fyi).

> “Design and implement a program that can compress and decompress text files using Huffman coding. The goal is to reduce file size while ensuring full fidelity on decompression.”

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/go-compression-tool.git
cd go-compression-tool
```

Build the project

```bash
go build -o hufftool main.go
```

## Usage

Compress the file. Generates input.bin

```bash
./hufftool -file=input.txt -compress
```

Decompress the file.

```bin
./hufftool -file=input.bin -decompress
```
