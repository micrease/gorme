
### GORME介绍

gorme是一个对gorm进行泛型封装的orm工具,可以让代码变的更简洁和优雅

#### 安装方法
```shell
go get github.com/micrease/gorme
```

#### 使用示例
更多方法请参考gorm官方文档,https://gorm.io/zh_CN/docs/

简单使用
```go
query := repo.NewQuery()
// SELECT * FROM `tb_example` WHERE age IN(20,21)  AND age >10  AND (age =20  OR age=23 OR (age=1 AND age=2)) AND `tb_example`.`deleted_at` IS NULL LIMIT 10
pageList, err := query.WhereIn("age", []any{20, 21}).Gt("age", 10).Where(func() {
    query.Eq("age", 20).Or("age=?", 23).Or(func() {
        query.Where("age", 1).Where("age", 2)
    })
}).Paginate(1, 10)
fmt.Println(pageList, err)
```
多种写法
```go
//相等条件,你认为合理的写法就是对的
repo.NewQuery().Where("age", 20).First()
repo.NewQuery().Where("age", "=", 20).First()
repo.NewQuery().Where("age=?", 20).First()
repo.NewQuery().Eq("age", 20).First()

repo.NewQuery().Where("age", ">", 20).First()
repo.NewQuery().Where("name", "IN", "张三,李四" ).First()
...
```
动态条件查询
```go
//case条件成立时,执行闭包中的方法
pageList, err := repo.NewQuery().Case(req.UserId > 0, func() {
    query.Where("user_id", req.UserId)
}).Case(len(req.GoodsName) > 0, func() {
    query.Like("goods_name", req.GoodsName)
}).Case(req.Amount > 0, func() {
    query.Where("amount", ">", req.Amount)
}).Paginate(1, 10)
```
联表查询,聚合查询
[更多见tests](https://github.com/micrease/gorme/blob/master/tests/gorme_test.go)

***
