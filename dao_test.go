package dao

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"strings"
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
delete from user;
`

func TestParseStruct(t *testing.T) {

	user := User{}
	structInfo, err := ParseStruct(user)
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
			"pwd":FieldInfo{Name:"Pwd", ColumnName:"pwd",Type:reflect.TypeOf(user.Pwd)},
		},
	}
	if !reflect.DeepEqual(structInfo, expected) {
		t.Fatalf("ParseStruct fail, \ngot:     %s\nexpected:%v\n", structInfo, expected)
	}
}

func TestCustomMapper(t *testing.T) {

	Mapper = func(s string) string {
		return s[1:]
	}

	user := User{}
	structInfo, err := ParseStruct(user)
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
			"wd":FieldInfo{Name:"Pwd", ColumnName:"wd",Type:reflect.TypeOf(user.Pwd)},
		},
	}
	if !reflect.DeepEqual(structInfo, expected) {
		t.Fatalf("ParseStruct fail, \ngot:     %s\nexpected:%v\n", structInfo, expected)
	}

	//reset Mapper
	Mapper = strings.ToLower
}

func TestFieldsValue(t *testing.T) {

	user := User{1, "tom", "tom123"}
	fieldsValue, err := FieldValue(user)
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
	err = db.List(&userList, nil)
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	for _, user := range userList {
		if user != m[user.Id] {
			t.Fatalf("List fail expedcted %v, but got :%v\n", m[user.Id], user)
		}
	}
}