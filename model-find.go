package dorm

import (
	"errors"
	"reflect"
	"strings"
)


type ModelScopeFind struct {
	*ModelScope
	stmt      string
	args      []interface{}
	fetchFetchFunc FetchFetchFunc
}

func (s *ModelScope) BuildFind(options ...interface{}) (t *ModelScopeFind) {
	t = &ModelScopeFind{ModelScope: s}

	for _, option := range options {
		switch o := option.(type) {
		case FetchFetchFunc:
			t.fetchFetchFunc = o
		}
	}

	//t.stmt = "select " + s.VPK + " from " + s.TableName + s.limitation(s.UPK, &t.args)
	if s.partialSet {
		t.stmt = "select " + strings.Join(s.fields, ",") + " from `" + s.tableName + "` " +
			s.limitation("id", &t.args)
		if t.fetchFetchFunc == nil {
			t.fetchFetchFunc, t.Error = s.fieldsFetch(s.fields)
		}
	} else {
		t.stmt = "select " + s.fullFields() + " from `" + s.tableName + "` " +
			s.limitation("id", &t.args)
		if t.fetchFetchFunc == nil {
			t.fetchFetchFunc, t.Error = s.fullFetch()
		}
	}

	return t
}

func (s *ModelScopeFind) Model(obj ORMObject) *ModelScopeFind {
	if reflect.TypeOf(obj) != s.addressableType {
		s.Error = errors.New("type error")
		return s
	}
	s.model = obj
	s.args[0] = obj.GetID()
	return s
}

func (s *ModelScopeFind) ID(id interface{}) *ModelScopeFind {
	s.args[0] = id
	return s
}

func (s *ModelScopeFind) Limit(sizeP interface{}) *ModelScopeFind {
	s.args[LimitPosition] = sizeP
	return s
}

func (s *ModelScopeFind) Offset(offsetP interface{}) *ModelScopeFind {
	s.args[OffsetPosition] = offsetP
	return s
}

func (s *ModelScopeFind) Rebind(offset int64, offsetP interface{}) *ModelScopeFind {
	s.args[offset] = offsetP
	return s
}

func (s *ModelScopeFind) GroupBy(groupBy string) *ModelScopeFind {
	s.groupBy = groupBy
	return s
}

func (s *ModelScopeFind) Order(order string) *ModelScopeFind {
	s.order = order
	return s
}

func (s *ModelScopeFind) Find(elem interface{}, args... interface{}) (err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	return s.db.QueryStatement(s.decide(s.stmt), s.fetchFetchFunc(elem), append(args, s.args...)...)
}
