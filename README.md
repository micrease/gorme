
### GORME介绍

gorme是一个小巧又实用的gorm查询辅助工具,使用go1.18最新特性泛型封装,可以容易的实现分页查询,列表查询和单一结果查询。
### 安装方法
```
go get github.com/micrease/gorme
```
由于使用了泛型，因此需要go最新版本1.18  
### 使用说明
提供了两种使用方式:  
1.一种是基于gorm写法，仅对结果集进行轻度封装，无耦合。适合简单改造。    
2.另一种是通过继承Repository可以实现链式调用。代码更加优雅。更符合工程化。

### 1.函数调用
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
### 2.继承Repository链式调用
可以像gorm一样的调用方式。Repository继承了几乎所有gorm方法。新增加方法有:  
* First() 查询单条记录。不同与gorm中的First，此方法无参数。  
* List()  查询列表  
* Paginate(pageNo,PageSize) 分页函数  
当然也可以使用gorm原生操作。u.DB即是*gorm.DB,你可以通过u.DB.Where(...).Find(...)形式进行原生查询.
```go
type UserModel struct {
	gorm.Model
	UserName string
	Age      int
}

type UserRepo struct {
	gorme.Repository[UserModel]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	repo := UserRepo{}
	repo.DB = db
	return &repo
}

//对First方法进行了重写，注意返回值
func (u *UserRepo) GetFirst() (*UserModel, error) {
	result, err := u.Select("age").Offset(1).Limit(1).Order("id desc").Where("id<?", 20).Where("age > ?", 1).First()
	fmt.Println(result.Age, result.UserName, err)
	return result, err
}

//查询列表
func (u *UserRepo) GetList() (*[]UserModel, error) {
    result, err := u.Select("age").Offset(1).Limit(5).Order("id desc").Where("id<?", 40).Where("age > ?", 1).List()
    return result, err
}


//分页查询
func (u *UserRepo) Paginate() (*gorme.PageResult[UserModel], error) {
	result, err := u.Select("age").Order("id desc").Where("id<?", 20).Where("age > ?", 1).Paginate(1, 2)
	fmt.Println(result.TotalSize, err)
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
