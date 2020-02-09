package main

import (
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func TestGetPossibleNodes(t *testing.T) {
	var tests = []struct {
		node     Node
		expected []string
	}{
		{Node{id: 1}, []string{"10000"}},
		{Node{id: 4}, []string{"10000", "10001", "10002", "10003"}},
	}

	for _, tt := range tests {
		actual := tt.node.getPossibleNodes()
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("got %q, want %q", actual, tt.expected)
		}
	}
}

func TestGetColor(t *testing.T) {
	var tests = []struct {
		colorStr string
		expected *color.Color
	}{
		{"red", color.New(color.FgRed)},
		{"blue", color.New(color.FgBlue)},
		{"asfdasdfasdfa not a color", color.New(color.FgWhite)},
	}

	for _, tt := range tests {
		actual := getColor(tt.colorStr)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("got %q, want %q", actual, tt.expected)
		}
	}
}
