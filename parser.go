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

// ParseSource is responsible for parsing an entire source file
func (p *Parser) ParseSource() AST {

	result := p.ParseCodeBlock()

	// I realise this is more of an expression thing,
	// but let's include it just in case.
	if p.tokeniser.currToken != EOF_TOKEN {
		panic("Unexpected characters at end of source file")
	}

	return result
}

// ParseCodeBlock is responsible for parsing a specific code block
func (p *Parser) ParseCodeBlock() AST {
	statements := make([]AST, 0)
	for p.tokeniser.currToken != EOF_TOKEN && p.tokeniser.currToken != CLOSE_CODE_BLOCK_TOKEN {
		nextStatement, _ := p.ParseStatement()
		statements = append(statements, nextStatement)
		if p.tokeniser.currToken == STATEMENT_ENDING_TOKEN {
			p.tokeniser.ReadToken()
		}
	}
	return NewCodeBlock(statements)
}

// ParseStatement is responsible for parsing any arbitrary statement.
// This may be any individual line of code.
func (p *Parser) ParseStatement() (tree AST, parenthesized bool) {
	precedenceListElement := p.opctx.precedenceList.Front()
	return p.ParsePrecedenceLevel(precedenceListElement)
}

func (p *Parser) ParsePrecedenceLevel(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	if precedenceListElement == nil {
		return p.ParseLeaf()
	}
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	// check bitmask in precedence_level.go
	switch precedenceLevel.properties & 0b111 {
	case 0b000: // prefix
		return p.ParsePrefix(precedenceListElement)
	case 0b001: // postfix
		return p.ParsePostfix(precedenceListElement)
	case 0b010: // infix left associative
		return p.ParseLeftAssociative(precedenceListElement)
	case 0b011: // infix right associative
		return p.ParseRightAssociative(precedenceListElement)
	case 0b110: // implied operation infix left associative
		return p.ParseImpliedLeftAssociative(precedenceListElement)
	case 0b111: // implied operation infix right associative
		return p.ParseImpliedRightAssociative(precedenceListElement)
	default:
		panic("invalid configuration")
	}
}

func (p *Parser) ParseImpliedLeftAssociative(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	lhs, lp := p.ParsePrecedenceLevel(precedenceListElement.Next())

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return lhs, lp
	}
	p.tokeniser.ReadToken()

	args := p.getInfixArguments(precedenceListElement.Next(), lhs, opProperties)

	expr := NewStatement(args, opProperties)

	impliedOpProp := precedenceLevel.operators[NIL_TOKEN]

	for {
		opProperties = precedenceLevel.operators[p.tokeniser.currToken]
		if opProperties == nil {
			return expr, false
		}
		p.tokeniser.ReadToken()

		lhs = args[len(args)-1]
		args = p.getInfixArguments(precedenceListElement.Next(), lhs, opProperties)
		expr = NewStatement([]AST{expr, NewStatement(args, opProperties)}, impliedOpProp)
	}
}

func (p *Parser) ParseImpliedRightAssociative(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	lhs, lp := p.ParsePrecedenceLevel(precedenceListElement.Next())

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return lhs, lp
	}
	p.tokeniser.ReadToken()

	args := p.getInfixArguments(precedenceListElement, lhs, opProperties)
	rhs := args[len(args)-1]

	binRight, ok := rhs.(*Statement)

	impliedOpProp := precedenceLevel.operators[NIL_TOKEN]

	// if rhs is an operator
	if ok {
		// if it's an operator in this precedence level
		if precedenceLevel.OperatorExists(binRight.properties) {
			binRightLeft, ok := binRight.terms[0].(*Statement)
			// if the right-hand side is the implied operation (this "ok &&" might cause some unexpected behaviour for weirdly structured trees)
			if ok && binRight.properties == impliedOpProp {
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
			return NewStatement([]AST{NewStatement(args, opProperties), binRight}, impliedOpProp), false
		}
	}
	/*
		  #
		 / \
		a   b
	*/
	return NewStatement([]AST{lhs, rhs}, opProperties), false
}

func (p *Parser) ParseLeftAssociative(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	lhs, lp := p.ParsePrecedenceLevel(precedenceListElement.Next())
	for {
		opProperties := precedenceLevel.operators[p.tokeniser.currToken]

		// if this operator isn't defined:
		if opProperties == nil {

			// we need to check if an implied operation exists for this precedence level
			nilOpProp := precedenceLevel.operators[NIL_TOKEN]

			// we know the next symbol isn't in this precedence level,
			// but still check if it's a control token in case of higher precedence operators.
			if nilOpProp == nil || p.tokeniser.currToken > NIL_TOKEN || p.tokeniser.currToken == STATEMENT_ENDING_TOKEN || p.tokeniser.currToken == EOF_TOKEN {
				return lhs, lp
			} else {
				opProperties = nilOpProp
			}
		} else {
			p.tokeniser.ReadToken()
		}

		// we've parsed the first argument, now we use the operator properties to deduce subsequent symbols to expect
		args := p.getInfixArguments(precedenceListElement.Next(), lhs, opProperties)

		lhs = NewStatement(args, opProperties)
	}
}

