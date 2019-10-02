package dorm

//type Expr struct {
//	concation string
//	variableCount uint8
//	Left *Expr
//	Right *Expr
//}
//
//func (e *Expr) String() string {
//	if e.Left == nil {
//		return e.concation
//	}
//	return "(" + e.Left.String() + ")" + e.concation + "(" + e.Right.String() + ")"
//}
//
//func (s *ManyToManyRelationshipScope) IdExp(exp string) *Expr {
//	return &Expr{concation:exp, variableCount:1}
//}
//
//func (s *ManyToManyRelationshipScope) exp(a, b interface{}, o string) *Expr {
//	if _, ok := a.(string); ok {
//		a = s.IdExp(a.(string))
//	}
//	if _, ok := b.(string); ok {
//		b = s.IdExp(b.(string))
//	}
//	return &Expr{concation:o, variableCount:a.(*Expr).variableCount + b.(*Expr).variableCount,
//		Left: a.(*Expr), Right: b.(*Expr)}
//}
//
//func (s *ManyToManyRelationshipScope) AndExp(a, b interface{}) *Expr {
//	return s.exp(a, b, "and")
//}
//
//func (s *ManyToManyRelationshipScope) OrExp(a, b interface{}) *Expr {
//	return s.exp(a, b, "or")
//}

