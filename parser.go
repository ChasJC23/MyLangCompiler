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
	for p.tokeniser.currToken != EOF_TOKEN && p.tokeniser.currToken != CLOSE_CODE_BLOCK_TOKEN {
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
	switch preclvl.properties & 0b111 {
	case 0b000: // prefix
		return p.ParsePrefix(preclvlel)
	case 0b001: // postfix
		return p.ParsePostfix(preclvlel)
	case 0b010: // infix left associative
		return p.ParseLeftAssociative(preclvlel)
	case 0b011: // infix right associative
		return p.ParseRightAssociative(preclvlel)
	case 0b110: // implied operation infix left associative
		return p.ParseImpliedLeftAssociative(preclvlel)
	case 0b111: // implied operation infix right associative
		return p.ParseImpliedRightAssociative(preclvlel)
	default:
		panic("invalid configuration")
	}
}

func (p *Parser) ParseImpliedLeftAssociative(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(*PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	lhs := p.ParsePrecedenceLevel(preclvlel.Next())

	opProperties := preclvl.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return lhs
	}
	p.tokeniser.ReadToken()

	args := p.getInfixArguments(preclvlel.Next(), lhs, opProperties)

	expr := NewStatement(args, opProperties)

	impliedOpProp := preclvl.operators[NIL_TOKEN]

	for {
		opProperties = preclvl.operators[p.tokeniser.currToken]
		if opProperties == nil {
			return expr
		}
		p.tokeniser.ReadToken()

		lhs = args[len(args)-1]
		args = p.getInfixArguments(preclvlel.Next(), lhs, opProperties)
		expr = NewStatement([]AST{expr, NewStatement(args, opProperties)}, impliedOpProp)
	}
}

func (p *Parser) ParseImpliedRightAssociative(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(*PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	lhs := p.ParsePrecedenceLevel(preclvlel.Next())

	opProperties := preclvl.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return lhs
	}
	p.tokeniser.ReadToken()

	args := p.getInfixArguments(preclvlel, lhs, opProperties)
	rhs := args[len(args)-1]

	binRight, err := rhs.(*Statement)

	impliedOpProp := preclvl.operators[NIL_TOKEN]

	// if rhs is an operator
	if !err {
		// if it's an operator in this precedence level
		if preclvl.OperatorExists(binRight.properties) {
			binRightLeft, err := binRight.terms[0].(*Statement)
			// if the right hand side is the implied operation (this "!err &&" might cause some unexpected behaviour for weirdly structured trees)
			if !err && binRight.properties == impliedOpProp {
				/*
						  &
						 / \
						#   &
					   /|  / \
					  / | #  ...
					 /  |/ \ / \
					a   b  ... ...
				*/
				args[len(args)-1] = binRightLeft.terms[0]
			} else {
				/*
						&
					   / \
					  #   #
					 / \ / \
					a   b   c
				*/
				args[len(args)-1] = binRight.terms[0]
			}
			return NewStatement([]AST{NewStatement(args, opProperties), binRight}, impliedOpProp)
		}
	}
	/*
		  #
		 / \
		a   b
	*/
	return NewStatement([]AST{lhs, rhs}, opProperties)
}

