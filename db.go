package dorm

import (
	"database/sql"
	"fmt"
	"io"

	_ "github.com/go-sql-driver/mysql"
)

type Escaper string
type DBType string

type SQLCommon interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type DB struct {
	SQLCommon
	escaper string
	logger Logger
}

func (db *DB) Clone() *DB {
	return &DB {
		SQLCommon: db.SQLCommon,
		escaper: db.escaper,
		logger: db.logger,
	}
}

func (db *DB) FixSqlDB(rdb SQLCommon) *DB {
	db.SQLCommon = rdb
	return db
}

func (db *DB) parseOption(options []interface{}) *DB {
	db.escaper = "`"
	db.logger = nl
	for _, option := range options {
		switch o := option.(type) {
		case Logger:
			db.logger = o
		case Escaper:
			db.escaper = string(o)
		}
	}
	return db
}

func getDBType(options []interface{}) string {
	for _, option := range options {
		switch o := option.(type) {
		case DBType:
			return string(o)
		}
	}
	return "mysql"
}


func Open(dsn string, options ...interface{}) (*DB, error) {
	rawDB, err := sql.Open(getDBType(options), dsn)
	if err != nil {
		return nil, err
	}
	return (&DB{SQLCommon: rawDB, logger: nl}).parseOption(options), nil
}

func IdleOpen(options ...interface{}) (*DB, error) {
	return (&DB{logger: nl}).parseOption(options), nil
}

func FromRaw(rdb SQLCommon, options ...interface{}) (*DB, error) {
	return (&DB{SQLCommon: rdb}).parseOption(options), nil
}


func (db *DB) Close() error {
	return db.SQLCommon.(io.Closer).Close()
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

func (db *DB) Statement(tables ...*TableObject) (t *Statement) {
	t = &Statement{db:db}
	t.registerTables(tables)
	t.tables = tables
	return
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

func (db *DB) SetEscaper(str string) {
	db.escaper = str
}

