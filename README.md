# go-dao

一个Go语言的ORM简单库, 主要功能就是通过反射Struct的属性推导出sql语句.

类似的库有: [sqlx](https://github.com/jmoiron/sqlx) ,但是它是把查询结果集*Row 装配成[]Sturct, sql语句还是需要手写


### Example

```go
db, _ := dao.Open("sqlite3", dbFile)

tom := User{1, "tom", "tom123"}
jake := User{2, "jake", "jake123"}
db.Save(tom)
db.Save(jake)

var list []User
err = db.List(&list)

fmt.Printf("%#v\n", list)
//Output: []dao.User{dao.User{Id:0x1, Name:"tom", Pwd:"tom123"}, dao.User{Id:0x2, Name:"jake", Pwd:"jake123"}}

tom.Name = "tom99"
db.Update(tom,"where id=?",1)

err = db.List(&list)
fmt.Printf("%#v\n", list)
//Output: []dao.User{dao.User{Id:0x1, Name:"tom99", Pwd:"tom123"}, dao.User{Id:0x2, Name:"jake", Pwd:"jake123"}}
```

### Benchmark

```
# 分析Struct不用cache的速度
BenchmarkParseStruct-2            500000              3102 ns/op
# 使用cache
BenchmarkParseStructUseCache-2  10000000               142 ns/op

# 原生手写sql的速度
BenchmarkInsertReference-2          5000            417416 ns/op
# 框架库
BenchmarkInsert-2                   5000            429126 ns/op
```