package struct2map

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/frankban/quicktest"
)

type testErr struct{}

func (te *testErr) Error() string {
	return "error happened"
}
func testError() *testErr {
	return nil
}

func badNil() error {
	//if err := testError(); err != nil {
	//	return err
	//}
	//return nil  // this make test return real nil
	return testError()
}

func Test_Convert_nilNotEqualNil(t *testing.T) {
	err := badNil()
	if err == nil {
		t.Fatal("err should not eq to nil")
	}
	c := NewConvertorAny(err)
	v := c.Convert()
	if v != nil {
		t.Fatal(err, "convert to", v, "not nil")
	}
}

func Test_Convert_nil(t *testing.T) {
	c2 := NewConvertorAny(nil)
	v := c2.Convert()
	if v != nil {
		t.Fatal("nil convert to", v, "not nil")
	}

	var t2 *testing.T
	c2 = NewConvertorAny(t2)
	v = c2.Convert()
	if v != nil {
		t.Fatal("*testing.T(nil) convert to", v, "not nil")
	}
}

func Test_Convert_notStruct(t *testing.T) {
	var vs = []any{1, 2.1, map[string]any{}, []any{}, true, uint32(32)}
	var panicNum int
	for _, v := range vs {
		last := panicNum
		func() {
			defer func() {
				if p := recover(); p != nil {
					panicNum++
				}
			}()
			NewConvertorAny(v).Convert()
		}()
		if last == panicNum {
			t.Fatalf("convert %v no panic", v)
		}
	}
}

func Test_Convert(t *testing.T) {
	type kc struct {
		a   int
		b   string
		c   bool
		kcc bool
		KCC bool
	}
	type BC struct {
		bb string
		cc bool
		DD bool
	}
	type kkc struct {
		kc
		BC
		d   bool
		DcD bool
		//DD bool
		VS []any
	}

	kv := []any{11, true, kc{KCC: true, a: 1232324}}
	vv := kkc{VS: kv}

	test4(quicktest.New(t), vv, [4]string{
		"map[BC:map[DD:false] DcD:false VS:[11 true map[KCC:true]]]",
		"map[BC:map[DD:false bb: cc:false] DcD:false VS:[11 true map[KCC:true a:1232324 b: c:false kcc:false]] d:false kc:map[KCC:false a:0 b: c:false kcc:false]]",
		"map[DD:false DcD:false VS:[11 true map[KCC:true]]]",
		"map[DD:false DcD:false KCC:false VS:[11 true map[KCC:true a:1232324 b: c:false kcc:false]] a:0 b: bb: c:false cc:false d:false kcc:false]",
	})
}

func Test_kindInterface(t *testing.T) {
	type BBC struct {
		VS []any
	}
	bbc := BBC{}
	bbc.VS = append(bbc.VS, 1, nil, map[string]any{}, map[any]any{"tu": 1, "1": 2, 3: 4, "5": map[any]any{"a": "b"}}, 2, "test1")

	inx0 := reflect.ValueOf(bbc).FieldByName("VS").Index(0)
	dkind := inx0.Kind() //  reflect.Interface
	if dkind != reflect.Interface {
		t.Fatalf("error idx0 kind")
	}
	of := reflect.ValueOf(inx0.Interface())
	wrapKind := of.Kind()        //  reflect.Int
	if wrapKind != reflect.Int { // kind turn back !
		t.Fatalf("error idx0 kind")
	}
}

func Test_Convert_withSlice(t *testing.T) {
	type BBC struct {
		VS []any
	}
	var vs2 []any
	var fp *float32
	var f32 = 1.22
	var fp2 = &f32
	bbc := BBC{}
	var mp2 map[string]any
	bbc.VS = append(bbc.VS,
		1,
		nil,
		map[string]any{},
		map[any]any{"tu": 1, "1": 2, 3: 4, "5": map[any]any{"a": "b"}},
		2,
		"test1",
		mp2,
		[]any{},
		vs2,
		fp2,
		fp,
	)

	test4(quicktest.New(t), bbc, [4]string{
		"map[VS:[1 <nil> map[] map[1:2 5:map[a:b] tu:1] 2 test1 map[] [] <nil> 1.22 <nil>]]",
		"map[VS:[1 <nil> map[] map[1:2 5:map[a:b] tu:1] 2 test1 map[] [] <nil> 1.22 <nil>]]",
		"map[VS:[1 <nil> map[] map[1:2 5:map[a:b] tu:1] 2 test1 map[] [] <nil> 1.22 <nil>]]",
		"map[VS:[1 <nil> map[] map[1:2 5:map[a:b] tu:1] 2 test1 map[] [] <nil> 1.22 <nil>]]",
	})
}

