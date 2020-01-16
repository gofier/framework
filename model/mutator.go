package model

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
)

func callMutator(scope *gorm.Scope, isGetter bool) {
	reflectValue := scope.IndirectValue()

	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		reflectValue = reflectValue.Addr()
	}

	switch reflectValue.Type().Elem().Kind() {
	case reflect.Struct:
	case reflect.Slice:
	default:
		panic("cannot use mutator in type:" + reflectValue.Type().Elem().Kind().String())
	}
}

func structReflect(reflectValue *reflect.Value, isGetter bool) {
	wg := &sync.WaitGroup{}

	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		tmp := reflectValue.Addr()
		reflectValue = &tmp
	}

	for i := 0; i < reflectValue.Type().Elem().NumField(); i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, index int) {
			fieldName := reflectValue.Type().Elem().Field(index).Name
			fieldValue := reflectValue.Elem().Field(index).Interface()

			if isGetter {
				getter(reflectValue, fieldName, fieldValue)
			} else {
				setter(*reflectValue, fieldName, fieldValue)
			}
		}(wg, i)
	}
	wg.Wait()
}

func setter(reflectValue reflect.Value, fieldName string, fieldValue interface{}) {
	const setterMethodTemplate = "Set%sAttribute"
	methodName := fmt.Sprintf(setterMethodTemplate, strcase.ToCamel(fieldName))

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		methodValue.Interface().(func(value interface{}))(fieldValue)
	}
}

func getter(reflectValue *reflect.Value, fieldName string, fieldValue interface{}) {
	const getterMethodTemplate = "Get%sAttribute"
	methodName := fmt.Sprintf(getterMethodTemplate, strcase.ToCamel(fieldName))

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		newGetData := methodValue.Interface().(func(value interface{}) interface{})(fieldValue)
		reflectValue.Elem().FieldByName(fieldName).Set(reflect.ValueOf(newGetData))
	}
}
