package dorm

import (
	"fmt"
	"reflect"
)

type ManyToManyRelationship struct {
	*RCommon
	uCommon *common
	vCommon *common
	db *DB
}

func (r *ManyToManyRelationship) Partial(u ORMObject) (pr *ManyToManyRelationship, err error) {
	switch reflect.TypeOf(u) {
	case r.uCommon.addressableType:
		return r, nil
	case r.vCommon.addressableType:
		return &ManyToManyRelationship{
			uCommon: r.vCommon,
			vCommon: r.uCommon,
			RCommon: r.RCommon.rotate(),
			db:      r.db,
		}, nil
	default:
		return nil, fmt.Errorf("cant partial %v object to relationship(%v, %v)",
			reflect.TypeOf(u), r.uCommon.addressableType, r.vCommon.addressableType)
	}
}

type ManyToManyRelationshipScope struct {
	modelCommon
	*ManyToManyRelationship
	Error error
}

func (r *ManyToManyRelationship) Anchor(u ORMObject) (s *ManyToManyRelationshipScope) {
	s = &ManyToManyRelationshipScope{
		ManyToManyRelationship: r,
	}
	s.model = u
	s.limit = -1
	s.offset = 0
	return
}

func (r *ManyToManyRelationship) MustAnchor(u ORMObject) (s *ManyToManyRelationshipScope) {
	if reflect.TypeOf(u) != r.uCommon.addressableType {
		s = &ManyToManyRelationshipScope{
			Error: fmt.Errorf("cant anchor %v object to relationship(%v, %v)",
				reflect.TypeOf(u), r.uCommon.addressableType, r.vCommon.addressableType),
			ManyToManyRelationship: r,
		}
		s.model = u
		return
	}
	return r.Anchor(u)
}

func (r *ManyToManyRelationship) RotateAndAnchor(u ORMObject) (s *ManyToManyRelationshipScope) {
	pr, err := r.Partial(u)
	if err != nil {
		return &ManyToManyRelationshipScope{Error: err, ManyToManyRelationship: r}
	}
	return pr.Anchor(u)
}

const (
	OffsetPosition = 0
	LimitPosition = 1
)

func (s *ManyToManyRelationshipScope) ID(id interface{}) *ManyToManyRelationshipScope {
	s.id = id
	return s
}

func (s *ManyToManyRelationshipScope) Where(exp string) *ManyToManyRelationshipScope {
	if len(s.whereExp) == 0 {
		s.whereExp = exp
	} else {
		s.whereExp = "(" + s.whereExp + ") and " + exp
	}
	return s
}

func (s *ManyToManyRelationshipScope) Limit(sizeP interface{}) *ManyToManyRelationshipScope {
	s.limit = sizeP
	return s
}

func (s *ManyToManyRelationshipScope) Offset(offsetP interface{}) *ManyToManyRelationshipScope {
	s.offset = offsetP
	return s
}

func (s *ManyToManyRelationshipScope) GroupBy(groupBy string) *ManyToManyRelationshipScope {
	s.groupBy = groupBy
	return s
}

func (s *ManyToManyRelationshipScope) Order(order string) *ManyToManyRelationshipScope {
	s.order = order
	return s
}

func (s *ManyToManyRelationshipScope) Find(result *[]uint, args ...interface{}) (aff int, err error) {
	return s.BuildFind().Find(result, args...)
}

func (s *ManyToManyRelationshipScope) Count(args ...interface{}) (count int, err error) {
	return s.BuildCount().Count(args...)
}
