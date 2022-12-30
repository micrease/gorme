package tests

import (
	"fmt"
	"testing"
	"time"
)

// 查询
func TestQuery(t *testing.T) {
	repo := NewOrderRepo()
	//链式where
	query := repo.NewQuery().Where("id > ?", 20).Where("id != ?", 23)
	//动态条件
	amount := 1
	if amount >= 18 {
		query = query.Where("amount>=?", amount)
	}

	//查询列表
	list, err := query.List(2)
	fmt.Println(err)
	for _, row := range list {
		fmt.Println(row.ID, row.UserId, row.Amount)
	}
}

// 查询
func TestQueryPamount(t *testing.T) {
	repo := NewOrderRepo()
	page, err := repo.NewQuery().Where("id", ">", 20).Order("id desc").Paginate(1, 10)
	fmt.Println(page)
	for _, row := range page.List {
		fmt.Println(row.ID, row.UserId, row.Amount, err)
	}
}

// 查询单值集合
func TestValues(t *testing.T) {
	repo := NewOrderRepo()
	//链式where
	query := repo.NewQuery().Limit(10).Where("id>?", 20).Where("id!=?", 23)
	//查询所有值
	// SELECT `amount` FROM `tb_example` WHERE id>20 AND id!=23 AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	amounts, err := query.Values("amount")
	fmt.Println(amounts, err)

	//查询时去重
	//SELECT DISTINCT `amount` FROM `tb_example` WHERE id>20 AND id!=23 AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	amounts, err = query.Limit(10).WhereRaw("amount>?", 20).DistinctValues("amount")
	fmt.Println(amounts, err)

	//不使用builder
	//SELECT DISTINCT `amount` FROM `tb_example` WHERE amount>20 AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	amounts, err = repo.NewQuery().Limit(10).WhereRaw("amount>?", 20).Distinct("amount").Values()
	fmt.Println(amounts, err)
}

// Where条件查询
func TestWhere(t *testing.T) {
	repo := NewOrderRepo()
	//查询所有值
	//SELECT * FROM `tb_example` WHERE amount=20 AND `tb_example`.`deleted_at` IS NULL ORDER BY `tb_example`.`id` LIMIT 1
	row, err := repo.NewQuery().Limit(1).Where("amount", 20).First()
	fmt.Println(row, err)

	//SELECT * FROM `tb_example` WHERE amount=20 AND `tb_example`.`deleted_at` IS NULL ORDER BY `tb_example`.`id` LIMIT 1
	row, err = repo.NewQuery().Limit(1).Where("amount", "=", 20).First()
	fmt.Println(row, err)

	//SELECT * FROM `tb_example` WHERE amount>20 AND `tb_example`.`deleted_at` IS NULL ORDER BY `tb_example`.`id` LIMIT 1
	row, err = repo.NewQuery().Limit(1).Where("amount", ">", 20).First()
	fmt.Println(row, err)

	//SELECT * FROM `tb_example` WHERE amount IN('21','23') AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := repo.NewQuery().Where("amount", "in", "21,23").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE amount IN(21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Where("amount", "in", 21).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE amount IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Where("amount", "in", []any{20, 21}).List(2)

	//SELECT * FROM `tb_example` WHERE amount IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().WhereIn("amount", []any{20, 21}).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE amount IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().WhereIn("amount", []any{20, 21}).List(2)
	fmt.Println(rows, err)
}

func TestWhereFunc(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery()
	// SELECT * FROM `tb_example` WHERE amount IN(20,21) AND amount<10 AND (amount=20 AND amount=23 AND (amount=1 AND amount=2)) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := query.WhereIn("amount", []any{20, 21}).Lt("amount", 10).Where(func() {
		query.Where("amount", 20).Where("amount=?", 23).Where(func() {
			query.Where("amount", 1).Where("amount", 2)
		})
	}).List(2)
	fmt.Println(rows, err)
}

func TestFindInSet(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery()
	//SELECT * FROM `tb_example` WHERE FIND_IN_SET(amount,'20') AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := query.FindInSet("user_name", "N2").List(2)
	fmt.Println(rows, err)
}

func TestOrFunc(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery()
	// SELECT * FROM `tb_example` WHERE amount IN(20,21)  AND amount >10  AND (amount =20  OR amount=23 OR (amount=1 AND amount=2)) AND `tb_example`.`deleted_at` IS NULL LIMIT 10
	page, err := query.WhereIn("amount", []any{20, 21}).Gt("amount", 10).Where(func() {
		query.Eq("amount", 20).Or("amount=?", 23).Or(func() {
			query.Where("amount", 1).Where("amount", 2)
		})
	}).Paginate(1, 10)

	fmt.Println(page, err)
	for _, item := range page.List {
		fmt.Println(item)
	}
}

func TestQuery2(t *testing.T) {
	repo := NewOrderRepo()
	//SELECT * FROM `tb_example` WHERE amount IN(20,21) AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err := repo.NewQuery().In("amount", []any{20, 21}).List(2)
	fmt.Println(rows, err)

	// SELECT * FROM `tb_example` WHERE amount Not IN('20','21')  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().NotIn("amount", "20,21").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE amount = 20  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Eq("amount", 20).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE amount > 20  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Gt("amount", 20).List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE user_name like '%name%'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().Like("goods_name", "goods_name").List(2)
	fmt.Println(rows, err)

	// SELECT * FROM `tb_example` WHERE user_name like '%name'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().LikeLeft("goods_name", "name").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE user_name not like '%name%'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().NotLike("goods_name", "name").List(2)
	fmt.Println(rows, err)

	//SELECT * FROM `tb_example` WHERE user_name IS Not NULL  AND user_name not like '%name%'  AND `tb_example`.`deleted_at` IS NULL LIMIT 2
	rows, err = repo.NewQuery().IsNotNull("goods_name").NotLike("goods_name", "name").List(2)
	fmt.Println(rows, err)
}

