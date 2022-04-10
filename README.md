
### GORME介绍

gorme是一个小巧又实用的gorm查询辅助工具,使用go1.18最新特性泛型封装,可以容易的实现分页查询,列表查询和单一结果查询。
### 安装方法
```
go get github.com/micrease/gorme
```
由于使用了泛型，因此需要go最新版本1.18  
### 使用说明
gorme只对查询结果进行了处理，因此你可以像往常一样使用gorm构建查询和排序等操作。只需要对获取结果进行很小的修改即可  
以下以Product为例,[详细代码](https://github.com/micrease/gorme/blob/master/example/example.go)
```go
type Product struct {
	gorm.Model
	Code  string `json:"code"`
	Price uint   `json:"price"`
}
```
#### 单一数据查询
```go
//使用gorm的常规方法构建查询
query := db.Where("id=?", 1)
//执行查询并返回Product类型的结果
product, err := gorme.GetOne[Product](query)
```

#### 列表查询
```go
//执行查询并返回Product类型的结果
products, err := gorme.List[Product](query)
```

#### 分页查询
分页查询结果数据结构
```go
type PageResult[T any] struct {
	PageSize  int    `json:"page_size"`  //每页条数
	PageNo    uint64 `json:"page_no"`    //当前页码
	TotalPage uint64 `json:"total_page"` //总页数
	TotalSize int64  `json:"total_size"` //总条数
	List      []T    `json:"list"`       //数据列表
}
```
分页查询方法
```go
//执行查询并返回一个分页结果集
result, err := gorme.PaginateSimple[Product](query, pageNo, pageSize)

fmt.Println("totalSize", result.TotalSize)
var products []Product
//可以看到result.List是一个[]Product类型，可以直接赋值给products
products = result.List
```
### 设计思想
1.通过赋值方式代码可读性更好  
通常查询结果通过传引用的方式使用反射机制赋值如:
```go
var product Product
db.Find(&product)
```
这种方式对于人类自然逻辑并不友好，并且每次查询都要提前声明一个变量接收值。显然通过赋值方式代码可读性更好,如:
```go
product := db.Find[Product]()
```
2.在比较复杂的结构中如分页查询，其结果集中除了列表外，还需要总条数，页码等信息。列表中的数据类型可以动态指定，因此使用泛型实现是一个很好的方式。

### features
1,增加更完善的查询条件构建方法。不需要关心where,limit,order等顺序问题。比如原生limit写在后面是无效的。  
2,update时传结构体会忽略0值。需要增加一个tag来控制是否在结构体中使用0值
