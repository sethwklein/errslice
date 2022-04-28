package errslice_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/diff"

	"github.com/sethwklein/errslice"
)

type other []error

func (errs other) Error() string {
	var b []byte
	b = append(b, "other: ("...)
	for i, err := range errs {
		if i > 0 {
			b = append(b, ", "...)
		}
		b = append(b, err.Error()...)
	}
	b = append(b, ')')
	return string(b)
}

var (
	e1 = errors.New("one")
	e2 = errors.New("two")
	e3 = errors.New("three")
)

func TestAppend(t *testing.T) {
	want := `
[<nil>]: <nil>
[<nil> <nil>]: <nil>
[<nil> <nil> <nil>]: <nil>
[one]: one
[one two]: one and two
[one two three]: one, two, and three
[one <nil>]: one
[<nil> one]: one
[one two <nil>]: one and two
[one and two three]: one, two, and three
[other: (one, two) three]: one, two, and three
[one two and three]: one, two, and three
`
	var got bytes.Buffer
	fmt.Fprintln(&got)
	for _, have := range [][]error{
		{nil},
		{nil, nil},
		{nil, nil, nil},
		{e1},
		{e1, e2},
		{e1, e2, e3},
		{e1, nil},
		{nil, e1},
		{e1, e2, nil},
		{errslice.Slice{e1, e2}, e3},
		{other{e1, e2}, e3},
		{e1, errslice.Slice{e2, e3}},
	} {
		fmt.Fprintf(&got, "%v: %v\n", have, errslice.Append(have...))
	}
	if differences := diff.Diff(want, got.String()); len(differences) > 0 {
		t.Error(differences)
	}
}

var fromTests = []error{
	nil,
	errslice.Slice{e1, e2},
	other{e1, e2},
	e1,
}

func testFrommer(t *testing.T, f func(error) errslice.Slice, want string) {
	var got bytes.Buffer
	fmt.Fprintln(&got)
	for _, have := range fromTests {
		fmt.Fprintf(&got, "%v: %v\n", have, f(have))
	}
	if differences := diff.Diff(want, got.String()); len(differences) > 0 {
		t.Error(differences)
	}
}

func TestFrom(t *testing.T) {
	want := `
<nil>: 
one and two: one and two
other: (one, two): one and two
one: one
`
	testFrommer(t, errslice.From, want)
}

func TestFromFast(t *testing.T) {
	want := `
<nil>: 
one and two: one and two
other: (one, two): other: (one, two)
one: one
`
	testFrommer(t, errslice.FromFast, want)
}
