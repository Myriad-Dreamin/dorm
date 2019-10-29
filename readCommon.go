package dorm

import (
	"math"
	"strings"
)

type modelCommon struct {
	model ORMObject
	id       interface{}
	limit    interface{}
	offset   interface{}
	order string
	groupBy string
	whereSize int
	whereExp string
	limitType uint8
	args []interface{}
}

func (s *modelCommon) limitation(upk string) (exp string) {

	s.whereSize = strings.Count(s.whereExp, "?")
	s.args = make([]interface{}, s.whereSize)

	if s.id != nil {
		if len(s.whereExp) != 0 {
			exp = "(" + s.whereExp + ") and " + upk + " = ?"
		} else {
			exp = upk + " = ?"
		}
		s.whereSize++
		s.args = append(s.args, s.id)
	} else if id := s.model.GetID(); id > 0 {
		if len(s.whereExp) != 0{
			exp = "(" + s.whereExp + ") and " + upk + " = ?"
		} else {
			exp = upk + " = ?"
		}
		s.whereSize++
		s.args = append(s.args, id)
	}

	if len(exp) != 0 {
		exp = " where " + exp
	}

	//s.args = append(s.args, s.offset, s.limit)
	//s.limitType = 2
	return exp
}


func (s *modelCommon) decideArgs(i []interface{}) []interface{} {
	var j = i
	if s.limit == int64(math.MaxInt64) {
		s.limit = nil
	}
	if s.offset == 0 {
		s.offset = nil
	}
	if s.offset != nil {
		j = append(j, s.offset)
	}
	if s.limit != nil {
		j = append(j, s.limit)
	}
	return j
}

func (s *modelCommon) decide(exp string) string {
	if len( s.groupBy ) != 0 {
		exp = exp + " group by " + s.groupBy
	}
	if len( s.order ) != 0 {
		exp = exp + " order by " + s.order
	}
	if s.limit == int64(math.MaxInt64) {
		s.limit = nil
	}
	if s.offset == 0 {
		s.offset = nil
	}
	if s.limit != nil {
		if s.offset != nil {
			return exp + " limit ?, ? "
		} else {
			return exp + " limit ? "
		}
	} else if s.offset != nil {
		return exp + " offset ? "
	} else {
		return exp
	}
}

