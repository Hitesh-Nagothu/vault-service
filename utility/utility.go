package utility

import (
	"reflect"
)

func IsStructEmpty(data interface{}) bool {
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue
		}
		return false //not empty
	}

	return true
}
