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

	decoder.RegisterConverter(Unknown, func(input string) reflect.Value {
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

	decoder.RegisterConverter(time.Time{}, func(input string) reflect.Value {
		time, _ := time.Parse(time.RFC3339, input)

		return reflect.ValueOf(time)
	})

	decoder.RegisterConverter(Project{}, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(nil)
		}

		project := Project{
			ID: input,
		}

		return reflect.ValueOf(project)
	})
}
