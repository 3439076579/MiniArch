package session

import (
	"MiniArch/golangorm/clause"
	"MiniArch/golangorm/dialect"
	"MiniArch/golangorm/logger"
	"MiniArch/golangorm/schema"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type CommonDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Session struct {
	Db *sql.DB
	tx *sql.Tx

	sql strings.Builder

	sqlVars []interface{}

	clause  *clause.Clause
	table   *schema.Schema
	Dialect dialect.Dialect
}

func New(db *sql.DB, datasource string) *Session {

	dial := dialect.GetDialect(datasource)

	return &Session{
		Db:      db,
		Dialect: dial,
		clause:  new(clause.Clause),
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() CommonDB {

	if s.tx != nil {
		return s.tx
	}
	return s.Db
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

	row := s.DB().QueryRow(s.sql.String(), s.sqlVars...)
	return row

}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	logger.Info("执行查询", s.sql.String(), s.sqlVars)

	rows, err := s.DB().Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		logger.Error(err)
	}

	return rows, err

}

func (s *Session) Model(value interface{}) *Session {

	if s.table == nil || s.table.ModelType != reflect.TypeOf(value).String() {
		s.table = schema.Parse(value, s.Dialect)
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

	_, err := s.Raw(s.Dialect.TableExistSql(s.table.TableName)).Exec()
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func (s *Session) Insert(values ...interface{}) (int64, error) {

	table := s.Model(reflect.New(reflect.Indirect(reflect.ValueOf(values[0])).Type()).Interface()).Table()
	// BeforeInsert Hook
	for i := 0; i < len(values); i++ {
		s.CallMethod(BeforeInsert, values[i])
	}

	s.clause.Set(clause.INSERT, table.TableName, table.FieldName, values)
	str := s.clause.Build(clause.INSERT)
	//fmt.Println(str)
	//return 0, nil

	result, err := s.Raw(str).Exec()
	if err != nil {
		return 0, err
	}
	// AfterInsert Hook
	for i := 0; i < len(values); i++ {
		s.CallMethod(AfterInsert, values[i])
	}

	return result.RowsAffected()
}

// Find 函数会将查找到的结果放置到dest变量中，Find能够查询多个变量放入dest中
// dest可以是切片，也可以是结构体指针
func (s *Session) Find(dest interface{}) error {

	var destType reflect.Type
	destSet := reflect.Indirect(reflect.ValueOf(dest))
	// 如果传入进来的不是一个slice，而是一个结构体或者结构体指针，执行if逻辑
	if reflect.ValueOf(dest).Kind() != reflect.Slice {
		destType = reflect.Indirect(reflect.ValueOf(dest)).Type()
		fmt.Println(destType)
	} else {
		destType = destSet.Type().Elem()
	}

	table := s.Model(reflect.New(destType).Interface()).Table()
	// BeforeQuery Hook
	s.CallMethod(BeforeQuery, nil)

	s.clause.Set(clause.SELECT, table.TableName, table.FieldName)
	selectSql := s.clause.Build(clause.SELECT, clause.WHERE, clause.LIMIT)

	rows, err := s.Raw(selectSql).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {

		// 创建一个目标结构体来存放数据
		targetStruct := reflect.Indirect(reflect.New(destType))
		var values []interface{}
		for i := 0; i < len(table.FieldName); i++ {
			values = append(values, targetStruct.FieldByName(table.FieldName[i]).Addr().Interface())
		}
		err := rows.Scan(values...)
		if err != nil {
			return err
		}
		// AfterQuery
		s.CallMethod(AfterQuery, targetStruct.Addr().Interface())

		destSet.Set(reflect.Append(destSet, targetStruct))

	}
	return rows.Close()

}

// Update 支持两种传递参数的方式:1.传入一个map[string]interface{}
// 2.传入列表:field1,value1,field2,value2...
func (s *Session) Update(values ...interface{}) (int64, error) {
	Map, ok := values[0].(map[string]interface{})
	if !ok {
		// 代表使用传入列表的方式
		Map = make(map[string]interface{})
		for i := 0; i < len(values)-1; i += 2 {
			Map[values[i].(string)] = values[i+1]
		}
	}
	// BeforeUpdate Hook
	s.CallMethod(BeforeUpdate, nil)

	s.clause.Set(clause.UPDATE, s.Table().TableName, Map)
	UpdateSql := s.clause.Build(clause.UPDATE, clause.WHERE)

	result, err := s.Raw(UpdateSql).Exec()
	if err != nil {
		return 0, err
	}
	// AfterUpdate Hook
	s.CallMethod(AfterUpdate, nil)

	return result.RowsAffected()

}

func (s *Session) Delete() (int64, error) {

	table := s.table.TableName
	if table == "" {
		logger.Error("should use Model() before Limit")
		return 0, errors.New("should use Model() before Limit")
	}
	// BeforeDelete Hook
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, table)
	DeleteSql := s.clause.Build(clause.DELETE, clause.WHERE)

	result, err := s.Raw(DeleteSql).Exec()
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	// AfterDelete Hook
	s.CallMethod(AfterDelete, nil)

	return result.RowsAffected()

}

func (s *Session) Where(query string) *Session {

	s.clause.Set(clause.WHERE, s.table.TableName, query)

	return s

}
func (s *Session) Limit(num int64) *Session {

	table := s.table.TableName
	if table == "" {
		logger.Error("should use Model() before Limit")
		return nil
	}

	s.clause.Set(clause.LIMIT, table, num)

	return s
}
