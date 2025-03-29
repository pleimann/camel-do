package utils

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/pleimann/camel-do/model"
)

var decoder *schema.Decoder

func Decoder() *schema.Decoder {
	return decoder
}

func init() {
	decoder = schema.NewDecoder()

	decoder.RegisterConverter(model.Zinc, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(model.Zinc)
		}

		color, _ := model.ParseColorString(input)

		return reflect.ValueOf(color)
	})

	decoder.RegisterConverter(model.Unknown, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf("")
		}

		color, _ := model.ParseIconString(input)

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

	decoder.RegisterConverter(model.Project{}, func(input string) reflect.Value {
		if input == "" {
			return reflect.ValueOf(model.Project{})
		}

		uuid, _ := uuid.Parse(input)

		project := model.Project{
			ID: uuid,
		}

		return reflect.ValueOf(project)
	})
}
