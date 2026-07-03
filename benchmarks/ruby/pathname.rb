# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference Ruby workload mirroring go/main.go for the go-ruby-pathname parity
# benchmark. Same path strings, same arguments, same iteration counts. All ops
# are PURE-PATH (lexical) — no filesystem access. `ruby pathname.rb verify`
# prints the canonical outputs run.sh diffs against the Go driver before timing.
require "pathname"
require_relative "_harness"

MAIN_PATH  = "/usr/local/lib/ruby/4.0.0/pathname.rb"
DIRTY_PATH = "/usr/local/../local/./lib//ruby/4.0.0/../4.0.0/pathname.rb"
PLUS_ARG   = "ruby/4.0.0/pathname.rb"

MAIN      = Pathname.new(MAIN_PATH)
JOIN_BASE = Pathname.new("/usr/local")
JOIN_ARGS = ["lib", "ruby", "4.0.0", "pathname.rb"]
PLUS_BASE = Pathname.new("/usr/local/lib")
REL_BASE  = Pathname.new("/usr/local/share/doc")

def outputs
  [
    "join|#{JOIN_BASE.join(*JOIN_ARGS)}",
    "plus|#{PLUS_BASE + PLUS_ARG}",
    "basename|#{MAIN.basename}",
    "dirname|#{MAIN.dirname}",
    "extname|#{MAIN.extname}",
    "cleanpath|#{Pathname.new(DIRTY_PATH).cleanpath}",
    "relative_path_from|#{MAIN.relative_path_from(REL_BASE)}",
    "split|#{MAIN.split.map(&:to_s).join("|")}",
    "each_filename|#{MAIN.each_filename.to_a.join(",")}",
  ]
end

if ARGV[0] == "verify"
  puts outputs
  exit
end

bench("join",               1000) { JOIN_BASE.join(*JOIN_ARGS) }
bench("plus",               1000) { PLUS_BASE + PLUS_ARG }
bench("basename",           1000) { MAIN.basename }
bench("dirname",            1000) { MAIN.dirname }
bench("extname",            1000) { MAIN.extname }
bench("cleanpath",          1000) { Pathname.new(DIRTY_PATH).cleanpath }
bench("relative_path_from", 1000) { MAIN.relative_path_from(REL_BASE) }
bench("split",              1000) { MAIN.split }
bench("each_filename",      1000) { MAIN.each_filename.to_a }
