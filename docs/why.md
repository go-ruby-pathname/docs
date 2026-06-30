# Why pure Go

`go-ruby-pathname/pathname` reimplements the path-manipulation surface of Ruby's
`Pathname` in **pure Go, with cgo disabled**. The slice of `Pathname` it covers is
**deterministic and interpreter-independent**: given a path string and arguments, the
result is a pure function of those inputs — no live binding, no filesystem access, no
evaluation of arbitrary Ruby. That is exactly the part that can — and should — live as a
standalone Go library, separate from the interpreter.

## What it is — and isn't

Ruby's `Pathname` is two libraries in one:

- a pure **path algebra** — lexical string surgery (`cleanpath`, `relative_path_from`,
  `basename`/`dirname`/`extname`, `join`/`split`/`ascend`), which never touches disk;
- a thin set of **delegations** to `File`/`Dir`/`IO` (`read`, `write`, `exist?`,
  `children`, …), which do.

This package is the **first half**. The filesystem-touching half stays host-side in
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby), where it forwards to the
host's `File` class; it is deliberately **out of scope** here.

## Extracted from rbgo, reusable by anyone

The path algebra has been **extracted into a reusable standalone library** so that:

- any Go program can import `github.com/go-ruby-pathname/pathname` directly, with no
  Ruby runtime;
- the dependency runs the *other* way — `rbgo` binds this module as a native module (the
  same pattern as [go-ruby-regexp](https://github.com/go-ruby-regexp) and
  [go-ruby-erb](https://github.com/go-ruby-erb)), rather than this module depending on
  the interpreter;
- the behaviour is pinned by a **differential oracle** against the system `ruby`,
  independent of any one consumer.

## Lexical, `/`-based, platform-independent

Ruby's `Pathname` hardcodes `"/"` as the component separator for **every** lexical
operation on **every** platform — it never consults the OS path separator. This package
does the same: it uses `"/"` unconditionally, so the behaviour and the test suite are
identical on Linux, macOS and Windows, and the module builds and passes `GOOS=windows`
with no OS-specific path code. (Ruby's own `File.basename` is `\`-aware on Windows;
`Pathname`'s lexical methods are not, and neither is this port.)

## Why pure Go matters here

Because the library is CGO-free and dependency-free, it:

- cross-compiles to every Go target with no C toolchain, and links into a single static
  binary;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `ruby` binary wherever one is on `PATH`,
  while the cross-arch and Windows lanes (where `ruby` is absent) still validate the
  library itself.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is in scope.
