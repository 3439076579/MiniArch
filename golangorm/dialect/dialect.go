package dialect

import (
	"MiniArch/golangorm/logger"
	"reflect"
)

var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataType(typ reflect.Value) string
	TableExistSql(tableName string) string
}

func RegisterDialect(dataSource string, dialect Dialect) {
	dialectMap[dataSource] = dialect
}

func GetDialect(database string) Dialect {
	dialect, ok := dialectMap[database]
	if !ok {
		logger.Error("database has not been registered")
		panic("database has not been registered")
	}

	return dialect
}
