package main

import "sort"

type Counter struct {
	Frequency int
	Character rune
}

type CounterMap map[rune]int

func GenerateFrequency(content string) ([]Counter, CounterMap) {
	var counterMap CounterMap
	counterMap = make(map[rune]int)
	for _, ch := range content {
		counterMap[ch]++
	}

	counter := make([]Counter, 0, len(counterMap))
	for r, count := range counterMap {
		counter = append(counter, Counter{Frequency: count, Character: r})
	}

	sort.Slice(counter, func(i, j int) bool {
		return counter[i].Frequency < counter[j].Frequency
	})

	return counter, counterMap
}
