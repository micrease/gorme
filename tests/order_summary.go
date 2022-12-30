package tests

import (
	"github.com/micrease/gorme"
)

// 这个一个例子
// model/example.go
type OrderSummaryModel struct {
	UserId   int64
	Username string
	Count    int
}

// 自定义表名
func (model OrderSummaryModel) TableName() string {
	return "tb_user as u left join tb_order as o on o.user_id=u.id"
}

// 实现Model接口中获取主键的方法
func (model OrderSummaryModel) GetID() any {
	return 0
}

// 举一个例子，ExampleRepo(可以换成你自己定义的Repo)继承gorme.Repository[T]
// repo/example.go
type OrderSummaryRepo struct {
	gorme.Repository[OrderSummaryModel]
}

func NewOrderSummaryRepo() *OrderSummaryRepo {
	repo := OrderSummaryRepo{}
	db := GetDB()
	repo.SetDB(db)
	return &repo
}
