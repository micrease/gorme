package main

import (
	"fmt"
	"github.com/micrease/gorme"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

func (u *UserRepo) GetFirst() (*UserModel, error) {
	//gorm原生查询
	//var m1 UserModel
	//err1:=u.DB.First(&m1).Error
	//gorme封装后的查询
	//m2,err2:= u.First()

	result, err := u.Select("age").Offset(1).Limit(1).Order("id desc").Where("id<?", 20).Where("age > ?", 1).First()
	fmt.Println(result.Age, result.UserName, err)
	return result, err
}

func (u *UserRepo) Paginate() (*gorme.PageResult[UserModel], error) {
	result, err := u.Select("age").Order("id desc").Where("id<?", 20).Where("age > ?", 1).Paginate(1, 2)
	fmt.Println(result.TotalSize, err)
	return result, err
}

func main() {
	dsn := "gorme:123456@tcp(127.0.0.1:3306)/gorme?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserModel{})
	for i := 0; i < 100; i++ {
		//name := fmt.Sprintf("Name%d", i)
		//age := rand.Intn(40)
		//db.Create(&UserModel{UserName: name, Age: age})
	}

	userRepo := NewUserRepo(db)
	userRepo.Paginate()
}
