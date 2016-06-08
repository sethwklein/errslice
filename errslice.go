// Package errslice is for working with multiple errors, especially for
// returning them to callers that may be expecting only one.
//
// See the AppendCall example for typical usage.
package errslice

import "reflect"

// Slice is needed because it is not possible to add methods to []error.

// Slice is a []error that implements error.
type Slice []error

// Error implements error. Better results are generally achieved by using From
// and a loop.
func (errs Slice) Error() string {
	var b []byte
	for i, err := range errs {
		b = append(b, err.Error()...)
		switch len(errs) - i {
		case 1:
			// end of line
		case 2:
			// oxford comma for the win
			if len(errs) > 2 {
				b = append(b, ',')
			}
			b = append(b, " and "...)
		default:
			b = append(b, ", "...)
		}
	}
	return string(b)
}

// FromFast is a like From but trades away detection of []errors named something
// other than Slice to gain performance.
func FromFast(err error) Slice {
	if err == nil {
		return nil
	}
	if errs, ok := err.(Slice); ok {
		return errs
	}
	return Slice{err}
}

// BUG(sk): should I expose fromByReflection? I used it for higher performance in Append.

// fromByReflection uses reflection to detect if err contains a named type of
// []error. It is its own function because it uses things that panic.
func fromByReflection(err error) (slice []error, ok bool) {
	defer func() {
		// there's no way to detect if this panic is from type conversion
		// or other causes, so this just acts as though it's from type
		// conversion.
		ok = recover() == nil
	}()
	slice = reflect.ValueOf(err).Convert(reflect.TypeOf(slice)).Interface().([]error)
	return // ok will be set by the defer
}

// From wraps err in a Slice unless it already is a Slice or []error by another
// name, in which case it returns it unaltered or uses reflection to convert
// it to a Slice, respectively. If err is nil, From returns nil.
//
// From is useful for iterating over errors, such as when printing them. See
// AppendCall's example for typical usage.
func From(err error) Slice {
	if err == nil {
		return nil
	}
	if errs, ok := err.(Slice); ok {
		return errs
	}
	if errs, ok := fromByReflection(err); ok {
		return Slice(errs)
	}
	return Slice{err}
}

// Append concatenates errors and Slices, dropping nils. It flattens Slices,
// but not recursively. If not passed at least one non-nil argument, it will
// return nil, and if passed only one non-nil argument, it will return that
// one, unchanged.
func Append(errs ...error) error {
	// common case
	if len(errs) == 2 && errs[0] == nil && errs[1] == nil {
		return nil
	}

	p := 0

	// first non-nil error
	for ; ; p++ {
		if p >= len(errs) {
			return nil
		}
		if errs[p] != nil {
			break
		}
	}
	first := errs[p]
	p++

	// second non-nil error
	for ; ; p++ {
		if p >= len(errs) {
			return first
		}
		if errs[p] != nil {
			break
		}
	}
	var acc Slice
	var ok bool
	if acc, ok = first.(Slice); ok {
		// it's good
	} else if acc, ok = fromByReflection(first); ok {
		// it's good
	} else {
		acc = Slice{first}
	}

	// even more things
	for ; p < len(errs); p++ {
		e := errs[p]
		if e == nil {
			continue
		}
		if es, ok := e.(Slice); ok {
			acc = append(acc, es...)
		} else {
			acc = append(acc, e)
		}
	}
	return acc
}

// AppendCall calls appends the result of calling f to errp. AppendCall is
// useful when deferring calls that return an error.
func AppendCall(errp *error, f func() error) {
	*errp = Append(*errp, f())
}
