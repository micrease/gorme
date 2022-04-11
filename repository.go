package gorme

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] struct {
	//如果继承的方式，无法形成链式结构x().y().Paginate(),因为x(),y()返回的是gorm.DB,而这个对象不具有Repository中的方法
	DB *gorm.DB
}

func (r *Repository[T]) First() (*T, error) {
	var t T
	err := r.DB.First(&t).Error
	return &t, err
}

func (r *Repository[T]) List() ([]T, error) {
	var t []T
	err := r.DB.Find(&t).Error
	return t, err
}

func (r *Repository[T]) Paginate(pageNo uint64, pageSize int) (*PageResult[T], error) {
	return PaginateSimple[T](r.DB, pageNo, pageSize)
}

//=========================================以下对DB原生方法套壳==================================================

func (r *Repository[T]) Where(query interface{}, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Where(query, args...)
	return r
}

func (r *Repository[T]) Order(value interface{}) *Repository[T] {
	r.DB = r.DB.Order(value)
	return r
}

func (r *Repository[T]) Model(value interface{}) *Repository[T] {
	r.DB = r.DB.Model(value)
	return r
}

func (r *Repository[T]) Create(value interface{}) *Repository[T] {
	r.DB = r.DB.Create(value)
	return r
}

func (r *Repository[T]) Or(value interface{}) *Repository[T] {
	r.DB = r.DB.Or(value)
	return r
}

func (r *Repository[T]) Limit(limit int) *Repository[T] {
	r.DB = r.DB.Limit(limit)
	return r
}

func (r *Repository[T]) Distinct(args ...interface{}) *Repository[T] {
	r.DB = r.DB.Distinct(args...)
	return r
}

func (r *Repository[T]) Offset(offset int) *Repository[T] {
	r.DB = r.DB.Offset(offset)
	return r
}

func (r *Repository[T]) Select(query interface{}, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Select(query, args...)
	return r
}

func (r *Repository[T]) Attrs(attrs ...interface{}) *Repository[T] {
	r.DB = r.DB.Attrs(attrs...)
	return r
}

func (r *Repository[T]) Joins(query string, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Joins(query, args...)
	return r
}

func (r *Repository[T]) Save(value interface{}) *Repository[T] {
	r.DB = r.DB.Save(value)
	return r
}

func (r *Repository[T]) Updates(values interface{}) *Repository[T] {
	r.DB = r.DB.Updates(values)
	return r
}

func (r *Repository[T]) Update(column string, value interface{}) *Repository[T] {
	r.DB = r.DB.Update(column, value)
	return r
}

func (r *Repository[T]) UpdateColumn(column string, value interface{}) *Repository[T] {
	r.DB = r.DB.UpdateColumn(column, value)
	return r
}

func (r *Repository[T]) UpdateColumns(values interface{}) *Repository[T] {
	r.DB = r.DB.UpdateColumns(values)
	return r
}

func (r *Repository[T]) Group(name string) *Repository[T] {
	r.DB = r.DB.Group(name)
	return r
}

func (r *Repository[T]) Having(query interface{}, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Having(query, args...)
	return r
}

func (r *Repository[T]) Delete(value interface{}, conds ...interface{}) *Repository[T] {
	r.DB = r.DB.Delete(value, conds...)
	return r
}

func (r *Repository[T]) Debug() *Repository[T] {
	r.DB = r.DB.Debug()
	return r
}

func (r *Repository[T]) Begin(opts ...*sql.TxOptions) *Repository[T] {
	r.DB = r.DB.Begin(opts...)
	return r
}

func (r *Repository[T]) Commit() *Repository[T] {
	r.DB = r.DB.Commit()
	return r
}

func (r *Repository[T]) Rollback() *Repository[T] {
	r.DB = r.DB.Rollback()
	return r
}

func (r *Repository[T]) Assign(attrs ...interface{}) *Repository[T] {
	r.DB = r.DB.Assign(attrs...)
	return r
}

func (r *Repository[T]) Clauses(conds ...clause.Expression) *Repository[T] {
	r.DB = r.DB.Clauses(conds...)
	return r
}

func (r *Repository[T]) Count(count *int64) *Repository[T] {
	r.DB = r.DB.Count(count)
	return r
}

func (r *Repository[T]) FirstOrCreate(dest interface{}, conds ...interface{}) *Repository[T] {
	r.DB = r.DB.FirstOrCreate(dest, conds...)
	return r
}

func (r *Repository[T]) FirstOrInit(dest interface{}, conds ...interface{}) *Repository[T] {
	r.DB = r.DB.FirstOrInit(dest, conds...)
	return r
}

func (r *Repository[T]) Pluck(column string, dest interface{}) *Repository[T] {
	r.DB = r.DB.Pluck(column, dest)
	return r
}

func (r *Repository[T]) Table(name string, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Table(name, args...)
	return r
}

func (r *Repository[T]) Session(config *gorm.Session) *Repository[T] {
	r.DB = r.DB.Session(config)
	return r
}

func (r *Repository[T]) Preload(query string, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Preload(query, args...)
	return r
}

func (r *Repository[T]) Omit(columns ...string) *Repository[T] {
	r.DB = r.DB.Omit(columns...)
	return r
}

func (r *Repository[T]) Not(query interface{}, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Not(query, args...)
	return r
}

func (r *Repository[T]) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.DB.Transaction(fc, opts...)
}

func (r *Repository[T]) Raw(sql string, values ...interface{}) *Repository[T] {
	r.DB = r.DB.Raw(sql, values...)
	return r
}

func (r *Repository[T]) Unscoped() *Repository[T] {
	r.DB = r.DB.Unscoped()
	return r
}
