package main

import (
	"fmt"
	"github.com/micrease/gorme"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
)

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
	result, err := e.Select("id,age").Where("id>10").GetOne()
	fmt.Println(result.ID, result.Age, result.UserName, err)
	return result, err
}

func (e *ExampleRepo) GetList() ([]ExampleModel, error) {
	//result, err := e.Select("age").Offset(1).Limit(5).Order("id desc").Where("id<?", 40).Where("age > ?", 1).List()
	//result, err := e.Limit(5).List()
	result, err := e.Order("id desc").List(3)
	printList(result)
	return result, err
}

func (e *ExampleRepo) GetPaginateList() (*gorme.PageResult[ExampleModel], error) {
	result, err := e.Select("*").Paginate(1, 10)
	fmt.Println("result list len=", len(result.List))
	fmt.Println(result.TotalSize, err)
	printList(result.List)
	return result, err
}

func printList(list []ExampleModel) {
	for _, v := range list {
		println(v.ID, v.UserName, v.Age)
	}
}

func main() {
	dsn := "gorme:123456@tcp(127.0.0.1:3306)/gorme?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&ExampleModel{})
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Name%d", i)
		age := rand.Intn(40)
		db.Create(&ExampleModel{UserName: name, Age: age})
	}

	exampleRepo := NewExampleRepo(db)
	exampleRepo.GetFirst()
	exampleRepo.GetList()
	exampleRepo.GetPaginateList()

	u, _ := exampleRepo.First()
	fmt.Println(u)
}
