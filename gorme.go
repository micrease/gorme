package gorme

import (
	"gorm.io/gorm"
)

type PageQuery struct {
	PageSize int    `json:"page_size"` //每页条数
	PageNo   uint64 `json:"page_no"`   //当前页码
}

type PageResult[T any] struct {
	PageSize  int    `json:"page_size"`  //每页条数
	PageNo    uint64 `json:"page_no"`    //当前页码
	TotalPage uint64 `json:"total_page"` //总页数
	TotalSize int64  `json:"total_size"` //总条数
	List      []T    `json:"list"`       //数据列表
}

type QueryBuilder struct {
	query *gorm.DB
	PageQuery
}

type Option func(*QueryBuilder)

func WithPageSize(pageSize int) Option {
	return func(builder *QueryBuilder) {
		builder.PageSize = pageSize
	}
}

func WithPageNo(pageNo uint64) Option {
	return func(builder *QueryBuilder) {
		builder.PageNo = pageNo
	}
}

func WithPage(page PageQuery) Option {
	return func(builder *QueryBuilder) {
		builder.PageNo = page.PageNo
		builder.PageSize = page.PageSize
	}
}

func WithQuery(query *gorm.DB) Option {
	return func(builder *QueryBuilder) {
		builder.query = query
	}
}

func Paginate[T any](opts ...Option) (*PageResult[T], error) {
	builder := &QueryBuilder{}
	for _, o := range opts {
		o(builder)
	}
	return PaginateQuery[T](builder.query, builder.PageQuery)
}

func PaginateSimple[T any](query *gorm.DB, pageNo uint64, pageSize int) (*PageResult[T], error) {
	return PaginateQuery[T](query, PageQuery{
		PageNo:   pageNo,
		PageSize: pageSize,
	})
}

func PaginateQuery[T any](query *gorm.DB, page PageQuery) (*PageResult[T], error) {
	result := new(PageResult[T])
	result.PageNo = page.PageNo
	result.PageSize = page.PageSize
	offset := int(page.PageNo-1) * page.PageSize
	var rows []T
	err := query.Limit(page.PageSize).Offset(int(offset)).Find(&rows).Count(&result.TotalSize).Error
	if (int(result.TotalSize) % page.PageSize) > 0 {
		result.TotalPage = uint64(result.TotalSize)/uint64(page.PageSize) + 1
	} else {
		result.TotalPage = uint64(result.TotalSize) / uint64(page.PageSize)
	}

	result.List = rows
	return result, err
}

//查询列表
func List[T any](query *gorm.DB) ([]T, error) {
	var rows []T
	err := query.Find(&rows).Error
	return rows, err
}

//查询单个对象
func GetOne[T any](query *gorm.DB) (T, error) {
	var row T
	err := query.Limit(1).Find(&row).Error
	return row, err
}
