package dorm

import (
	"errors"
	"math"
)

type Model struct {
	*common
	ori ORMObject
	db *DB
}

type ModelScope struct {
	modelCommon
	*Model
	partialSet bool
	fields []string
	Error error
}

func (m *Model) Anchor(u ORMObject) (s *ModelScope) {
	s = &ModelScope{Model: m,}
	s.model = u
	s.limit = int64(math.MaxInt64)
	s.offset = 0
	return
}

func (m *Model) Scope() (s *ModelScope) {
	s = &ModelScope{Model: m,}
	s.model = m.ori
	s.limit = int64(math.MaxInt64)
	s.offset = 0
	return
}

func (s *ModelScope) BindModel(obj ORMObject) *ModelScope {
	s.model = obj
	return s
}

func (s *ModelScope) ID(id interface{}) *ModelScope {
	s.id = id
	return s
}

func (s *ModelScope) Where(exp string) *ModelScope {
	if len(s.whereExp) == 0 {
		s.whereExp = exp
	} else {
		s.whereExp = "(" + s.whereExp + ") and " + exp
	}
	return s
}

func (s *ModelScope) Limit(sizeP interface{}) *ModelScope {
	s.limit = sizeP
	return s
}

func (s *ModelScope) Offset(offsetP interface{}) *ModelScope {
	s.offset = offsetP
	return s
}

func (s *ModelScope) GroupBy(groupBy string) *ModelScope {
	s.groupBy = groupBy
	return s
}

func (s *ModelScope) Order(order string) *ModelScope {
	s.order = order
	return s
}

func (s *ModelScope) Select (fields ...string) *ModelScope {
	if len(fields) == 0 {
		s.partialSet = false
	} else {
		s.partialSet = true
		s.fields = fields
	}
	return s
}

func (s *ModelScope) UpdateFields(args ...interface{}) (int64, error) {
	return s.BuildUpdateFields().UpdateFields(args...)
}

func (s *ModelScope) Insert() (int64, error) {
	return s.BuildInsert().Insert()
}

func (s *ModelScope) Delete(args ...interface{}) (int64, error) {
	return s.BuildDelete().Delete(args...)
}

func (s *ModelScope) Find(elem interface{}, args ...interface{}) error {

	return s.BuildFind().Find(elem, args...)
}

func (s *ModelScope) tagsTemplate() (string, string, error) {
	return "", "", errors.New("todo")
	//var b = new(bytes.Buffer)
	//for i, field := range t.typeInfo {
	//}
	//
	//return b.String(), "", nil
}


func (s *ModelScope) generateCreateSQL() (string, error) {
	tags, options, err := s.tagsTemplate()
	if err != nil {
		return "", err
	}

	return "create table" + s.tableName + " (" + tags + ")" + options, nil
}