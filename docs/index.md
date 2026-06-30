# go-ruby-pathname documentation

**Ruby's `Pathname` path algebra in pure Go — MRI-compatible, filesystem-free, no cgo.**

`go-ruby-pathname/pathname` is a faithful, pure-Go (zero cgo) reimplementation of the
**path-manipulation surface** of Ruby's `Pathname` standard library, matching reference
Ruby (MRI) byte-for-byte. The module path is `github.com/go-ruby-pathname/pathname`.

Ruby's `Pathname` is two libraries in one: a pure **path algebra** (lexical string
surgery, no I/O) and a thin set of delegations to `File`/`Dir`/`IO`. This package is the
**first half** — fully deterministic, needs no interpreter, lives here as pure Go.
The filesystem-touching half stays host-side in
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby), where it forwards to the
host's `File` class. The dependency runs the other way: this library has **no dependency
on the Ruby runtime**, and `rbgo` binds it as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb).

!!! success "Status: path algebra complete — MRI byte-exact"
    Faithful port of `Pathname`'s pure path methods: **`Basename`** / **`Dirname`** / **`Extname`** / **`Cleanpath`**, **`RelativePathFrom`** (with its two `ArgumentError` cases), **`Plus`/`Join`**, **`Split`/`EachFilename`**, **`Ascend`/`Descend`**, the **`Absolute`/`Relative`/`Root`** predicates and **`Cmp`/`Eql`/`Hash`**. Validated by a **differential oracle** against the system `ruby` — every method compared against MRI's `Pathname` byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes (incl. Windows).

## Lexical, `/`-based, platform-independent

Ruby's `Pathname` hardcodes `"/"` as the component separator for **every** lexical
operation on **every** platform — it never consults the OS path separator for
`join`/`split`/`cleanpath`/`basename`. This package does the same, so the behaviour
(and the test suite) is identical on Linux, macOS and Windows. It builds and passes
`GOOS=windows` with no OS-specific path code.

## Quick taste

```go
p := pathname.New("a/./b/../c")
fmt.Println(p.Cleanpath().ToS()) // a/c

fmt.Println(pathname.New("/usr/bin/ruby.rb").BasenameSuffix(".*").ToS()) // ruby
fmt.Println(pathname.New("foo.tar.gz").Extname())                        // .gz

rel, _ := pathname.New("/a/b/c").RelativePathFrom(pathname.New("/a/x/y"))
fmt.Println(rel.ToS()) // ../../b/c
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`pathname`](https://github.com/go-ruby-pathname/pathname) | the library — Ruby's `Pathname` path algebra in pure Go |
| [`docs`](https://github.com/go-ruby-pathname/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-pathname.github.io`](https://github.com/go-ruby-pathname/go-ruby-pathname.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-pathname/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **Lexical & filesystem-free** — only the deterministic path algebra, `"/"`-based
  on every platform; the `File`/`Dir` delegations stay host-side.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** No dependency on the Ruby runtime — the dependency
  runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-pathname/pathname](https://github.com/go-ruby-pathname/pathname).
