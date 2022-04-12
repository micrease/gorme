
### GORME介绍

gorme是一个小巧又实用的gorm查询辅助工具,使用go1.18最新特性泛型封装,可以容易的实现分页查询,列表查询和单一结果查询。
```go
query := db.Where("age>?", 3)
result,err := gorme.Paginate[Product](query), 1, 10)
or:
result,err := repo.Paginate( 1, 10)
```
### 安装方法
```
go get github.com/micrease/gorme
```
> 由于使用了泛型，因此需要新版本go1.18  
> Goland对泛型支持尚不完善,更好的支持可以选择2022.1RC版本  
> 使用方法可参考example中的代码
### 使用说明
提供了两种使用方式:  
1.一种是基于gorm写法，仅对结果集进行轻度封装，无耦合。适合简单改造。    
2.另一种是通过继承Repository可以实现链式调用。代码更加优雅。更符合工程化。

### 1.函数调用
gorme只对查询结果进行了处理，因此你可以像往常一样使用gorm构建查询和排序等操作。只需要对获取结果进行很小的修改即可，通常是替换Fist,Find等结果获取方法  

| gorme方法                                                        | 功能说明                | 对应gorm                |
|:---------------------------------------------------------------|:--------------------|:----------------------|
| gorme.First\[T](query)(T,error)                                | 查询第一条记录             | db.First(&t)          |                
| gorme.Last\[T](query)(T,error)                                 | 查询最后一条记录            | db.Last(&t)           |                
| gorme.Take\[T](query)(T,error)                                 | 查询一条记录              | db.Take(&t)           |                
| gorme.GetOne\[T](query)(T,error)                               | 查询一条记录,Take的别名      | db.Take(&t)           |
| gorme.List\[T](query)([]T,error)                               | 查询多条记录,返回列表         | db.Find(&t)           |                
| gorme.Paginate\[T](query,pageNo,pageSize)(*PageResult[T],error)| 分页查询,返回PageResult结构 | db.Find(&t).Count(&c) |                

以下以Product为例,[详细代码](https://github.com/micrease/gorme/blob/master/example/example.go)
```go
type Product struct {
	gorm.Model
	Code  string `json:"code"`
	Price uint   `json:"price"`
}
```
构建query参数
```go
//使用gorm的常规方法构建查询,如
query := db.Where("id=?", 1)
query := db.Model(&Product{}).Select("age").Offset(1).Limit(1).Order("id desc").Where("id<?", 20).Where("price > ?", 1)
```
#### 单一数据查询
```go
//执行查询并返回Product类型的结果
product, err := gorme.GetOne[Product](query)
product, err := gorme.First[Product](query)
product, err := gorme.Last[Product](query)
product, err := gorme.Take[Product](query)
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
result, err := gorme.Paginate[Product](query, pageNo, pageSize)

fmt.Println("totalSize", result.TotalSize)
var products []Product
//可以看到result.List是一个[]Product类型，可以直接赋值给products
products = result.List
```
### 2.继承Repository链式调用
可以像gorm一样的调用方式。Repository继承了几乎所有gorm方法。

| Repository方法                                         | 功能说明                | 对应gorm                |
|:-----------------------------------------------------|:--------------------|:----------------------|
| repo.First()(T,error)                                | 查询第一条记录             | db.First(&t)          |                
| repo.Last()(T,error)                                 | 查询最后一条记录            | db.Last(&t)           |                
| repo.Take()(T,error)                                 | 查询一条记录              | db.Take(&t)           |                
| repo.GetOne()(T,error)                               | 查询一条记录,Take的别名      | db.Take(&t)           |
| repo.List(num)([]T,error)                            | 查询多条记录,返回列表         | db.Find(&t)           |                
| repo.Paginate(pageNo,pageSize)(*PageResult[T],error) | 分页查询,返回PageResult结构 | db.Find(&t).Count(&c) |

当然也可以使用gorm原生操作。u.DB即是*gorm.DB,你可以通过u.DB.Where(...).Find(...)形式进行原生查询.
```go
//这个一个例子
type ExampleModel struct {
    gorm.Model
    UserName string
    Age      int
}

//举一个例子，ExampleRepo(可以换成你自己定义的Repo)继承gorme.Repository[T]
type ExampleRepo struct {
    gorme.Repository[ExampleModel]
}

func NewExampleRepo(db *gorm.DB) *ExampleRepo {
    repo := ExampleRepo{}
    repo.SetDB(db)
    return &repo
}

func (e *ExampleRepo) GetFirst() (ExampleModel, error) {
    result, err := e.Select("id,age").Where("id>10").First()
    fmt.Println(result.ID, result.Age, result.UserName, err)
    return result, err
}

func (e *ExampleRepo) GetList() ([]ExampleModel, error) {
    //result, err := e.Select("age").Offset(1).Limit(5).Order("id desc").Where("id<?", 40).Where("age > ?", 1).List()
    //result, err := e.Limit(5).List()
    result, err := e.List(3)
    printList(result)
    return result, err
}

func (e *ExampleRepo) GetPaginateList() (*gorme.PageResult[ExampleModel], error) {
    result, err := e.Paginate(1, 10)
    fmt.Println("result list len=", len(result.List))
    fmt.Println(result.TotalSize, err)
    printList(result.List)
    return result, err
}

```
调用方式
```go
dsn := "gorme:123456@tcp(127.0.0.1:3306)/gorme?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
Logger: logger.Default.LogMode(logger.Info)})
if err != nil {
    panic("failed to connect database")
}
//分页查询
userRepo := NewUserRepo(db)
result,err:=userRepo.Paginate()
```
