package dorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type FetchFetchFunc func(value interface{}) FetchFunc

func (s *ModelScope) fetchArgsFunc(args []int) FetchFetchFunc {
	var retrieves = make([]interface{}, len(args))
	return func(obj interface{}) FetchFunc {
		var err error
		objValue := reflect.ValueOf(obj)
		objType := objValue.Type()
		if objType.Kind() != reflect.Ptr {
			err = errors.New("not ptr input")
		} else if objValue.IsNil() {
			err = errors.New("nil object")
		}


		if objType == s.addressableType {
			objValue = objValue.Elem()
			return func(rows *sql.Rows) error {
				if err != nil {
					return err
				}

				for i, arg := range args {
					retrieves[i] = objValue.Field(int(arg)).Addr().Interface()
				}

				return rows.Scan(retrieves...)
			}
		} else {
			objValue = objValue.Elem()
			sliceType := objType.Elem().Elem()
			if sliceType == s.addressableType {
				return func(rows *sql.Rows) error {
					if err != nil {
						return err
					}
					value := reflect.New(s.valueType)
					valueElem := value.Elem()
					for i, arg := range args {
						retrieves[i] = valueElem.Field(int(arg)).Addr().Interface()
					}
					err = rows.Scan(retrieves...)
					if err != nil {
						return err
					}
					objValue.Set(reflect.Append(objValue, value))
					return nil
				}
			} else if sliceType == s.valueType {
				return func(rows *sql.Rows) error {
					if err != nil {
						return err
					}
					value := reflect.New(s.valueType).Elem()
					for i, arg := range args {
						retrieves[i] = value.Field(int(arg)).Addr().Interface()
					}
					err = rows.Scan(retrieves...)
					if err != nil {
						return err
					}
					objValue.Set(reflect.Append(objValue, value))
					return nil
				}
			} else {
				return func(rows *sql.Rows) error {
					return errors.New("invalid slice type")
				}
			}
		}
	}
}

func (s *ModelScope) fieldsFetch(fields []string) (FetchFetchFunc, error) {

	var args = make([]int, len(fields))
	for i, field := range fields {
		if fieldType, ok := s.typeInfo[field]; !ok {
			return nil, fmt.Errorf("field %v not find", field)
		} else {
			args[i] = int(fieldType.FieldOffset)
		}
	}

	return s.fetchArgsFunc(args), nil
}

func (s *ModelScope) fullFetch() (FetchFetchFunc, error) {

	var args = make([]int, len(s.typeInfoSlice))
	for i, fieldType := range s.typeInfoSlice {
		args[i] = int(fieldType.FieldOffset)
	}

	return s.fetchArgsFunc(args), nil
}
