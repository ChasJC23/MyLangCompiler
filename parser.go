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
func (p *Parser) ParseSource() AST {

	result := p.ParseCodeBlock()

	// I realise this is more of an expressions thing,
	// but let's include it just in case.
	if p.tokeniser.currToken != EOF_TOKEN {
		panic("Unexpected characters at end of source file")
	}

	return result
}

// function responsible for parsing a specific code block
func (p *Parser) ParseCodeBlock() AST {
	statements := make([]AST, 0)
	for p.tokeniser.currToken != EOF_TOKEN && p.tokeniser.currToken != CLOSE_CODE_BLOCK {
		statements = append(statements, p.ParseStatement())
	}
	return NewCodeBlock(statements)
}

// function responsible for parsing any arbitrary statement.
// This may be any individual line of code.
func (p *Parser) ParseStatement() AST {
	preclvlel := p.opctx.precList.Front()
	result := p.ParsePrecedenceLevel(preclvlel)
	return result
}

func (p *Parser) ParsePrecedenceLevel(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(*PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	// check bitmask in preclevel.go
	switch preclvl.properties & 0b1111 {
	case 0b0000: // prefix
		return p.ParsePrefix(preclvl)
	case 0b0001: // postfix
		return p.ParsePostfix(preclvl)
	case 0b0010: // infix left associative
		return p.ParseLeftAssociative(preclvl)
	case 0b0011: // infix right associative
		return p.ParseRightAssociative(preclvl)
	case 0b0100: // repeatable prefix
		return p.ParseRepeatablePrefix(preclvl)
	case 0b0101: // repeatable postfix
		return p.ParseRepeatablePostfix(preclvl)
	case 0b0110: // repeatable infix left associative
		return p.ParseRepeatableLeftAssociative(preclvl)
	case 0b0111: // repeatable infix right associative
		return p.ParseRepeatableRightAssociative(preclvl)
	case 0b1010: // implied operation infix left associative
		return p.ParseImpliedLeftAssociative(preclvl)
	case 0b1011: // implied operation infix right associative
		return p.ParseImpliedRightAssociative(preclvl)
	case 0b1110: // implied operation repeatable infix left associative
		return p.ParseImpliedRepeatableLeftAssociative(preclvl)
	case 0b1111: // implied operation repeatable infix right associative
		return p.ParseImpliedRepeatableRightAssociative(preclvl)
	default:
		panic("invalid configuration")
	}
}

func (p *Parser) ParseImpliedLeftAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseImpliedRepeatableLeftAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseImpliedRightAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseImpliedRepeatableRightAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseLeftAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseRepeatableLeftAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseRightAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseRepeatableRightAssociative(preclvl *PrecedenceLevel) AST
func (p *Parser) ParsePrefix(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseRepeatablePrefix(preclvl *PrecedenceLevel) AST
func (p *Parser) ParsePostfix(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseRepeatablePostfix(preclvl *PrecedenceLevel) AST
func (p *Parser) ParseLeaf() AST
