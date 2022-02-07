package main

type Parser struct {
	tokeniser *Tokeniser
}

func NewParser(tokeniser *Tokeniser) *Parser {
	parser := new(Parser)
	parser.tokeniser = tokeniser
	return parser
}

func (p *Parser) ParseSource() {

}

func (p *Parser) ParseCodeBlock() {

}

func (p *Parser) ParseExpression() {

}
