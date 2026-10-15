package assembler

import (
	"fmt"
	"slices"
	"strings"
)

type Tokentype int

var instructions = []string{
	"nop",
	"hlt",
	"ldi",
	"sti",
	"ld",
	"st",
	"mov",
	"str",
	"inc",
	"dec",
	"cmp",
	"add",
	"sub",
	"mul",
	"div",
	"jmp",
	"jnz",
}

const (
	INSTRUCTION Tokentype = iota
	NUMBER
	REGISTER
	COMMA
	LABEL
	IDENT
)

type Token struct {
	Type  Tokentype
	Value string
}

type Lexer struct {
	Code    string
	Pointer int
	Tokens  []Token
}

func NewLexer(code string) Lexer {
	return Lexer{
		Code:    code,
		Pointer: 0,
		Tokens:  []Token{},
	}
}

func (l *Lexer) Lex() []Token {
	for l.Pointer < len(l.Code) {
		switch l.Code[l.Pointer] {
		case ';':
			fmt.Println("comment detected")
			for l.Peek() != '\n' {
				l.Advance()
			}

		case '$':
			fmt.Println("register detected")
			var tmp strings.Builder
			for !l.IsWhitespace(l.Peek()) && l.Peek() != ',' {
				tmp.WriteByte(l.Peek())
				l.Advance()
			}

			l.AddToken(REGISTER, tmp.String())

		case ',':
			l.AddToken(COMMA, ",")

		case ' ':
		case '\n':
		case '\t':

		default:
			if l.IsNumber(l.Code[l.Pointer]) {
				var tmp strings.Builder
				tmp.WriteByte(l.Code[l.Pointer])
				for l.IsNumber(l.Peek()) {
					tmp.WriteByte(l.Peek())
					l.Advance()
				}

				number := tmp.String()
				l.AddToken(NUMBER, number)
			} else {
				var tmp strings.Builder
				tmp.WriteByte(l.Code[l.Pointer])
				for !l.IsWhitespace(l.Peek()) {
					tmp.WriteByte(l.Peek())
					l.Advance()
				}

				str := tmp.String()
				if strings.HasSuffix(str, ":") {
					l.AddToken(LABEL, strings.TrimSuffix(str, ":"))
				} else if slices.Contains(instructions, strings.ToLower(str)) {
					l.AddToken(INSTRUCTION, str)
				} else {
					l.AddToken(IDENT, str)
				}
			}
		}
		l.Advance()
	}
	return l.Tokens
}

func (l *Lexer) Advance() {
	l.Pointer++
}

func (l *Lexer) Peek() byte {
	return l.Code[l.Pointer+1]
}

func (l *Lexer) IsWhitespace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t'
}

func (l *Lexer) AddToken(t Tokentype, value string) {
	l.Tokens = append(l.Tokens, Token{t, value})
}

func (l *Lexer) IsNumber(value byte) bool {
	return ('0' <= value && value <= '9')
}
