package dorm

import (
	"errors"
	"reflect"
)

type ModelScopeUpdateFields struct {
	*ModelScope
	stmt string
}

func (s *ModelScope) BuildUpdateFields() (t *ModelScopeUpdateFields) {
	t = &ModelScopeUpdateFields{ModelScope: s}

	if s.partialSet {
		t.stmt = "update " + s.db.escaper + s.tableName + s.db.escaper + " set " + s.fieldsTemplate(s.fields) + s.limitation("id")
	} else {
		t.stmt = "update " + s.db.escaper + s.tableName + s.db.escaper + " set " + s.fullFieldsTemplate() + s.limitation("id")
	}

	return t
}

func (s *ModelScopeUpdateFields) Model(obj ORMObject) *ModelScopeUpdateFields {
	if reflect.TypeOf(obj) != s.addressableType {
		s.Error = errors.New("type error")
		return s
	}
	s.model = obj
	s.args[0] = obj.GetID()
	return s
}

func (s *ModelScopeUpdateFields) ID(id interface{}) *ModelScopeUpdateFields {
	s.args[0] = id
	return s
}

func (s *ModelScopeUpdateFields) Limit(sizeP interface{}) *ModelScopeUpdateFields {
	s.limit = sizeP
	return s
}

func (s *ModelScopeUpdateFields) Offset(offsetP interface{}) *ModelScopeUpdateFields {
	s.offset = offsetP
	return s
}

func (s *ModelScopeUpdateFields) Rebind(offset int64, offsetP interface{}) *ModelScopeUpdateFields {
	s.args[offset] = offsetP
	return s
}

func (s *ModelScopeUpdateFields) UpdateFields(args ...interface{}) (aff int64, err error) {
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

	return s.db.ExecStatement(s.decide(s.stmt), append(append(fvs, args...), s.decideArgs(s.args)...)...)
}
