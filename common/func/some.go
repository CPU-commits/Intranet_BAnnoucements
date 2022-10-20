package common

import (
	"fmt"
	"reflect"
)

func Some(validorFunc func(x interface{}) bool, slice interface{}) (bool, error) {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return false, fmt.Errorf("the second argument must be a slice")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		v := sliceValue.Index(i).Interface()

		if !validorFunc(v) {
			return false, nil
		}
	}

	return true, nil
}
