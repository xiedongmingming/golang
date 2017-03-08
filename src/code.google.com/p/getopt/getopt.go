package getopt

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

var DisplayWidth = 80
var HelpColumn = 20

func (s *Set) PrintUsage(w io.Writer) { // 将使用说明打印到指定输出中

	sort.Sort(s.options) // 对选项列表进行排序

	flags := "" // 表示定义并赋值

	max := 4

	for _, opt := range s.options {

		if opt.name == "" {
			opt.name = "value"
		}

		if opt.uname == "" {
			opt.uname = opt.usageName()
		}

		if max < len(opt.uname) && len(opt.uname) <= HelpColumn-3 {
			max = len(opt.uname)
		}

		if opt.flag && opt.short != 0 && opt.short != '-' {
			flags += string(opt.short)
		}
	}

	var opts []string

	if s.shortOptions['-'] != nil {
		opts = append(opts, "-")
	}

	if flags != "" {
		opts = append(opts, "-"+flags)
	}

	for _, opt := range s.options {
		if opt.flag {
			if opt.short != 0 {
				continue
			}
			flags = "--" + opt.long
		} else if opt.short != 0 {
			flags = "-" + string(opt.short) + " " + opt.name
		} else {
			flags = "--" + string(opt.long) + " " + opt.name
		}
		opts = append(opts, flags)
	}

	flags = strings.Join(opts, "] [")

	if flags != "" {
		flags = " [" + flags + "]"
	}

	if s.parameters != "" {
		flags += " " + s.parameters
	}
	fmt.Fprintf(w, "Usage: %s%s\n", s.program, flags)

	for _, opt := range s.options {
		if opt.uname != "" {
			opt.help = strings.TrimSpace(opt.help)
			if len(opt.help) == 0 {
				fmt.Fprintf(w, " %s\n", opt.uname)
				continue
			}
			help := strings.Split(opt.help, "\n")
			if len(help) == 1 {
				help = breakup(help[0], DisplayWidth-HelpColumn)
			}
			if len(opt.uname) <= max {
				fmt.Fprintf(w, " %-*s  %s\n", max, opt.uname, help[0])
				help = help[1:]
			} else {
				fmt.Fprintf(w, " %s\n", opt.uname)
			}
			for _, s := range help {
				fmt.Fprintf(w, " %-*s  %s\n", max, " ", s)
			}
		}
	}
}

func breakup(s string, max int) []string { // breakup breaks s up into strings no longer than max bytes.

	var a []string

	for {

		for len(s) > 0 && s[0] == ' ' { // strip leading spaces
			s = s[1:]
		}

		if len(s) <= max { // If the option is no longer than the max just return it
			if len(s) != 0 {
				a = append(a, s)
			}
			return a
		}
		x := max
		for s[x] != ' ' {
			// the first word is too long?!
			if x == 0 {
				x = max
				for x < len(s) && s[x] != ' ' {
					x++
				}
				if x == len(s) {
					x--
				}
				break
			}
			x--
		}
		for s[x] == ' ' {
			x--
		}
		a = append(a, s[:x+1])
		s = s[x+1:]
	}
	panic("unreachable")
}

func (s *Set) Parse(args []string) { // 对参数进行解析

	// Parse uses Getopt to parse args using the options set for s.  The first
	// element of args is used to assign the program for s if it is not yet set.  On
	// error, Parse displays the error message as well as a usage message on
	// standard error and then exits the program.

	if err := s.Getopt(args, nil); err != nil {

		fmt.Fprintln(os.Stderr, err)

		s.usage()

		os.Exit(1)
	}

}

// Parse uses Getopt to parse args using the options set for s.  The first
// element of args is used to assign the program for s if it is not yet set.
// Getop calls fn, if not nil, for each option parsed.
//
// Getopt returns nil when all options have been processed (a non-option
// argument was encountered, "--" was encountered, or fn returned false).
//
// On error getopt returns a refernce to an InvalidOption (which implements
// the error interface).
func (s *Set) Getopt(args []string, fn func(Option) bool) (err error) {
	s.State = InProgress
	defer func() {
		if s.State == InProgress {
			switch {
			case err != nil:
				s.State = Failure
			case len(s.args) == 0:
				s.State = EndOfArguments
			default:
				s.State = Unknown
			}
		}
	}()
	if fn == nil {
		fn = func(Option) bool { return true }
	}
	if len(args) == 0 {
		return nil
	}

	if s.program == "" {
		s.program = path.Base(args[0])
	}
	args = args[1:]
Parsing:
	for len(args) > 0 {
		arg := args[0]
		s.args = args
		args = args[1:]

		// end of options?
		if arg == "" || arg[0] != '-' {
			s.State = EndOfOptions
			return nil
		}

		if arg == "-" {
			goto ShortParsing
		}

		// explicitly request end of options?
		if arg == "--" {
			s.args = args
			s.State = DashDash
			return nil
		}

		// Long option processing
		if len(s.longOptions) > 0 && arg[1] == '-' {
			e := strings.IndexRune(arg, '=')
			var value string
			if e > 0 {
				value = arg[e+1:]
				arg = arg[:e]
			}
			opt := s.longOptions[arg[2:]]
			// If we are processing long options then --f is -f
			// if f is not defined as a long option.
			// This lets you say --f=false
			if opt == nil && len(arg[2:]) == 1 {
				opt = s.shortOptions[rune(arg[2])]
			}
			if opt == nil {
				return unknownOption(arg[2:])
			}
			opt.isLong = true
			// If we require an option and did not have an =
			// then use the next argument as an option.
			if !opt.flag && e < 0 && !opt.optional {
				if len(args) == 0 {
					return missingArg(opt)
				}
				value = args[0]
				args = args[1:]
			}
			opt.count++

			if err := opt.value.Set(value, opt); err != nil {
				return setError(opt, value, err)
			}

			if !fn(opt) {
				s.State = Terminated
				return nil
			}
			continue Parsing
		}

		// Short option processing
		arg = arg[1:] // strip -
	ShortParsing:
		for i, c := range arg {
			opt := s.shortOptions[c]
			if opt == nil {
				// In traditional getopt, if - is not registered
				// as an option, a lone - is treated as
				// if there were a -- in front of it.
				if arg == "-" {
					s.State = Dash
					return nil
				}
				return unknownOption(c)
			}
			opt.isLong = false
			opt.count++
			var value string
			if !opt.flag {
				value = arg[1+i:]
				if value == "" && !opt.optional {
					if len(args) == 0 {
						return missingArg(opt)
					}
					value = args[0]
					args = args[1:]
				}
			}
			if err := opt.value.Set(value, opt); err != nil {
				return setError(opt, value, err)
			}
			if !fn(opt) {
				s.State = Terminated
				return nil
			}
			if !opt.flag {
				continue Parsing
			}
		}
	}
	s.args = []string{}
	return nil
}
