package dorm

type ManyToManyRelationshipScopeInsertSet struct {
	*ManyToManyRelationshipScope
	stmt string
}

func (s *ManyToManyRelationshipScope) BuildInsertSet() (t *ManyToManyRelationshipScopeInsertSet) {
	t = &ManyToManyRelationshipScopeInsertSet{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "insert into " + s.db.escaper + s.TableName + s.db.escaper + "(" + s.uCommon.fieldsColumns([]string{s.UPK, s.VPK}) +
		") values " + twoPlaceHolder
	return
}

func (s *ManyToManyRelationshipScopeInsertSet) ID(id interface{}) *ManyToManyRelationshipScopeInsertSet {
	s.args[0] = id
	return s
}

func (s *ManyToManyRelationshipScopeInsertSet) InsertSet(v uint) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	return s.db.ExecStatement(s.stmt, s.id, v)
}