package schema

import (
	"MiniArch/golangorm/dialect"
	"go/ast"
	"reflect"
	"strings"
)

type Tabler interface {
	TableName() string
}

type Field struct {
	Name string
	Type string //对应数据库表中的field的Type
	Tag  string
}

type Schema struct {
	Model     interface{}
	ModelType string
	Fields    []*Field
	FieldName []string
	TableName string
	filedMap  map[string]*Field
}

// 将一个结构体解析为Schema

func Parse(model interface{}, dialect dialect.Dialect) *Schema {

	schema := new(Schema)

	//step1:解析model的type
	// reflect.indirect可以把指针或者非指针的对象的值提取出来
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	schema.ModelType = modelType.String()
	schema.Model = model
	schema.filedMap = make(map[string]*Field)

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {

			field := &Field{
				Name: p.Name,
				Type: dialect.DataType(reflect.Indirect(reflect.New(p.Type))),
			}

			if v, ok := p.Tag.Lookup("golorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldName = append(schema.FieldName, field.Name)
			schema.filedMap[field.Name] = field
		}

	}

	method, exist := reflect.TypeOf(modelType).MethodByName("TableName")

	if !exist {
		StringSet := strings.Split(schema.ModelType, ".")
		schema.TableName = StringSet[len(StringSet)-1] + "s"
	} else {
		tmp := method.Func.Call([]reflect.Value{reflect.ValueOf(model)})
		schema.TableName = tmp[0].String()
	}

	return schema

}

func (s *Schema) GetFieldByName(name string) *Field {
	return s.filedMap[name]
}
