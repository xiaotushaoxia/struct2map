# struct2map

convert struct to map[string]any


## Usage

``` go
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
```