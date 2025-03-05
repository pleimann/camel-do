package model

//go:generate ../go-enum --marshal --names --values

/*
ENUM(

	Cat,
	Dog,
	Rabbit,
	Snail,
	Squirrel,
	Turtle,
	Bird,
	Bug,
	Fish,
	Rat,
	Worm,

)
*/
type Icon string
