package gohm

import(
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var NoStructError error = errors.New(`model is not a struct`)
var NoIDError error = errors.New(`model does not have an ohm:"id" tagged field`)
var NonExportedAttrError error = errors.New(`can't put ohm tags in unexported fields`)

// If you plan on calling any of the Model helpers available in this package
// make sure you always run ValidateModel on your model, or you run a pretty
// big risk of raising a panic: gohm uses *a lot* of reflection, which is very
// prone to panics when the type received doesn't follow certain assumptions.
func ValidateModel(model interface{}) error {
	var hasID bool
	modelData := reflect.ValueOf(model).Elem()
	modelType := modelData.Type()

	if modelData.Kind().String() != `struct` {
		return NoIDError
	}

	for i := 0; i < modelData.NumField(); i++ {
		if !modelData.Field(i).CanSet() {
			return NonExportedAttrError
		}

		if modelType.Field(i).Tag.Get(`ohm`) == `id` {
			hasID = true
		}
	}

	if !hasID {
		return NoIDError
	}

	return nil
}

func ModelAttrIndexMap(model interface{}) (map[string]int) {
	attrs := map[string]int{}
	typeData := reflect.TypeOf(model).Elem()
	for i := 0; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		tag := field.Tag.Get(`ohm`)
		if tag != `` && tag != `-` && tag != "id" {
			attrs[tag] = i
		}
	}

	return attrs
}

func ModelKey(model interface{}) (key string) {
	modelType := reflect.TypeOf(model).Elem()
	key = fmt.Sprintf("%v:%v", modelType.Name(), ModelID(model))

	return
}

func ModelID(model interface{}) (id string) {
	modelData := reflect.ValueOf(model).Elem()
	idFieldName := ModelIDFieldName(model)
	id = modelData.FieldByName(idFieldName).String()

	return
}

func ModelHasAttribute(model interface{}, attribute string) bool {
	attrIndexMap := ModelAttrIndexMap(model)
	for attr, _ := range attrIndexMap {
		if attribute == attr {
			return true
		}
	}

	return false
}

func ModelIDFieldName(model interface{}) (fieldName string) {
	modelData := reflect.ValueOf(model).Elem()
	modelType := modelData.Type()

	for i := 0; i < modelData.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get(`ohm`)
		if tag == `id` {
			fieldName = field.Name
			break
		}
	}

	return
}

func ModelIndices(model interface{}) []int {
	indices := []int{}

	typeData := reflect.TypeOf(model).Elem()
	for i := 0; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		tag := field.Tag.Get(`ohm`)
		if strings.Contains(tag, `index`) {
			indices = append(indices, i)
		}
	}

	return indices
}

func ModelUniques(model interface{}) []int {
	uniques := []int{}

	typeData := reflect.TypeOf(model).Elem()
	for i := 0; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		tag := field.Tag.Get(`ohm`)
		if strings.Contains(tag, `unique`) {
			uniques = append(uniques, i)
		}
	}

	return uniques
}

func ModelSetID(id string, model interface{}) {
	reflect.ValueOf(model).Elem().FieldByName(ModelIDFieldName(model)).SetString(id)
}
