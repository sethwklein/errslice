package errslice_test

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sethwklein/errslice"
)

type angryWriteCloser struct{}

func (_ angryWriteCloser) Write(b []byte) (int, error) {
	return 0, errors.New("no writing!")
}
func (_ angryWriteCloser) Close() error {
	return errors.New("no closing!")
}

// BE CAREFUL! The return value must be named. Else the error returned from the
// deferred call will be silently discarded with no compiler error.

func exampleWrite(w io.WriteCloser) (err error) {
	defer errslice.AppendCall(&err, w.Close)

	for c := byte('a'); c <= 'z'; c++ {
		_, err = w.Write([]byte{c})
		if err != nil {
			return err
		}
	}
	return nil
}

func ExampleAppendCall() {
	err := exampleWrite(angryWriteCloser{})
	for _, e := range errslice.From(err) {
		// In real code, use os.Stderr, but that's broken in examples.
		fmt.Fprintln(os.Stdout, "Error:", e)
	}
	// Output:
	// Error: no writing!
	// Error: no closing!
}
