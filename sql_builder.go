package dao

import (
	"strings"
	"reflect"
	"errors"
	"bytes"
)

const SQL_INSERT = "insert"
const SQL_SELECT = "select"
const SQL_UPDATE = "update"

var sqls = map[string]func(tableInfo TableInfo) (sql string, err error){
	SQL_INSERT:func(tableInfo TableInfo) (sql string, err error){
		buf := bytes.NewBufferString("INSERT INTO ")
		buf.WriteString(tableInfo.TableName)
		buf.WriteString("(")
		buf.WriteString(strings.Join(tableInfo.ColumnNames, ","))
		buf.WriteString(") VALUES (")

		bindStr := strings.Repeat("?,", len(tableInfo.ColumnNames))
		buf.WriteString(bindStr[:len(bindStr) - 1])
		buf.WriteString(")")

		sql = buf.String()

		return sql,nil
	},
	SQL_SELECT: func(tableInfo TableInfo)(sql string, err error){
		buf := bytes.NewBufferString("SELECT ")
		buf.WriteString(strings.Join(tableInfo.ColumnNames, ","))
		buf.WriteString(" FROM ")
		buf.WriteString(tableInfo.TableName)

		sql = buf.String()

		return sql, nil
	},
	SQL_UPDATE: func(tableInfo TableInfo)(sql string, err error){
		buf := bytes.NewBufferString("UPDATE ")
		buf.WriteString(tableInfo.TableName)
		buf.WriteString(" SET ")

		for _,column := range tableInfo.ColumnNames {
			buf.WriteString(column)
			buf.WriteString("=?,")
		}

		sql = buf.String()
		sql = sql[:len(sql)-1]
		return sql, nil
	},
}

type SqlBuilder struct {
	mapper     func(s string) string
	cacheable  bool
	tableCache map[reflect.Type]TableInfo
}

func DefaultBuilder() SqlBuilder {
	return SqlBuilder{
		mapper:strings.ToLower,
		cacheable:true,
		tableCache:make(map[reflect.Type]TableInfo),
	}
}

type ColumnInfo struct {
	FieldName  string
	ColumnName string
	FieldType  reflect.Type
}

type TableInfo struct {
	StructName  string
	StructType  reflect.Type
	TableName   string
	Columns     map[string]ColumnInfo
	ColumnNames []string
	Sqls   map[string]string
}

func (builder *SqlBuilder) ParseStruct(v interface{}) (tableInfo TableInfo, err error) {
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		err = errors.New("v is not a struct")
		return
	}

	return builder.ParseType(value.Type())
}

func (builder *SqlBuilder) ParseType(t reflect.Type) (tableInfo TableInfo, err error) {

	if builder.cacheable {
		tableInfo, ok := builder.tableCache[t]
		if ok {
			return tableInfo, nil
		}
	}
	tableInfo = TableInfo{}
	tableInfo.StructName = t.Name()
	tableInfo.StructType = t
	tableInfo.TableName = builder.mapper(t.Name())
	tableInfo.Columns = make(map[string]ColumnInfo, t.NumField())
	tableInfo.ColumnNames = make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		columnInfo := ColumnInfo{}
		columnInfo.FieldName = fieldName
		columnInfo.ColumnName = builder.mapper(fieldName)
		columnInfo.FieldType = t.Field(i).Type

		tableInfo.Columns[columnInfo.ColumnName] = columnInfo
		tableInfo.ColumnNames[i] = columnInfo.ColumnName
	}

	tableInfo.Sqls = make(map[string]string)
	for key,creater := range sqls {
		sql, err := creater(tableInfo)
		if err != nil {
			return tableInfo, err
		}
		tableInfo.Sqls[key] = sql
	}

	if builder.cacheable {
		builder.tableCache[t] = tableInfo
	}
	return
}

func (builder *SqlBuilder) CleanStructCache() {
	builder.tableCache = make(map[reflect.Type]TableInfo)
}

func (builder *SqlBuilder) SetMapper(m func(s string) string) {
	builder.mapper = m
	builder.CleanStructCache()
}