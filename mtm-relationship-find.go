package dorm

import "database/sql"

type ManyToManyRelationshipScopeFind struct {
	*ManyToManyRelationshipScope
	stmt string
	args []interface{}
}

func (s *ManyToManyRelationshipScope) BuildFind() (t *ManyToManyRelationshipScopeFind) {
	t = &ManyToManyRelationshipScopeFind{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "select " + s.VPK + " from `" + s.TableName + "`" + s.limitation(s.UPK, &t.args)
	return
}

func (s *ManyToManyRelationshipScopeFind) ID(id interface{}) *ManyToManyRelationshipScopeFind {
	// assert id in the where exp
	s.args[s.whereSize - 1] = id
	s.id = id
	return s
}

func (s *ManyToManyRelationshipScopeFind) Limit(sizeP interface{}) *ManyToManyRelationshipScopeFind {
	s.args[s.whereSize + LimitPosition] = sizeP
	s.limit = sizeP
	return s
}

func (s *ManyToManyRelationshipScopeFind) Offset(offsetP interface{}) *ManyToManyRelationshipScopeFind {
	s.args[s.whereSize + OffsetPosition] = offsetP
	s.offset = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeFind) Rebind(offset int, offsetP interface{}) *ManyToManyRelationshipScopeFind {
	s.args[offset] = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeFind) GroupBy(groupBy string) *ManyToManyRelationshipScopeFind {
	s.groupBy = groupBy
	return s
}

func (s *ManyToManyRelationshipScopeFind) Order(order string) *ManyToManyRelationshipScopeFind {
	s.order = order
	return s
}

func (s *ManyToManyRelationshipScopeFind) Find(result *[]uint, args ...interface{}) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}
	s.args[s.whereSize + LimitPosition] = s.limit
	s.args[s.whereSize + OffsetPosition] = s.offset
	copy(s.args, args)
	err = s.db.QueryStatement(s.decide(s.stmt), func(row *sql.Rows) error {
		if int(aff) < len(*result) {
			err := row.Scan(&(*result)[aff])
			if err != nil {
				return err
			}
		} else {
			var i uint
			err := row.Scan(&i)
			if err != nil {
				return err
			}
			*result = append(*result, i)
		}
		aff++
		return nil
	}, s.args...)
	return
}