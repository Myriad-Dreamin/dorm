package dorm

import (
	"errors"
	"reflect"
)

type ModelScopeDelete struct {
	*ModelScope
	stmt string
}

func (s *ModelScope) BuildDelete() (t *ModelScopeDelete) {
	t = &ModelScopeDelete{ModelScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "delete from " + s.db.escaper + s.tableName + s.db.escaper + s.limitation("id")
	return
}

func (s *ModelScopeDelete) Model(obj ORMObject) *ModelScopeDelete {
	if reflect.TypeOf(obj) != s.addressableType {
		s.Error = errors.New("type error")
		return s
	}
	s.model = obj
	s.id = obj.GetID()
	return s
}

func (s *ModelScopeDelete) ID(id interface{}) *ModelScopeDelete {
	s.args[0] = id
	return s
}

func (s *ModelScopeDelete) Limit(sizeP interface{}) *ModelScopeDelete {
	s.limit = sizeP
	return s
}

func (s *ModelScopeDelete) Offset(offsetP interface{}) *ModelScopeDelete {
	s.offset = offsetP
	return s
}

func (s *ModelScopeDelete) Rebind(offset int64, offsetP interface{}) *ModelScopeDelete {
	s.args[offset] = offsetP
	return s
}

func (s *ModelScopeDelete) Delete(args ...interface{}) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	return s.db.ExecStatement(s.decide(s.stmt), append(args, s.decideArgs(s.args)...)...)
}
