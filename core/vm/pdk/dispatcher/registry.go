package dispatcher

import (
	"reflect"
	"sync"
)

// registry maps all precompiled state types to their dispatchers
var (
	registry = sync.Map{}
)

// RegisterDispatcher creates and registers a dispatcher for a given logic interface.
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

// GetDispatcher retrieves the registered dispatcher for a given logic interface.
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
