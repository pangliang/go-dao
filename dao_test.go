package dao

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"os"
	"reflect"
	"fmt"
)

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
		t.Fatalf("fieldsValue name got %#v \n", fieldsValue["name"])
	}

	if fieldsValue["Pwd"].Interface() != reflect.ValueOf("tom123").Interface() {
		t.Fatalf("fieldsValue pwd got %#v \n", fieldsValue["pwd"])
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
			t.Fatalf("expected RowsAffected 1, but got :%#v\n", rowAffected)
		}
	}

	var userList []User
	err = db.List(&userList)
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	for _, user := range userList {
		if user != m[user.Id] {
			t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[user.Id], user)
		}
	}

	err = db.List(&userList, "where name=? or pwd=?", "tom", "jake123")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	for _, user := range userList {
		if user != m[user.Id] {
			t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[user.Id], user)
		}
	}

	err = db.List(&userList, "order by id")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	fmt.Printf("%#v\n", userList)
	if userList[0] != m[1] {
		t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[1], userList[0])
	}
	if userList[1] != m[2] {
		t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[2], userList[1])
	}

	err = db.List(&userList, "order by id desc")
	fmt.Printf("%#v\n", userList)
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}
	if userList[0] != m[2] {
		t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[2], userList[0])
	}
	if userList[1] != m[1] {
		t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[1], userList[1])
	}

	var one []User
	err = db.List(&one, "where name=?", "tom")
	if err != nil {
		t.Fatalf("error:%s\n", err)
	}

	if len(one) != 1 {
		t.Fatalf("List fail expedcted 1 obj, but got :%#v\n", len(one))
	}

	if one[0] != m[1] {
		t.Fatalf("List fail expedcted %#v, but got :%#v\n", m[1], one[0])
	}
}

func TestUpdate(t *testing.T) {
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
		db.Save(user)
	}

	user := m[1]
	user.Name = "tom99"
	db.Update(user, "where id=?", 1)

	var list []User
	db.List(&list, "order by id")
	if list[0] != user {
		t.Fatalf("Update fail expedcted %#v, but got :%#v\n", user, list[0])
	}
	if list[1] != m[2] {
		t.Fatalf("Update fail expedcted %#v, but got :%#v\n", m[2], list[1])
	}

	_, err = db.Update(user)
	if err == nil {
		t.Fatalf("Update not where condition")
	}
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

