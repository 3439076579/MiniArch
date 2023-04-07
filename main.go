package main

import (
	"MiniArch/golangorm"
	"MiniArch/golangorm/session"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type student struct {
	Name string `golorm:"PrimaryKey;default:not null"`
	Age  int64
}

func (s *student) BeforeInsert(session *session.Session) {

	s.Age += 100

	fmt.Println("Hello,-->BeforeInsert")
}

//func Hello(value ...interface{}) {
//	fmt.Println("Hello world", value)
//}

func main() {
	//[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	db := golangorm.Open("mysql", "root:wjb20031205@tcp(localhost:3306)/douyin_projoect")
	defer db.Close()

	s := session.New(db.DB(), "mysql")

	var u student

	s.Model(&student{}).Insert(&student{Name: "张三", Age: 10}, &student{Name: "李四", Age: 15})
	//s.Model(&student{}).Where("Age > 91").Find(&u)

	fmt.Println(&u)

}
