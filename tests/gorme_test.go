package tests

import (
	"fmt"
	"github.com/micrease/gorme"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"testing"
	"time"
)

// 这个一个例子
// model/example.go
type ExampleModel struct {
	gorm.Model
	UserName string
	Age      int
}

// 自定义表名
func (model ExampleModel) TableName() string {
	return "tb_example"
}

// 实现Model接口中获取主键的方法
func (model ExampleModel) GetID() any {
	return model.ID
}

// 举一个例子，ExampleRepo(可以换成你自己定义的Repo)继承gorme.Repository[T]
// repo/example.go
type ExampleRepo struct {
	gorme.Repository[ExampleModel]
}

func NewExampleRepo() *ExampleRepo {
	repo := ExampleRepo{}
	db := GetDB()
	repo.SetDB(db)
	return &repo
}

// 查询
func TestQuery(t *testing.T) {
	repo := NewExampleRepo()
	//链式where
	builder := repo.NewQueryBuilder().Where("id>?", 20).Where("id!=?", 23)

	//动态条件
	age := 19
	if age >= 18 {
		builder = builder.Where("age>=?", age)
	}

	//查询列表
	list, err := repo.QueryWithBuilder(builder).List(2)
	for _, row := range list {
		fmt.Println(row.ID, row.UserName, row.Age)
	}

	//查询一条数据
	builder = repo.NewQueryBuilder().Where("id=?", 20)
	model, err := repo.QueryWithBuilder(builder).First()
	fmt.Println(model.ID, model.UserName, err)

	//查询分页,第1页，每页10条
	builder = repo.NewQueryBuilder().Where("id>?", 20).Order("id desc")
	page, err := repo.QueryWithBuilder(builder).Paginate(1, 10)
	fmt.Println(page.PageNo, page.PageSize, page.TotalPage, page.TotalSize)
	for _, row := range page.List {
		fmt.Println(row.ID, row.UserName, row.Age)
	}
}

// 查询单值集合
func TestValues(t *testing.T) {
	repo := NewExampleRepo()
	//链式where
	builder := repo.NewQueryBuilder().Limit(10).Where("id>?", 20).Where("id!=?", 23)
	//查询所有值
	// SELECT `age` FROM `tb_example` WHERE id>20 AND id!=23 AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	ages, err := repo.QueryWithBuilder(builder).Values("age")
	fmt.Println(ages, err)

	//查询时去重
	//SELECT DISTINCT `age` FROM `tb_example` WHERE id>20 AND id!=23 AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	ages, err = repo.QueryWithBuilder(builder).DistinctValues("age")
	fmt.Println(ages, err)

	//不使用builder
	//SELECT DISTINCT `age` FROM `tb_example` WHERE age>20 AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	ages, err = repo.NewQuery().Limit(10).WhereRaw("age>?", 20).Distinct("age").Values()
	fmt.Println(ages, err)
}

// Where条件查询
func TestWhere(t *testing.T) {
	repo := NewExampleRepo()
	//查询所有值
	//SELECT * FROM `tb_example` WHERE age=20 AND `tb_example`.`deleted_at` IS NULL ORDER BY `tb_example`.`id` LIMIT 1
	row, err := repo.NewQuery().Limit(1).Where("age", 20).First()
	fmt.Println(row, err)

	//SELECT * FROM `tb_example` WHERE age=20 AND `tb_example`.`deleted_at` IS NULL ORDER BY `tb_example`.`id` LIMIT 1
	row, err = repo.NewQuery().Limit(1).Where("age", "=", 20).First()
	fmt.Println(row, err)

	//SELECT * FROM `tb_example` WHERE age>20 AND `tb_example`.`deleted_at` IS NULL ORDER BY `tb_example`.`id` LIMIT 1
	row, err = repo.NewQuery().Limit(1).Where("age", ">", 20).First()
	fmt.Println(row, err)

	//SELECT * FROM `tb_example` WHERE age IN('21','23') AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := repo.NewQuery().Where("age", "in", "21,23").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE age IN(21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Where("age", "in", 21).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE age IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Where("age", "in", []any{20, 21}).List(2)

	//SELECT * FROM `tb_example` WHERE age IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().WhereIn("age", []any{20, 21}).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE age IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().WhereIn("age", []any{20, 21}).List(2)
	fmt.Println(rows, err)
}

func TestWhereFunc(t *testing.T) {
	repo := NewExampleRepo()
	query := repo.NewQuery()
	// SELECT * FROM `tb_example` WHERE age IN(20,21) AND age<10 AND (age=20 AND age=23 AND (age=1 AND age=2)) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := query.WhereIn("age", []any{20, 21}).Lt("age", 10).Where(func() {
		query.Where("age", 20).Where("age=?", 23).Where(func() {
			query.Where("age", 1).Where("age", 2)
		})
	}).List(2)
	fmt.Println(rows, err)
}

func TestFindInSet(t *testing.T) {
	repo := NewExampleRepo()
	query := repo.NewQuery()
	//SELECT * FROM `tb_example` WHERE FIND_IN_SET(age,'20') AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := query.FindInSet("user_name", "N2").List(2)
	fmt.Println(rows, err)
}

func TestOrFunc(t *testing.T) {
	repo := NewExampleRepo()
	query := repo.NewQuery()
	// SELECT * FROM `tb_example` WHERE age IN(20,21)  AND age >10  AND (age =20  OR age=23 OR (age=1 AND age=2)) AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	pageList, err := query.WhereIn("age", []any{20, 21}).
		Gt("age", 10).
		Where(func() {
			query.Eq("age", 20).
				Or("age=?", 23).Or(func() {
				query.Where("age", 1).Where("age", 2)
			})
		}).Paginate(1, 10)
	fmt.Println(pageList, err)
}

