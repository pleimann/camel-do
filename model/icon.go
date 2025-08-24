package model

//go:generate go tool go-enum -type=Icon

type Icon int

const (
	Unknown Icon = iota
	Bear
	Bee
	Bird
	Bug
	Butterfly
	Cat
	Crab
	Cow
	Dog
	Elephant
	Fish
	Frog
	Hedgehog
	Horse
	Lion
	Narwhal
	Owl
	Panda
	Pig
	Rabbit
	Rat
	Snail
	Squirrel
	Turtle
	Worm
	Shark
	Spider
	Whale
)
