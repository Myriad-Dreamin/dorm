package dorm

import (
	"errors"
	"reflect"
	"time"
)

type ORMObject interface {
	TableName() string
	GetID() uint
}

var (
	timeType = reflect.TypeOf(time.Time{})
)

func isBaseKind(k reflect.Type) bool {
	switch k.Kind() {
	case reflect.UnsafePointer, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Chan, reflect.Array,
		reflect.Uintptr, reflect.Ptr,
		reflect.Complex64, reflect.Complex128, reflect.Invalid:
		return false
	case reflect.Struct:
		return k == timeType
	default:
		return true
	}
}

func isBaseType(fieldType reflect.Type) bool {
	return (fieldType.Kind() == reflect.Ptr && isBaseKind(fieldType.Elem())) || isBaseKind(fieldType)
}

type SType struct {
	*reflect.StructField
	FieldOffset int
}

type TypeInfo = map[string]*SType

func generateTypeInfo(obj ORMObject) (TypeInfo, error) {
	var typeInfo = make(map[string]*SType)
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}

	if objType.Kind() != reflect.Struct {
		return nil, errors.New("not available struct object")
	}

	for i := objType.NumField() - 1; i >= 0; i-- {
		field := objType.Field(i)
		if tag, ok := field.Tag.Lookup("dorm"); ok {

			if !isBaseType(field.Type) {
				return nil, errors.New("not base kind")
			}

			typeInfo[tag] = &SType{
				StructField: &field,
				FieldOffset: i,
			}
		}
	}

	return typeInfo, nil
}
