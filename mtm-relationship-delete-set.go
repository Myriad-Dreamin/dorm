package dorm

type ManyToManyRelationshipScopeDeleteSet struct {
	*ManyToManyRelationshipScope
	stmt string
}

func (s *ManyToManyRelationshipScope) BuildDeleteSet() (t *ManyToManyRelationshipScopeDeleteSet) {
	t = &ManyToManyRelationshipScopeDeleteSet{ManyToManyRelationshipScope: s.Where(s.VPK + " = ?")}
	if t.Error != nil {
		return
	}

	t.stmt = "delete from " + s.db.escaper + t.TableName + s.db.escaper + t.limitation(t.UPK)
	return
}

func (s *ManyToManyRelationshipScopeDeleteSet) ID(id interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.args[0] = id
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) Limit(sizeP interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.limit = sizeP
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) Offset(offsetP interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.offset = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) Rebind(offset int64, offsetP interface{}) *ManyToManyRelationshipScopeDeleteSet {
	s.args[offset] = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeDeleteSet) DeleteSet(args ...interface{}) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}
	sargs := s.decideArgs(s.args)
	copy(sargs, args)
	return s.db.ExecStatement(s.decide(s.stmt), sargs...)
}