func Test_Convert_withPtrSliceMap(t *testing.T) {
	type BBC struct {
		VS  *[]any
		int *map[string]any
	}

	c := quicktest.New(t)

	result := NewConvertorAny(BBC{}, KeepUnexported()).Convert()
	// nil with type
	c.Assert(result["VS"], quicktest.IsNil)
	c.Assert(reflect.TypeOf(result["VS"]), quicktest.Equals, reflect.TypeOf([]any{}))
	c.Assert(result["int"], quicktest.IsNil)
	c.Assert(reflect.TypeOf(result["int"]), quicktest.Equals, reflect.TypeOf(map[string]any{}))

	test4(c, BBC{}, [4]string{
		"map[VS:[]]", // output [], but VS is nil
		"map[VS:[] int:map[]]",
		"map[VS:[]]",
		"map[VS:[] int:map[]]",
	})
	mp := map[string]any{"a": 1, "b": 2}
	test4(c, BBC{VS: &([]any{1, 2, 3}), int: &mp}, [4]string{
		"map[VS:[1 2 3]]",
		"map[VS:[1 2 3] int:map[a:1 b:2]]",
		"map[VS:[1 2 3]]",
		"map[VS:[1 2 3] int:map[a:1 b:2]]",
	})
}

func Test_Convert_withPtrAnyInSlice(t *testing.T) {
	type BBC struct {
		VS []any
		Va *any
	}
	bbc := BBC{}
	var nm any
	var nilSlice []int
	var mk any = 12323
	bbc.VS = append(bbc.VS, &nm, 1, 2, &nilSlice, nilSlice, &mk)

	test4(quicktest.New(t), bbc, [4]string{
		"map[VS:[<nil> 1 2 <nil> <nil> 12323] Va:<nil>]",
		"map[VS:[<nil> 1 2 <nil> <nil> 12323] Va:<nil>]",
		"map[VS:[<nil> 1 2 <nil> <nil> 12323] Va:<nil>]",
		"map[VS:[<nil> 1 2 <nil> <nil> 12323] Va:<nil>]",
	})
}

func Test_Convert_withNilSlice(t *testing.T) {
	type BBC struct {
		VS []any
	}
	bbc := BBC{}
	pbbc := &bbc
	test4(quicktest.New(t), pbbc, [4]string{
		"map[VS:<nil>]",
		"map[VS:<nil>]",
		"map[VS:<nil>]",
		"map[VS:<nil>]",
	})
	test4(quicktest.New(t), bbc, [4]string{
		"map[VS:<nil>]",
		"map[VS:<nil>]",
		"map[VS:<nil>]",
		"map[VS:<nil>]",
	})
	test4(quicktest.New(t), &pbbc, [4]string{
		"map[VS:<nil>]",
		"map[VS:<nil>]",
		"map[VS:<nil>]",
		"map[VS:<nil>]",
	})
}

func Test_Convert_withFuncChan(t *testing.T) {
	type BBC struct {
		VS func()
		C  chan int
		c2 chan int
		ii int
	}

	test4(quicktest.New(t), BBC{}, [4]string{
		"map[]",
		"map[ii:0]",
		"map[]",
		"map[ii:0]",
	})
}

func Test_Convert_embedBaseType(t *testing.T) {
	type tt struct {
		int
		float64
		*float32
	}
	// for base type embed, same as follows
	//type tt struct {
	//	int int
	//	float64 float64
	//}

	vv := tt{
		int:     11,
		float64: 0.2,
	}

	c := quicktest.New(t)

	test4(c, vv, [4]string{
		"map[]", "map[float32:<nil> float64:0.2 int:11]", "map[]", "map[float32:<nil> float64:0.2 int:11]",
	})

	pp := float32(1.1)
	vv.float32 = &pp

	test4(c, vv, [4]string{
		"map[]", "map[float32:1.1 float64:0.2 int:11]", "map[]", "map[float32:1.1 float64:0.2 int:11]",
	})
}

