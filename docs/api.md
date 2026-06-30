# Usage & API

The public API lives at the module root (`github.com/go-ruby-pathname/pathname`). It is
**Ruby-shaped but Go-idiomatic**: the methods mirror Ruby's `Pathname` (`basename`,
`cleanpath`, `relative_path_from`, …), while the surface follows Go conventions — value
types, an explicit `error` where Ruby raises, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-pathname/pathname`, bound into `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-pathname/pathname
```

## Worked example

```go
package main

import (
	"fmt"

	"github.com/go-ruby-pathname/pathname"
)

func main() {
	p := pathname.New("a/./b/../c")
	fmt.Println(p.Cleanpath().ToS()) // a/c

	fmt.Println(pathname.New("/usr/bin/ruby.rb").BasenameSuffix(".*").ToS()) // ruby
	fmt.Println(pathname.New("foo.tar.gz").Extname())                        // .gz

	rel, _ := pathname.New("/a/b/c").RelativePathFrom(pathname.New("/a/x/y"))
	fmt.Println(rel.ToS()) // ../../b/c

	fmt.Println(pathname.New("a").JoinStrings("b", "/c", "d").ToS()) // /c/d
}
```

## Shape

`*Pathname` is an **immutable wrapper** over a path string; every transforming method
returns a fresh `*Pathname`, `string`, `[]*Pathname` or `bool`, mirroring the
corresponding Ruby method.

- **`New` / `ToS` / `Inspect`** — wrap a path, render it, `#<Pathname:…>`.
- **`Basename` / `BasenameSuffix`** — last component, with optional suffix strip
  (`".*"` strips any extension).
- **`Dirname` / `Parent`** — all but the last component (`.`, `/` edge cases).
- **`Extname` / `SubExt`** — MRI's `File.extname` rule **byte-for-byte**, including the
  fiddly trailing-dot behaviour (`"foo." → "."`, `"a..b" → ".b"`) and the dotfile rule
  (`".foo" → ""`).
- **`Cleanpath`** — collapse `.`, `..` and redundant separators; a `..` past the root of
  an absolute path is dropped, a `..` escaping a relative path is kept
  (`"a/./b/../c" → "a/c"`, `"/a/../../b" → "/b"`, `"a/../../b" → "../b"`).
- **`RelativePathFrom`** — MRI's lexical `relative_path_from` algorithm, returning
  `(*Pathname, error)` where the error is an `*ArgumentError` carrying MRI's exact
  message (the two cases: mixing absolute/relative; a `..` in the base).
- **`Plus` / `Join`** (`+`, `join`) — append components; an absolute component resets to
  the root.
- **`Split` / `EachFilename` / `Filenames`** — `[dirname, basename]` and the non-empty
  component list (`Ascend`, `Descend` and `EachFilename` take a `func` callback — the Go
  analogue of Ruby's block; `Filenames` returns the slice form).
- **`Ascend` / `Descend`** — the path and each parent up to the root (or the first
  relative component), and the reverse.
- **`Absolute` / `Relative` / `Root`** — the predicate trio.
- **`Cmp` / `Eql` / `Hash`** — `<=>`, `==`/`eql?`, a stable FNV hash.

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** shells out to the
`ruby` binary (gated on `RUBY_VERSION >= "4.0"`) and checks each method against MRI's
`Pathname` over a corpus of edge cases — the trailing-slash / `.` / `..` / root-escape
cleanpath cases, the trailing-dot extname cases, and the `relative_path_from`
`ArgumentError` messages — comparing results **byte-for-byte**, not approximated from
memory. The oracle skips itself where `ruby` is absent (the qemu cross-arch lanes and the
Windows lane); the deterministic, ruby-free tests alone hold coverage at 100% there.

## Relationship to Ruby

`go-ruby-pathname/pathname` is **standalone and reusable**, and is the path-algebra
backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo`
as a native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs the other
way: this library has no dependency on the Ruby runtime.
