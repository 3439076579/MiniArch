package session

import (
	"MiniArch/golangorm/logger"
	"reflect"
)

/*
	钩子所需要的入参为Session，返回值为error
*/

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

func (s *Session) CallMethod(method string, value interface{}) {

	function := reflect.ValueOf(s.table.Model).MethodByName(method)
	if value != nil {
		function = reflect.ValueOf(value).MethodByName(method)
	}

	params := []reflect.Value{reflect.ValueOf(s)}

	if function.IsValid() {
		result := function.Call(params)
		if len(result) > 0 {
			err, ok := result[0].Interface().(error)
			if ok {
				logger.Error(err)
			}
		}
	}
	return

}
