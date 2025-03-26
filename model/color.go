package model

//go:generate go tool go-enum -type=Color

type Color int

const (
	Red Color = iota
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
	Zinc
)