// 插入或更新,ID存在就更新，不存在就插入,用Save()方法
func TestInsertOrUpdate(t *testing.T) {
	repo := NewOrderRepo()
	insertModel := repo.NewModel()
	insertModel.ID = 0
	insertModel.Amount = 12
	insertModel.UserId = 100
	insertModel.CreatedAt = time.Now()
	err := repo.Save(insertModel).Error
	fmt.Println(insertModel.ID, err)

	updateModel := repo.NewModel()
	updateModel.Amount = 12
	updateModel.ID = 400
	updateModel.UserId = 100
	updateModel.CreatedAt = time.Now()
	err = repo.Save(&updateModel).Error
	fmt.Println(updateModel.ID, err)
}

// 插入数据Create
func TestInsert(t *testing.T) {
	repo := NewOrderRepo()
	model := repo.NewModel()
	model.Amount = 12
	//model.ID = 1
	model.UserId = 100
	model.CreatedAt = time.Now()
	err := repo.Create(model).Error
	fmt.Println(model.ID, err)
}

// 更新，建议用Save方法
func TestUpdate(t *testing.T) {
	repo := NewOrderRepo()
	//常规Updates更新方式,0值不会更新
	model := repo.NewModel()
	model.ID = 1
	model.Amount = 0 //0值不会更新
	model.UserId = 0 //空值不会更新
	model.CreatedAt = time.Now()
	err := repo.Updates(model).Error
	fmt.Println(model.ID, err)

	//更新0值方式一,使用Updates方法+Select字段名
	repo.Select("amount").Updates(model)
	//更新0值方式二,Save方法会更新所有字段,不限0值
	repo.Save(model)
	//Save方法指定Select字段时，等价于Updates
	repo.Select("amount").Save(model)

	fmt.Println(model.ID, err)
}

// 根据查询条件,更新单个字段
// UPDATE `tb_order` SET `amount`=11,`updated_at`='2022-12-30 14:43:01.11' WHERE id=1 AND `tb_order`.`deleted_at` IS NULL
func TestUpdateSingleColumnWithQueryBuilder(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery().Where("id=?", 1)
	err := query.Update("amount", 11).Error
	fmt.Println(err)
}

// 根据查询条件,更新Model中的多个字段
// UPDATE `tb_order` SET `updated_at`='2022-12-30 14:44:39.426',`amount`=100 WHERE id=1 AND `tb_order`.`deleted_at` IS NULL
func TestUpdateModelColumnWithQueryBuilder(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery().Where("id=?", 1)

	model := repo.NewModel()
	model.Amount = 100 //0值不会更新
	model.UserId = 100 //空值不会更新
	model.CreatedAt = time.Now()
	err := query.Select("amount", "user_name").Updates(model)
	fmt.Println(err)
	//UPDATE `example_models` SET `updated_at`='2022-07-14 17:00:28.43',`user_name`='test100',`amount`=100 WHERE id=1 AND `example_models`.`deleted_at` IS NULL
}

// 根据查询条件,更新动态多个字段
// UPDATE `tb_order` SET `amount`=12,`goods_name`='gggg',`updated_at`='2022-12-30 14:46:25.922' WHERE id=1 AND `tb_order`.`deleted_at` IS NULL
func TestUpdateBySetterWithQueryBuilder(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery().Where("id=?", 1)
	setter := query.NewSetter().Set("amount", 12).Set("goods_name", "gggg")
	err := query.Updates(setter).Error
	fmt.Println(err)
}

// 根据查询条件,删除
// DELETE FROM `tb_order` WHERE id=2
func TestDeleteByQueryBuilder(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery().Where("id=?", 2)
	err := query.Delete().Error
	fmt.Println(err)
}

// 根据ID删除
func TestDeleteID(t *testing.T) {
	repo := NewOrderRepo()
	//单个ID
	err := repo.Delete(2).Error
	fmt.Println(err)

	//多个ID
	err = repo.Delete([]int{2, 3}).Error
	fmt.Println(err)
}

// 根据查询条件,软删除
// UPDATE `tb_order` SET `deleted_at`='2022-12-30 14:47:55.504' WHERE id=2 AND `tb_order`.`deleted_at` IS NULL
func TestSoftDeleteByQueryBuilder(t *testing.T) {
	repo := NewOrderRepo()
	query := repo.NewQuery().Where("id=?", 2)
	err := query.DeleteSoft().Error
	fmt.Println(err)
}

func TestLeftJoin(t *testing.T) {
	list, err := NewOrderSummaryRepo().NewQuery().Select("u.id as user_id,u.user_name as username").Where("u.id=?", 10).List(10)
	fmt.Println(list, err)
}

// SELECT u.id as user_id,u.user_name as username,count(*) as count FROM tb_user as u left join tb_order as o on o.user_id=u.id WHERE u.id=10 GROUP BY `u`.`id` HAVING count<10 LIMIT 10
func TestSummary(t *testing.T) {
	list, err := NewOrderSummaryRepo().NewQuery().
		Select("u.id as user_id,u.user_name as username,count(*) as count").Where("u.id=?", 10).
		Group("u.id").Having("count<10").
		List(10)
	fmt.Println(list, err)
}
