package dorm

type ManyToManyRelationshipScopeInsertSet struct {
	*ManyToManyRelationshipScope
	stmt string
	args []interface{}
}

func (s *ManyToManyRelationshipScope) BuildInsertSet() (t *ManyToManyRelationshipScopeInsertSet) {
	t = &ManyToManyRelationshipScopeInsertSet{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "insert into `" + s.TableName + "`(" + s.uCommon.fieldsColumns([]string{s.UPK, s.VPK}) +
		") values " + twoPlaceHolder
	return
}

func (s *ManyToManyRelationshipScopeInsertSet) ID(id interface{}) *ManyToManyRelationshipScopeInsertSet {
	s.id = id
	return s
}

func (s *ManyToManyRelationshipScopeInsertSet) InsertSet(v uint) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	return s.db.ExecStatement(s.stmt, s.id, v)
}