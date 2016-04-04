package dao

import (
	"database/sql"
	"reflect"
	"strings"
	"bytes"
	"errors"
	"log"
)

type DB struct {
	*sql.DB
	driverName string
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName}, nil
}

func (db *DB) Save(v interface{}) (result sql.Result, err error) {

	structInfo, err := ParseStruct(v)
	if err != nil {
		return
	}

	sql := bytes.NewBufferString("insert into ")
	sql.WriteString(structInfo.TableName)
	sql.WriteString("(")
	sql.WriteString(strings.Join(structInfo.ColumnNames, ","))
	sql.WriteString(") values (")

	bindStr := strings.Repeat("?,", len(structInfo.ColumnNames))
	sql.WriteString(bindStr[:len(bindStr) - 1])
	sql.WriteString(")")

	fieldValues, err := FieldValue(v)
	if err != nil {
		return
	}

	args := make([]interface{}, 0, len(structInfo.ColumnNames))
	for _, column := range structInfo.ColumnNames {
		fieldInfo, ok := structInfo.Columns[column]
		if !ok {
			err = errors.New("column " + column + " mismatch")
			return
		}
		args = append(args, fieldValues[fieldInfo.Name].Interface())
	}

	result, err = db.Exec(sql.String(), args...)

	return
}

func (db *DB) List(v interface{}, args...interface{}) error {

	if (reflect.TypeOf(v).Kind() != reflect.Ptr) {
		return errors.New("must pass a slice pointer, like &[]xxx")
	}

	listPtr := reflect.Indirect(reflect.ValueOf(v))
	listValue := reflect.MakeSlice(listPtr.Type(), 0, 1)

	structInfo, err := ParseType(listValue.Type().Elem())
	if err != nil {
		return err
	}

	buffer := bytes.NewBufferString("select ")
	buffer.WriteString(strings.Join(structInfo.ColumnNames, ","))
	buffer.WriteString(" from ")
	buffer.WriteString(structInfo.TableName)
	if args != nil {
		buffer.WriteString(" ")
		buffer.WriteString(args[0].(string))
		args = args[1:]
	}

	sql := buffer.String()
	log.Printf("%s\n", sql)
	rows, err := db.Query(sql, args...)
	if err != nil {
		return err

	}
	defer rows.Close()

	for rows.Next() {
		columns, _ := rows.Columns()
		fieldsSlice := make([]interface{}, len(columns))

		for i, column := range columns {
			fieldInfo, ok := structInfo.Columns[column]
			if !ok {
				return errors.New("column " + column + " mismatch")
			}
			fieldValue := reflect.New(fieldInfo.Type)
			fieldsSlice[i] = fieldValue.Interface()
		}

		err = rows.Scan(fieldsSlice...)
		if err != nil {
			return err
		}

		obj := reflect.New(structInfo.Type).Elem()
		for i, column := range columns {
			fieldInfo, ok := structInfo.Columns[column]
			if !ok {
				return errors.New("column " + column + " mismatch")
			}
			obj.FieldByName(fieldInfo.Name).Set(reflect.Indirect(reflect.ValueOf(fieldsSlice[i])))
		}
		listValue = reflect.Append(listValue, obj)

	}
	listPtr.Set(listValue)
	return nil

}
