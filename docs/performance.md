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

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `Pathname`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; the Go driver's output was checked
**byte-identical to MRI** (the resulting path strings for every op) before any
timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5 — **date 2026-07-03**.
  All runtimes ran natively on the host (no VM).
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Ops:** all **pure-path (lexical)** — no filesystem access, so the workload is
  deterministic and reproducible: `Pathname#join`, `#+` (plus), `#basename`,
  `#dirname`, `#extname`, `#cleanpath`, `#relative_path_from`, `#split`, and
  `#each_filename`. Ruby's `Pathname` is pure-Ruby stdlib wrapping `File` path
  methods, so parity is algorithmically reachable — and, being compiled, the
  pure-Go library realizes it with room to spare. No op was unmeasurable: every
  one of the nine has a direct, output-equal Go counterpart.
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop (1000 ops each), timed with a monotonic clock; the **best**
  pass is reported as **ns/op** (lower is better). `vs MRI` < 1.00× means
  *faster than MRI*. Interpreter start-up is outside the timed region, so these
  are operation costs, not `ruby file.rb` process costs. Harness + workload:
  [`benchmarks/`](https://github.com/go-ruby-pathname/docs/tree/main/benchmarks)
  (`go/` pins the published library by pseudo-version; `ruby/pathname.rb`);
  reproduce with `bash benchmarks/run.sh`.

#### join

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 112.1 | 0.01× |
| MRI | 11018.0 | 1.00× |
| MRI + YJIT | 10161.0 | 0.92× |
| JRuby | 8448.9 | 0.77× |
| TruffleRuby | 6322.3 | 0.57× |

#### plus

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 25.9 | 0.01× |
| MRI | 2822.0 | 1.00× |
| MRI + YJIT | 2688.0 | 0.95× |
| JRuby | 2011.2 | 0.71× |
| TruffleRuby | 1524.3 | 0.54× |

#### basename

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 84.7 | 0.25× |
| MRI | 336.0 | 1.00× |
| MRI + YJIT | 277.0 | 0.82× |
| JRuby | 210.6 | 0.63× |
| TruffleRuby | 442.7 | 1.32× |

#### dirname

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 14.3 | 0.05× |
| MRI | 298.0 | 1.00× |
| MRI + YJIT | 227.0 | 0.76× |
| JRuby | 265.7 | 0.89× |
| TruffleRuby | 1973.8 | 6.62× |

#### extname

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 105.9 | 0.52× |
| MRI | 203.0 | 1.00× |
| MRI + YJIT | 163.0 | 0.80× |
| JRuby | 313.0 | 1.54× |
| TruffleRuby | 493.3 | 2.43× |

#### cleanpath

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 250.0 | 0.04× |
| MRI | 6044.0 | 1.00× |
| MRI + YJIT | 6058.0 | 1.00× |
| JRuby | 4995.4 | 0.83× |
| TruffleRuby | 3353.3 | 0.55× |

#### relative_path_from

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 578.2 | 0.05× |
| MRI | 12150.0 | 1.00× |
| MRI + YJIT | 11913.0 | 0.98× |
| JRuby | 9125.1 | 0.75× |
| TruffleRuby | 3648.9 | 0.30× |

#### split

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 132.5 | 0.19× |
| MRI | 703.0 | 1.00× |
| MRI + YJIT | 544.0 | 0.77× |
| JRuby | 389.7 | 0.55× |
| TruffleRuby | 881.9 | 1.25× |

#### each_filename

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 166.1 | 0.06× |
| MRI | 2885.0 | 1.00× |
| MRI + YJIT | 2711.0 | 0.94× |
| JRuby | 2173.6 | 0.75× |
| TruffleRuby | 4788.2 | 1.66× |

The pure-Go library is **faster than MRI on every one of the nine ops**. The
widest margins are on the allocation-heavy path builders — `join` and `+`
(~100× and ~110× faster than MRI, since MRI's `Pathname#join` chains `#+`, each
allocating a fresh `Pathname` and string) and `cleanpath` /
`relative_path_from` (~24× and ~21×, which split, normalise and re-join). The
narrowest are the already-tiny single-scan ops `extname` (~1.9×) and `basename`
(~4×), where MRI has little Ruby overhead to shed. MRI's `Pathname` is pure-Ruby
stdlib, so each op pays interpreter dispatch the compiled Go code does not.
Parity is not just reached, it is beaten across the board. Outputs were checked
byte-identical to MRI before timing, so the speed is not bought with divergent
behavior.

!!! note "Cold-JIT caveat"
    JRuby (JVM) and TruffleRuby (GraalVM) get only the harness's 3 warm-up
    passes, so they are effectively **cold** — TruffleRuby's `dirname` row
    (1973.8 ns, 6.62×) and `extname` row (2.43×) are dominated by Graal
    compilation kicking in on just-touched paths, not steady-state throughput;
    on the heavier `relative_path_from`/`cleanpath`/`join` ops, where it has more
    work to amortise over, TruffleRuby is already the fastest Ruby column. JRuby
    is closer to warm but still not fully JITed. Read the JVM/Graal columns as
    cold single-shot costs, the same footing as the top-of-page table. These are
    **real measured numbers** from the 2026-07-03 run (Apple M4 Max, arm64;
    `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`, `truffleruby 34.0.1`) — nothing is
    fabricated or cherry-picked.
