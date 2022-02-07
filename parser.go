package main

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

func (p *Parser) ParseSource() {
	p.ParseCodeBlock()

	// I realise this is more of an expressions thing,
	// but let's include it just in case.
	if p.tokeniser.opctx.opToken != EOF {
		panic("Unexpected characters at end of source file")
	}
}

func (p *Parser) ParseCodeBlock() {

}

func (p *Parser) ParseExpression() {

}
