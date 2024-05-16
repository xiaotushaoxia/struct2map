package struct2map

import (
	"fmt"
	"reflect"
)

type ConvertorOption func(*Convertor)

func KeepUnexported() ConvertorOption {
	return func(c *Convertor) {
		c.KeepUnexported = true
	}
}
func FlattenEmbed() ConvertorOption {
	return func(c *Convertor) {
		c.FlattenEmbed = true
	}
}

func NewConvertorAny(v any, options ...ConvertorOption) *Convertor {
	return NewConvertorValue(reflect.ValueOf(v), options...)
}

func NewConvertorValue(v reflect.Value, options ...ConvertorOption) *Convertor {
	c := &Convertor{v: v, options: options}
	for _, option := range options {
		option(c)
	}
	return c
}

type Convertor struct {
	v reflect.Value

	KeepUnexported bool
	FlattenEmbed   bool

	options []ConvertorOption
}

func (c *Convertor) Convert() map[string]any {
	v := indirectForValue(c.v)
	if !v.IsValid() {
		return nil
	}
	// Why Golang Nil Is Not Always Nil? Nil Explained  https://codefibershq.com/blog/golang-why-nil-is-not-always-nil
	// for reflect.Ptr. call v.IsNil() to ensure that v is really nil
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%v not struct or ptr to struct", v.Kind()))
	}
	var result = map[string]any{}

	var copied bool
	ta := v.Type()
	for i := 0; i < ta.NumField(); i++ {
		sf := ta.Field(i)
		if !sf.IsExported() && !c.KeepUnexported {
			continue
		}
		field := v.Field(i)
		if field.Kind() == reflect.UnsafePointer || field.Kind() == reflect.Chan || field.Kind() == reflect.Func {
			continue
		}
		if !field.CanInterface() {
			if !copied {
				v = copyReflectValue(v)
				field, copied = v.Field(i), true
			}
			field = getUnexportedField(field)
		}
		if !(sf.Anonymous && c.FlattenEmbed) {
			result[sf.Name] = c.convertSingle(field)
			continue
		}
		vv := c.convertSingle(field)
		switch tv := vv.(type) {
		case map[string]any:
			mergeMap(result, tv)
		default:
			result[sf.Name] = tv
		}
	}
	return result
}

// convertSingle if f cant be convert, return nil
func (c *Convertor) convertSingle(f reflect.Value) any {
	switch f.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128, reflect.String: // 17
		return f.Interface()
	case reflect.Struct:
		return c.copyConfig(f).Convert()
	case reflect.Ptr:
		elem := f.Elem()
		if elem.Kind() == reflect.Invalid {
			return getEmptyPtrDefaultValue(newElemOfPtr(f.Interface()).Elem().Kind())
		}
		return c.convertSingle(elem)
	case reflect.Array, reflect.Slice:
		if f.IsNil() {
			var vv []any
			return vv
		}
		l := make([]any, f.Len())
		// at first convertSingle return (result any, ok bool)
		// but I don't know how to convert a slice containing non-convertible elements.
		// if I skip non-convertible, it makes len(l) not equal to f.Len()
		// so if f is non-convertible, return nil
		for i := 0; i < len(l); i++ {
			l[i] = c.convertSingle(f.Index(i))
		}
		return l
	case reflect.Map:
		if f.IsNil() {
			var mp map[string]any
			return mp
		}
		m := map[string]any{}
		for _, k := range f.MapKeys() {
			if k.Kind() == reflect.Interface {
				// if f is map[any]T, kind of k is Interface.  this line make kind of k turn back to String
				k = reflect.ValueOf(k.Interface())
			}
			if k.Kind() != reflect.String {
				continue
			}
			m[k.String()] = c.convertSingle(f.MapIndex(k))
		}
		return m
	case reflect.Invalid:
		return nil
	case reflect.Interface:
		// I noticed the f.Type().Name() is empty, it is weird (for me).
		// Having no other options, I tried the following line of code and miraculously got type of f
		f = reflect.ValueOf(f.Interface()) // see Test_kindInterface
		return c.convertSingle(f)
	//case reflect.UnsafePointer: // fixme what is this ?
	//	return nil
	//case reflect.Chan, reflect.Func: // won't reach
	//	return nil
	default: // won't reach
		return nil
	}
}

func (c *Convertor) copyConfig(v reflect.Value) *Convertor {
	return NewConvertorValue(v, c.options...)
}

func getEmptyPtrDefaultValue(kind reflect.Kind) any {
	switch kind {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128, reflect.String: // 17
		return nil
	case reflect.Struct:
		var mp map[string]any
		return mp
	case reflect.Map:
		var mp map[string]any
		return mp
	case reflect.Slice, reflect.Array:
		var vv []any
		return vv
	case reflect.Interface:
		return nil
	case reflect.Invalid:
		return nil
	//case reflect.UnsafePointer: // fix me what is this ?
	//	return nil
	//case reflect.Chan, reflect.Func: // won't reach
	//	return nil
	default: // won't reach
		return nil
	}
}
