package comm

import (
	"errors"
	"reflect"
)

var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
)

type Funcs map[string]reflect.Value

func NewFuncs(size int) Funcs {
	return make(Funcs, size)
}

func (f Funcs) Bind(name string, fn interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(name + " is not callable.")
		}
	}()
	v := reflect.ValueOf(fn)
	v.Type().NumIn()
	f[name] = v
	return
}

func (f Funcs) Call(name string, params ... interface{}) (result []reflect.Value, err error) {
	if _, ok := f[name]; !ok {
		err = errors.New(name + " does not exist.")
		return
	}
	if len(params) != f[name].Type().NumIn() {
		err = ErrParamsNotAdapted
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f[name].Call(in)
	return
}


func FuncBind(funcMap map[string]interface{}, funcs Funcs) error{

	for k, v := range funcMap {
		err := funcs.Bind(k, v)
		if k[:3] == "err" {
			if err == nil {
				return errors.New("Bind:" + k + "an error should be paniced.")

			}
		} else {
			if err != nil {
				return errors.New("Bind:" + k + ",err:" + err.Error())
			}
		}
	}
	return nil
}
