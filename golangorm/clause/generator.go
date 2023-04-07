package clause

import (
	"fmt"
	"reflect"
	"strings"
)

type Type int

const (
	INSERT = iota
	SELECT
	LIMIT
	UPDATE
	WHERE
	DELETE
)

type generator func(tableName string, values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {

	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[UPDATE] = _update
	generators[LIMIT] = _limit
	generators[DELETE] = _delete

}

// INSERT INTO $tableName ($field1,$field2...) VALUES ($value1,$value2...)
// 预期接收的values为0 values[0]为tableName,后面的value为一个tableName对应的结构
func _insert(tableName string, values ...interface{}) (string, []interface{}) {

	var fieldAddr []string

	fieldAddr = values[0].([]string)

	fields := strings.Join(fieldAddr, ",")

	var Valueresult []string
	var tmp strings.Builder
	var v = values[1].([]interface{})
	for j := 0; j < len(v); j++ {
		for i := 0; i < len(fieldAddr); i++ {
			Value := reflect.Indirect(reflect.ValueOf(v[j])).FieldByName(fieldAddr[i])
			if i == 0 {
				tmp.WriteString("(")
			}
			if i == len(fieldAddr)-1 {
				if Value.Kind() == reflect.String {
					tmp.WriteString(fmt.Sprintf("'%v'", Value) + ")")
					break
				}
				tmp.WriteString(fmt.Sprintf("%v", Value) + ")")
				break
			}

			if Value.Kind() == reflect.String {
				tmp.WriteString(fmt.Sprintf("'%v',", Value))
			} else {
				tmp.WriteString(fmt.Sprintf("%v,", Value))
			}

		}
		Valueresult = append(Valueresult, tmp.String())
		tmp.Reset()
	}
	valueStr := strings.Join(Valueresult, ",")

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, fields, valueStr), []interface{}{}

}

// SELECT ($field1,$field2...) FROM $tables
func _select(tableName string, values ...interface{}) (string, []interface{}) {

	var fields string

	fmt.Println(values[0])
	var valueAddr []string

	for _, value := range values[0].([]string) {
		valueAddr = append(valueAddr, value)
	}

	fields = strings.Join(valueAddr, ",")

	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []interface{}{}
}

// WHERE ?
func _where(tableName string, values ...interface{}) (string, []interface{}) {

	return fmt.Sprintf("WHERE %s", values[0]), []interface{}{}

}

// 期望values[0]是field的集合，values[1]是var的集合
// UPDATE tableName SET $field1=$value1,...
func _update(tableName string, values ...interface{}) (string, []interface{}) {

	Map := values[0].(map[string]interface{})
	var strSet []string

	for key, value := range Map {
		strSet = append(strSet, fmt.Sprintf("%s = %v", key, value))
	}

	str := strings.Join(strSet, ",")

	return fmt.Sprintf("UPDATE %s SET %s", tableName, str), []interface{}{}
}

func _limit(tableName string, values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("LIMIT %v", values[0]), []interface{}{}
}

func _delete(tableName string, values ...interface{}) (string, []interface{}) {

	return fmt.Sprintf("DELETE FROM %s", tableName), []interface{}{}
}
