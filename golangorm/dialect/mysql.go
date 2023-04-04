package dialect

import (
	"MiniArch/golangorm/logger"
	"fmt"
	"reflect"
	"time"
)

// mysql ->Dialect
type mysql struct{}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (m *mysql) DataType(typ reflect.Value) string {

	switch typ.Kind() {
	case reflect.Bool:
		return "TINYINT"
	case reflect.Uint, reflect.Uint32, reflect.Int16, reflect.Int,
		reflect.Uintptr, reflect.Uint8, reflect.Int8, reflect.Uint16,
		reflect.Int32:
		return "INTEGER"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.Float32:
		return "FLOAT"
	case reflect.String:
		return "TEXT"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Timer); ok {
			return "DATATIME(3)"
		}
	}
	logger.Error("invalid typ")
	panic("invalid typ")

}

func (m *mysql) TableExistSql(tableName string) string {
	return fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)

}
