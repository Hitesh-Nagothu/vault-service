package utility

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func IntersectionOfIds(existing []primitive.ObjectID, new []primitive.ObjectID) []primitive.ObjectID {
	lookup := make(map[primitive.ObjectID]bool)
	for _, id := range existing {
		lookup[id] = true
	}

	intersection := []primitive.ObjectID{}
	for _, id := range new {
		if lookup[id] {
			intersection = append(intersection, id)
		}
	}

	return intersection
}
