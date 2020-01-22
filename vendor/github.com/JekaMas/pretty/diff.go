package pretty

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
)

type sbuf []string

func (s *sbuf) Write(b []byte) (int, error) {
	*s = append(*s, string(b))
	return len(b), nil
}

const (
	defaultMessageHeader     = "Obtained:\t\tExpected:"
	defaultDiffMessageFormat = "\n%v"
)

//DiffMessage - get fine diff message for two terms
//Get two format strings. First for header message, second for each diff message
func DiffMessage(obtained, expected interface{}, format ...string) string {
	var (
		failMessage       bytes.Buffer
		messageHeader     = defaultMessageHeader
		diffMessageFormat = defaultDiffMessageFormat
	)

	diffs := Diff(obtained, expected)

	if len(format) > 0 {
		messageHeader = format[0]
	}

	if len(format) > 1 {
		diffMessageFormat = format[1]
	}

	if len(diffs) > 0 {
		failMessage.WriteString(messageHeader)
		for _, singleDiff := range diffs {
			failMessage.WriteString(fmt.Sprintf(diffMessageFormat, singleDiff))
		}
	}

	return failMessage.String()
}

// Diff returns a slice where each element describes
// a difference between a and b.
func Diff(a, b interface{}) (desc []string) {
	Fdiff((*sbuf)(&desc), a, b)
	return desc
}

// Fdiff writes to w a description of the differences between a and b.
func Fdiff(w io.Writer, a, b interface{}) {
	diffWriter{w: w}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
}

type diffWriter struct {
	w io.Writer
	l string // label
}

func (w diffWriter) printf(f string, a ...interface{}) {
	var l string
	if w.l != "" {
		l = w.l + ": "
	}
	fmt.Fprintf(w.w, l+f, a...)
}

func (w diffWriter) diff(av, bv reflect.Value) {
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %#v", bv.Interface())
		return
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%#v != nil", av.Interface())
		return
	}
	if !av.IsValid() && !bv.IsValid() {
		return
	}

	at := av.Type()
	bt := bv.Type()
	if at != bt {
		w.printf("%v != %v", at, bt)
		return
	}

	// numeric types, including bool
	if at.Kind() < reflect.Array {
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%#v != %#v", a, b)
		}
		return
	}

	switch at.Kind() {
	case reflect.String:
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%q != %q", a, b)
		}
	case reflect.Ptr:
		switch {
		case av.IsNil() && !bv.IsNil():
			w.printf("nil != %v", bv.Interface())
		case !av.IsNil() && bv.IsNil():
			w.printf("%v != nil", av.Interface())
		case !av.IsNil() && !bv.IsNil():
			w.diff(av.Elem(), bv.Elem())
		}
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			// If a field is exported. See: https://golang.org/src/reflect/type.go?s=20439:20967#L726
			if at.Field(i).PkgPath == "" {
				w.relabel(at.Field(i).Name).diff(av.Field(i), bv.Field(i))
			}
		}
	case reflect.Slice:
		lenA := av.Len()
		lenB := bv.Len()
		if lenA != lenB {
			w.printf("%s[%d] != %s[%d]", av.Type(), lenA, bv.Type(), lenB)
		}
		for i := 0; i < int(math.Min(float64(lenA), float64(lenB))); i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i))
		}
		if lenA > lenB {
			for i := lenB; i < lenA; i++ {
				w.relabel(fmt.Sprintf("[%d]", i)).printf("%v != (missing)", av.Index(i))
			}
		} else if lenA < lenB {
			for i := lenA; i < lenB; i++ {
				w.relabel(fmt.Sprintf("[%d]", i)).printf("(missing) != %v", bv.Index(i))
			}
		}

	case reflect.Map:
		ak, both, bk := keyDiff(av.MapKeys(), bv.MapKeys())
		for _, k := range ak {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("%v{%#+v} != (missing)", av.MapIndex(k).Kind(), av.MapIndex(k).Interface())
		}
		for _, k := range both {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.diff(av.MapIndex(k), bv.MapIndex(k))
		}
		for _, k := range bk {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("(missing) != %v{%#+v}", bv.MapIndex(k).Kind(), bv.MapIndex(k).Interface())

		}
	case reflect.Interface:
		w.diff(reflect.ValueOf(av.Interface()), reflect.ValueOf(bv.Interface()))
	default:
		if !reflect.DeepEqual(av.Interface(), bv.Interface()) {
			w.printf("%# v != %# v", Formatter(av.Interface()), Formatter(bv.Interface()))
		}
	}
}

func (d diffWriter) relabel(name string) (d1 diffWriter) {
	d1 = d
	if d.l != "" && name[0] != '[' {
		d1.l += "."
	}
	d1.l += name
	return d1
}

func keyDiff(a, b []reflect.Value) (ak, both, bk []reflect.Value) {
	for _, av := range a {
		inBoth := false
		for _, bv := range b {
			if reflect.DeepEqual(av.Interface(), bv.Interface()) {
				inBoth = true
				both = append(both, av)
				break
			}
		}
		if !inBoth {
			ak = append(ak, av)
		}
	}
	for _, bv := range b {
		inBoth := false
		for _, av := range a {
			if reflect.DeepEqual(av.Interface(), bv.Interface()) {
				inBoth = true
				break
			}
		}
		if !inBoth {
			bk = append(bk, bv)
		}
	}
	return
}
