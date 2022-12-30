package tests

import (
	"github.com/micrease/gorme"
	"gorm.io/gorm"
)

// 这个一个例子
// model/example.go
type UserModel struct {
	gorm.Model
	UserId int64
	Age    int
}

// 自定义表名
func (model UserModel) TableName() string {
	return "tb_user"
}

// 实现Model接口中获取主键的方法
func (model UserModel) GetID() any {
	return model.ID
}

// 举一个例子，ExampleRepo(可以换成你自己定义的Repo)继承gorme.Repository[T]
// repo/example.go
type UserRepo struct {
	gorme.Repository[UserModel]
}

func NewUserRepo() *UserRepo {
	repo := UserRepo{}
	db := GetDB()
	repo.SetDB(db)
	return &repo
}
