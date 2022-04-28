# errslice

`errslice` is a [Go](https://golang.org/) package for handling multiple
errors and returning them to callers that may not be aware they're dealing
with more than one.

## Usage

`go get github.com/sethwklein/errslice`

`import "github.com/sethwklein/errslice"`

https://pkg.go.dev/github.com/sethwklein/errslice

## Rationale

When writing a file, either `Write` or `Close` may return an error. While
Unix files don't, technically nothing stops some weird `io.WriteCloser`
from returning interesting and unique errors from both. This library makes
it trivial to just return however many errors you get. It also comes in
handy in a surprising number of situations where you want to return all
the errors, not just the first. Your callers need not be aware that they
may be receiving more than one error, although they can be.

## History

This is the fourth iteration of this concept.

The first was an experiment in an unreleased project. I liked it enough to
release the second.

The second was a drop in replacement for the standard `errors` package. This
was unpopular with the community, confusing, and an extra maintenance burden
so I dropped it to get the third.

The third hid the implementation in an opaque type. Experience using it
revealed few advantages and the occasional case where I had to reimplement
it using a plain old slice.

Hence this, fourth, version that just uses a good old fashioned slice.
