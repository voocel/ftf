package main

import "fmt"

const (
	reset = iota
	bold
)

const (
	black = iota + 30
	red
	green
	yellow
	blue
	pink
	cyan
	gray

	white   = 97
	unknown = 999
)

var colorMap = map[int]string{
	bold:    "bold",
	black:   "black",
	red:     "red",
	green:   "green",
	yellow:  "yellow",
	blue:    "blue",
	pink:    "pink",
	cyan:    "cyan",
	gray:    "gray",
	white:   "white",
	unknown: "unknown",
}

func SetColor(text string, conf, bg, color int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, color, text, 0x1B)
}

func Bold(s string) string {
	return SetColor(s, 0, 0, bold)
}

func Black(s string) string {
	return SetColor(s, 0, 0, black)
}

func Red(s string) string {
	return SetColor(s, 0, 0, red)
}

func Green(s string) string {
	return SetColor(s, 0, 0, green)
}

func Yellow(s string) string {
	return SetColor(s, 0, 0, yellow)
}

func Blue(s string) string {
	return SetColor(s, 0, 0, blue)
}

func Pink(s string) string {
	return SetColor(s, 0, 0, pink)
}

func Cyan(s string) string {
	return SetColor(s, 0, 0, cyan)
}

func Gray(s string) string {
	return SetColor(s, 0, 0, gray)
}

func White(s string) string {
	return SetColor(s, 0, 0, white)
}

func PrintBold(s string) {
	println(Bold(s))
}

func PrintBoldf(format, s string) {
	PrintBold(fmt.Sprintf(format, s))
}

func PrintBlack(s string) {
	println(Black(s))
}

func PrintBlackf(format, s string) {
	PrintBlack(fmt.Sprintf(format, s))
}

func PrintRed(s string) {
	println(Red(s))
}

func PrintRedf(format, s string) {
	PrintRed(fmt.Sprintf(format, s))
}

func PrintGreen(s string) {
	println(Green(s))
}

func PrintGreenf(format, s string) {
	PrintGreen(fmt.Sprintf(format, s))
}

func PrintYellow(s string) {
	println(Yellow(s))
}

func PrintYellowf(format, s string) {
	PrintYellow(fmt.Sprintf(format, s))
}

func PrintBlue(s string) {
	println(Blue(s))
}

func PrintBluef(format, s string) {
	PrintBlue(fmt.Sprintf(format, s))
}

func PrintPink(s string) {
	println(Pink(s))
}

func PrintPinkf(format, s string) {
	PrintPink(fmt.Sprintf(format, s))
}

func PrintCyan(s string) {
	println(Cyan(s))
}
func PrintCyanf(format, s string) {
	PrintCyan(fmt.Sprintf(format, s))
}

func PrintGray(s string) {
	println(Gray(s))
}

func PrintGrayf(format, s string) {
	PrintGray(fmt.Sprintf(format, s))
}

func PrintWhite(s string) {
	println(White(s))
}

func PrintWhitef(format, s string) {
	PrintWhite(fmt.Sprintf(format, s))
}

func codeReason(code int) string {
	v, ok := colorMap[code]
	if !ok {
		v = colorMap[unknown]
	}
	return v
}

func colorToCode(s string) int {
	for k := range colorMap {
		if colorMap[k] == s {
			return k
		}
	}
	return unknown
}
