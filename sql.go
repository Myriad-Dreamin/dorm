package dorm

import (
	"fmt"
)

type Table struct {
	typeInfo TypeInfo
	db *DB
}

func (db *DB) callSQL(cmd string, args ...interface{}) (int64, error) {
	stmt, err := db.db.Prepare(cmd)
	fmt.Println(cmd, args)

	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (db *DB) Executor(d ORMObject) (t *Table, err error) {
	t = new(Table)
	t.typeInfo, err = generateTypeInfo(d)
	if err != nil {
		return nil, err
	}
	t.db = db
	return
}