func (p *Parser) ParseLeftAssociative(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(*PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	lhs := p.ParsePrecedenceLevel(preclvlel.Next())
	for {
		opProperties := preclvl.operators[p.tokeniser.currToken]

		// if this operator isn't defined:
		if opProperties == nil {

			// we need to check if an implied operation exists for this precedence level
			nilOpProp := preclvl.operators[NIL_TOKEN]

			// we know the next symbol isn't in this precedence level,
			// but still check if it's a control token in case of higher precedence operators.
			if nilOpProp == nil || p.tokeniser.currToken > NIL_TOKEN {
				return lhs
			} else {
				opProperties = nilOpProp
			}
		} else {
			p.tokeniser.ReadToken()
		}

		// we've parsed the first argument, now we use the operator properties to deduce subsequent symbols to expect
		args := p.getInfixArguments(preclvlel.Next(), lhs, opProperties)

		lhs = NewStatement(args, opProperties)
	}
}

func (p *Parser) ParseRightAssociative(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	lhs := p.ParsePrecedenceLevel(preclvlel.Next())
	opProperties := preclvl.operators[p.tokeniser.currToken]
	// same logic as left associative, just a bit of recursion to get the associativity right
	// private function might be useful, quite a lot of redundancy here
	if opProperties == nil {
		nilOpProp := preclvl.operators[NIL_TOKEN]
		if nilOpProp == nil || p.tokeniser.currToken > NIL_TOKEN {
			return lhs
		} else {
			opProperties = nilOpProp
		}
	} else {
		p.tokeniser.ReadToken()
	}
	args := p.getInfixArguments(preclvlel, lhs, opProperties)

	return NewStatement(args, opProperties)
}

func (p *Parser) ParsePrefix(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(*PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	opProperties := preclvl.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return p.ParsePrecedenceLevel(preclvlel.Next())
	}
	argumentCount := opProperties.argumentCount
	argumentSlice := make([]AST, argumentCount)
	// uh... that works I guess
	for argumentIndex, bit := 0, uint(1); argumentIndex < argumentCount; argumentIndex, bit = argumentIndex+1, bit<<1 {
		p.tokeniser.ReadToken()
		isCodeBlock := opProperties.codeBlockArguments&bit != 0
		var argument AST
		if isCodeBlock {
			if p.tokeniser.currToken == OPEN_CODE_BLOCK_TOKEN {
				p.tokeniser.ReadToken()
				argument = p.ParseCodeBlock()
				if p.tokeniser.currToken != CLOSE_CODE_BLOCK_TOKEN {
					panic("missing close code block symbol")
				}
				p.tokeniser.ReadToken()
			} else {
				argument = p.ParseStatement()
			}
		} else {
			argument = p.ParsePrefix(preclvlel)
		}
		argumentSlice[argumentIndex] = argument
	}
	return NewStatement(argumentSlice, opProperties)
}

func (p *Parser) ParsePostfix(preclvlel *list.Element) AST {
	preclvl, err := preclvlel.Value.(*PrecedenceLevel)
	if err {
		return p.ParseLeaf()
	}
	// stack based parsing
	stack := make(ASTStack, 0)

	// The stack should be the deciding factor for when we complete this precedence level
	for len(stack) != 1 {
		opProperties := preclvl.operators[p.tokeniser.currToken]
		// for symbols we don't recognise here, pass onto higher precedence parsing and add to the stack
		if opProperties == nil {
			stack.Push(p.ParsePrecedenceLevel(preclvlel.Next()))
		} else
		// for symbols we do recognise, replace top few elements of the stack with the parsed result
		{
			argCount := opProperties.argumentCount
			argumentSlice := make([]AST, argCount)
			for i := argCount - 1; i >= 0; i++ {
				argumentSlice[i] = stack.Pop()
			}
			stack.Push(NewStatement(argumentSlice, opProperties))
		}
		p.tokeniser.ReadToken()
	}
	return stack[0]
}

func (p *Parser) ParseLeaf() AST {
	var result AST
	switch p.tokeniser.currToken {
	case INT_LITERAL:
		result = IntLiteral{p.tokeniser.intLiteral}
	case FLOAT_LITERAL:
		result = FloatLiteral{p.tokeniser.floatLiteral}
	case CHAR_LITERAL:
		result = CharLiteral{p.tokeniser.charLiteral}
	case STRING_LITERAL:
		result = StringLiteral{p.tokeniser.stringLiteral}
	case OPEN_PARENS_TOKEN:
		p.tokeniser.ReadToken()
		result = p.ParseStatement()
		if p.tokeniser.currToken != CLOSE_PARENS_TOKEN {
			panic("missing parentheses")
		}
	case IDENTIFIER_TOKEN:
		result = Identifier{p.tokeniser.identifier}
	}
	p.tokeniser.ReadToken()
	return result
}

func (p *Parser) getInfixArguments(term_preclvlel *list.Element, firstArg AST, opProperties *OpProp) []AST {
	args := make([]AST, 2, len(opProperties.subsequentSymbols)+1)
	args[0] = firstArg
	args[1] = p.ParsePrecedenceLevel(term_preclvlel)
	for i := 0; i < len(opProperties.subsequentSymbols); i++ {
		nextSymbol := opProperties.subsequentSymbols[i]
		if p.tokeniser.currToken == nextSymbol || p.tokeniser.currToken < NIL_TOKEN && nextSymbol == NIL_TOKEN {
			if nextSymbol != NIL_TOKEN {
				p.tokeniser.ReadToken()
			}
			args = append(args, p.ParsePrecedenceLevel(term_preclvlel))
		} else {
			panic("Unexpected symbol")
		}
	}
	return args
}
