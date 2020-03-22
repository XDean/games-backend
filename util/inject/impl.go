package inject

import (
	"errors"
	"fmt"
	"reflect"
)

type impl struct {
	objs []interface{}
}

func (t *impl) Register(obj interface{}) {
	t.objs = append(t.objs, obj)
}

func (t *impl) Get(rec interface{}) error {
	err, _ := t.get(rec, 0)
	return err
}

func (t *impl) GetAll(recs interface{}) error {
	panic("not implement")
}

func (t *impl) Inject(obj interface{}) error {
	objValue := reflect.ValueOf(obj)
	for objValue.Kind() == reflect.Ptr || objValue.Kind() == reflect.Interface {
		objValue = objValue.Elem()
	}
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Type().Field(i)
		fieldValue := objValue.Field(i)
		if _, ok := field.Tag.Lookup("inject"); ok {
			if err := t.Get(fieldValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *impl) Refresh() error {
	for i := range t.objs {
		err := t.Inject(&t.objs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *impl) get(rec interface{}, offset int) (err error, pos int) {
	var recValue reflect.Value
	if v, ok := rec.(reflect.Value); ok {
		recValue = v
	} else {
		recValue = reflect.ValueOf(rec)
	}
	var elemValue reflect.Value
	var elemType reflect.Type
	if recValue.CanSet() {
		elemValue = recValue
		elemType = recValue.Type()
	} else if recValue.Kind() == reflect.Ptr || recValue.Kind() == reflect.Interface {
		elemValue = recValue.Elem()
		elemType = elemValue.Type()
	} else {
		return errors.New("must give addr-able receiver"), -1
	}
	for i := offset; i < len(t.objs); i++ {
		objValue := reflect.ValueOf(t.objs[i])
		objType := objValue.Type()
		if objType.AssignableTo(elemType) {
			elemValue.Set(objValue)
			return nil, i
		}
		if objValue.Kind() == reflect.Ptr {
			objPtrValue := objValue.Elem()
			objPtrType := objPtrValue.Type()
			fmt.Println(objPtrType, elemType)
			if objPtrType.AssignableTo(elemType) {
				elemValue.Set(objPtrValue)
				return nil, i
			}
		}
	}
	return errors.New("no such bean: " + elemType.String()), -1
}
