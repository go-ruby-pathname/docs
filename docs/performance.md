# Performance

`go-ruby-pathname/pathname` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `Pathname` path
algebra. This page records the **methodology** for the comparative benchmark of that
module against the reference Ruby runtimes — part of the ecosystem-wide per-module parity
suite. No numbers are quoted here until they are measured on a pinned host and committed
to the benchmark harness.

## What is measured

The **same** Ruby script — a representative mix of `Pathname` lexical operations
(`cleanpath`, `relative_path_from`, `basename`/`dirname`/`extname`, `join`/`split`,
`ascend`) over a corpus of paths — is run under every runtime. `rbgo`'s number reflects
**this pure-Go library doing the work**; every other column is that interpreter's own
`pathname` stdlib. So the comparison is the **Ruby-visible operation**, apples-to-apples
across interpreters. The script prints a deterministic checksum and its output is checked
**byte-identical to MRI** before any timing is taken.

## Method

- **Best-of-N wall time** (best, not mean, to suppress scheduler noise); single-shot
  processes, no warm-up beyond the script's own loop.
- **Runtimes:** MRI (the oracle) and MRI `--yjit`; JRuby (on OpenJDK); TruffleRuby
  (GraalVM CE Native) — each timed cold, single-shot, so JVM / Graal startup is on every
  run; read those as one-shot `ruby file.rb` costs, not steady-state JIT numbers.
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules).
  Reproduce with `RBGO=./rbgo bash bench/modules/run.sh N`.

!!! note "Honest framing"
    Because `Pathname`'s lexical operations are individually cheap, a round-trip
    completes fast and the relative noise is high; treat any ratio that completes in well
    under ~200 ms as order-of-magnitude. Numbers will be added here only once they are
    real, measured on a pinned host, and reproducible from the committed harness — nothing
    is estimated or cherry-picked.
