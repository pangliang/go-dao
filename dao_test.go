package dao

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"database/sql"
	"os"
)

type User struct {
	Id   uint32
	Name string
	Pwd  string
}

const dbFile = "./test.db"

const ddl = `
drop table if exists user;
CREATE TABLE IF NOT EXISTS USER
(
id INTEGER NOT NULL,
pwd TEXT DEFAULT '' NOT NULL,
name TEXT DEFAULT '' NOT NULL
);
`

func TestParseStruct(t *testing.T) {
	reader := DefaultReader()
	user := User{}
	structInfo, err := reader.ParseStruct(user)
	if err != nil {
		t.Fatal(err)
	}

	expected := StructInfo{
		Name:"User",
		TableName:"user",
		Type       :reflect.TypeOf(user),
		ColumnNames :[]string{"id", "name", "pwd"},
		Columns     :map[string]FieldInfo{
			"id":FieldInfo{Name:"Id", ColumnName:"id", Type:reflect.TypeOf(user.Id)},
			"name":FieldInfo{Name:"Name", ColumnName:"name", Type:reflect.TypeOf(user.Name)},
			"pwd":FieldInfo{Name:"Pwd", ColumnName:"pwd", Type:reflect.TypeOf(user.Pwd)},
		},
	}
	if !reflect.DeepEqual(structInfo, expected) {
		t.Fatalf("ParseStruct fail, \ngot:     %s\nexpected:%v\n", structInfo, expected)
	}
}

func TestCustomMapper(t *testing.T) {
	reader := DefaultReader()

	reader.SetMapper(func(s string) string {
		return s[1:]
	})

	user := User{}
	structInfo, err := reader.ParseStruct(user)
	if err != nil {
		t.Fatal(err)
	}

	expected := StructInfo{
		Name:"User",
		TableName:"ser",
		Type       :reflect.TypeOf(user),
		ColumnNames :[]string{"d", "ame", "wd"},
		Columns     :map[string]FieldInfo{
			"d":FieldInfo{Name:"Id", ColumnName:"d", Type:reflect.TypeOf(user.Id)},
			"ame":FieldInfo{Name:"Name", ColumnName:"ame", Type:reflect.TypeOf(user.Name)},
			"wd":FieldInfo{Name:"Pwd", ColumnName:"wd", Type:reflect.TypeOf(user.Pwd)},
		},
	}
	if !reflect.DeepEqual(structInfo, expected) {
		t.Fatalf("ParseStruct fail, \ngot:     %s\nexpected:%v\n", structInfo, expected)
	}
}

func TestFieldsValue(t *testing.T) {
	reader := DefaultReader()
	user := User{1, "tom", "tom123"}
	fieldsValue, err := reader.FieldValue(user)
	if err != nil {
		t.Fatal(err)
	}

	if fieldsValue["Id"].Interface() != reflect.ValueOf(uint32(1)).Interface() {
		t.Fatalf("fieldsValue id got %p \n", fieldsValue["id"].Interface())
	}

	if fieldsValue["Name"].Interface() != reflect.ValueOf("tom").Interface() {
		t.Fatalf("fieldsValue name got %v \n", fieldsValue["name"])
	}

	if fieldsValue["Pwd"].Interface() != reflect.ValueOf("tom123").Interface() {
		t.Fatalf("fieldsValue pwd got %v \n", fieldsValue["pwd"])
	}
}

func TestDaoList(t *testing.T) {
	os.Remove(dbFile)
	db, err := Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	defer db.Close()
	_, err = db.Exec(ddl);
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}

	m := map[uint32]User{
		1:User{1, "tom", "tom123"},
		2:User{2, "jake", "jake123"},
	}
	for _, user := range m {
		result, err := db.Save(user)
		if err != nil {
			t.Fatalf("error:%s\n", err)
		}

		rowAffected, _ := result.RowsAffected()
		if rowAffected != 1 {
			t.Fatalf("expected RowsAffected 1, but got :%v\n", rowAffected)
		}
	}

	var userList []User
	err = db.List(&userList)
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	for _, user := range userList {
		if user != m[user.Id] {
			t.Fatalf("List fail expedcted %v, but got :%v\n", m[user.Id], user)
		}
	}

	err = db.List(&userList, "where name=? or pwd=?", "tom", "jake123")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	for _, user := range userList {
		if user != m[user.Id] {
			t.Fatalf("List fail expedcted %v, but got :%v\n", m[user.Id], user)
		}
	}

	err = db.List(&userList, "order by id")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	if userList[0] != m[1] {
		t.Fatalf("List fail expedcted %v, but got :%v\n", m[1], userList[0])
	}
	if userList[1] != m[2] {
		t.Fatalf("List fail expedcted %v, but got :%v\n", m[2], userList[1])
	}

	err = db.List(&userList, "order by id desc")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	if userList[0] != m[2] {
		t.Fatalf("List fail expedcted %v, but got :%v\n", m[2], userList[0])
	}
	if userList[1] != m[1] {
		t.Fatalf("List fail expedcted %v, but got :%v\n", m[1], userList[1])
	}

	var one []User
	err = db.List(&one, "where name=?", "tom")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}

	if len(one) != 1 {
		t.Fatalf("List fail expedcted 1 obj, but got :%v\n", len(one))
	}

	if one[0] != m[1] {
		t.Fatalf("List fail expedcted %v, but got :%v\n", m[1], one[0])
	}
}

func BenchmarkParseStruct(b *testing.B) {
	reader := DefaultReader()
	reader.cacheable = false
	type Apple struct {
		p1  int
		p2  int
		p3  int
		p4  int
		p5  int
		p6  int
		p7  int
		p8  int
		p9  int
		p10 int
	}
	apple := Apple{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reader.ParseStruct(apple)
		if err != nil {
			b.Fatalf("error:%s\n", err)
		}
	}
}

func BenchmarkParseStructUseCache(b *testing.B) {
	reader := DefaultReader()
	type Apple struct {
		p1  int
		p2  int
		p3  int
		p4  int
		p5  int
		p6  int
		p7  int
		p8  int
		p9  int
		p10 int
	}
	apple := Apple{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reader.ParseStruct(apple)
		if err != nil {
			b.Fatalf("error:%s\n", err)
		}
	}
	b.StopTimer()
}

// 使用原生sql的基准测试, 作为对照
func BenchmarkInsertReference(b *testing.B) {
	os.Remove(dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		b.Fatalf("error:%s\n", err)
	}
	defer db.Close()
	_, err = db.Exec(ddl);
	if err != nil {
		b.Fatalf("error:%s\n", err)
	}
	b.ResetTimer()

	sql := "insert into user (id, name, pwd) values (? , ?, ?)"
	for i := 0; i < b.N; i++ {
		_, err = db.Exec(sql, i, i, i)
		if err != nil {
			b.Fatalf("error:%s\n", err)
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	os.Remove(dbFile)
	db, err := Open("sqlite3", dbFile)
	if err != nil {
		b.Fatalf("error:%s\n", err)
	}
	defer db.Close()
	_, err = db.Exec(ddl);
	if err != nil {
		b.Fatalf("error:%s\n", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var user = User{uint32(i), string(i), string(i)}
		_, err = db.Save(user)
		if err != nil {
			b.Fatalf("error:%s\n", err)
		}
	}
}

