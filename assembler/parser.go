package assembler

import (
	"regexp"

	"github.com/annaxiomm/moa/errors"
)

type StatementType int
type OperandType int

const (
	PLABEL StatementType = iota
	PINSTRUCTION
)

const (
	OREGISTER OperandType = iota
	OLABEL
	ONUMBER
)

type Parser struct {
	Tokens  []Token
	Pointer int
}

type Program struct {
	Statements []Statement
}

type Statement struct {
	Type     StatementType
	Name     string
	Operands []Operand
}

type Operand struct {
	Type  OperandType
	Value string
}

func NewParser(tokens []Token) Parser {
	return Parser{
		Tokens:  tokens,
		Pointer: 0,
	}
}

func NewStatement(t StatementType, name string, operands []Operand) Statement {
	return Statement{
		Type:     t,
		Name:     name,
		Operands: operands,
	}
}

func (p *Parser) Parse() (Program, error) {
	program := Program{Statements: []Statement{}}
	for p.Cur().Type != EOF {
		if p.Cur().Type == NEWLINE {
			p.Advance()
			continue
		}
		statement, err := p.ParseStatement()
		if err != nil {
			return program, err
		}

		program.Statements = append(program.Statements, statement)
		p.Advance()
	}

	return program, nil
}

func (p *Parser) Peek() Token {
	return p.Tokens[p.Pointer+1]
}

func (p *Parser) IsEnd(tok Token) bool {
	return tok.Type == NEWLINE || tok.Type == EOF
}

func (p *Parser) Cur() Token {
	return p.Tokens[p.Pointer]
}

func (p *Parser) Advance() {
	p.Pointer++
}

func (p *Parser) ParseStatement() (Statement, error) {
	token := p.Tokens[p.Pointer]
	switch token.Type {
	case LABEL:
		return NewStatement(PLABEL, token.Value, []Operand{}), nil
	case INSTRUCTION:
		return p.ParseInstruction()
	default:
		return Statement{}, &errors.ParseError{
			Type:     errors.UnexpectedToken,
			Tok:      token.Value,
			Expected: "instruction or label",
		}
	}
}

func (p *Parser) ParseInstruction() (Statement, error) {
	instruction := p.Cur()
	statement := Statement{Type: PINSTRUCTION, Name: instruction.Value}
	operands := []Operand{}

	if p.IsEnd(p.Peek()) {
		return statement, nil
	}

	for !p.IsEnd(p.Peek()) {
		p.Advance()
		cur := p.Cur()
		switch cur.Type {
		case NUMBER:
			operands = append(operands, Operand{Type: ONUMBER, Value: cur.Value})
		case REGISTER:
			if x, _ := regexp.MatchString("^R[0-7]$", cur.Value); !x {
				return Statement{}, &errors.ParseError{Type: errors.InvalidRegister, Tok: "$" + cur.Value}
			}
			operands = append(operands, Operand{Type: OREGISTER, Value: string(cur.Value[1])})
		case LABEL:
			operands = append(operands, Operand{Type: OLABEL, Value: cur.Value})
		default:
			return Statement{}, &errors.ParseError{Type: errors.UnexpectedToken, Tok: cur.Value, Expected: "number, register, or label"}
		}

		if p.IsEnd(p.Peek()) {
			break
		}

		p.Advance()
		if p.Cur().Type != COMMA {
			return Statement{}, &errors.ParseError{Type: errors.MissingComma}
		}

		if p.IsEnd(p.Peek()) {
			return Statement{}, &errors.ParseError{Type: errors.MissingOperand}
		}
	}

	statement.Operands = operands

	return statement, nil
}
