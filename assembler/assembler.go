package assembler

import "fmt"

func Assemble(input string) error {
	l := NewLexer(input)
	tokens, err := l.Lex()
	if err != nil {
		return err
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		return err
	}
	for _, s := range program.Statements {
		fmt.Println(s.Name, s.Operands)
	}
	return nil
}
