package object_parser

import (
	"reflect"
	"strings"
)

type ObjectParser struct {
	object      interface{}
	objectType  reflect.Type
	objectValue reflect.Value
}

func NewObjectParser(object interface{}) *ObjectParser {
	return &ObjectParser{
		object:      object,
		objectType:  reflect.TypeOf(object),
		objectValue: reflect.ValueOf(object),
	}
}

// Ignore field
// 1. field wasn't granted target tag
// 2. field is ptr & nil
// 3. target tag contains omitempty & field is zero value
func (objectParser *ObjectParser) TagValueMap(targetTag string, castMap ...map[reflect.Type]reflect.Type) map[string]interface{} {
	namedParam := make(map[string]interface{})
	t, v := objectParser.getTypeAndValue()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := objectParser.getValue(v.Field(i))

		if field.Type.Kind() == reflect.Ptr && v.Field(i).IsNil() {
			continue
		}

		tag := field.Tag.Get(targetTag)
		if tag == "" {
			continue
		}

		tagValues := strings.Split(field.Tag.Get(targetTag), ",")
		if objectParser.tagContains(tagValues, "omitempty") && fieldValue.IsZero() {
			continue
		}

		fieldValueType := fieldValue.Type()
		if len(castMap) > 0 && castMap[0][fieldValueType] != nil {
			namedParam[tagValues[0]] = fieldValue.Convert(castMap[0][fieldValueType]).Interface()
			continue
		}

		namedParam[tagValues[0]] = fieldValue.Interface()
	}
	return namedParam
}

func (objectParser *ObjectParser) getTypeAndValue() (reflect.Type, reflect.Value) {
	if objectParser.objectType.Kind() == reflect.Ptr {
		return objectParser.objectType.Elem(), objectParser.objectValue.Elem()
	}
	return objectParser.objectType, objectParser.objectValue
}

func (objectParser *ObjectParser) getValue(value reflect.Value) reflect.Value {
	if value.Type().Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

func (objectParser *ObjectParser) tagContains(tagValues []string, element string) bool {
	for _, tag := range tagValues {
		if tag == element {
			return true
		}
	}
	return false
}