package dorm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type TableObject struct {
	ORMObject
	*common
	alias string
}

func Table(table ORMObject) (obj *TableObject) {
	obj = &TableObject{
		ORMObject: table,
	}
	var err error
	obj.common, err = commonFrom(table)
	if err != nil {
		panic(fmt.Errorf("type %T is not avaliable table: %v", obj, err))
	}
	return
}

func (t *TableObject) String() string {
	return t.tableName
}

func (t *TableObject) Comma(b interface{}) *Expr {
	return exp(t, b, ",")
}

func (t *TableObject) Union(b interface{}) *Expr {
	return exp(t, b, "union")
}

func (t *TableObject) UnionAll(b interface{}) *Expr {
	return exp(t, b, "union all")
}

func (t *TableObject) Except(b interface{}) *Expr {
	return exp(t, b, "except")
}

func (t *TableObject) Intersect(b interface{}) *Expr {
	return exp(t, b, "intersect")
}

func (t *TableObject) InnerJoin(b interface{}) *Expr {
	return exp(t, b, "inner join")
}

func (t *TableObject) LeftJoin(b interface{}) *Expr {
	return exp(t, b, "left join")
}

func (t *TableObject) RightJoin(b interface{}) *Expr {
	return exp(t, b, "right join")
}

func (t *TableObject) On(b string) *Expr {
	return exp(t, b, "on")
}

func (t *TableObject) As(alias string) *TableObject {
	t.alias = alias
	return t
}

func (s *Statement) registerTables(tables []*TableObject) {
	for _, table := range tables {
		s.registerTable(table)
		if s.Error != nil {
			return
		}
	}
	return
}

func (s *Statement) registerTable(table *TableObject) {
	if s.tablesVars == nil {
		s.tablesVars = make(map[interface{}]*TSType)
	}
	tableValue := reflect.ValueOf(table.ORMObject).Elem()
	for _, fieldType := range table.common.typeInfoSlice {
		addr := tableValue.Field(fieldType.FieldOffset).Addr().Interface()
		if _, ok := s.tablesVars[addr]; ok {
			s.Error = errors.New("duplicate registered address")
			return
		}
		s.tablesVars[addr] = &TSType{
			SType:  fieldType,
			parent: table,
		}
	}
}

type Expr struct {
	repr  interface{}
	Left  *Expr
	Right *Expr
}

func (e *Expr) maybeParens() (string, bool) {
	if e == nil {
		return "", false
	}
	if e.Left == nil {
		if s, ok := e.repr.(string); ok {
			return s, false
		} else if s, ok := e.repr.(fmt.Stringer); ok {
			return s.String(), false
		} else {
			panic(fmt.Errorf("bad expression type %T", e.repr))
		}
	}

	return e.applyMaybeParens(), true
}

func (e *Expr) applyMaybeParens() string {
	if e == nil || e.Left == nil {
		a, _ := e.maybeParens()
		return a
	}
	var leftExp, rightExp string
	if s, ok := e.repr.(string); !(ok && len(s) == 1 && s == ",") {
		var p bool
		if leftExp, p = e.Left.maybeParens(); p {
			leftExp = "(" + leftExp + ")"
		}
		if rightExp, p = e.Right.maybeParens(); p {
			rightExp = "(" + rightExp + ")"
		}
		return leftExp + " " + e.repr.(string) + " " + rightExp
	} else {
		leftExp, _ = e.Left.maybeParens()
		rightExp, _ = e.Right.maybeParens()
		return leftExp + "" + e.repr.(string) + " " + rightExp
	}
}

func (e *Expr) String() string {
	return e.applyMaybeParens()
}

type FunctionObject struct {
	f      string
	params *Expr
}

func Func(f string, params ...interface{}) *FunctionObject {
	return &FunctionObject{
		f:      f,
		params: Commas(params...),
	}
}

func (f FunctionObject) String() string {
	return f.f + "(" + f.params.String() + ")"
}

func (f *FunctionObject) Apply(params ...interface{}) *FunctionObject {
	return &FunctionObject{
		f:      f.f,
		params: f.params.Commas(params...),
	}
}

func FuncE(f string, params ...interface{}) *Expr {
	return IdExp(Func(f, params...))
}

func (f *FunctionObject) ApplyE(params ...interface{}) *Expr {
	return IdExp(f.Apply(params...))
}

func Commas(params ...interface{}) *Expr {
	if len(params) == 0 {
		return nil
	}
	return IdExp(params[0]).Commas(params[1:]...)
}

