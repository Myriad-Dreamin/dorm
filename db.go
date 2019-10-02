package dorm

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
	logger Logger
}

func Open(dsn string, options ...interface{}) (*DB, error) {
	rawDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db := &DB{DB: rawDB, logger: nl}

	for _, option := range options {
		if t, ok := option.(Logger); ok {
			db.logger = t
		}

		//switch o := option.(type) {
		//case Logger:
		//
		//}
	}

	return db, nil
}

func IdleOpen(options ...interface{}) (*DB, error) {

	db := &DB{logger: nl}

	for _, option := range options {
		if t, ok := option.(Logger); ok {
			db.logger = t
		}
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) ExecStatement(statement string, args ...interface{}) (int64, error) {
	//fmt.Println(statement, args)
	db.logger.Debug("exec ", "statement", statement, "args", args)

	stmt, err := db.Prepare(statement)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (db *DB) ExecStatementR(statement string, args ...interface{}) (sql.Result, error) {
	db.logger.Debug("exec ", "statement", statement, "args", args)

	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(args...)
}

type FetchFunc func(*sql.Rows) error

func (db *DB) closeRow(row *sql.Rows) {
	if err := row.Close(); err != nil {
		db.logger.Error("close connection error", "error", err)
	}
}

func (db *DB) QueryStatement(statement string, fetch FetchFunc, args ...interface{}) error {
	db.logger.Debug("exec ", "statement", statement, "args", args)

	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}

	row, err := stmt.Query(args...)
	if err != nil {
		return err
	}

	for row.Next() {
		err = fetch(row)
		if err != nil {
			db.closeRow(row)
			return err
		}
	}

	db.closeRow(row)
	return row.Err()
}

func (db *DB) QueryRowStatement(statement string, fetchHandler []interface{}, args ...interface{}) error {
	db.logger.Debug("exec ", "statement", statement, "args", args)

	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	return stmt.QueryRow(args...).Scan(fetchHandler...)
}

func (db *DB) Model(d ORMObject) (t *Model, err error) {
	t = &Model{db: db, ori: d}
	t.common, err = commonFrom(d)
	return
}

func (db *DB) ManyToManyRelation(u ORMObject, v ORMObject, options ...interface{}) (r *ManyToManyRelationship, err error) {
	r = &ManyToManyRelationship{db: db,}
	r.uCommon, err = commonFrom(u)
	if err != nil {
		return
	}
	r.vCommon, err = commonFrom(v)

	for _, option := range options {
		switch o := option.(type) {
		case *RCommon:
			r.RCommon = o
		case RCommon:
			r.RCommon = &o
		default:
			return nil, fmt.Errorf("invalid option of type %T", o)
		}
	}

	if r.RCommon == nil {
		r.RCommon = NewRCommon(
			r.uCommon.tableName+"_id", r.vCommon.tableName+"_id",
			r.uCommon.tableName+"_"+r.vCommon.tableName+"s")
	}

	return
}
