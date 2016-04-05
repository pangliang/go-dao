package dao

import (
	"testing"
	"reflect"
)

type User struct {
	Id   uint32
	Name string
	Pwd  string
}

func TestParseStruct(t *testing.T) {
	reader := DefaultBuilder()
	user := User{}
	tableInfo, err := reader.ParseStruct(user)
	if err != nil {
		t.Fatal(err)
	}

	expected := TableInfo{
		StructName:"User",
		TableName:"user",
		StructType       :reflect.TypeOf(user),
		ColumnNames :[]string{"id", "name", "pwd"},
		Columns     :map[string]ColumnInfo{
			"id":ColumnInfo{FieldName:"Id", ColumnName:"id", FieldType:reflect.TypeOf(user.Id)},
			"name":ColumnInfo{FieldName:"Name", ColumnName:"name", FieldType:reflect.TypeOf(user.Name)},
			"pwd":ColumnInfo{FieldName:"Pwd", ColumnName:"pwd", FieldType:reflect.TypeOf(user.Pwd)},
		},
		Sqls:map[string]string{
			SQL_INSERT:"INSERT INTO user(id,name,pwd) VALUES (?,?,?)",
			SQL_SELECT:"SELECT id,name,pwd FROM user",
			SQL_UPDATE:"UPDATE user SET id=?,name=?,pwd=?",
		},
	}
	if !reflect.DeepEqual(tableInfo, expected) {
		t.Fatalf("ParseStruct fail, \ngot:     %s\nexpected:%v\n", tableInfo, expected)
	}
}

func TestCustomMapper(t *testing.T) {
	reader := DefaultBuilder()

	reader.SetMapper(func(s string) string {
		return s[1:]
	})

	user := User{}
	tableInfo, err := reader.ParseStruct(user)
	if err != nil {
		t.Fatal(err)
	}

	expected := TableInfo{
		StructName:"User",
		TableName:"ser",
		StructType       :reflect.TypeOf(user),
		ColumnNames :[]string{"d", "ame", "wd"},
		Columns     :map[string]ColumnInfo{
			"d":ColumnInfo{FieldName:"Id", ColumnName:"d", FieldType:reflect.TypeOf(user.Id)},
			"ame":ColumnInfo{FieldName:"Name", ColumnName:"ame", FieldType:reflect.TypeOf(user.Name)},
			"wd":ColumnInfo{FieldName:"Pwd", ColumnName:"wd", FieldType:reflect.TypeOf(user.Pwd)},
		},
		Sqls:map[string]string{
			SQL_INSERT:"INSERT INTO ser(d,ame,wd) VALUES (?,?,?)",
			SQL_SELECT:"SELECT d,ame,wd FROM ser",
			SQL_UPDATE:"UPDATE ser SET d=?,ame=?,wd=?",
		},
	}
	if !reflect.DeepEqual(tableInfo, expected) {
		t.Fatalf("ParseStruct fail, \ngot:     %s\nexpected:%v\n", tableInfo, expected)
	}
}

func BenchmarkParseStruct(b *testing.B) {
	reader := DefaultBuilder()
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
	reader := DefaultBuilder()
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
