package Helpers

import (
	"errors"
	"reflect"
)

func ArrayIndexOf(array interface{}, search interface{}) (int, error) {
	if reflect.TypeOf(array).Kind() == reflect.Array || reflect.TypeOf(array).Kind() == reflect.Slice {

		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == search {
				return i, nil
			}
		}

		return -1, nil
	} else {
		return -1, errors.New("ArrayIndexOf needs an array")
	}
}
