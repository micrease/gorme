package main

import "gorm.io/gorm"

type Repository[T any] struct {
	db *gorm.DB
}

func (r *Repository[T]) Where(query interface{}, args ...interface{}) *Repository[T] {
	r.db.Where(query, args...)
	return r
}

func (r *Repository[T]) GetOne() (*T, error) {
	var t T
	err := r.db.First(&t).Error
	return &t, err
}
