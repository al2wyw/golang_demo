package main

import (
	"log"
	"reflect"
	"testing"
)

// User 结构体表示一个用户
type User struct {
	ID   int
	Name string
}

type Aspect struct {
	Before []func([]interface{}) error
	After  []func([]interface{}) error
}

// Apply 利用反射执行切面的前置和后置处理函数
func (a *Aspect) Apply(targetFunc interface{}, args ...interface{}) ([]interface{}, error) {
	for _, beforeFunc := range a.Before {
		if err := beforeFunc(args); err != nil {
			return nil, err
		}
	}

	vArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vArgs[i] = reflect.ValueOf(arg)
	}
	result := reflect.ValueOf(targetFunc).Call(vArgs)
	var rets []interface{}
	for _, v := range result {
		if v.Type().AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
			if v.IsNil() {
				continue
			}
			if err, ok := v.Interface().(error); ok {
				if err != nil {
					return nil, err
				}
			} else {
				log.Panicln("Failed to convert return value to error")
			}
		} else {
			rets = append(rets, v.Interface())
		}
	}

	for _, afterFunc := range a.After {
		if err := afterFunc(args); err != nil {
			return nil, err
		}
	}

	return rets, nil
}

func BeforeAspect1(args []interface{}) error {
	id := args[0].(int)
	log.Println("BeforeAspect1 get id:", id)
	// change the first argument value
	args[0] = id + 1
	return nil
}

func BeforeAspect2(args []interface{}) error {
	name := args[1].(string)
	log.Println("BeforeAspect2 get name:", name)
	return nil
}

func BeforeAspect3(args []interface{}) error {
	user := args[2].(*User)
	log.Println("BeforeAspect3 get user:", user.ID, user.Name)
	return nil
}

func AfterAspect1(args []interface{}) error {
	log.Println("AfterAspect1")
	return nil
}

func TargetFunction(id int, name string, user *User) (*User, error) {
	user.ID = 13
	user.Name = "John Doe"
	log.Printf("Executing target function with ID: %d, Name: %s, User: %+v\n", id, name, user)
	return user, nil
}

func TestAop(f *testing.T) {
	myAspect := &Aspect{
		Before: []func([]interface{}) error{
			BeforeAspect1,
			BeforeAspect2,
			BeforeAspect3,
		},
		After: []func([]interface{}) error{
			AfterAspect1,
		},
	}

	// 构造目标函数的参数值
	id := 1
	name := "John"
	user := &User{ID: 100, Name: "Alice"}

	result, err := myAspect.Apply(TargetFunction, id, name, user)
	if err != nil {
		log.Println("have err:", err)
		return
	}

	log.Println("Target function result:", len(result), result[0])
}