func (p *Parser) ParseRightAssociative(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	lhs, lp := p.ParsePrecedenceLevel(precedenceListElement.Next())
	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	// same logic as left associative, just a bit of recursion to get the associativity right
	// private function might be useful, quite a lot of redundancy here
	if opProperties == nil {
		nilOpProp := precedenceLevel.operators[NIL_TOKEN]
		if nilOpProp == nil || p.tokeniser.currToken > NIL_TOKEN || p.tokeniser.currToken == STATEMENT_ENDING_TOKEN || p.tokeniser.currToken == EOF_TOKEN {
			return lhs, lp
		} else {
			opProperties = nilOpProp
		}
	} else {
		p.tokeniser.ReadToken()
	}
	args := p.getInfixArguments(precedenceListElement, lhs, opProperties)

	return NewStatement(args, opProperties), false
}

func (p *Parser) ParsePrefix(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return p.ParsePrecedenceLevel(precedenceListElement.Next())
	}
	argumentCount := opProperties.argumentCount
	argumentSlice := make([]AST, argumentCount)
	p.tokeniser.ReadToken()
	// uh... that works I guess
	for argumentIndex, bit := 0, uint(1); argumentIndex < argumentCount; argumentIndex, bit = argumentIndex+1, bit<<1 {
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
				argument, _ = p.ParseStatement()
			}
		} else {
			argument, _ = p.ParsePrefix(precedenceListElement)
		}
		argumentSlice[argumentIndex] = argument
	}
	return NewStatement(argumentSlice, opProperties), false
}

func (p *Parser) ParsePostfix(precedenceListElement *list.Element) (tree AST, parenthesized bool) {
	precedenceLevel, ok := precedenceListElement.Value.(*PrecedenceLevel)
	if !ok {
		return p.ParseLeaf()
	}
	// stack based parsing
	stack := make(ASTStack, 0)

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]

	// The stack should be the deciding factor for when we complete this precedence level
	for len(stack) != 1 && p.tokeniser.currToken != STATEMENT_ENDING_TOKEN || opProperties != nil {
		// for the symbols we don't recognise here, pass onto higher precedence parsing and add to the stack
		if opProperties == nil {
			arg, _ := p.ParsePrecedenceLevel(precedenceListElement.Next())
			stack.Push(arg)
		} else
		// for the symbols we do recognise, replace top few elements of the stack with the parsed result
		{
			argCount := opProperties.argumentCount
			argumentSlice := make([]AST, argCount)
			for i := argCount - 1; i >= 0; i++ {
				argumentSlice[i] = stack.Pop()
			}
			stack.Push(NewStatement(argumentSlice, opProperties))
			p.tokeniser.ReadToken()
		}
		opProperties = precedenceLevel.operators[p.tokeniser.currToken]
	}
	if len(stack) != 1 {
		panic("unused arguments in postfix expression")
	}
	return stack[0], false
}

func (p *Parser) ParseLeaf() (tree AST, parenthesized bool) {
	var result AST
	parenthesized = false
	switch p.tokeniser.currToken {
	case INT_LITERAL:
		result = IntLiteral{p.tokeniser.intLiteral}
	case FLOAT_LITERAL:
		result = FloatLiteral{p.tokeniser.floatLiteral}
	case CHAR_LITERAL:
		result = CharLiteral{p.tokeniser.charLiteral}
	case STRING_LITERAL:
		result = StringLiteral{p.tokeniser.stringLiteral}
	case TRUE_LITERAL:
		result = BoolLiteral{true}
	case FALSE_LITERAL:
		result = BoolLiteral{false}
	case OPEN_PARENS_TOKEN:
		parenthesized = true
		p.tokeniser.ReadToken()
		result, _ = p.ParseStatement()
		if p.tokeniser.currToken != CLOSE_PARENS_TOKEN {
			panic("missing parentheses")
		}
	case IDENTIFIER_TOKEN:
		result = Identifier{p.tokeniser.identifier}
	}
	p.tokeniser.ReadToken()
	return result, parenthesized
}

func (p *Parser) getInfixArguments(termPrecedenceListElement *list.Element, firstArg AST, opProperties *OpProp) []AST {
	args := make([]AST, 2, len(opProperties.subsequentSymbols)+2)
	args[0] = firstArg
	args[1], _ = p.ParsePrecedenceLevel(termPrecedenceListElement)
	for i := 0; i < len(opProperties.subsequentSymbols); i++ {
		nextSymbol := opProperties.subsequentSymbols[i]
		if p.tokeniser.currToken == nextSymbol || p.tokeniser.currToken < NIL_TOKEN && nextSymbol == NIL_TOKEN {
			if nextSymbol != NIL_TOKEN {
				p.tokeniser.ReadToken()
			}
			nextArg, _ := p.ParsePrecedenceLevel(termPrecedenceListElement)
			args = append(args, nextArg)
		} else {
			panic("Unexpected symbol")
		}
	}
	return args
}
