package dorm

type modelCommon struct {
	model ORMObject
	id       interface{}
	limit    interface{}
	offset   interface{}
	order string
	groupBy string
	whereExp string
}

func (s modelCommon) limitation(upk string, args *[]interface{}) (exp string) {
	if s.id != nil {
		if len(s.whereExp) != 0 {
			exp = "(" + s.whereExp + ") and " + upk + " = ?"
		} else {
			exp = upk + " = ?"
		}
		*args = append(*args, s.id)
	} else if id := s.model.GetID(); id > 0 {
		if len(s.whereExp) != 0{
			exp = "(" + s.whereExp + ") and " + upk + " = ?"
		} else {
			exp = upk + " = ?"
		}
		*args = append(*args, id)
	}
	*args = append(*args, s.offset, s.limit)

	if len(exp) != 0 {
		exp = " where " + exp
	}

	return exp
}

func (s modelCommon) whereLimitation(upk string, args *[]interface{}) (exp string) {
	if s.id != nil {
		if len(s.whereExp) != 0 {
			exp = "(" + s.whereExp + ") and " + upk + " = ?"
		} else {
			exp = upk + " = ?"
		}
		*args = append(*args, s.id)
	} else if id := s.model.GetID(); id > 0 {
		if len(s.whereExp) != 0{
			exp = "(" + s.whereExp + ") and " + upk + " = ?"
		} else {
			exp = upk + " = ?"
		}
		*args = append(*args, id)
	} else {
		exp = s.whereExp
	}

	if len(exp) != 0 {
		exp = " where " + exp
	}

	return exp
}


func (s modelCommon) decide(exp string) string {
	if len( s.groupBy ) != 0 {
		exp = exp + " group by " + s.groupBy
	}
	if len( s.order ) != 0 {
		exp = exp + " order by " + s.order
	}
	return exp + " limit ?, ? "
}

