package dorm

import (
	"fmt"
	"math"
	"reflect"
)

type ManyToManyRelationship struct {
	*RCommon
	uCommon *common
	vCommon *common
	db      *DB
}

func (r *ManyToManyRelationship) Clone() *ManyToManyRelationship {
	return &ManyToManyRelationship{
		RCommon: r.RCommon,
		uCommon: r.uCommon,
		vCommon: r.vCommon,
		db:      r.db,
	}
}

func (r *ManyToManyRelationship) FixDB(db *DB) *ManyToManyRelationship {
	r.db = db
	return r
}

func (r *ManyToManyRelationship) FixSqlDB(db SQLCommon) *ManyToManyRelationship {
	r.db = r.db.Clone().FixSqlDB(db)
	return r
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
	partialSet bool
	fields     []string
	Error      error
}

func (r *ManyToManyRelationship) Anchor(u ORMObject) (s *ManyToManyRelationshipScope) {
	s = &ManyToManyRelationshipScope{
		ManyToManyRelationship: r,
	}
	s.model = u
	s.id = u.GetID()
	s.limit = int64(math.MaxInt64)
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
	LimitPosition  = 1
)

func (s *ManyToManyRelationshipScope) ID(id interface{}) *ManyToManyRelationshipScope {
	s.id = id
	return s
}

func (s *ManyToManyRelationshipScope) Select(fields ...string) *ManyToManyRelationshipScope {
	if len(fields) == 0 {
		s.partialSet = false
	} else {
		s.partialSet = true
		s.fields = fields
	}
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

func (s *ManyToManyRelationshipScope) Find(result *[]uint, args ...interface{}) (aff int64, err error) {
	return s.BuildFind().Find(result, args...)
}

func (s *ManyToManyRelationshipScope) InsertSet(idSet uint) (aff int64, err error) {
	return s.BuildInsertSet().InsertSet(idSet)
}

func (s *ManyToManyRelationshipScope) DeleteSet(idSet uint) (aff int64, err error) {
	return s.BuildDeleteSet().DeleteSet(idSet)
}

func (s *ManyToManyRelationshipScope) InsertsSet(idSet ...uint) (aff int64, err error) {
	return s.BuildInsertsSet().InsertsSet(idSet...)
}

func (s *ManyToManyRelationshipScope) Count(args ...interface{}) (count int64, err error) {
	return s.BuildCount().Count(args...)
}
