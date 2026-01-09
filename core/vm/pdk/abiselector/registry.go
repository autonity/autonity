package abiselector

import (
	"reflect"
	"sync"
)

// registry maps all precompiled state types to their dispatchers
var (
	registry = sync.Map{}
)

func RegisterDispatcher(logic interface{}) *Dispatcher {
	typ := reflect.TypeOf(logic)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	d := newDispatcher()
	err := InferABIMethods(d, reflect.ValueOf(logic))
	if err != nil {
		panic(err)
	}
	registry.Store(typ, d)
	return d
}

func GetDispatcher(logic interface{}) *Dispatcher {
	typ := reflect.TypeOf(logic)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if dispatcher, ok := registry.Load(typ); ok {
		return dispatcher.(*Dispatcher)
	}
	// note: we could choose to auto register here, but to avoid unintended consequences, we return nil
	return nil
}

//func GetOrRegisterDispatcher(logic interface{}) *Dispatcher {
//	typ := reflect.TypeOf(logic)
//	if typ.Kind() == reflect.Ptr {
//		typ = typ.Elem()
//	}
//	if dispatcher, ok := registry.Load(typ); ok {
//		return dispatcher.(*Dispatcher)
//	}
//	d := newDispatcher()
//	err := InferABIMethods(d, reflect.ValueOf(logic))
//	if err != nil {
//		panic(err)
//	}
//
//	registry.Store(typ, d)
//	return d
//}