func TestQuery2(t *testing.T) {
	repo := NewExampleRepo()
	//SELECT * FROM `tb_example` WHERE age IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := repo.NewQuery().In("age", []any{20, 21}).List(2)
	fmt.Println(rows, err)

	// SELECT * FROM `tb_example` WHERE age Not IN('20','21')  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().NotIn("age", "20,21").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE age = 20  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Eq("age", 20).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE age > 20  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Gt("age", 20).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE user_name like '%name%'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Like("user_name", "name").List(2)
	fmt.Println(rows, err)

	// SELECT * FROM `tb_example` WHERE user_name like '%name'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().LikeLeft("user_name", "name").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE user_name not like '%name%'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().NotLike("user_name", "name").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE user_name IS Not NULL  AND user_name not like '%name%'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().IsNotNull("user_name").NotLike("user_name", "name").List(2)
	fmt.Println(rows, err)
}

// 插入或更新,ID存在就更新，不存在就插入,用Save()方法
func TestInsertOrUpdate(t *testing.T) {
	repo := NewExampleRepo()
	insertModel := repo.NewModelPtr()
	insertModel.Age = 12
	insertModel.UserName = "test1"
	insertModel.CreatedAt = time.Now()
	err := repo.Save(insertModel).Error
	fmt.Println(insertModel.ID, err)

	updateModel := repo.NewModelPtr()
	updateModel.Age = 12
	updateModel.ID = 1
	updateModel.UserName = "test1"
	updateModel.CreatedAt = time.Now()
	err = repo.Save(updateModel).Error
	fmt.Println(updateModel.ID, err)
}

// 插入数据Create
func TestInsert(t *testing.T) {
	repo := NewExampleRepo()
	model := repo.NewModelPtr()
	model.Age = 12
	//model.ID = 1
	model.UserName = "test1"
	model.CreatedAt = time.Now()
	err := repo.Create(model).Error
	fmt.Println(model.ID, err)
}

// 更新，建议用Save方法
func TestUpdate(t *testing.T) {
	repo := NewExampleRepo()
	//常规Updates更新方式,0值不会更新
	model := repo.NewModelPtr()
	model.ID = 1
	model.Age = 0       //0值不会更新
	model.UserName = "" //空值不会更新
	model.CreatedAt = time.Now()
	err := repo.Updates(model).Error
	fmt.Println(model.ID, err)

	//更新0值方式一,使用Updates方法+Select字段名
	repo.Select("age").Updates(model)
	//更新0值方式二,Save方法会更新所有字段,不限0值
	repo.Save(model)
	//Save方法指定Select字段时，等价于Updates
	repo.Select("age").Save(model)

	fmt.Println(model.ID, err)
}

// 根据查询条件,更新单个字段
func TestUpdateSingleColumnWithQueryBuilder(t *testing.T) {
	repo := NewExampleRepo()
	builder := repo.NewQueryBuilder().Where("id=?", 1)
	err := repo.QueryWithBuilder(builder).Update("age", 11).Error
	fmt.Println(err)
}

// 根据查询条件,更新Model中的多个字段
func TestUpdateModelColumnWithQueryBuilder(t *testing.T) {
	repo := NewExampleRepo()
	builder := repo.NewQueryBuilder().Where("id=?", 1)

	model := repo.NewModelPtr()
	model.Age = 100            //0值不会更新
	model.UserName = "test100" //空值不会更新
	model.CreatedAt = time.Now()
	err := repo.QueryWithBuilder(builder).Select("age", "user_name").Updates(model)
	fmt.Println(err)
	//UPDATE `example_models` SET `updated_at`='2022-07-14 17:00:28.43',`user_name`='test100',`age`=100 WHERE id=1 AND `example_models`.`deleted_at` IS NULL
}

// 根据查询条件,更新动态多个字段
func TestUpdateBySetterWithQueryBuilder(t *testing.T) {
	repo := NewExampleRepo()
	builder := repo.NewQueryBuilder().Where("id=?", 1)
	setter := repo.NewSetter().Set("age", 12).Set("user_name", "gggg")
	err := repo.QueryWithBuilder(builder).Updates(setter).Error
	fmt.Println(err)
}

// 根据查询条件,删除
func TestDeleteByQueryBuilder(t *testing.T) {
	repo := NewExampleRepo()
	builder := repo.NewQueryBuilder().Where("id=?", 2)
	err := repo.QueryWithBuilder(builder).Delete().Error
	fmt.Println(err)
}

// 根据ID删除
func TestDeleteID(t *testing.T) {
	repo := NewExampleRepo()
	//单个ID
	err := repo.Delete(2).Error
	fmt.Println(err)

	//多个ID
	err = repo.Delete([]int{2, 3}).Error
	fmt.Println(err)
}

// 根据查询条件,软删除
func TestSoftDeleteByQueryBuilder(t *testing.T) {
	repo := NewExampleRepo()
	builder := repo.NewQueryBuilder().Where("id=?", 2)
	err := repo.QueryWithBuilder(builder).DeleteSoft().Error
	fmt.Println(err)
}

func GetDB() *gorm.DB {
	dsn := "gorme:123456@tcp(127.0.0.1:3306)/gorme?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalln("连接数据库失败")
	}
	return db
}
