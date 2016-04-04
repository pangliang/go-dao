package dao

import (
	"testing"
	"reflect"
)

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
