# Roadmap

`go-ruby-pathname/pathname` is grown **test-first**, each capability differential-tested
against MRI rather than built in isolation. The path-algebra slice of Ruby's `Pathname` —
the deterministic, interpreter-independent, filesystem-free core — is **complete**.

| Stage | What | Status |
| --- | --- | --- |
| Basename / dirname / extname | `Basename`/`BasenameSuffix`, `Dirname`/`Parent`, `Extname`/`SubExt` — MRI's `File.extname` rule byte-for-byte, including trailing-dot (`"foo." → "."`) and dotfile (`".foo" → ""`) cases. | **Done** |
| Cleanpath | Collapse `.`, `..` and redundant separators exactly as MRI does, with the absolute-root-drop / relative-escape-keep distinction. | **Done** |
| RelativePathFrom | MRI's lexical `relative_path_from`, with its two `ArgumentError` cases and message text matching MRI. | **Done** |
| Join / split / traversal | `Plus`/`Join` (absolute component resets to root); `Split`/`EachFilename`/`Filenames`; `Ascend`/`Descend` via Go-callback blocks. | **Done** |
| Predicates & comparison | `Absolute`/`Relative`/`Root`; `Cmp` (`<=>`), `Eql` (`==`/`eql?`), stable FNV `Hash`. | **Done** |
| Differential oracle & coverage | Every method checked against the system `ruby`/`Pathname` over a corpus of edge cases; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No filesystem.** The library is the lexical path algebra only; `read`, `write`,
  `exist?`, `children` and the other `File`/`Dir`/`IO` delegations stay host-side in
  go-embedded-ruby, where they forward to the host's `File` class. The library never
  touches disk.
- **Lexical, `"/"`-based, every platform.** Like MRI's `Pathname`, the separator is
  always `"/"`; there is no OS-specific path code, and the behaviour is identical on
  Linux, macOS and Windows.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour, pinned by the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime; the
  dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic / filesystem split.
