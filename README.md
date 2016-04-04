# go-dao

一个Go语言的ORM简单库, 主要功能就是通过反射Struct的属性推导出sql语句.
类似的库有: [sqlx](https://github.com/jmoiron/sqlx) ,但是它是把查询结果集*Row 装配成[]Sturct, sql语句还是需要手写