package utils

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/url"
	"reflect"
)

func SetFormValues(values url.Values, structPtr interface{}) error {
	// Get the pointer of struct
	ptr := reflect.ValueOf(structPtr)

	// Get the value of struct
	value := ptr.Elem()

	slog.Debug("utils.SetFormData", "form data", values)

	// Set value to struct field
	valueType := value.Type()
	for i := range value.NumField() {
		field := value.Field(i)
		jsonTag := valueType.Field(i).Tag.Get("json")

		if !values.Has(jsonTag) {
			continue
		}

		if !field.CanSet() {
			return errors.New("cannot set value to field")
		}

		if field.IsValid() && field.Type() == reflect.TypeOf((*Duration)(nil)).Elem() {
			var duration Duration
			json.Unmarshal([]byte(values.Get(jsonTag)), &duration)
			field.Set(reflect.ValueOf(duration))
			continue
		}

		field.Set(reflect.ValueOf(values.Get(jsonTag)).Convert(field.Type()))
	}

	return nil
}
