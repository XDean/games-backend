package util

import (
	"fmt"
	"reflect"
)

func ReflectSetReceiver(data interface{}) func(receiver interface{}) error {
	dataValue := reflect.ValueOf(data)
	return func(receiver interface{}) error {
		receiverValue := reflect.ValueOf(receiver)
		if receiverValue.Kind() == reflect.Interface || receiverValue.Kind() == reflect.Ptr {
			receiveValue := receiverValue.Elem()
			if receiveValue.Type() == dataValue.Type() {
				receiveValue.Set(dataValue)
				return nil
			}
		}
		return fmt.Errorf("Wrong receiver type, except *(%T), actual %T", data, receiver)
	}
}
