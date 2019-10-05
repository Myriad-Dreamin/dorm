package dorm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)


func defaultValue(fieldType reflect.Type) string {
	if fieldType == timeType {
		return "date('9999-12-31')"
	}
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Float32, reflect.Float64:
		return "0"
	case  reflect.Bool:
		return "false"
	case reflect.String:
		return "''"
	case reflect.Ptr:
		return defaultValue(fieldType.Elem())
	}
	panic(fmt.Errorf("unknown type of value %v %v", fieldType.Name(), fieldType))
	return ""
}

func (s *common) fullFields() string {
	var fs = make([]string, len(s.typeInfoSlice))
	for i, fieldType := range s.typeInfoSlice {
		fs[i] = "ifnull("+ fieldType.ColumnName +", "+ defaultValue(fieldType.Type) +")"
	}

	return strings.Join(fs, ",")
}

func (s *common) fieldsS(fields []string) string {
	var fs = make([]string, len(s.typeInfoSlice))
	for i, field := range fields {
		if fieldType, ok := s.typeInfo[field]; !ok {
			panic(fmt.Errorf("bad field %v into %v", field, s.addressableType))
		} else {
			// if type not equal ...
			fs[i] = "ifnull("+ fieldType.ColumnName +", "+ defaultValue(fieldType.Type) +")"
		}
	}

	return strings.Join(fs, ",")
}

func (s *common) fullFieldsTemplate() string {
	var fs = make([]string, len(s.typeInfoSlice))
	for i, fieldType := range s.typeInfoSlice {
		fs[i] = fieldType.ColumnName + " = ?"
	}

	return strings.Join(fs, ",")
}

func (s *common) fieldsTemplate(fields []string) string {
	var b = new(bytes.Buffer)
	for i, field := range fields {
		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(field)
		b.WriteString(" = ?")
	}

	return b.String()
}

func (s *common) fieldsColumns(fields []string) string {
	var b = new(bytes.Buffer)
	for i, field := range fields {
		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(field)
	}

	return b.String()
}

func (s *common) fullValues(obj ORMObject) ([]interface{}, error) {

	if reflect.TypeOf(obj) != s.addressableType {
		return nil, errors.New("type error")
	}
	objValue := reflect.ValueOf(obj)
	if objValue.IsNil() {
		return nil, errors.New("nil pointer object")
	}
	objValue = objValue.Elem()
	var argLength = len(s.typeInfo)
	var args = make([]interface{}, argLength)
	for i, fieldType := range s.typeInfoSlice {
		args[i] = objValue.Field(fieldType.FieldOffset).Interface()
	}
	return args, nil
}

func (s *common) fieldsValues(obj ORMObject, fields []string) ([]interface{}, error) {

	if reflect.TypeOf(obj) != s.addressableType {
		return nil, errors.New("type error")
	}
	objValue := reflect.ValueOf(obj)
	if objValue.IsNil() {
		return nil, errors.New("nil pointer object")
	}
	objValue = objValue.Elem()

	var args = make([]interface{}, len(fields))
	for i, field := range fields {
		if fieldType, ok := s.typeInfo[field]; !ok {
			return nil, fmt.Errorf("field %v not find", field)
		} else {
			// if type not equal ...
			args[i] = objValue.Field(fieldType.FieldOffset).Interface()
		}
	}

	return args, nil
}

func (s *common) fieldsTemplateI(obj ORMObject, fields []string) (string, []interface{}, reflect.Value, error) {

	if reflect.TypeOf(obj) != s.addressableType {
		return "", nil, reflect.Value{}, errors.New("type error")
	}
	objValue := reflect.ValueOf(obj)
	if objValue.IsNil() {
		return "", nil, reflect.Value{}, errors.New("nil pointer object")
	}
	objValue = objValue.Elem()

	var b = new(bytes.Buffer)
	var args = make([]interface{}, len(fields))
	for i, field := range fields {
		if fieldType, ok := s.typeInfo[field]; !ok {
			return "", nil, reflect.Value{}, errors.New("not find")
		} else {
			// if type not equal ...
			args[i] = objValue.Field(fieldType.FieldOffset).Interface()
		}

		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(field)
		b.WriteString(" = ?")
	}

	if fieldType, ok := s.typeInfo["id"]; ok {
		return b.String(), args, objValue.Field(fieldType.FieldOffset).Addr(), nil
	} else {
		return b.String(), args, reflect.Value{}, nil
	}
}

func (s *common) fieldsTemplateV(valueObj ORMObject, fields []string) (string, []interface{}, error) {

	if reflect.TypeOf(valueObj) != s.valueType {
		return "", nil, errors.New("type error")
	}
	objValue := reflect.ValueOf(valueObj)

	var b = new(bytes.Buffer)
	var args = make([]interface{}, len(fields))
	for i, field := range fields {
		if fieldType, ok := s.typeInfo[field]; !ok {
			return "", nil, errors.New("not find")
		} else {
			// if type not equal ...
			args[i] = objValue.Field(fieldType.FieldOffset).Interface()
		}

		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(field)
		b.WriteString(" = ?")
	}

	return b.String(), args, nil
}
