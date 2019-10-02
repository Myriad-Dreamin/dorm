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
	timePtrType = reflect.TypeOf(&time.Time{})
)

func isBaseKind(k reflect.Type) bool {
	switch k.Kind() {
	case reflect.UnsafePointer, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func, reflect.Chan, reflect.Array,
		reflect.Uintptr,
		reflect.Complex64, reflect.Complex128, reflect.Invalid:
		return false
	case reflect.Struct:
		return k == timeType
	case reflect.Ptr:
		return isBaseKind(k.Elem())
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
	ColumnName string
}

type TypeInfo = map[string]*SType
type TypeInfoSlice = []*SType

func generateTypeInfo(obj ORMObject) (TypeInfo, TypeInfoSlice, reflect.Type, error) {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	if objType.Kind() != reflect.Ptr {
		return nil, nil, nil, errors.New("not available pointer struct object")
	}
	objType = objType.Elem()
	objValue = objValue.Elem()

	if objType.Kind() != reflect.Struct {
		return nil, nil, nil, errors.New("not available pointer struct object")
	}

	var typeInfo, typeSlice = make(map[string]*SType), make(TypeInfoSlice, 0)

	l := objType.NumField()
	for i := 0; i < l; i++ {
		field := objType.Field(i)
		if tag, ok := field.Tag.Lookup("dorm"); ok {

			if !isBaseType(field.Type) {
				return nil, nil, nil, errors.New("not base kind")
			}
			st := &SType{
				StructField: &field,
				FieldOffset: i,
				ColumnName: tag,
			}
			typeInfo[tag] = st
			typeSlice = append(typeSlice, st)
		}
	}

	return typeInfo, typeSlice, objType, nil
}


type common struct {
	typeInfo        TypeInfo
	typeInfoSlice        TypeInfoSlice
	addressableType reflect.Type
	valueType       reflect.Type
	tableName       string
}

type RCommon struct {
	UPK       string
	VPK       string
	TableName string
}

func NewRCommon(uPK, vPK, tableName string) *RCommon {
	return &RCommon{
		UPK:       vPK,
		VPK:       uPK,
		TableName: tableName,
	}
}

func (c *RCommon) rotate() *RCommon {
	return &RCommon{
		UPK:       c.VPK,
		VPK:       c.UPK,
		TableName: c.TableName,
	}
}

func commonFrom(d ORMObject) (c *common, err error) {
	c = new(common)
	c.typeInfo, c.typeInfoSlice, c.valueType, err = generateTypeInfo(d)
	if err != nil {
		return nil, err
	}
	c.addressableType = reflect.TypeOf(d)
	c.tableName = d.TableName()
	return
}