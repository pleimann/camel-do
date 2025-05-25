package model

//go:generate go tool go-enum -type=Color

type Color int

const (
	Zinc Color = iota
	Red
	Orange
	Amber
	Yellow
	Lime
	Green
	Emerald
	Teal
	Cyan
	Sky
	Violet
	Purple
	Fuchsia
	Pink
	Rose
)
