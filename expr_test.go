package dorm

import (
	"fmt"
	"testing"
)

func TestExpr_String(t *testing.T) {
	fmt.Println(IdExp("a").And("b"))
	Date := Func("date")
	integer := Func("integer")
	fmt.Println(Date)
	fmt.Println(IdExp("a").Apply(Date))
	fmt.Println(IdExp("a, b").Apply(Date))
	fmt.Println(IdExp("a").Comma("b").Apply(Date))

	logger := NewFmtLogger()
	db, err := IdleOpen(logger)
	if err != nil {
		logger.Error("error open", "error", err)
	}

	fmt.Println()
	s := db.Statement()
	fmt.Println(s.Select(s.Var("hh"), s.Var("hh")))

	s = db.Statement()
	fmt.Println(s.Select(s.Var("hh"), s.Var("hh")).From(IdExp("table")))
	s = db.Statement()
	fmt.Println(s.Select(s.Var("hh"), s.Var("hh")).From(
		IdExp("table").As("a").Union(IdExp("table").As("b")).Union(IdExp("table").As("c"))))

	var table = new(User)

	s = db.Statement(Table(table).As("QAQ"))
	fmt.Println(s.Select(s.Var("hh"), s.VarHandler(&table.ID).As("'super hh'").Apply(integer)).From(
		IdExp("table").As("a").Union(IdExp("table").As("b")).Union(IdExp("table").As("c"))))


	s = db.Statement(Table(table).As("QAQ"))
	fmt.Println(s.Select(s.Var("hh"), s.VarHandler(&table.ID).As("'super hh'").Apply(integer)).From(
		s.Table(table).As("QwQ").Union(IdExp("table").As("b")).Union(IdExp("table").As("c"))))

	s = db.Statement(Table(table).As("QAQ"))
	fmt.Println(s.Select(s.Var("hh"), s.VarHandler(&table.ID).As("'super hh'").Apply(integer)).From(
		s.Table(table).As("QwQ").Union(IdExp("table").As("b")).Union(IdExp("table").As("c"))).Limit(3).Offset(5))

	//var id uint
	s = db.Statement(Table(table).As("QAQ"))
	s.Update(s.Table(table)).Set().Where(s.VarReceiver(&table.ID).Equal())
	fmt.Println(s, s.receivers)


	var id uint
	s = db.Statement(Table(table).As("QAQ"))
	s.Update(s.Table(table)).Set(s.VarReceiver(&table.ID), s.VarReceiver(&table.UserName)).
		Where(s.VarReceiver(&table.ID, &id))
	fmt.Println(s, s.receivers, []interface{}{&table.ID, &table.UserName, &id})
}


