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

	decoder.RegisterConverter(Zinc, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(Zinc)
		}

		color, _ := ParseColorString(input)

		return reflect.ValueOf(color)
	})

	decoder.RegisterConverter(Bear, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf("")
		}

		color, _ := ParseIconString(input)

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
