package model

import (
	"errors"
	"fmt"

	"github.com/gofier/framework/database"

	"github.com/jinzhu/gorm"
)

type IBaseModel interface {
	Query() *gorm.DB
}

type BaseModel struct {
	db *gorm.DB
}

func (bm *BaseModel) Query(m interface{}) *gorm.DB {
	if bm.db == nil {
		return database.DB().Model(m)
	}
	return bm.db.Model(m)
}

func (bm *BaseModel) BeforeSave(scope *gorm.Scope) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			if __err, ok := _err.(error); ok {
				err = __err
				return
			}
			err = errors.New(fmt.Sprint(_err))
			return
		}
	}()

	callMutator(scope, false)

	return nil
}

func (bm *BaseModel) BeforeUpdate(scope *gorm.Scope) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			if __err, ok := _err.(error); ok {
				err = __err
				return
			}
			err = errors.New(fmt.Sprint(_err))
			return
		}
	}()

	callMutator(scope, false)
	return nil
}

func (bm *BaseModel) AfterFind(scope *gorm.Scope) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			if __err, ok := _err.(error); ok {
				err = __err
				return
			}
			err = errors.New(fmt.Sprint(_err))
			return
		}
	}()

	callMutator(scope, true)
	return nil
}
