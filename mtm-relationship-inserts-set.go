package dorm

import (
	"bytes"
)

type ManyToManyRelationshipScopeInsertsSet struct {
	*ManyToManyRelationshipScope
	stmt string
}

func (s *ManyToManyRelationshipScope) BuildInsertsSet() (t *ManyToManyRelationshipScopeInsertsSet) {
	t = &ManyToManyRelationshipScopeInsertsSet{ManyToManyRelationshipScope: s}
	if s.Error != nil {
		return
	}

	t.stmt = "insert into `" + s.TableName + "`(" + s.uCommon.fieldsColumns([]string{s.UPK, s.VPK}) +
		") values " + twoPlaceHolder
	return
}

func (s *ManyToManyRelationshipScopeInsertsSet) ID(id interface{}) *ManyToManyRelationshipScopeInsertsSet {
	s.args[0] = id
	return s
}

const dotTwoPlaceHolder = "," + twoPlaceHolder

func (s *ManyToManyRelationshipScopeInsertsSet) InsertsSet(vs ...uint) (aff int64, err error) {
	if s.Error != nil {
		err = s.Error
		return
	}

	var fvs = make([]interface{}, len(vs)<<1)

	var stmtBuf = bytes.NewBufferString(s.stmt)
	for i := len(vs) - 1; i >= 0; i-- {
		fvs[i<<1] = s.id
		fvs[i<<1|1] = vs[i]
	}
	for i := len(vs) - 1; i > 0; i-- {
		stmtBuf.WriteString(dotTwoPlaceHolder)
	}
	return s.db.ExecStatement(stmtBuf.String(), fvs...)
}
