package dao

import (
	"database/sql"
	"reflect"
	"errors"
	"log"
)

type DB struct {
	*sql.DB
	driverName string
	builder    SqlBuilder
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName, builder:DefaultBuilder()}, nil
}

func (db *DB) Save(v interface{}) (result sql.Result, err error) {

	tableInfo, err := db.builder.ParseStruct(v)
	if err != nil {
		return
	}

	fieldValues, err := FieldValue(v)
	if err != nil {
		return
	}

	args := make([]interface{}, 0, len(tableInfo.ColumnNames))
	for _, column := range tableInfo.ColumnNames {
		columnInfo, ok := tableInfo.Columns[column]
		if !ok {
			err = errors.New("column " + column + " mismatch")
			return
		}
		args = append(args, fieldValues[columnInfo.FieldName].Interface())
	}

	sql := tableInfo.Sqls[SQL_INSERT]
	result, err = db.Exec(sql, args...)

	return
}

func (db *DB) List(v interface{}, args...interface{}) error {

	if (reflect.TypeOf(v).Kind() != reflect.Ptr) {
		return errors.New("must pass a slice pointer, like &[]xxx")
	}

	listPtr := reflect.Indirect(reflect.ValueOf(v))
	listValue := reflect.MakeSlice(listPtr.Type(), 0, 1)

	tableInfo, err := db.builder.ParseType(listValue.Type().Elem())
	if err != nil {
		return err
	}

	sql := tableInfo.Sqls[SQL_SELECT]
	if args != nil {
		sql = sql + " " + args[0].(string)
		args = args[1:]
	}

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
			columnInfo, ok := tableInfo.Columns[column]
			if !ok {
				return errors.New("column " + column + " mismatch")
			}
			fieldValue := reflect.New(columnInfo.FieldType)
			fieldsSlice[i] = fieldValue.Interface()
		}

		err = rows.Scan(fieldsSlice...)
		if err != nil {
			return err
		}

		obj := reflect.New(tableInfo.StructType).Elem()
		for i, column := range columns {
			columnInfo, ok := tableInfo.Columns[column]
			if !ok {
				return errors.New("column " + column + " mismatch")
			}
			obj.FieldByName(columnInfo.FieldName).Set(reflect.Indirect(reflect.ValueOf(fieldsSlice[i])))
		}
		listValue = reflect.Append(listValue, obj)

	}
	listPtr.Set(listValue)
	return nil

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
