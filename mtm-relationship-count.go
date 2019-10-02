package dorm

type ManyToManyRelationshipScopeCount struct {
	*ManyToManyRelationshipScope
	stmt string
	args []interface{}
}

func (s *ManyToManyRelationshipScope) BuildCount() (t *ManyToManyRelationshipScopeCount) {
	t = &ManyToManyRelationshipScopeCount{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}
	t.stmt = "select count(" + s.VPK + ") from `" + s.TableName + "` " + s.limitation(s.UPK, &t.args)
	return
}

func (s *ManyToManyRelationshipScopeCount) Limit(sizeP interface{}) *ManyToManyRelationshipScopeCount {
	s.args[LimitPosition] = sizeP
	return s
}

func (s *ManyToManyRelationshipScopeCount) Offset(offsetP interface{}) *ManyToManyRelationshipScopeCount {
	s.args[OffsetPosition] = offsetP
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

func (s *ManyToManyRelationshipScopeCount) Rebind(offset int, offsetP interface{}) *ManyToManyRelationshipScopeCount {
	s.args[offset] = offsetP
	return s
}

func (s *ManyToManyRelationshipScopeCount) Count(args ...interface{}) (count int, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	err = s.db.QueryRowStatement(s.decide(s.stmt), []interface{}{count}, append(args, s.args...)...)
	return
}
