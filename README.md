# go-hash

Wrap smaller hash libraries with interfaces similar to those in standard library hash interface.

## Background

The Go standard library defines the [`hash`](https://pkg.go.dev/hash) package.
This package defines a [`Hash`](https://pkg.go.dev/hash#Hash) interface
and then specific size interfaces are derived for
[`Hash32`](https://pkg.go.dev/hash#Hash32)
and [`Hash64`](https://pkg.go.dev/hash#Hash64).
Subdirectories provide implementations of these interfaces for specific hash functions
as well as the hashing function used by Go maps.

One feature of the `Hash` interface is that it derives from `io.Writer`.
This allows types that implement `Hash` to be used with
[`io.Copy()`](https://pkg.go.dev/io#Copy).

## Implementation

This project provides missing interfaces:

* `Hash16`
* `Hash8`

and wrappers around pre-existing hash libraries that implement these interfaces.
