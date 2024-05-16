package struct2map

import (
	"reflect"
	"unsafe"
)

func getUnexportedField(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}

func copyReflectValue(v reflect.Value) reflect.Value {
	elem := reflect.New(v.Type()).Elem()
	elem.Set(v)
	return elem
}

func indirectForValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func newElemOfPtr(p any) reflect.Value {
	// p must be ptr to struct
	elemType := reflect.TypeOf(p).Elem()
	zeroVal := reflect.Zero(elemType)
	ptr := reflect.New(elemType)
	ptr.Elem().Set(zeroVal)
	return ptr
}
