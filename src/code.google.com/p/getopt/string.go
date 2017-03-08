package getopt

type stringValue string

func (s *stringValue) Set(value string, opt Option) error {
	*s = stringValue(value)
	return nil
}

func (s *stringValue) String() string {
	return string(*s)
}

// String returns a value option that stores is value as a string.  The
// initial value of the string is passed in value.
func String(name rune, value string, helpvalue ...string) *string {
	return CommandLine.String(name, value, helpvalue...)
}
func StringLong(name string, short rune, value string, helpvalue ...string) *string {
	return CommandLine.StringLong(name, short, value, helpvalue...)
}
func StringVar(p *string, name rune, helpvalue ...string) Option {
	return CommandLine.StringVar(p, name, helpvalue...)
}
func StringVarLong(p *string, name string, short rune, helpvalue ...string) Option {
	return CommandLine.StringVarLong(p, name, short, helpvalue...)
}
func (s *Set) String(name rune, value string, helpvalue ...string) *string {
	p := value
	s.StringVarLong(&p, "", name, helpvalue...)
	return &p
}
func (s *Set) StringLong(name string, short rune, value string, helpvalue ...string) *string {
	s.StringVarLong(&value, name, short, helpvalue...)
	return &value
}
func (s *Set) StringVar(p *string, name rune, helpvalue ...string) Option {
	return s.VarLong((*stringValue)(p), "", name, helpvalue...)
}
func (s *Set) StringVarLong(p *string, name string, short rune, helpvalue ...string) Option {
	return s.VarLong((*stringValue)(p), name, short, helpvalue...)
}
