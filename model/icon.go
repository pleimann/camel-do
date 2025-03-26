package model

//go:generate go tool go-enum -type=Icon

type Icon int

const (
	Bear Icon = iota
	Bird
	Bug
	Butterfly
	Cat
	Cow
	Crab
	Elephant
	Fish
	Frog
	Hedgehog
	Horse
	Lion
	Narwhal
	Owl
	Pig
	Rabbit
	Shark
	Snail
	Squirrel
	Rat
	Turtle
	Worm
	Whale
)
