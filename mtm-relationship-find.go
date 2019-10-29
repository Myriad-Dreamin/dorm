package dorm

import "database/sql"

type ManyToManyRelationshipScopeFind struct {
	*ManyToManyRelationshipScope
	stmt string
}

func (s *ManyToManyRelationshipScope) BuildFind() (t *ManyToManyRelationshipScopeFind) {
	t = &ManyToManyRelationshipScopeFind{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "select " + s.VPK + " from `" + s.TableName + "`" + s.limitation(s.UPK)
	return
}

func (s *ManyToManyRelationshipScopeFind) ID(id interface{}) *ManyToManyRelationshipScopeFind {
	s.args[0] = id
	return s
}

func (s *ManyToManyRelationshipScopeFind) Limit(sizeP interface{}) *ManyToManyRelationshipScopeFind {
	s.limit = sizeP
	return s
}

func (s *ManyToManyRelationshipScopeFind) Offset(offsetP interface{}) *ManyToManyRelationshipScopeFind {
	s.offset = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeFind) Rebind(offset int64, offsetP interface{}) *ManyToManyRelationshipScopeFind {
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
	sargs := s.decideArgs(s.args)
	copy(sargs, args)
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
	}, sargs...)
	return
}