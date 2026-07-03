// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level benchmark driver for the pure-Go go-ruby-pathname/pathname
// library. Mirrors ruby/pathname.rb exactly: same path strings, same arguments,
// same iteration counts. All ops are PURE-PATH (lexical) — no filesystem access.
// Run with the single arg "verify" to print the canonical outputs run.sh diffs
// against MRI before any timing.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ruby-pathname/pathname"
)

// Inputs are byte-identical to ruby/pathname.rb.
const (
	mainPath  = "/usr/local/lib/ruby/4.0.0/pathname.rb"
	dirtyPath = "/usr/local/../local/./lib//ruby/4.0.0/../4.0.0/pathname.rb"
	plusArg   = "ruby/4.0.0/pathname.rb"
)

var (
	mainPN   = pathname.New(mainPath)
	joinBase = pathname.New("/usr/local")
	joinArgs = []string{"lib", "ruby", "4.0.0", "pathname.rb"}
	plusBase = pathname.New("/usr/local/lib")
	relBase  = pathname.New("/usr/local/share/doc")
)

func splitStr() string {
	s := mainPN.Split()
	return s[0].String() + "|" + s[1].String()
}

func eachStr() string {
	var b strings.Builder
	mainPN.EachFilename(func(f string) {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(f)
	})
	return b.String()
}

func relStr() string {
	r, err := mainPN.RelativePathFrom(relBase)
	if err != nil {
		return "ERR:" + err.Error()
	}
	return r.String()
}

// outputs renders every op's result as "op|value" lines, matching pathname.rb.
func outputs() []string {
	return []string{
		"join|" + joinBase.JoinStrings(joinArgs...).String(),
		"plus|" + plusBase.PlusString(plusArg).String(),
		"basename|" + mainPN.Basename().String(),
		"dirname|" + mainPN.Dirname().String(),
		"extname|" + mainPN.Extname(),
		"cleanpath|" + pathname.New(dirtyPath).Cleanpath().String(),
		"relative_path_from|" + relStr(),
		"split|" + splitStr(),
		"each_filename|" + eachStr(),
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "verify" {
		for _, l := range outputs() {
			fmt.Println(l)
		}
		return
	}
	bench("join", 1000, func() { sink = joinBase.JoinStrings(joinArgs...) })
	bench("plus", 1000, func() { sink = plusBase.PlusString(plusArg) })
	bench("basename", 1000, func() { sink = mainPN.Basename() })
	bench("dirname", 1000, func() { sink = mainPN.Dirname() })
	bench("extname", 1000, func() { sink = mainPN.Extname() })
	bench("cleanpath", 1000, func() { sink = pathname.New(dirtyPath).Cleanpath() })
	bench("relative_path_from", 1000, func() { r, _ := mainPN.RelativePathFrom(relBase); sink = r })
	bench("split", 1000, func() { sink = mainPN.Split() })
	bench("each_filename", 1000, func() {
		var fs []string
		mainPN.EachFilename(func(f string) { fs = append(fs, f) })
		sink = fs
	})
}