func (e *Expr) Commas(params ...interface{}) *Expr {
	if e == nil {
		return Commas(params...)
	}
	if len(params) == 0 {
		return e
	}
	var t = new(Expr)
	t = e.Comma(params[0])
	for _, param := range params[1:] {
		t = t.Comma(param)
	}
	return t
}

func IdExp(exp interface{}) *Expr {
	return &Expr{repr: exp}
}

func exp(a, b interface{}, o string) *Expr {
	switch ax := a.(type) {
	case string, *FunctionObject, FunctionObject, *TableObject, TableObject:
		a = IdExp(ax)
	}
	switch bx := b.(type) {
	case string, *FunctionObject, FunctionObject, *TableObject, TableObject:
		b = IdExp(bx)
	}
	return &Expr{repr: o, Left: a.(*Expr), Right: b.(*Expr)}
}

func travelAs(a *Expr, alias string) {
	if a == nil {
		return
	}
	switch ax := a.repr.(type) {
	case *TableObject:
		ax.alias = alias
		return
	case string:
		if len(ax) == 2 && ax == "as" {
			travelAs(a.Left, alias)
			*a = *a.Left
			return
		}
		travelAs(a.Left, alias)
		travelAs(a.Right, alias)
	default:
		return
	}
}

func asExp(a interface{}, alias string) *Expr {
	switch ax := a.(type) {
	case *TableObject:
		ax.alias = alias
		a = IdExp(ax)
	case *Expr:
		travelAs(ax, alias)
	}
	return &Expr{repr: "as", Left: a.(*Expr), Right: IdExp(alias)}
}

func (e *Expr) And(b interface{}) *Expr {
	if e == nil {
		panic(fmt.Errorf("nil expression"))
	}
	return exp(e, b, "and")
}

func (e *Expr) Or(b interface{}) *Expr {
	return exp(e, b, "or")
}

func (e *Expr) Equal() *Expr {
	return exp(e, "?", "=")
}

func (e *Expr) Comma(b interface{}) *Expr {
	return exp(e, b, ",")
}

func (e *Expr) Union(b interface{}) *Expr {
	return exp(e, b, "union")
}

func (e *Expr) UnionAll(b interface{}) *Expr {
	return exp(e, b, "union all")
}

func (e *Expr) Except(b interface{}) *Expr {
	return exp(e, b, "except")
}

func (e *Expr) Intersect(b interface{}) *Expr {
	return exp(e, b, "intersect")
}

func (e *Expr) InnerJoin(b interface{}) *Expr {
	return exp(e, b, "inner join")
}

func (e *Expr) LeftJoin(b interface{}) *Expr {
	return exp(e, b, "left join")
}

func (e *Expr) RightJoin(b interface{}) *Expr {
	return exp(e, b, "right join")
}

func (e *Expr) On(b string) *Expr {
	return exp(e, b, "on")
}

func (e *Expr) As(b string) *Expr {
	if t, ok := e.repr.(*TableObject); ok {
		return IdExp(t.As(b))
	}
	return asExp(e, b)
}

func (e *Expr) Apply(f *FunctionObject) *Expr {
	return f.ApplyE(e)
}

type limitStatement struct {
	i interface{}
}

func (s limitStatement) String() string {
	return "limit " + stringify(s.i)
}

type offsetStatement struct {
	i interface{}
}

func (s offsetStatement) String() string {
	return "offset " + stringify(s.i)
}

type orderByStatement string

func (s orderByStatement) String() string {
	return "order by " + string(s)
}

type groupByStatement string

func (s groupByStatement) String() string {
	return "group by " + string(s)
}

type selectExp []interface{}

func (s selectExp) String() string {
	return "select " + Commas(s...).String()
}

type setExp []interface{}

func (s setExp) String() string {
	return "set " + Commas(s...).String()
}

type updateExp []interface{}

func (s updateExp) String() string {
	return "update " + Commas(s...).String()
}

type deleteExp struct{}

func (s deleteExp) String() string {
	return "delete"
}

type insertExp struct{}

func (s insertExp) String() string {
	return "insert"
}

type whereExp struct {
	e *Expr
}

func (s whereExp) String() string {
	return "where " + s.e.String()
}

type fromExp struct {
	e *Expr
}

func (s fromExp) String() string {
	return "from " + s.e.String()
}

type intoExp struct {
	e *Expr
}

func (s intoExp) String() string {
	return "into " + s.e.String()
}

type Statement struct {
	total      *Statement
	tablesVars map[interface{}]*TSType
	tables     []*TableObject
	db         *DB
	conds      []interface{}
	receivers  []interface{}
	Error      error
}

type TSType struct {
	*SType
	parent *TableObject
}

