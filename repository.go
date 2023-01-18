package gorme

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Model interface {
	TableName() string
	GetID() any
}

type Repository[T Model] struct {
	//一个初始干净的db
	DB *gorm.DB
	//如果继承的方式，无法形成链式结构x().y().Paginate(),因为x(),y()返回的是gorm.DB,而这个对象不具有Repository中的方法
	Query *sql.DB
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
	return r
}

func (r *Repository[T]) Reset() *Repository[T] {
	tmp := r.DB.Statement
	oldStatement := *tmp
	r.DB.Statement = &oldStatement
	r.DB.Statement.Clauses = map[string]clause.Clause{}
	r.DB.Statement.Preloads = map[string][]interface{}{}
	return r
}

func (r *Repository[T]) First() (T, error) {
	var t T
	err := r.DB.First(&t).Error
	//把DB初始化
	r.Reset()
	return t, r.IgnoreError(err)
}

func (r *Repository[T]) IgnoreError(err error) error {
	if err == nil {
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (r *Repository[T]) Last() (T, error) {
	var t T
	err := r.DB.Last(&t).Error
	//把DB初始化
	r.Reset()
	return t, r.IgnoreError(err)
}

func (r *Repository[T]) GetOne() (T, error) {
	return r.Take()
}

func (r *Repository[T]) Take() (T, error) {
	var t T
	err := r.DB.Take(&t).Error
	//把DB初始化
	r.Reset()
	return t, r.IgnoreError(err)
}

func (r *Repository[T]) Values(column ...string) ([]any, error) {
	var pluckColumn string
	if len(column) > 0 {
		pluckColumn = column[0]
	}

	var values []any
	err := r.DB.Pluck(pluckColumn, &values).Error
	//把DB初始化
	r.Reset()
	return values, r.IgnoreError(err)
}

func (r *Repository[T]) DistinctValues(column string) ([]any, error) {
	var values []any
	err := r.DB.Distinct(column).Pluck(column, &values).Error
	//把DB初始化
	r.Reset()
	return values, r.IgnoreError(err)
}

func (r *Repository[T]) Pluck(column string) ([]any, error) {
	var values []any
	err := r.DB.Pluck(column, &values).Error
	r.Reset()
	return values, r.IgnoreError(err)
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
	return t, r.IgnoreError(err)
}

func (r *Repository[T]) Paginate(pageNo int, pageSize int) (*PageResult[T], error) {
	result, err := Paginate[T](r.DB, pageNo, pageSize)
	//把DB初始化
	r.Reset()
	return result, r.IgnoreError(err)
}

// ======================================Query Builder=====================================
func (r *Repository[T]) NewQueryBuilder() *gorm.DB {
	var t T
	r.DB.Statement.Clauses = map[string]clause.Clause{}
	r.DB = r.DB.Model(&t).Table(t.TableName())
	return r.DB
}

func (r *Repository[T]) NewQuery() *Repository[T] {
	var t T
	r.DB.Statement.Clauses = map[string]clause.Clause{}
	r.DB = r.DB.Model(&t).Table(t.TableName())
	return r
}

func (r *Repository[T]) QueryWithBuilder(builder *gorm.DB) *Repository[T] {
	r.DB = builder
	return r
}

func (r *Repository[T]) NewModelValue() T {
	var t T
	return t
}

func (r *Repository[T]) NewSetter() Setter {
	return Setter{}
}

func (r *Repository[T]) NewModel() *T {
	var t = new(T)
	return t
}

//======================================最后调用的方法返回*gorm.DB,这样获取结果中的信息更方便一些=====================================

func (r *Repository[T]) Create(value interface{}) *gorm.DB {
	tx := r.DB.Create(value)
	return tx
}

func (r *Repository[T]) Save(value interface{}) *gorm.DB {
	tx := r.DB.Save(value)
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

// 软删除,前提是有 Deleted gorm.DeletedAt
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

// -------------------以下Where查询方式-------------------------
func (r *Repository[T]) Or(query any, args ...interface{}) *Repository[T] {
	return r.OrWhere(query, args...)
}

func (r *Repository[T]) OrWhere(query any, args ...interface{}) *Repository[T] {
	switch query.(type) {
	case string:
		argsLen := len(args)
		if argsLen == 0 {
			return r.OrRaw(query, args...)
		}

		queryStr, _ := query.(string)
		if strings.Contains(queryStr, "?") {
			return r.OrRaw(queryStr, args...)
		}

		expr := "="
		var value any
		if argsLen == 1 {
			value = args[0]
		} else if argsLen > 1 {
			value = args[1]
			expr = fmt.Sprintf("%v", args[0])
			if strings.ToUpper(expr) == "IN" {
				return r.OrWhereIn(queryStr, args[1])
			}
		}
		queryExpr := queryStr + expr + "?"
		r.DB = r.DB.Or(queryExpr, value)
	case func():
		f, _ := query.(func())
		tmp := *r.DB
		oldDB := &tmp
		r.Reset()
		f()
		r.DB = oldDB.Or(r.DB)
	}
	return r
}

func (r *Repository[T]) Case(isTrue bool, handleFunc func()) *Repository[T] {
	if isTrue {
		handleFunc()
	}
	return r
}

func (r *Repository[T]) Where(query any, args ...interface{}) *Repository[T] {
	switch query.(type) {
	case string:
		argsLen := len(args)
		if argsLen == 0 {
			return r.WhereRaw(query, args...)
		}

		queryStr, _ := query.(string)
		if strings.Contains(queryStr, "?") {
			return r.WhereRaw(queryStr, args...)
		}

		expr := "="
		var value any
		if argsLen == 1 {
			value = args[0]
		} else if argsLen > 1 {
			value = args[1]
			expr = fmt.Sprintf("%v", args[0])
			if strings.ToUpper(expr) == "IN" {
				return r.WhereIn(queryStr, args[1])
			}
		}
		queryExpr := queryStr + expr + "?"
		r.DB = r.DB.Where(queryExpr, value)
	case func():
		f, _ := query.(func())
		tmp := *r.DB
		oldDB := &tmp
		r.Reset()
		fmt.Println(oldDB, r.DB)

		f()
		r.DB = oldDB.Where(r.DB)
	}
	return r
}

func (r *Repository[T]) WhereIn(column string, args interface{}) *Repository[T] {
	var values any
	switch args.(type) {
	case string:
		values = strings.Split(args.(string), ",")
	default:
		values = args
	}
	r.DB = r.DB.Where(column+" IN(?) ", values)
	return r
}

func (r *Repository[T]) OrWhereIn(column string, args interface{}) *Repository[T] {
	var values any
	switch args.(type) {
	case string:
		values = strings.Split(args.(string), ",")
	default:
		values = args
	}
	r.DB = r.DB.Or(column+" IN(?) ", values)
	return r
}

func (r *Repository[T]) WhereNotIn(column string, args interface{}) *Repository[T] {
	var values any
	switch args.(type) {
	case string:
		values = strings.Split(args.(string), ",")
	default:
		values = args
	}
	r.DB = r.DB.Where(column+" NOT IN(?) ", values)
	return r
}

func (r *Repository[T]) NotIn(column string, args interface{}) *Repository[T] {
	return r.WhereNotIn(column, args)
}

func (r *Repository[T]) In(column string, args interface{}) *Repository[T] {
	return r.WhereIn(column, args)
}

func (r *Repository[T]) FindInSet(column string, set any) *Repository[T] {
	r.DB = r.DB.Where("FIND_IN_SET(?,`"+column+"`)", set)
	return r
}

func (r *Repository[T]) Between(column string, value1, value2 any) *Repository[T] {
	r.DB = r.DB.Where(column+" BETWEEN ? AND ? ", value1, value2)
	return r
}

func (r *Repository[T]) NotBetween(column string, value1, value2 any) *Repository[T] {
	r.DB = r.DB.Where(column+" NOT BETWEEN ? AND ? ", value1, value2)
	return r
}

func (r *Repository[T]) Eq(key string, value any) *Repository[T] {
	r.DB = r.DB.Where(key+" =? ", value)
	return r
}

func (r *Repository[T]) Neq(key string, value any) *Repository[T] {
	r.DB = r.DB.Where(key+" !=? ", value)
	return r
}

func (r *Repository[T]) Gt(key string, value any) *Repository[T] {
	r.DB = r.DB.Where(key+" >? ", value)
	return r
}

func (r *Repository[T]) Ge(key string, value any) *Repository[T] {
	r.DB = r.DB.Where(key+" >=? ", value)
	return r
}

func (r *Repository[T]) Lt(key string, value any) *Repository[T] {
	r.DB = r.DB.Where(key+" <? ", value)
	return r
}

func (r *Repository[T]) Le(key string, value any) *Repository[T] {
	r.DB = r.DB.Where(key+" <=? ", value)
	return r
}

func (r *Repository[T]) Like(key string, value string) *Repository[T] {
	r.DB = r.DB.Where(key+" LIKE ? ", "%"+value+"%")
	return r
}

func (r *Repository[T]) LikeLeft(key string, value string) *Repository[T] {
	r.DB = r.DB.Where(key+" LIKE ? ", "%"+value)
	return r
}

func (r *Repository[T]) LikeRight(key string, value string) *Repository[T] {
	r.DB = r.DB.Where(key+" LIKE ? ", value+"%")
	return r
}

func (r *Repository[T]) NotLike(key string, value string) *Repository[T] {
	r.DB = r.DB.Where(key+" NOT LIKE ? ", "%"+value+"%")
	return r
}

func (r *Repository[T]) NotLikeLeft(key string, value string) *Repository[T] {
	r.DB = r.DB.Where(key+" NOT LIKE ? ", "%"+value)
	return r
}
func (r *Repository[T]) NotLikeRight(key string, value string) *Repository[T] {
	r.DB = r.DB.Where(key+" NOT LIKE ? ", value+"%")
	return r
}

func (r *Repository[T]) IsNull(key string) *Repository[T] {
	r.DB = r.DB.Where(key + " IS NULL ")
	return r
}

func (r *Repository[T]) IsNotNull(key string) *Repository[T] {
	r.DB = r.DB.Where(key + " IS NOT NULL ")
	return r
}

//-------------------以下对DB原生方法套壳-------------------------

func (r *Repository[T]) WhereRaw(query interface{}, args ...interface{}) *Repository[T] {
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

func (r *Repository[T]) OrRaw(query interface{}, args ...interface{}) *Repository[T] {
	r.DB = r.DB.Or(query, args...)
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
