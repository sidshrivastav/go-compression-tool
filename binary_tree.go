package main

type HuffmanNode struct {
	Val   int
	Left  *HuffmanNode
	Right *HuffmanNode
}

func NewNode(val int) *HuffmanNode {
	return &HuffmanNode{Val: val}
}
