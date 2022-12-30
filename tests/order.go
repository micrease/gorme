package tests

import (
	"github.com/micrease/gorme"
	"gorm.io/gorm"
)

// 这个一个例子
// model/example.go
type OrderModel struct {
	gorm.Model
	UserId    int64
	Amount    int
	GoodsName string
}

// 自定义表名
func (model OrderModel) TableName() string {
	return "tb_order"
}

// 实现Model接口中获取主键的方法
func (model OrderModel) GetID() any {
	return model.ID
}

// 举一个例子，ExampleRepo(可以换成你自己定义的Repo)继承gorme.Repository[T]
// repo/example.go
type OrderRepo struct {
	gorme.Repository[OrderModel]
}

func NewOrderRepo() *OrderRepo {
	repo := OrderRepo{}
	db := GetDB()
	repo.SetDB(db)
	return &repo
}
