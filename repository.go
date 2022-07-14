package gorme

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model interface {
	TableName() string
	GetID() uint
}

type Repository[T Model] struct {
	//一个初始干净的db
	_initDB *gorm.DB
	//如果继承的方式，无法形成链式结构x().y().Paginate(),因为x(),y()返回的是gorm.DB,而这个对象不具有Repository中的方法
	DB *gorm.DB
	//在增删改查时，存放的待处理数据
	Data map[string]any
}

type Setter struct {
	Data map[string]any
}

func (s Setter) Set(key string, value any) Setter {
	if s.Data == nil {
		s.Data = map[string]any{}
	}
	s.Data[key] = value
	return s
}

func (r *Repository[T]) SetDB(db *gorm.DB) *Repository[T] {
	r.DB = db
	r._initDB = db
	return r
}

func (r *Repository[T]) Reset() *Repository[T] {
	r.DB = r._initDB
	return r
}

func (r *Repository[T]) First() (T, error) {
	var t T
	err := r.DB.First(&t).Error
	//把DB初始化
	r.Reset()
	return t, err
}

func (r *Repository[T]) Last() (T, error) {
	var t T
	err := r.DB.Last(&t).Error
	//把DB初始化
	r.Reset()
	return t, err
}

func (r *Repository[T]) GetOne() (T, error) {
	return r.Take()
}

func (r *Repository[T]) Take() (T, error) {
	var t T
	err := r.DB.Take(&t).Error
	//把DB初始化
	r.Reset()
	return t, err
}

func (r *Repository[T]) List(args ...int) ([]T, error) {
	if len(args) > 1 {
		panic("the number of args cannot exceed 1")
	}

	var t []T
	if len(args) == 1 {
		limit := args[0]
		r.DB = r.DB.Limit(limit)
	}
	err := r.DB.Find(&t).Error
	//DB初始化
	r.Reset()
	return t, err
}

func (r *Repository[T]) Paginate(pageNo int, pageSize int) (*PageResult[T], error) {
	result, err := Paginate[T](r.DB, pageNo, pageSize)
	//把DB初始化
	r.Reset()
	return result, err
}

//======================================Query Builder=====================================
func (r *Repository[T]) NewQueryBuilder() *gorm.DB {
	r.Reset()
	var t T
	r.DB = r.DB.Model(&t)
	return r.DB
}

func (r *Repository[T]) QueryWithBuilder(builder *gorm.DB) *Repository[T] {
	r.DB = builder
	return r
}

func (r *Repository[T]) NewModel() T {
	var t T
	return t
}

func (r *Repository[T]) NewSetter() Setter {
	return Setter{}
}

func (r *Repository[T]) NewModelPtr() *T {
	var t = new(T)
	return t
}

//======================================最后调用的方法返回*gorm.DB,这样获取结果中的信息更方便一些=====================================

func (r *Repository[T]) Create(value interface{}) *gorm.DB {
	tx := r.DB.Create(value)
	r.Reset()
	return tx
}

func (r *Repository[T]) Save(value interface{}) *gorm.DB {
	tx := r.DB.Save(value)
	r.Reset()
	return tx
}

func (r *Repository[T]) Updates(values interface{}) *gorm.DB {
	var tx *gorm.DB
	if setter, ok := values.(Setter); ok {
		tx = r.DB.Updates(setter.Data)
	} else {
		tx = r.DB.Updates(values)
	}
	r.Reset()
	return tx
}

func (r *Repository[T]) Update(column string, value interface{}) *gorm.DB {
	var t T
	tx := r.DB.Update(column, value).Model(&t)
	r.Reset()
	return tx
}

func (r *Repository[T]) UpdateColumn(column string, value interface{}) *gorm.DB {
	var t T
	tx := r.DB.Model(&t).UpdateColumn(column, value)
	r.Reset()
	return tx
}

func (r *Repository[T]) UpdateColumns(values interface{}) *gorm.DB {
	tx := r.DB.UpdateColumns(values)
	r.Reset()
	return tx
}

func (r *Repository[T]) Delete(conds ...interface{}) *gorm.DB {
	var t T
	tx := r.DB.Unscoped().Delete(&t, conds...)
	r.Reset()
	return tx
}

//软删除,前提是有 Deleted gorm.DeletedAt
func (r *Repository[T]) DeleteSoft(conds ...interface{}) *gorm.DB {
	var t T
	tx := r.DB.Delete(&t, conds...)
	r.Reset()
	return tx
}

func (r *Repository[T]) Begin(opts ...*sql.TxOptions) *gorm.DB {
	tx := r.DB.Begin(opts...)
	r.Reset()
	return tx
}

func (r *Repository[T]) Commit() *gorm.DB {
	tx := r.DB.Commit()
	r.Reset()
	return tx
}

func (r *Repository[T]) Rollback() *gorm.DB {
	tx := r.DB.Rollback()
	r.Reset()
	return tx
}

func (r *Repository[T]) Scan(dest interface{}) *gorm.DB {
	tx := r.DB.Scan(dest)
	r.Reset()
	return tx
}

func (r *Repository[T]) ScanRows(rows *sql.Rows, dest interface{}) error {
	err := r.DB.ScanRows(rows, dest)
	r.Reset()
	return err
}

func (r *Repository[T]) Exec(sql string, values ...interface{}) *gorm.DB {
	tx := r.DB.Exec(sql, values...)
	r.Reset()
	return tx
}

func (r *Repository[T]) Raw(sql string, values ...interface{}) *gorm.DB {
	tx := r.DB.Raw(sql, values...)
	r.Reset()
	return tx
}

func (r *Repository[T]) Row() *sql.Row {
	row := r.DB.Row()
	r.Reset()
	return row
}

func (r *Repository[T]) Rows() (*sql.Rows, error) {
	rows, err := r.DB.Rows()
	r.Reset()
	return rows, err
}

//-------------------以下对DB原生方法套壳-------------------------

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

func (r *Repository[T]) Group(name string) *Repository[T] {
	r.DB = r.DB.Group(name)
	return r
}

func (r *Repository[T]) Having(query interface{}, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Having(query, args...)
	return r
}

func (r *Repository[T]) Debug() *Repository[T] {
	r.DB = r.DB.Debug()
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

func (r *Repository[T]) Unscoped() *Repository[T] {
	r.DB = r.DB.Unscoped()
	return r
}

func (r *Repository[T]) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *Repository[T] {
	r.DB = r.DB.Scopes(funcs...)
	return r
}
