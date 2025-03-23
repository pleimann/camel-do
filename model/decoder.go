package model

import (
	"reflect"
	"time"

	"github.com/gorilla/schema"
)

var decoder *schema.Decoder

func Decoder() *schema.Decoder {
	return decoder
}

func init() {
	decoder = schema.NewDecoder()

	decoder.RegisterConverter(ColorZinc, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(ColorZinc)
		}

		color, _ := ParseColor(input)

		return reflect.ValueOf(color)
	})

	decoder.RegisterConverter(IconBear, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf("")
		}

		color, _ := ParseIcon(input)

		return reflect.ValueOf(color)
	})

	decoder.RegisterConverter(time.Second, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(nil)
		}

		duration, _ := time.ParseDuration(input)

		return reflect.ValueOf(duration)
	})
}
