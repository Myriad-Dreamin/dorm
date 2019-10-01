package dorm

import (
	"bytes"
	"errors"
	"reflect"
)

func (t *Table) fieldsTemplate(obj ORMObject, fields []string) (string, []interface{}, error) {
	objValue := reflect.ValueOf(obj)

	var b = new(bytes.Buffer)
	var args = make([]interface{}, 0, len(fields))
	for i, field := range fields {
		if fieldType, ok := t.typeInfo[field]; !ok {
			return "", nil, errors.New("not find")
		} else {
			// if type not equal ...
			args = append(args, objValue.Field(fieldType.FieldOffset).Interface())
		}

		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(field)
		b.WriteString(" = ?")
	}

	return b.String(), append(args, obj.GetID()), nil
}

func (t *Table) generateUpdateSQL(d ORMObject, fields []string) (string, []interface{}, error) {
	template, params, err := t.fieldsTemplate(d, fields)
	if err != nil {
		return "", nil, err
	}

	return "UPDATE " + d.TableName() + " set " + template + " where id = ?", params, nil
}
