package main

import (
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func (node *Node) getPossibleNodes() []string {
	var ret = make([]string, node.id)
	for i := 0; i < node.id; i++ {
		ret[i] = strconv.Itoa(SERVER_BASE_PORT + i)
	}
	return ret
}

func getColor(colorStr string) *color.Color {
	colorStr = strings.ToLower(colorStr)
	if colorStr == "red" {
		return color.New(color.FgRed)
	} else if colorStr == "green" {
		return color.New(color.FgGreen)
	} else if colorStr == "blue" {
		return color.New(color.FgBlue)
	} else {
		return color.New(color.FgWhite)
	}
}
