package assembler

func Assemble(input string) ([]uint8, error) {
	l := NewLexer(input)
	tokens := l.Lex()

	p := NewParser(tokens)
	p.Parse()
	return []uint8{}, nil
}
