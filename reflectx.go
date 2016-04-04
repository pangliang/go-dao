package dao

import (
	"strings"
	"reflect"
	"errors"
	"fmt"
)

var Mapper = strings.ToLower
var Cacheable = true
var structInfoCache = make(map[string]StructInfo)

type FieldInfo struct {
	Name       string
	ColumnName string
	Type       reflect.Type
}

type StructInfo struct {
	Name        string
	Type        reflect.Type
	TableName   string
	//Fields      map[string]FieldInfo
	Columns     map[string]FieldInfo
	ColumnNames []string
}

func FieldValue(v interface{}) (fieldValue map[string]reflect.Value, err error) {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		err = errors.New("v is not a struct")
		return
	}

	t := value.Type()

	fieldValue = make(map[string]reflect.Value, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fname := t.Field(i).Name
		fieldValue[fname] = value.Field(i)
	}

	return
}

func ParseStruct(v interface{}) (structInfo StructInfo, err error) {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		err = errors.New("v is not a struct")
		return
	}

	t := value.Type()

	return ParseType(t)
}

func ParseType(t reflect.Type) (structInfo StructInfo, err error) {

	pkgPath := fmt.Sprintf("%v", t)
	if Cacheable {
		structInfo, ok := structInfoCache[pkgPath]
		if ok {
			return structInfo,nil
		}
	}
	structInfo = StructInfo{}
	structInfo.Name = t.Name()
	structInfo.Type = t
	structInfo.TableName = Mapper(t.Name())
	//structInfo.Fields = make(map[string]FieldInfo, t.NumField())
	structInfo.Columns = make(map[string]FieldInfo, t.NumField())
	structInfo.ColumnNames = make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fname := t.Field(i).Name
		fieldInfo := FieldInfo{}
		fieldInfo.Name = fname
		fieldInfo.ColumnName = Mapper(fname)
		fieldInfo.Type = t.Field(i).Type

		//structInfo.Fields[fname] = fieldInfo
		structInfo.Columns[fieldInfo.ColumnName] = fieldInfo
		structInfo.ColumnNames[i] = fieldInfo.ColumnName
	}

	if Cacheable {
		structInfoCache[pkgPath] = structInfo
	}
	return
}

func CleanStructCache() {
	structInfoCache = make(map[string]StructInfo)
}
