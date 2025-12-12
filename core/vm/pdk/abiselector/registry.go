package abiselector

import (
	"reflect"
	"sync"
)

// registry maps all precompiled state types to their dispatchers
var (
	registry = sync.Map{}
)

func GetOrRegisterDispatcher(logic interface{}) *Dispatcher {
	typ := reflect.TypeOf(logic)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if dispatcher, ok := registry.Load(typ); ok {
		return dispatcher.(*Dispatcher)
	}
	d := newDispatcher()
	// since we are passing address of struct, type will be pointer type,
	// we can create a zero value type, which can be used for lookup later
	val := reflect.New(typ)
	err := InferABIMethods(d, val)
	if err != nil {
		panic(err)
	}

	registry.Store(typ, d)
	return d
}