func (s *TSType) String() string {
	if len(s.parent.alias) != 0 {
		return s.parent.alias + "." + s.ColumnName
	}
	return s.ColumnName
}

func stringify(c interface{}) string {
	switch c := c.(type) {
	case fmt.Stringer:
		return c.String()
	case string:
		return c
	case []byte:
		return string(c)
	case uint8:
		return strconv.FormatUint(uint64(c), 10)
	case uint16:
		return strconv.FormatUint(uint64(c), 10)
	case uint32:
		return strconv.FormatUint(uint64(c), 10)
	case uint64:
		return strconv.FormatUint(c, 10)
	case uint:
		return strconv.FormatUint(uint64(c), 10)
	case int:
		return strconv.FormatInt(int64(c), 10)
	case int8:
		return strconv.FormatInt(int64(c), 10)
	case int16:
		return strconv.FormatInt(int64(c), 10)
	case int32:
		return strconv.FormatInt(int64(c), 10)
	case int64:
		return strconv.FormatInt(c, 10)
	case float32:
		return strconv.FormatFloat(float64(c), 'f', 15, 64)
	case float64:
		return strconv.FormatFloat(c, 'f', 8, 32)
	default:
		panic(fmt.Errorf("cant stringify %T", c))
	}
}

func (s Statement) String() string {
	var b = new(bytes.Buffer)
	for _, cond := range s.conds {
		b.WriteString(stringify(cond))
		b.WriteByte(' ')
	}
	return b.String()
}

func (s *Statement) Var(a interface{}) *Expr {
	return IdExp(a)
}

func (s *Statement) Table(a ORMObject) *Expr {
	for _, t := range s.tables {
		if t.ORMObject == a {
			if len(t.alias) != 0 {
				return asExp(t, t.alias)
			}
			return IdExp(t)
		}
	}
	t := Table(a)
	s.registerTable(t)
	return IdExp(t)
}

func (s *Statement) VarHandler(a interface{}) *Expr {
	if tVar, ok := s.tablesVars[a]; ok {
		return IdExp(tVar)
	} else {
		panic("var handler is not in any table registered")
	}
}

func (s *Statement) VarReceiver(a interface{}, options ...interface{}) *Expr {
	return exp(s.VarArgs(a, options...), "?", "=")
}

func (s *Statement) VarArgs(a interface{}, options ...interface{}) *Expr {
	if tVar, ok := s.tablesVars[a]; ok {
		if len(options) == 0 {
			s.receivers = append(s.receivers, a)
		} else if len(options) == 1 {
			s.receivers = append(s.receivers, options[0])
		} else if len(options) == 2 {
			i := options[0].(int)
			if len(s.receivers) < i {
				s.receivers = append(s.receivers, make([]interface{}, i - len(s.receivers))...)
				s.receivers[i] = options[1]
			} else if len(s.receivers) == i {
				s.receivers = append(s.receivers, options[1])
			}
		} else {
			panic("var receiver need options count of 0 ~ 2")
		}
		return IdExp(tVar)
	} else {
		panic("var handler is not in any table registered")
	}
}

func (s *Statement) Select(vars ...interface{}) *Statement {
	s.conds = append(s.conds, selectExp(vars))
	return s
}

func (s *Statement) Update(vars ...interface{}) *Statement {
	s.conds = append(s.conds, updateExp(vars))
	return s
}

func (s *Statement) Set(vars ...interface{}) *Statement {
	s.conds = append(s.conds, setExp(vars))
	return s
}

func (s *Statement) Delete() *Statement {
	s.conds = append(s.conds, deleteExp{})
	return s
}

func (s *Statement) Insert() *Statement {
	s.conds = append(s.conds, insertExp{})
	return s
}

func (s *Statement) Where(exp *Expr) *Statement {
	s.conds = append(s.conds, whereExp{exp})
	return s
}

func (s *Statement) From(exp *Expr) *Statement {
	s.conds = append(s.conds, fromExp{exp})
	return s
}

func (s *Statement) Into(exp *Expr) *Statement {
	s.conds = append(s.conds, intoExp{exp})
	return s
}

func (s *Statement) Limit(limit interface{}) *Statement {
	s.conds = append(s.conds, limitStatement{limit})
	return s
}

func (s *Statement) Offset(offset interface{}) *Statement {
	s.conds = append(s.conds, offsetStatement{offset})
	return s
}

func (s *Statement) OrderBy(order orderByStatement) *Statement {
	s.conds = append(s.conds, order)
	return s
}

func (s *Statement) GroupBy(group groupByStatement) *Statement {
	s.conds = append(s.conds, group)
	return s
}
