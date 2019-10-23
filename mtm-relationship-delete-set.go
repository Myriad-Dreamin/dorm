package dorm

type ManyToManyRelationshipScopeDeleteSet struct {
	*ManyToManyRelationshipScope
	stmt string
	args []interface{}
}

func (s *ManyToManyRelationshipScope) BuildDeleteSet() (t *ManyToManyRelationshipScopeDeleteSet) {
	t = &ManyToManyRelationshipScopeDeleteSet{ManyToManyRelationshipScope: s.Where(s.VPK + " = ?")}
	if t.Error != nil {
		return
	}

	t.stmt = "delete from `" + t.TableName + "`" + t.limitation(t.UPK, &t.args)
	return
}

func (s *ManyToManyRelationshipScopeDeleteSet) ID(id interface{}) *ManyToManyRelationshipScopeDeleteSet {
	// assert id in the where exp
	s.args[s.whereSize - 1] = id
	s.id = id
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) Limit(sizeP interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.args[s.whereSize + LimitPosition] = sizeP
	s.limit = sizeP
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) Offset(offsetP interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.args[s.whereSize + OffsetPosition] = offsetP
	s.offset = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) Rebind(offset int64, offsetP interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.args[offset] = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) DeleteSet(args... interface{}) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	s.args[s.whereSize + LimitPosition] = s.limit
	s.args[s.whereSize + OffsetPosition] = s.offset
	copy(s.args, args)
	return s.db.ExecStatement(s.decide(s.stmt), s.args...)
}