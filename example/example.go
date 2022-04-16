package main

import (
	"fmt"
	"github.com/micrease/gorme"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	gorm.Model
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

type Other struct {
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

func main() {

	dsn := "gorme:123456@tcp(127.0.0.1:3306)/gorme?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})
	for i := 0; i < 100; i++ {
		//code := fmt.Sprintf("D%d", i)
		//db.Create(&Product{Code: code, Price: uint(rand.Int31n(100))})
	}

	testObject(db)
	testList(db)
	testPaginate(db)
}

//获取单一对象
func testObject(db *gorm.DB) Product {
	//val类型为Product
	val, _ := gorme.Last[Product](db)
	fmt.Println("Last", val.ID, val.Price, val.Code)

	val, _ = gorme.First[Product](db)
	fmt.Println("First", val.ID, val.Price, val.Code)

	db = db.Where("id=?", 3)
	val, _ = gorme.Take[Product](db)
	fmt.Println("Take", val.ID, val.Price, val.Code)
	//GetOne = Take
	val, _ = gorme.GetOne[Product](db)
	fmt.Println("GetOne", val.ID, val.Price, val.Code)

	query2 := db.Model(&Product{}).Where("id=?", 2)
	val2, _ := gorme.GetOne[Other](query2)
	fmt.Println(val2)

	return val
}

//获取指定对象列表
func testList(db *gorm.DB) []*Product {
	query := db.Limit(6)
	query.Select("id,price").Where("id>?", 3).
		Order("price desc").
		Order("id desc")

	//list类型为[]Product
	list, _ := gorme.List[*Product](query)
	for _, v := range list {
		//v是一个Product类型
		v.Price = 33
		fmt.Println(v.ID, v.Price, v.Code)
	}
	return list
}

//获取指定对象分页结构
func testPaginate(db *gorm.DB) *gorme.PageResult[*Product] {
	//example1
	query := db.Where("code = ?", "D42") // 查找 code 字段值为 D42 的记录
	query.Where("id>?", 3).
		Where("id<?", 100).
		Order("id desc").
		Order("price desc")

	pageNo := 1
	pageSize := 5
	result, _ := gorme.Paginate[*Product](query, pageNo, pageSize)

	printPageResult(result)
	// SELECT * FROM `products` WHERE code = 'D42' AND id>3 AND id<100 AND `products`.`deleted_at` IS NULL ORDER BY id desc,price desc LIMIT 5
	// SELECT count(*) FROM `products` WHERE code = 'D42' AND id>3 AND id<100 AND `products`.`deleted_at` IS NULL LIMIT 5

	//example2
	type QueryReq struct {
		gorme.PageQuery
		MinPrice float32
		MaxPrice float32
		Id       int64
	}

	req := &QueryReq{}
	req.PageNo = 1
	req.PageSize = 7
	req.MinPrice = 10
	query = db.Where("price>?", req.MinPrice)
	result, _ = gorme.PaginateByOptions[*Product]([]gorme.Option{
		gorme.WithQuery(query),
		gorme.WithPage(req.PageQuery),
	}...)
	// SELECT * FROM `products` WHERE price>10.000000 AND `products`.`deleted_at` IS NULL LIMIT 7
	// SELECT count(*) FROM `products` WHERE price>10.000000 AND `products`.`deleted_at` IS NULL LIMIT 7

	printPageResult(result)
	return result
}

func printPageResult(result *gorme.PageResult[*Product]) {
	fmt.Println("totalSize", result.TotalSize)
	fmt.Println("totalPage", result.TotalPage)
	fmt.Println("PageNo", result.PageNo)
	fmt.Println("PageSize", result.PageSize)
	fmt.Println("ID", "Price", "Code")

	var products []*Product
	//result.List是[]Product类型
	products = result.List
	for _, v := range products {
		//v是一个Product类型
		fmt.Println(v.ID, v.Price, v.Code)
		v.Price = 44
	}
}
