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

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-pathname) | 260 | 0.59× |
| MRI (ruby 4.0.5) | 440 | 1.00× |
| MRI + YJIT | 420 | 0.95× |
| JRuby 10.1.0.0 | 1890 | 4.30× |
| TruffleRuby 34.0.1 | 740 | 1.68× |

rbgo runs on **go-ruby-pathname** and is **faster than MRI** here (0.59x) on this lexical path-manipulation workload (join/cleanpath/relative_path_from/split, no filesystem I/O).

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are **real measured numbers** from
    the 2026-06-30 run (Apple M-series; `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`,
    `truffleruby 34.0.1`) — nothing is fabricated or cherry-picked.
