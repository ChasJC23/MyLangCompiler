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
		p.ParsePrefix(preclvl)
	case 0b0001: // postfix
		p.ParsePostfix(preclvl)
	case 0b0010: // infix left associative
		p.ParseLeftAssociative(preclvl)
	case 0b0011: // infix right associative
		p.ParseRightAssociative(preclvl)
	case 0b0100: // repeatable prefix
		p.ParseRepeatablePrefix(preclvl)
	case 0b0101: // repeatable postfix
		p.ParseRepeatablePostfix(preclvl)
	case 0b0110: // repeatable infix left associative
		p.ParseRepeatableLeftAssociative(preclvl)
	case 0b0111: // repeatable infix right associative
		p.ParseRepeatableRightAssociative(preclvl)
	case 0b1010: // implied operation infix left associative
		p.ParseImpliedLeftAssociative(preclvl)
	case 0b1011: // implied operation infix right associative
		p.ParseImpliedRightAssociative(preclvl)
	case 0b1110: // implied operation repeatable infix left associative
		p.ParseImpliedRepeatableLeftAssociative(preclvl)
	case 0b1111: // implied operation repeatable infix right associative
		p.ParseImpliedRepeatableRightAssociative(preclvl)
	default:
		panic("invalid configuration")
	}
	return nil
}

func (p *Parser) ParseImpliedLeftAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseImpliedRepeatableLeftAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseImpliedRightAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseImpliedRepeatableRightAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseLeftAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseRepeatableLeftAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseRightAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParseRepeatableRightAssociative(preclvl *PrecedenceLevel)
func (p *Parser) ParsePrefix(preclvl *PrecedenceLevel)
func (p *Parser) ParseRepeatablePrefix(preclvl *PrecedenceLevel)
func (p *Parser) ParsePostfix(preclvl *PrecedenceLevel)
func (p *Parser) ParseRepeatablePostfix(preclvl *PrecedenceLevel)
