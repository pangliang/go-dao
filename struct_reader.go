package dao

import (
	"strings"
	"reflect"
	"errors"
	"sync"
)

type StructReader struct {
	mapper    func(s string) string
	cacheable bool
	cache     map[reflect.Type]StructInfo
	mutex     sync.Mutex
}

func DefaultReader() StructReader {
	return StructReader{mapper:strings.ToLower, cacheable:true, cache:make(map[reflect.Type]StructInfo)}
}

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

func (reader *StructReader) FieldValue(v interface{}) (fieldValue map[string]reflect.Value, err error) {
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

func (reader *StructReader) ParseStruct(v interface{}) (structInfo StructInfo, err error) {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		err = errors.New("v is not a struct")
		return
	}

	t := value.Type()

	return reader.ParseType(t)
}

func (reader *StructReader) ParseType(t reflect.Type) (structInfo StructInfo, err error) {

	if reader.cacheable {
		reader.mutex.Lock()
		structInfo, ok := reader.cache[t]
		if ok {
			reader.mutex.Unlock()
			return structInfo, nil
		}
		reader.mutex.Unlock()
	}
	structInfo = StructInfo{}
	structInfo.Name = t.Name()
	structInfo.Type = t
	structInfo.TableName = reader.mapper(t.Name())
	//structInfo.Fields = make(map[string]FieldInfo, t.NumField())
	structInfo.Columns = make(map[string]FieldInfo, t.NumField())
	structInfo.ColumnNames = make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fname := t.Field(i).Name
		fieldInfo := FieldInfo{}
		fieldInfo.Name = fname
		fieldInfo.ColumnName = reader.mapper(fname)
		fieldInfo.Type = t.Field(i).Type

		//structInfo.Fields[fname] = fieldInfo
		structInfo.Columns[fieldInfo.ColumnName] = fieldInfo
		structInfo.ColumnNames[i] = fieldInfo.ColumnName
	}

	if reader.cacheable {
		reader.mutex.Lock()
		reader.cache[t] = structInfo
		reader.mutex.Unlock()
	}
	return
}

func (reader *StructReader) CleanStructCache() {
	reader.cache = make(map[reflect.Type]StructInfo)
}

func (reader *StructReader) SetMapper(m func(s string) string) {
	reader.mapper = m
	reader.CleanStructCache()
}
