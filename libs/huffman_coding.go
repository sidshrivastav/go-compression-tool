package huffman_coding

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"strings"
)

type HuffmanNode struct {
	Char        rune
	Freq        int
	Left, Right *HuffmanNode
}

type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*HuffmanNode))
}
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[:n-1]
	return node
}

// GenerateFrequencyMap counts the frequency of each character.
func GenerateFrequencyMap(input string) map[rune]int {
	freqMap := make(map[rune]int)
	for _, ch := range input {
		freqMap[ch]++
	}
	return freqMap
}

// BuildHuffmanTree builds the Huffman tree and returns the root.
func BuildHuffmanTree(freqMap map[rune]int) *HuffmanNode {
	pq := make(PriorityQueue, 0)
	for ch, freq := range freqMap {
		pq = append(pq, &HuffmanNode{Char: ch, Freq: freq})
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)

		merged := &HuffmanNode{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, merged)
	}

	return heap.Pop(&pq).(*HuffmanNode)
}

// GeneratePrefixCodeTable generate prefix code table from huffman binary tree
func GeneratePrefixCodeTable(root *HuffmanNode) map[rune]string {
	table := make(map[rune]string)

	var helper func(node *HuffmanNode, slate []string)
	helper = func(node *HuffmanNode, slate []string) {
		if node.Left == nil && node.Right == nil {
			table[node.Char] = strings.Join(slate, "")
			return
		}

		if node.Left != nil {
			slate = append(slate, "0")
			helper(node.Left, slate)
			slate = slate[:len(slate)-1]
		}

		if node.Right != nil {
			slate = append(slate, "1")
			helper(node.Right, slate)
			slate = slate[:len(slate)-1]
		}
	}

	if root != nil {
		helper(root, []string{})
	}

	return table
}

func CompressFile(content string, prefixCodeTable map[rune]string, encodingFile *os.File) error {
	var bitBufferBuilder strings.Builder

	// Convert each rune to its prefix code
	for _, ch := range content {
		code, ok := prefixCodeTable[ch]
		if !ok {
			return fmt.Errorf("character %q not found in prefixCodeTable", ch)
		}
		bitBufferBuilder.WriteString(code)
	}

	bitBuffer := bitBufferBuilder.String()

	// Padding to make it byte-aligned
	padding := (8 - len(bitBuffer)%8) % 8
	bitBuffer += strings.Repeat("0", padding)

	// Write padding info
	_, err := encodingFile.WriteString(fmt.Sprintf("Padding:%d\n", padding))
	if err != nil {
		return err
	}
	_, err = encodingFile.WriteString("--- Encoding-Start ---\n")
	if err != nil {
		return err
	}

	// Convert bit string to raw bytes
	var byteBuffer bytes.Buffer
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
	_, err = encodingFile.Write(byteBuffer.Bytes())
	if err != nil {
		return err
	}
	_, err = encodingFile.WriteString("\n--- Encoding-End ---\n")
	if err != nil {
		return err
	}

	return nil
}
