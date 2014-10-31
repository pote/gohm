package gohm

import(
	"errors"
	"reflect"
)

var NoStructError error = errors.New(`model is not a struct`)
var NoIDError error = errors.New(`model does not have an ohm:"id" tagged field`)

func ValidateModel(model interface{}) error {
	var hasID bool
	modelData := reflect.ValueOf(model).Elem()
	modelType := modelData.Type()

	if modelData.Kind().String() != `struct` {
		return NoIDError
	}

	for i := 0; i < modelData.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get(`ohm`)
		if !hasID && tag == `id` {
			hasID = true
		}
	}

	if !hasID {
		return NoIDError
	}

	return nil
}
