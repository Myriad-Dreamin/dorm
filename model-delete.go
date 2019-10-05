package dorm

import (
	"errors"
	"reflect"
)

type ModelScopeDelete struct {
	*ModelScope
	stmt string
	args []interface{}
}


func (s *ModelScope) BuildDelete() (t *ModelScopeDelete) {
	t = &ModelScopeDelete{ModelScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "delete from " + s.tableName + s.limitation("id", &t.args)
	return
}

func (s *ModelScopeDelete) Model(obj ORMObject) *ModelScopeDelete {
	if reflect.TypeOf(obj) != s.addressableType {
		s.Error = errors.New("type error")
		return s
	}
	s.model = obj
	s.args[0] = obj.GetID()
	return s
}

func (s *ModelScopeDelete) ID(id interface{}) *ModelScopeDelete {
	s.args[0] = id
	return s
}

func (s *ModelScopeDelete) Limit(sizeP interface{}) *ModelScopeDelete {
	s.args[LimitPosition] = sizeP
	return s
}

func (s *ModelScopeDelete) Offset(offsetP interface{}) *ModelScopeDelete {
	s.args[OffsetPosition] = offsetP
	return s
}

func (s *ModelScopeDelete) Rebind(offset int, offsetP interface{}) *ModelScopeDelete {
	s.args[offset] = offsetP
	return s
}

func (s *ModelScopeDelete) Delete(args...interface{}) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	return s.db.ExecStatement(s.stmt, append(args, s.args...)...)
}

