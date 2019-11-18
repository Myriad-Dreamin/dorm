package dorm

import (
	"errors"
	"reflect"
)

type ModelScopeInsert struct {
	*ModelScope
	stmt string
}


func (s *ModelScope) BuildInsert() (t *ModelScopeInsert) {
	t = &ModelScopeInsert{ModelScope: s}
	if s.partialSet {
		t.stmt = "insert into " + s.db.escaper + s.tableName + s.db.escaper + "(" + s.fieldsColumns(s.fields) + ") values " + nPlaceHolder(len(s.typeInfoSlice))
	} else {
		t.stmt = "insert into " + s.db.escaper + s.tableName + s.db.escaper + "(" + s.fullFieldColumns() + ") values " + nPlaceHolder(len(s.typeInfoSlice))
	}

	return t
}

func (s *ModelScopeInsert) Model(obj ORMObject) *ModelScopeInsert {
	if reflect.TypeOf(obj) != s.addressableType {
		s.Error = errors.New("type error")
		return s
	}
	s.model = obj
	return s
}

func (s *ModelScopeInsert) Insert() (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}
	var fvs []interface{}
	if s.partialSet {
		fvs, err = s.fieldsValues(s.model, s.fields)
		if err != nil {
			return
		}
	} else {
		fvs, err = s.fullValues(s.model)
		if err != nil {
			return
		}
	}
	res, err :=s.db.ExecStatementR(s.stmt, fvs...)
	if err != nil {
		return 0, err
	}


	if fieldType, ok := s.typeInfo["id"]; ok {
		idAddr := reflect.ValueOf(s.model).Elem().Field(fieldType.FieldOffset)
		if idAddr.CanSet() {
			id, _ := res.LastInsertId()
			idAddr.Set(reflect.ValueOf(id).Convert(idAddr.Type()))
		}
	}

	return res.RowsAffected()
}

