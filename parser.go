package main

import "container/list"

type Parser struct {
	tokeniser *Tokeniser
	opctx     *OpContext
}

func NewParser(tokeniser *Tokeniser) *Parser {
	parser := new(Parser)
	parser.tokeniser = tokeniser
	parser.opctx = tokeniser.opctx
	return parser
}

// function responsible for parsing an entire source file
func (p *Parser) ParseSource() interface{} {

	result := p.ParseCodeBlock()

	// I realise this is more of an expressions thing,
	// but let's include it just in case.
	if p.tokeniser.currToken != EOF_TOKEN {
		panic("Unexpected characters at end of source file")
	}

	return result
}

// function responsible for parsing a specific code block
func (p *Parser) ParseCodeBlock() interface{} {
	statements := make([]interface{}, 0)
	for p.tokeniser.currToken != EOF_TOKEN && p.tokeniser.currToken != CLOSE_CODE_BLOCK {
		statements = append(statements, p.ParseStatement())
	}
	return statements
}

// function responsible for parsing any arbitrary statement.
// This may be any individual line of code.
func (p *Parser) ParseStatement() interface{} {
	preclvlel := p.opctx.precList.Front()
	result := p.ParsePrecisionLevel(preclvlel)
	return result
}

func (p *Parser) ParsePrecisionLevel(preclvlel *list.Element) interface{} {
	preclvl := preclvlel.Value.(*PrecedenceLevel)
	// check bitmask in preclevel.go
	switch preclvl.properties & 0b1111 {
	case 0b0000: // prefix
	case 0b0001: // postfix
	case 0b0010: // infix left associative
	case 0b0011: // infix right associative
	case 0b0100: // repeatable prefix
	case 0b0101: // repeatable postfix
	case 0b0110: // repeatable infix left associative
	case 0b0111: // repeatable infix right associative
	case 0b1010: // implied operation infix left associative
	case 0b1011: // implied operation infix right associative
	case 0b1110: // implied operation repeatable infix left associative
	case 0b1111: // implied operation repeatable infix right associative
	default:
		panic("invalid configuration")
	}
	return nil
}
