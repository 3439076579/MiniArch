package golangorm

import (
	"MiniArch/golangorm/dialect"
	"MiniArch/golangorm/logger"
	"MiniArch/golangorm/schema"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Session struct {
	db *sql.DB

	sql strings.Builder

	sqlVars []interface{}

	table   *schema.Schema
	dialect dialect.Dialect
}

func New(db *sql.DB, datasource string) *Session {

	dial := dialect.GetDialect(datasource)

	return &Session{
		db:      db,
		dialect: dial,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {

	defer s.Clear()
	logger.Info("执行:", s.sql.String(), s.sqlVars)

	res, err := s.DB().Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, err

}

func (s *Session) QueryRow() *sql.Row {

	defer s.Clear()
	logger.Info("执行查询", s.sql.String(), s.sqlVars)

	row := s.DB().QueryRow(s.sql.String(), s.sqlVars)
	return row

}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	logger.Info("执行查询", s.sql.String(), s.sqlVars)

	rows, err := s.DB().Query(s.sql.String(), s.sqlVars)
	if err != nil {
		logger.Error(err)
	}

	return rows, err

}

func (s *Session) Model(value interface{}) *Session {

	if s.table == nil || s.table.ModelType != reflect.TypeOf(value).String() {
		s.table = schema.Parse(value, s.dialect)
	}
	return s

}

func (s *Session) Table() *schema.Schema {
	if s.table != nil {
		return s.table
	}
	logger.Error("table is not set")
	return s.table
}

func (s *Session) CreateTable() error {
	t := s.table

	var columns []string
	for i := 0; i < len(t.Fields); i++ {
		var column string
		column = t.Fields[i].Name + "  " + t.Fields[i].Type
		if i != len(t.Fields)-1 {
			column = column + ","
		}

		columns = append(columns, column)
	}
	dest := strings.Join(columns, "\n")

	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s ( %s)", t.TableName, dest)).Exec()

	return err
}

func (s *Session) DropTable() error {

	_, err := s.Raw(fmt.Sprintf("DROP TABLE %s ;", s.table.TableName)).Exec()

	return err
}

func (s *Session) HasTable() bool {

	_, err := s.Raw(s.dialect.TableExistSql(s.table.TableName)).Exec()
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}
