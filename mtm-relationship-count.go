package dorm

type ManyToManyRelationshipScopeCount struct {
	*ManyToManyRelationshipScope
	stmt string
}

func (s *ManyToManyRelationshipScope) BuildCount() (t *ManyToManyRelationshipScopeCount) {
	t = &ManyToManyRelationshipScopeCount{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}
	t.stmt = "select count(" + s.VPK + ") from `" + s.TableName + "` " + s.limitation(s.UPK)
	return
}

func (s *ManyToManyRelationshipScopeCount) ID(id interface{}) *ManyToManyRelationshipScopeCount {
	s.args[0] = id
	return s
}

func (s *ManyToManyRelationshipScopeCount) Limit(sizeP interface{}) *ManyToManyRelationshipScopeCount {
	s.limit = sizeP
	return s
}

func (s *ManyToManyRelationshipScopeCount) Offset(offsetP interface{}) *ManyToManyRelationshipScopeCount {
	s.offset = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeCount) Rebind(offset int64, offsetP interface{}) *ManyToManyRelationshipScopeCount {
	s.args[offset] = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeCount) GroupBy(groupBy string) *ManyToManyRelationshipScopeCount {
	s.groupBy = groupBy
	return s
}

func (s *ManyToManyRelationshipScopeCount) Order(order string) *ManyToManyRelationshipScopeCount {
	s.order = order
	return s
}

func (s *ManyToManyRelationshipScopeCount) Count(args ...interface{}) (count int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	sargs := s.decideArgs(s.args)
	copy(sargs, args)
	err = s.db.QueryRowStatement(s.decide(s.stmt), []interface{}{count}, sargs...)
	return
}
