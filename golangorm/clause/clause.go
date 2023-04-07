package clause

import (
	"MiniArch/golangorm/logger"
	"strings"
)

type Clause struct {
	sql map[Type]string
}

func (s *Clause) Set(Kind Type, values ...interface{}) {

	if s.sql == nil {
		s.sql = make(map[Type]string)
	}

	v, ok := generators[Kind]
	if !ok {
		logger.Error("invalid Type")
	}

	sqlClause, _ := v(values[0].(string), values[1:]...)

	s.sql[Kind] = sqlClause

}

// Build 函数会根据据传入的Order顺序来构建SQL语句
func (s *Clause) Build(orders ...Type) string {

	var SqlSet []string

	for _, order := range orders {
		v, ok := s.sql[order]
		if ok {
			SqlSet = append(SqlSet, v)
		}
	}

	return strings.Join(SqlSet, " ")

}