func Test_Convert_embedPointer(t *testing.T) {
	type tt2 struct {
		A int
		B bool
		c bool
	}
	type tt struct {
		*tt2
		D bool
		e bool
	}

	vv := tt{}

	c := quicktest.New(t)

	test4(c, vv, [4]string{
		"map[D:false]", "map[D:false e:false tt2:map[]]", "map[D:false]", "map[D:false e:false]",
	})

}

func test4(c *quicktest.C, v any, want [4]string) {
	vv := reflect.ValueOf(v)

	c1, c2, c3, c4 := NewConvertorValue(vv), NewConvertorValue(vv, KeepUnexported()),
		NewConvertorValue(vv, FlattenEmbed()), NewConvertorValue(vv, KeepUnexported(), FlattenEmbed())

	var cs []*Convertor
	cs = append(cs, c1, c2, c3, c4)

	var gots []string
	var output []string
	for _, convertor := range cs {
		convert := convertor.Convert()
		gi := fmt.Sprintf("%v", convert)
		gots = append(gots, gi)
		output = append(output, fmt.Sprintf("\"%s\"", gi))
	}
	//fmt.Println(strings.Join(output, ",\n") + ",")
	for i, got := range gots {
		c.Assert(got, quicktest.Equals, want[i], quicktest.Commentf(c.TB.Name()+":"+getComment(cs[i])))
	}
}

func getComment(c *Convertor) string {
	return fmt.Sprintf("KeepUnexported:%t, FlattenEmbed: %t, %#v", c.KeepUnexported, c.FlattenEmbed, c.v.Interface())
}

func Test_usage(t *testing.T) {
	type B struct {
		B1 int
		B2 string
		b3 bool
	}

	type C struct {
		C1 int
		C2 string
		C3 bool
	}
	type A struct {
		A1 int
		A3 bool
		B  B
		C
	}

	v := A{
		A1: 1,
		A3: false,
		B: B{
			B1: 99,
			B2: "b2",
			b3: true,
		},
		C: C{
			C1: 22,
			C2: "c2",
			C3: false,
		},
	}
	c1 := NewConvertorAny(v)
	c2 := NewConvertorAny(v, FlattenEmbed())   // difference between result of c1 and c2
	c3 := NewConvertorAny(v, KeepUnexported()) // export unexported fields
	c4 := NewConvertorAny(v, FlattenEmbed(), KeepUnexported())

	fmt.Println(c1.Convert()) // map[A1:1 A3:false B:map[B1:99 B2:b2] C:map[C1:22 C2:c2 C3:false]]

	fmt.Println(c2.Convert()) // map[A1:1 A3:false B:map[B1:99 B2:b2] C1:22 C2:c2 C3:false]

	fmt.Println(c3.Convert()) // map[A1:1 A3:false B:map[B1:99 B2:b2 b3:true] C:map[C1:22 C2:c2 C3:false]]

	fmt.Println(c4.Convert()) // map[A1:1 A3:false B:map[B1:99 B2:b2 b3:true] C1:22 C2:c2 C3:false]

	fmt.Println(NewConvertorAny(nil).Convert()) // map[]  (nil map)

	var a = struct{}{}

	fmt.Println(NewConvertorAny(a).Convert()) // map[]  (empty map)

	type D struct {
		*C
		V *C
	}

	c5 := NewConvertorAny(D{})
	c6 := NewConvertorAny(D{}, FlattenEmbed())
	r5 := c5.Convert()

	// r5["C"] is nil map but not nil
	fmt.Println(r5, reflect.ValueOf(r5["C"]).IsNil(), r5["C"] == nil) // map[C:map[] V:map[]] true false

	// enable FlattenEmbed, *C convert to nil map and merge to result.
	fmt.Println(c6.Convert()) // map[V:map[]]
}
