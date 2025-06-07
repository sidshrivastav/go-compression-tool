package main

import (
	"reflect"
	"testing"
)

func TestGenerateFrequency(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Counter
	}{
		{
			name:  "simple text",
			input: "hello",
			want: []Counter{
				{Frequency: 1, Character: 'h'},
				{Frequency: 1, Character: 'e'},
				{Frequency: 1, Character: 'o'},
				{Frequency: 2, Character: 'l'},
			},
		},
		{
			name:  "empty string",
			input: "",
			want:  []Counter{},
		},
		{
			name:  "with spaces",
			input: "a a b",
			want: []Counter{
				{Frequency: 1, Character: 'b'},
				{Frequency: 2, Character: 'a'},
				{Frequency: 2, Character: ' '},
			},
		},
		{
			name:  "unicode characters",
			input: "世世a",
			want: []Counter{
				{Frequency: 1, Character: 'a'},
				{Frequency: 2, Character: '世'},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateFrequency(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateFrequency(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}
